package moltbook

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrUnauthorized       = errors.New("unauthorized: invalid API key")
	ErrRateLimit         = errors.New("rate limit exceeded")
	ErrNotFound          = errors.New("resource not found")
	ErrInvalidRequest    = errors.New("invalid request")
	ErrPostCooldown      = errors.New("post cooldown active")
	ErrCommentCooldown   = errors.New("comment cooldown active")
	ErrBadRequest        = errors.New("bad request")
	ErrInternalServer    = errors.New("internal server error")
	ErrAlreadyFollowing  = errors.New("already following")
	ErrNotFollowing      = errors.New("not following")
	ErrAlreadySubscribed = errors.New("already subscribed")
	ErrNotSubscribed     = errors.New("not subscribed")
	ErrTooManyPins       = errors.New("too many pinned posts")
	ErrMaxSizeExceeded   = errors.New("file size exceeded")
)

type APIError struct {
	StatusCode int
	Message    string
	Success    bool   `json:"success"`
	ErrorMessage string `json:"error"`
	Hint       string `json:"hint"`
}

func (e APIError) Error() string {
	if e.ErrorMessage != "" {
		return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.ErrorMessage)
	}
	if e.Message != "" {
		return fmt.Sprintf("HTTP error (status %d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("HTTP error (status %d)", e.StatusCode)
}

func (e APIError) IsRateLimit() bool {
	return e.StatusCode == http.StatusTooManyRequests
}

func (e APIError) IsUnauthorized() bool {
	return e.StatusCode == http.StatusUnauthorized
}

func (e APIError) IsNotFound() bool {
	return e.StatusCode == http.StatusNotFound
}

type RateLimitError struct {
	APIError
	RetryAfterSeconds int `json:"retry_after_seconds"`
	RetryAfterMinutes int `json:"retry_after_minutes"`
	DailyRemaining    int `json:"daily_remaining"`
}

func (e RateLimitError) Error() string {
	if e.RetryAfterMinutes > 0 {
		return fmt.Sprintf("rate limit: please wait %d minutes before posting again", e.RetryAfterMinutes)
	}
	if e.RetryAfterSeconds > 0 {
		return fmt.Sprintf("rate limit: please wait %d seconds before commenting again", e.RetryAfterSeconds)
	}
	return e.APIError.Error()
}

type Response struct {
	StatusCode int          `json:"-"`
	Success    bool         `json:"success"`
	Data       interface{}  `json:"data"`
	ErrorMessage string       `json:"error"`
	Hint       string       `json:"hint"`
	Body       []byte       `json:"-"`
}

func (r *Response) UnmarshalData(v interface{}) error {
	if !r.Success {
		if r.ErrorMessage != "" {
			return fmt.Errorf("API error: %s", r.ErrorMessage)
		}
		if r.Hint != "" {
			return fmt.Errorf("API error: %s", r.Hint)
		}
		return fmt.Errorf("API error (status %d)", r.StatusCode)
	}

	if r.Data == nil {
		return fmt.Errorf("response success but no data")
	}

	dataBytes, err := json.Marshal(r.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := json.Unmarshal(dataBytes, v); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return nil
}

func structToURLValues(data interface{}) (url.Values, error) {
	values := url.Values{}
	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return values, nil
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		tagName := strings.Split(jsonTag, ",")[0]
		if tagName == "" {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			if field.String() != "" {
				values.Set(tagName, field.String())
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			values.Set(tagName, strconv.FormatInt(field.Int(), 10))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			values.Set(tagName, strconv.FormatUint(field.Uint(), 10))
		case reflect.Float32, reflect.Float64:
			values.Set(tagName, strconv.FormatFloat(field.Float(), 'f', -1, 64))
		case reflect.Bool:
			values.Set(tagName, strconv.FormatBool(field.Bool()))
		case reflect.Slice:
			if field.Type().Elem().Kind() == reflect.String {
				slice := field.Interface().([]string)
				for _, item := range slice {
					values.Add(tagName, item)
				}
			}
		}
	}

	return values, nil
}

func openFile(filePath string) (io.ReadCloser, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}
