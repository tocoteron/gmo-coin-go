package client

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/tocoteron/gmo-coin-go/api/v1/auth"
	"github.com/tocoteron/gmo-coin-go/api/v1/http/endpoint"
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

func requestPublic[P, B, E any](c *Client, e endpoint.Endpoint[P, B, E], p *P, b *B) (*E, *http.Response, error) {
	// Create request
	req, err := createReq(c, e.Method, "/public"+e.Path, p, b)
	if err != nil {
		return nil, nil, err
	}

	// Send request
	entity, resp, err := sendReq[E](c, req)
	if err != nil {
		return nil, nil, err
	}

	return entity, resp, nil
}

func requestPrivate[P, B, E any](c *Client, e endpoint.Endpoint[P, B, E], p *P, b *B) (*E, *http.Response, error) {
	// Check auth config
	if c.AuthConfig == nil {
		return nil, nil, errors.New("failed to request because auth config is required for private API")
	}

	// Create request
	req, err := createReq(c, e.Method, "/private"+e.Path, p, b)
	if err != nil {
		return nil, nil, err
	}

	// Sign request
	{
		// Get request body if it exists
		var reqBody []byte
		if req.Body != nil {
			reqBody, err = io.ReadAll(req.Body)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to read request body: %w", err)
			}
		}

		// Sign
		err = c.sign(req, e.Method, e.Path, reqBody)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to sign http request: %w", err)
		}
	}

	// Send request
	entity, resp, err := sendReq[E](c, req)
	if err != nil {
		return nil, nil, err
	}

	return entity, resp, nil
}

func createReq[P, B any](c *Client, method, path string, params *P, body *B) (*http.Request, error) {
	// Build endpoint
	endpoint, err := url.Parse(c.Host + path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url %s: %w", endpoint, err)
	}

	// Set query params
	query, err := query.Values(params)
	if err != nil {
		return nil, fmt.Errorf("failed to parse query: %w", err)
	}
	endpoint.RawQuery = query.Encode()

	// Set request body if it exists
	if body != nil {
		// Marshal request body to json
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body of http request: %w", err)
		}

		// Create new request
		req, err := http.NewRequest(method, endpoint.String(), bytes.NewReader(jsonBody))
		if err != nil {
			return nil, fmt.Errorf("failed to create request to %s: %w", endpoint, err)
		}
		return req, nil
	}

	// Create new request
	req, err := http.NewRequest(method, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to %s: %w", endpoint, err)
	}
	return req, nil
}

func sendReq[E any](c *Client, req *http.Request) (*E, *http.Response, error) {
	// Request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to request to %s: %w", req.URL.String(), err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to read body of http response: %w", err)
	}

	// Unmarshal response body to data structure
	var entity E
	err = json.Unmarshal(respBody, &entity)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to unmarshal body of http reponse: %w", err)
	}

	return &entity, resp, nil
}

func (c *Client) sign(req *http.Request, method, path string, body []byte) error {
	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	text := timestamp + method + path + string(body)

	hc := hmac.New(sha256.New, []byte(c.AuthConfig.APISecret))
	_, err := hc.Write([]byte(text))
	if err != nil {
		return err
	}

	sign := hex.EncodeToString(hc.Sum(nil))

	// Set headers
	req.Header.Set("API-KEY", c.AuthConfig.APIKey)
	req.Header.Set("API-TIMESTAMP", timestamp)
	req.Header.Set("API-SIGN", sign)

	return nil
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
