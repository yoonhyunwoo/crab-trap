package moltbook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	apiKey := "test_api_key"
	client := NewClient(apiKey)

	assert.NotNil(t, client)
	assert.Equal(t, apiKey, client.apiKey)
	assert.NotNil(t, client.httpClient)
	assert.NotNil(t, client.baseURL)
	assert.Equal(t, "https://www.moltbook.com/api/v1", client.baseURL.String())
}

func TestNewClientWithCustomBaseURL(t *testing.T) {
	apiKey := "test_api_key"
	customURL := "https://custom.api.com/v1"
	client := NewClient(apiKey).WithBaseURL(customURL)

	assert.NotNil(t, client)
	assert.Equal(t, customURL, client.baseURL.String())
}

func TestBuildEndpoint(t *testing.T) {
	client := NewClient("test_api_key")

	tests := []struct {
		name     string
		parts    []string
		expected string
	}{
		{"single", []string{"posts"}, "https://www.moltbook.com/api/v1/posts"},
		{"multiple", []string{"posts", "123", "comments"}, "https://www.moltbook.com/api/v1/posts/123/comments"},
		{"nested", []string{"agents", "me"}, "https://www.moltbook.com/api/v1/agents/me"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.buildEndpoint(tt.parts...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  APIError
		want string
	}{
		{
			name: "with error field",
			err: APIError{StatusCode: 400, ErrorMessage: "Bad request"},
			want: "API error (status 400): Bad request",
		},
		{
			name: "with message",
			err: APIError{StatusCode: 404, Message: "Not found"},
			want: "HTTP error (status 404): Not found",
		},
		{
			name: "minimal",
			err:  APIError{StatusCode: 500},
			want: "HTTP error (status 500)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestAPIError_IsRateLimit(t *testing.T) {
	err := APIError{StatusCode: 429}
	assert.True(t, err.IsRateLimit())

	err = APIError{StatusCode: 200}
	assert.False(t, err.IsRateLimit())
}

func TestAPIError_IsUnauthorized(t *testing.T) {
	err := APIError{StatusCode: 401}
	assert.True(t, err.IsUnauthorized())

	err = APIError{StatusCode: 200}
	assert.False(t, err.IsUnauthorized())
}

func TestAPIError_IsNotFound(t *testing.T) {
	err := APIError{StatusCode: 404}
	assert.True(t, err.IsNotFound())

	err = APIError{StatusCode: 200}
	assert.False(t, err.IsNotFound())
}

func TestRateLimitError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  RateLimitError
		want string
	}{
		{
			name: "with minutes",
			err:  RateLimitError{RetryAfterMinutes: 5},
			want: "rate limit: please wait 5 minutes before posting again",
		},
		{
			name: "with seconds",
			err:  RateLimitError{RetryAfterSeconds: 20},
			want: "rate limit: please wait 20 seconds before commenting again",
		},
		{
			name: "default",
			err:  RateLimitError{APIError: APIError{ErrorMessage: "rate limited"}},
			want: "API error (status 0): rate limited",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			assert.Contains(t, result, tt.want)
		})
	}
}

func TestStructToURLValues(t *testing.T) {
	type TestStruct struct {
		Name   string   `json:"name"`
		Age    int      `json:"age"`
		Active bool     `json:"active"`
		Tags   []string `json:"tags"`
		Ignore string   `json:"-"`
	}

	tests := []struct {
		name     string
		input    interface{}
		expected map[string]string
	}{
		{
			name: "full struct",
			input: TestStruct{
				Name:   "test",
				Age:    25,
				Active: true,
				Tags:   []string{"a", "b"},
			},
			expected: map[string]string{
				"name":   "test",
				"age":    "25",
				"active": "true",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := structToURLValues(tt.input)
			assert.NoError(t, err)

			for k, v := range tt.expected {
				assert.Equal(t, v, result.Get(k))
			}
		})
	}
}

func TestResponse_UnmarshalData(t *testing.T) {
	type TestType struct {
		Name string `json:"name"`
	}

	tests := []struct {
		name    string
		resp    Response
		want    TestType
		wantErr bool
	}{
		{
			name: "success",
			resp: Response{
				Success: true,
				Data: map[string]interface{}{
					"name": "test",
				},
			},
			want:    TestType{Name: "test"},
			wantErr: false,
		},
		{
			name: "no data",
			resp: Response{
				Success: false,
				ErrorMessage:   "some error",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result TestType
			err := tt.resp.UnmarshalData(&result)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}
		})
	}
}
