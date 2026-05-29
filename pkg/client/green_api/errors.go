package green_api

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrUnauthorized    = errors.New("green_api: unauthorized — check idInstance and apiTokenInstance")
	ErrInstanceBanned  = errors.New("green_api: instance is banned")
	ErrRateLimited     = errors.New("green_api: rate limit exceeded")
	ErrNotFound        = errors.New("green_api: resource not found")
	ErrPayloadTooLarge = errors.New("green_api: payload too large (max 100 KB)")

	ErrEmptyChatID    = errors.New("green_api: chatId must not be empty")
	ErrEmptyMessage   = errors.New("green_api: message must not be empty")
	ErrEmptyURLFile   = errors.New("green_api: urlFile must not be empty")
	ErrEmptyFileName  = errors.New("green_api: fileName must not be empty")
	ErrMessageTooLong = errors.New("green_api: message exceeds 20 000 characters")
	ErrCaptionTooLong = errors.New("green_api: caption exceeds 1 024 characters")
)

// APIError — структурированная ошибка, вернувшаяся от Green API
type APIError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	Info       string `json:"info"`
}

func (e *APIError) Error() string {
	if e.Info != "" {
		return fmt.Sprintf("green_api [%d]: %s — %s", e.StatusCode, e.Message, e.Info)
	}
	if e.Message != "" {
		return fmt.Sprintf("green_api [%d]: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("green_api: unexpected status %d", e.StatusCode)
}

// parseAPIError парсит тело ответа и возвращает подходящую ошибку
func parseAPIError(statusCode int, body []byte) error {
	apiErr := &APIError{StatusCode: statusCode}
	_ = json.Unmarshal(body, apiErr)

	switch statusCode {
	case 401:
		return fmt.Errorf("%w: %s", ErrUnauthorized, apiErr.Info)
	case 403:
		return fmt.Errorf("%w: %s", ErrInstanceBanned, apiErr.Info)
	case 404:
		return fmt.Errorf("%w", ErrNotFound)
	case 429:
		return fmt.Errorf("%w", ErrRateLimited)
	case 413:
		return fmt.Errorf("%w", ErrPayloadTooLarge)
	case 500:
		if apiErr.Info == "request entity too large" {
			return fmt.Errorf("%w", ErrPayloadTooLarge)
		}
		return apiErr
	}
	return apiErr
}
