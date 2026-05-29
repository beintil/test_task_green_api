package green_api

import (
	"context"
	"fmt"
)

const (
	methodSendMessage   = "sendMessage"
	methodSendFileByURL = "sendFileByUrl"

	maxMessageLength = 20_000
	maxCaptionLength = 1_024
)

func (c *Client) SendMessage(ctx context.Context, idInstance, apiToken string, req SendMessageRequest) (*SendResponse, error) {
	if err := validateSendMessage(req); err != nil {
		return nil, err
	}
	var out SendResponse
	if err := c.post(ctx, idInstance, apiToken, methodSendMessage, req, &out); err != nil {
		return nil, fmt.Errorf("SendMessage: %w", err)
	}
	return &out, nil
}

func (c *Client) SendFileByURL(ctx context.Context, idInstance, apiToken string, req SendFileByURLRequest) (*SendResponse, error) {
	if err := validateSendFileByURL(req); err != nil {
		return nil, err
	}
	var out SendResponse
	if err := c.post(ctx, idInstance, apiToken, methodSendFileByURL, req, &out); err != nil {
		return nil, fmt.Errorf("SendFileByURL: %w", err)
	}
	return &out, nil
}

func validateSendMessage(req SendMessageRequest) error {
	if req.ChatID == "" {
		return ErrEmptyChatID
	}
	if req.Message == "" {
		return ErrEmptyMessage
	}
	if len([]rune(req.Message)) > maxMessageLength {
		return ErrMessageTooLong
	}
	return nil
}

func validateSendFileByURL(req SendFileByURLRequest) error {
	if req.ChatID == "" {
		return ErrEmptyChatID
	}
	if req.URLFile == "" {
		return ErrEmptyURLFile
	}
	if req.FileName == "" {
		return ErrEmptyFileName
	}
	if len([]rune(req.Caption)) > maxCaptionLength {
		return ErrCaptionTooLong
	}
	return nil
}
