package client

import (
	"fmt"
	"net/http"

	"github.com/tocoteron/gmo-coin-go/api/v1/auth"
	"github.com/tocoteron/gmo-coin-go/api/v1/http/endpoint/private/account/margin"
	"github.com/tocoteron/gmo-coin-go/api/v1/http/endpoint/public/status"
)

const (
	API_HOST = "https://api.coin.z.com"
)

type Client struct {
	Host       string
	HTTPClient *http.Client
	AuthConfig *auth.AuthConfig
}

type ClientOpts struct {
	AuthConfig *auth.AuthConfig
}

func NewClient(opts *ClientOpts) *Client {
	return &Client{Host: API_HOST, HTTPClient: &http.Client{}, AuthConfig: opts.AuthConfig}
}

// Public API

func (c *Client) Status() (*status.Status, *http.Response, error) {
	entity, httpResp, err := requestPublic(c, *status.Endpoint, nil, nil)
	if err != nil {
		return nil, httpResp, fmt.Errorf("failed to call status API: %w", err)
	}
	return entity, httpResp, nil
}

// Private API

func (c *Client) Margin() (*margin.Margin, *http.Response, error) {
	entity, httpResp, err := requestPrivate(c, *margin.Endpoint, nil, nil)
	if err != nil {
		return nil, httpResp, fmt.Errorf("failed to call margin API: %w", err)
	}
	return entity, httpResp, nil
}
