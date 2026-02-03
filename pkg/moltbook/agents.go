package moltbook

func Register(name, description string) (*RegisterResponse, error) {
	client := NewClient("")

	req := RegisterRequest{
		Name:        name,
		Description: description,
	}

	resp, err := client.doRequestWithJSON("POST", client.buildEndpoint("agents", "register"), req)
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

	var result RegisterResponse
	if err := resp.UnmarshalData(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetMe() (*Agent, error) {
	resp, err := c.doRequest("GET", c.buildEndpoint("agents", "me"), nil)
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

	var result Agent
	if err := resp.UnmarshalData(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetProfile(name string) (*AgentProfileResponse, error) {
	endpoint := c.buildEndpoint("agents", "profile") + "?name=" + name

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

	var result AgentProfileResponse
	if err := resp.UnmarshalData(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) UpdateProfile(req UpdateAgentRequest) (*Agent, error) {
	resp, err := c.doRequestWithJSON("PATCH", c.buildEndpoint("agents", "me"), req)
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

	var result Agent
	if err := resp.UnmarshalData(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) UploadAvatar(filePath string) (*Agent, error) {
	files := map[string]string{"file": filePath}
	fields := map[string]string{}

	resp, err := c.doRequestWithMultipart("POST", c.buildEndpoint("agents", "me", "avatar"), files, fields)
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

	var result Agent
	if err := resp.UnmarshalData(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) DeleteAvatar() error {
	resp, err := c.doRequest("DELETE", c.buildEndpoint("agents", "me", "avatar"), nil)
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

func (c *Client) GetStatus() (*StatusResponse, error) {
	resp, err := c.doRequest("GET", c.buildEndpoint("agents", "status"), nil)
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

	var result StatusResponse
	if err := resp.UnmarshalData(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) Follow(moltyName string) (*FollowResponse, error) {
	resp, err := c.doRequest("POST", c.buildEndpoint("agents", moltyName, "follow"), nil)
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

	var result FollowResponse
	if err := resp.UnmarshalData(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) Unfollow(moltyName string) error {
	resp, err := c.doRequest("DELETE", c.buildEndpoint("agents", moltyName, "follow"), nil)
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
