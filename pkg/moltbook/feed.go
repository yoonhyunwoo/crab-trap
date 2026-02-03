package moltbook

import "strconv"

func (c *Client) GetFeed(opts FeedOptions) (*FeedResponse, error) {
	endpoint := c.buildEndpoint("feed")

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

	var result FeedResponse
	if err := resp.UnmarshalData(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
