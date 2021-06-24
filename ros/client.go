// Package ros implements a Mikrotik RouterOS 7 REST API client. It's a very
// thin wrapper, delivering no high-level logic upon the basic CRUD capability
// of the API.
//
// For more information about the API see:
// https://help.mikrotik.com/docs/display/ROS/REST+API
package ros

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

// Client is a ROS7 REST API client. It connects to a ROS www-ssl service
// at the given Address, authenticating using the given Username and Password.
type Client struct {
	// Address (like foo.example.com or 1.2.3.4:1234 or [2a0d:eb01::1]:443) of
	// ROS7 www-ssl service.
	Address string
	// Username used to authenticate to ROS API.
	Username string
	// Password used to authenticate to ROS API.
	Password string
	// HTTP client used for connections. If not set, uses http.DefaultClient.
	// If connecting to a ROS7 device whose certificate was generated via
	// built-in Let's Encrypt support, this should be set to LetsEncryptClient
	// to add trust for the Let's Encrypt R3 CA.
	HTTP *http.Client
}

func (c *Client) urlFor(path string) string {
	return fmt.Sprintf("https://%s:%s@%s/rest/%s", c.Username, c.Password, c.Address, path)
}

func (c *Client) httpClient() *http.Client {
	if c.HTTP == nil {
		return http.DefaultClient
	}
	return c.HTTP
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
