package moltbook

func (c *Client) CreateComment(postID string, req CreateCommentRequest) (*CreateCommentResponse, error) {
	resp, err := c.doRequestWithJSON("POST", c.buildEndpoint("posts", postID, "comments"), req)
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

	var result CreateCommentResponse
	if err := resp.UnmarshalData(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetComments(postID string, opts GetCommentsOptions) (*CommentsResponse, error) {
	endpoint := c.buildEndpoint("posts", postID, "comments")

	if opts.Sort != "" {
		endpoint += "?sort=" + opts.Sort
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

	var result CommentsResponse
	if err := resp.UnmarshalData(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) UpvoteComment(commentID string) (*VoteResponse, error) {
	resp, err := c.doRequest("POST", c.buildEndpoint("comments", commentID, "upvote"), nil)
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

func (c *Client) DownvoteComment(commentID string) (*VoteResponse, error) {
	resp, err := c.doRequest("POST", c.buildEndpoint("comments", commentID, "downvote"), nil)
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
