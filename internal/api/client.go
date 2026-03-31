package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"scale-test/cli/internal/model"
)

// Client is an authenticated HTTP client for the Scale-Test API.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new API client.
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func (c *Client) do(req *http.Request, out interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var apiErr model.ErrorResponse
		if jsonErr := json.Unmarshal(respBody, &apiErr); jsonErr == nil && apiErr.Error != "" {
			return fmt.Errorf("API error %d: %s", resp.StatusCode, apiErr.Error)
		}
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	if out != nil {
		if err := json.Unmarshal(respBody, out); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}
	return nil
}

// CreateRun calls POST /run/new.
func (c *Client) CreateRun(reqBody model.CreateRunRequest) (*model.CreateRunResponse, error) {
	req, err := c.newRequest("POST", "/run/new", reqBody)
	if err != nil {
		return nil, err
	}
	var out model.CreateRunResponse
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetRun calls GET /runs/{id}.
func (c *Client) GetRun(id string) (*model.Run, error) {
	req, err := c.newRequest("GET", "/runs/"+id, nil)
	if err != nil {
		return nil, err
	}
	var out model.GetRunResponse
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

// DeleteRun calls DELETE /runs/{id}.
func (c *Client) DeleteRun(id string) (*model.SuccessMessage, error) {
	req, err := c.newRequest("DELETE", "/runs/"+id, nil)
	if err != nil {
		return nil, err
	}
	var out model.SuccessMessage
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
