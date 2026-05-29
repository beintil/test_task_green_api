package green_api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultBaseURL = "https://api.green-api.com"
	defaultTimeout = 30 * time.Second
)

type Client struct {
	baseURL string
	http    *http.Client
}

type Config struct {
	BaseURL string
	Timeout time.Duration
}

type Interface interface {
	GetSettings(ctx context.Context, idInstance, apiToken string) (*Settings, error)
	GetStateInstance(ctx context.Context, idInstance, apiToken string) (*StateInstanceResponse, error)
	SendMessage(ctx context.Context, idInstance, apiToken string, req SendMessageRequest) (*SendResponse, error)
	SendFileByURL(ctx context.Context, idInstance, apiToken string, req SendFileByURLRequest) (*SendResponse, error)
}

var _ Interface = (*Client)(nil)

func NewClient(cfg Config) *Client {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}
	return &Client{
		baseURL: baseURL,
		http:    &http.Client{Timeout: timeout},
	}
}

func (c *Client) buildURL(idInstance, apiToken, method string) string {
	return fmt.Sprintf("%s/waInstance%s/%s/%s", c.baseURL, idInstance, method, apiToken)
}

func (c *Client) get(ctx context.Context, idInstance, apiToken, method string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.buildURL(idInstance, apiToken, method), nil)
	if err != nil {
		return fmt.Errorf("green_api: build GET request (%s): %w", method, err)
	}
	return c.do(req, out)
}

func (c *Client) post(ctx context.Context, idInstance, apiToken, method string, body, out any) error {
	raw, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("green_api: marshal request body (%s): %w", method, err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.buildURL(idInstance, apiToken, method), bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("green_api: build POST request (%s): %w", method, err)
	}
	req.Header.Set("Content-Type", "application/json")
	return c.do(req, out)
}

func (c *Client) do(req *http.Request, out any) error {
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("green_api: http: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("green_api: read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return parseAPIError(resp.StatusCode, body)
	}

	if out != nil {
		if err := json.Unmarshal(body, out); err != nil {
			return fmt.Errorf("green_api: decode response: %w", err)
		}
	}
	return nil
}
