package moltbook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

const (
	defaultBaseURL    = "https://www.moltbook.com/api/v1"
	defaultUserAgent  = "Moltbook-Go-SDK/1.0.0"
	defaultHTTPClient = 30 * time.Second
)

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	apiKey     string
}

type RequestOptions struct {
	Body   io.Reader
	Header http.Header
}

func NewClient(apiKey string) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: defaultHTTPClient},
		apiKey:     apiKey,
	}
}

func (c *Client) WithBaseURL(baseURL string) *Client {
	u, err := url.Parse(baseURL)
	if err != nil {
		return c
	}
	c.baseURL = u
	return c
}

func (c *Client) WithHTTPClient(httpClient *http.Client) *Client {
	c.httpClient = httpClient
	return c
}

func (c *Client) buildEndpoint(parts ...string) string {
	return c.baseURL.String() + "/" + path.Join(parts...)
}

func (c *Client) doRequest(method, endpoint string, opts *RequestOptions) (*Response, error) {
	req, err := c.newRequest(method, endpoint, opts)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	return c.parseResponse(resp)
}

func (c *Client) newRequest(method, endpoint string, opts *RequestOptions) (*http.Request, error) {
	var body io.Reader
	if opts != nil && opts.Body != nil {
		body = opts.Body
	}

	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", defaultUserAgent)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	if opts != nil && opts.Header != nil {
		for k, v := range opts.Header {
			req.Header[k] = v
		}
	}

	return req, nil
}

func (c *Client) doRequestWithJSON(method, endpoint string, data interface{}) (*Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	opts := &RequestOptions{
		Body:   bytes.NewBuffer(jsonData),
		Header: http.Header{"Content-Type": []string{"application/json"}},
	}

	return c.doRequest(method, endpoint, opts)
}

func (c *Client) doRequestWithFormData(method, endpoint string, data interface{}) (*Response, error) {
	values, err := structToURLValues(data)
	if err != nil {
		return nil, err
	}

	opts := &RequestOptions{
		Body:   strings.NewReader(values.Encode()),
		Header: http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}},
	}

	return c.doRequest(method, endpoint, opts)
}

func (c *Client) doRequestWithMultipart(method, endpoint string, files map[string]string, fields map[string]string) (*Response, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for fieldName, filePath := range files {
		file, err := openFile(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		part, err := writer.CreateFormFile(fieldName, filePath)
		if err != nil {
			return nil, err
		}

		if _, err := io.Copy(part, file); err != nil {
			return nil, err
		}
	}

	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			return nil, err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	opts := &RequestOptions{
		Body:   body,
		Header: http.Header{"Content-Type": []string{writer.FormDataContentType()}},
	}

	return c.doRequest(method, endpoint, opts)
}

func (c *Client) parseResponse(resp *http.Response) (*Response, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &Response{
		StatusCode: resp.StatusCode,
		Body:       body,
	}

	if err := json.Unmarshal(body, &response); err == nil {
		return response, nil
	}

	if resp.StatusCode >= 400 {
		return response, APIError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	return response, nil
}
