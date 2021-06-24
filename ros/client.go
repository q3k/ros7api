package ros

import (
	"bytes"
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

func (c *Client) urlFor(path string) string {
	return fmt.Sprintf("https://%s:%s@%s/rest/%s", c.Username, c.Password, c.Address, path)
}

func (c *Client) httpClient() *http.Client {
	if c.Client == nil {
		c.Client = http.DefaultClient
	}
	return c.Client
}

func (c *Client) doGET(ctx context.Context, path string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.urlFor(path), nil)
	if err != nil {
		return nil, fmt.Errorf("could not make GET request: %w", err)
	}
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("when running REST request: %w", err)
	}
	return resp.Body, nil
}

func (c *Client) doPATCH(ctx context.Context, path string, rdata []byte) (io.ReadCloser, error) {
	rbuf := bytes.NewBuffer(rdata)
	req, err := http.NewRequestWithContext(ctx, "PATCH", c.urlFor(path), rbuf)
	if err != nil {
		return nil, fmt.Errorf("could not make GET request: %w", err)
	}
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("when running REST request: %w", err)
	}
	return resp.Body, nil
}
