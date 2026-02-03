package moltbook

import "strconv"

func (c *Client) CreateSubmolt(req CreateSubmoltRequest) (*SubmoltResponse, error) {
	resp, err := c.doRequestWithJSON("POST", c.buildEndpoint("submolts"), req)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, APIError{
			StatusCode: resp.StatusCode,
			ErrorMessage:      resp.ErrorMessage,
			Hint:       resp.Hint,
		}
	}

	var result SubmoltResponse
	if err := resp.UnmarshalData(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetSubmolts() (*SubmoltsResponse, error) {
	resp, err := c.doRequest("GET", c.buildEndpoint("submolts"), nil)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, APIError{
			StatusCode: resp.StatusCode,
			ErrorMessage:      resp.ErrorMessage,
			Hint:       resp.Hint,
		}
	}

	var result SubmoltsResponse
	if err := resp.UnmarshalData(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetSubmolt(submoltName string) (*Submolt, error) {
	resp, err := c.doRequest("GET", c.buildEndpoint("submolts", submoltName), nil)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, APIError{
			StatusCode: resp.StatusCode,
			ErrorMessage:      resp.ErrorMessage,
			Hint:       resp.Hint,
		}
	}

	var result Submolt
	if err := resp.UnmarshalData(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetSubmoltFeed(submoltName string, opts GetPostsOptions) (*PostsResponse, error) {
	endpoint := c.buildEndpoint("submolts", submoltName, "feed")

	if opts.Sort != "" {
		endpoint += "?sort=" + opts.Sort
	}

	if opts.Limit > 0 {
		if opts.Sort != "" {
			endpoint += "&limit=" + strconv.Itoa(opts.Limit)
		} else {
			endpoint += "?limit=" + strconv.Itoa(opts.Limit)
		}
	}

	resp, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, APIError{
			StatusCode: resp.StatusCode,
			ErrorMessage:      resp.ErrorMessage,
			Hint:       resp.Hint,
		}
	}

	var result PostsResponse
	if err := resp.UnmarshalData(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) Subscribe(submoltName string) error {
	resp, err := c.doRequest("POST", c.buildEndpoint("submolts", submoltName, "subscribe"), nil)
	if err != nil {
		return err
	}

	if !resp.Success {
		return APIError{
			StatusCode: resp.StatusCode,
			ErrorMessage:      resp.ErrorMessage,
			Hint:       resp.Hint,
		}
	}

	return nil
}

func (c *Client) Unsubscribe(submoltName string) error {
	resp, err := c.doRequest("DELETE", c.buildEndpoint("submolts", submoltName, "subscribe"), nil)
	if err != nil {
		return err
	}

	if !resp.Success {
		return APIError{
			StatusCode: resp.StatusCode,
			ErrorMessage:      resp.ErrorMessage,
			Hint:       resp.Hint,
		}
	}

	return nil
}

func (c *Client) UpdateSubmoltSettings(submoltName string, req UpdateSubmoltRequest) (*Submolt, error) {
	resp, err := c.doRequestWithJSON("PATCH", c.buildEndpoint("submolts", submoltName, "settings"), req)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, APIError{
			StatusCode: resp.StatusCode,
			ErrorMessage:      resp.ErrorMessage,
			Hint:       resp.Hint,
		}
	}

	var result Submolt
	if err := resp.UnmarshalData(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) UploadSubmoltAvatar(submoltName, filePath string) (*Submolt, error) {
	files := map[string]string{"file": filePath}
	fields := map[string]string{"type": "avatar"}

	resp, err := c.doRequestWithMultipart("POST", c.buildEndpoint("submolts", submoltName, "settings"), files, fields)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, APIError{
			StatusCode: resp.StatusCode,
			ErrorMessage:      resp.ErrorMessage,
			Hint:       resp.Hint,
		}
	}

	var result Submolt
	if err := resp.UnmarshalData(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) UploadSubmoltBanner(submoltName, filePath string) (*Submolt, error) {
	files := map[string]string{"file": filePath}
	fields := map[string]string{"type": "banner"}

	resp, err := c.doRequestWithMultipart("POST", c.buildEndpoint("submolts", submoltName, "settings"), files, fields)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, APIError{
			StatusCode: resp.StatusCode,
			ErrorMessage:      resp.ErrorMessage,
			Hint:       resp.Hint,
		}
	}

	var result Submolt
	if err := resp.UnmarshalData(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) AddModerator(submoltName string, req AddModeratorRequest) error {
	resp, err := c.doRequestWithJSON("POST", c.buildEndpoint("submolts", submoltName, "moderators"), req)
	if err != nil {
		return err
	}

	if !resp.Success {
		return APIError{
			StatusCode: resp.StatusCode,
			ErrorMessage:      resp.ErrorMessage,
			Hint:       resp.Hint,
		}
	}

	return nil
}

func (c *Client) RemoveModerator(submoltName string, req RemoveModeratorRequest) error {
	resp, err := c.doRequestWithJSON("DELETE", c.buildEndpoint("submolts", submoltName, "moderators"), req)
	if err != nil {
		return err
	}

	if !resp.Success {
		return APIError{
			StatusCode: resp.StatusCode,
			ErrorMessage:      resp.ErrorMessage,
			Hint:       resp.Hint,
		}
	}

	return nil
}

func (c *Client) GetModerators(submoltName string) (*ModeratorsResponse, error) {
	resp, err := c.doRequest("GET", c.buildEndpoint("submolts", submoltName, "moderators"), nil)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, APIError{
			StatusCode: resp.StatusCode,
			ErrorMessage:      resp.ErrorMessage,
			Hint:       resp.Hint,
		}
	}

	var result ModeratorsResponse
	if err := resp.UnmarshalData(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
