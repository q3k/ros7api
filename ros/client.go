package ros

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	Username string
	Password string
	Address  string
	Client   *http.Client
}

func (c *Client) doGET(ctx context.Context, path string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://%s:%s@%s/rest/%s", c.Username, c.Password, c.Address, path), nil)
	if err != nil {
		return nil, fmt.Errorf("GET: %w", err)
	}
	cl := c.Client
	if cl == nil {
		cl = http.DefaultClient
	}
	resp, err := cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Do: %w", err)
	}
	return resp.Body, nil
}
