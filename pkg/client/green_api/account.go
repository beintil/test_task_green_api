package green_api

import (
	"context"
	"fmt"
)

const (
	methodGetSettings      = "getSettings"
	methodGetStateInstance = "getStateInstance"
)

func (c *Client) GetSettings(ctx context.Context, idInstance, apiToken string) (*Settings, error) {
	var out Settings
	if err := c.get(ctx, idInstance, apiToken, methodGetSettings, &out); err != nil {
		return nil, fmt.Errorf("GetSettings: %w", err)
	}
	return &out, nil
}

func (c *Client) GetStateInstance(ctx context.Context, idInstance, apiToken string) (*StateInstanceResponse, error) {
	var out StateInstanceResponse
	if err := c.get(ctx, idInstance, apiToken, methodGetStateInstance, &out); err != nil {
		return nil, fmt.Errorf("GetStateInstance: %w", err)
	}
	return &out, nil
}
