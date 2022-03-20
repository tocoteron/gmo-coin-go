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
	"github.com/tocoteron/gmo-coin-go/api/v1/http/private/account/margin"
	"github.com/tocoteron/gmo-coin-go/api/v1/http/public/status"
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

func (c *Client) requestPublic(method, path string, params, body, v interface{}) (*http.Response, error) {
	// Create request
	req, err := c.createReq(method, "/public"+path, params, body)
	if err != nil {
		return nil, err
	}

	// Send request
	resp, err := c.sendReq(req, v)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) requestPrivate(method, path string, params, body, v interface{}) (*http.Response, error) {
	// Check auth config
	if c.AuthConfig == nil {
		return nil, errors.New("failed to request because auth config is required for private API")
	}

	// Create request
	req, err := c.createReq(method, "/private"+path, params, body)
	if err != nil {
		return nil, err
	}

	// Sign request
	{
		// Get request body if it exists
		var reqBody []byte
		if req.Body != nil {
			reqBody, err = io.ReadAll(req.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read request body: %w", err)
			}
		}

		// Sign
		err = c.sign(req, method, path, reqBody)
		if err != nil {
			return nil, fmt.Errorf("failed to sign http request: %w", err)
		}
	}

	// Send request
	resp, err := c.sendReq(req, v)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) createReq(method, path string, params, body interface{}) (*http.Request, error) {
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

func (c *Client) sendReq(req *http.Request, v interface{}) (*http.Response, error) {
	// Request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return resp, fmt.Errorf("failed to request to %s: %w", req.URL.String(), err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, fmt.Errorf("failed to read body of http response: %w", err)
	}

	// Unmarshal response body to data structure
	err = json.Unmarshal(respBody, v)
	if err != nil {
		return resp, fmt.Errorf("failed to unmarshal body of http reponse: %w", err)
	}

	return resp, nil
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
	resp := &status.Status{}

	httpResp, err := c.requestPublic(status.Method, status.Path, nil, nil, resp)
	if err != nil {
		return nil, httpResp, err
	}

	return resp, httpResp, nil
}

// Private API

func (c *Client) Margin() (*margin.Margin, *http.Response, error) {
	resp := &margin.Margin{}

	httpResp, err := c.requestPrivate(margin.Method, margin.Path, nil, nil, resp)
	if err != nil {
		return nil, httpResp, err
	}

	return resp, httpResp, nil
}
