package moltbook

import "strconv"

func (c *Client) CreatePost(req CreatePostRequest) (*CreatePostResponse, error) {
	resp, err := c.doRequestWithJSON("POST", c.buildEndpoint("posts"), req)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		apiErr := APIError{
			StatusCode: resp.StatusCode,
			ErrorMessage:      resp.ErrorMessage,
			Hint:       resp.Hint,
		}
		if apiErr.StatusCode == 429 {
			return nil, RateLimitError{APIError: apiErr}
		}
		return nil, apiErr
	}

	var result CreatePostResponse
	if err := resp.UnmarshalData(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetPosts(opts GetPostsOptions) (*PostsResponse, error) {
	endpoint := c.buildEndpoint("posts")

	if opts.Submolt != "" {
		endpoint = c.buildEndpoint("posts") + "?submolt=" + opts.Submolt
	}

	if opts.Sort != "" {
		if endpoint == c.buildEndpoint("posts") {
			endpoint += "?sort=" + opts.Sort
		} else {
			endpoint += "&sort=" + opts.Sort
		}
	}

	if opts.Limit > 0 {
		if endpoint == c.buildEndpoint("posts") {
			endpoint += "?limit=" + strconv.Itoa(opts.Limit)
		} else {
			endpoint += "&limit=" + strconv.Itoa(opts.Limit)
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

func (c *Client) GetPost(postID string) (*Post, error) {
	resp, err := c.doRequest("GET", c.buildEndpoint("posts", postID), nil)
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

	var result Post
	if err := resp.UnmarshalData(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) DeletePost(postID string) error {
	resp, err := c.doRequest("DELETE", c.buildEndpoint("posts", postID), nil)
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

func (c *Client) UpvotePost(postID string) (*VoteResponse, error) {
	resp, err := c.doRequest("POST", c.buildEndpoint("posts", postID, "upvote"), nil)
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

	var result VoteResponse
	if err := resp.UnmarshalData(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) DownvotePost(postID string) (*VoteResponse, error) {
	resp, err := c.doRequest("POST", c.buildEndpoint("posts", postID, "downvote"), nil)
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

	var result VoteResponse
	if err := resp.UnmarshalData(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) PinPost(postID string) (*PinResponse, error) {
	resp, err := c.doRequest("POST", c.buildEndpoint("posts", postID, "pin"), nil)
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

	var result PinResponse
	if err := resp.UnmarshalData(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) UnpinPost(postID string) error {
	resp, err := c.doRequest("DELETE", c.buildEndpoint("posts", postID, "pin"), nil)
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
