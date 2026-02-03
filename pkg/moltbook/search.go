package moltbook

import "strconv"

func (c *Client) Search(query string, opts SearchRequest) (*SearchResponse, error) {
	endpoint := c.buildEndpoint("search") + "?q=" + query

	if opts.Type != "" {
		endpoint += "&type=" + opts.Type
	}

	if opts.Limit > 0 {
		endpoint += "&limit=" + strconv.Itoa(opts.Limit)
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

	var result SearchResponse
	if err := resp.UnmarshalData(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
