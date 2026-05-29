package greenapi

import (
	"context"

	"github.com/darkness/green_api/internal/domain"
	srverr "github.com/darkness/green_api/internal/shared/server_error"
	greenapiclient "github.com/darkness/green_api/pkg/client/green_api"
)

const (
	ServiceErrBadRequest srverr.ErrorTypeBadRequest          = "green_api_bad_request"
	ServiceErrInternal   srverr.ErrorTypeInternalServerError = "green_api_internal_error"
)

const (
	msgIDInstanceRequired = "idInstance must not be empty"
	msgAPITokenRequired   = "apiTokenInstance must not be empty"
	msgChatIDRequired     = "chatId must not be empty"
	msgMessageRequired    = "message must not be empty"
	msgURLFileRequired    = "urlFile must not be empty"
	msgFileNameRequired   = "fileName must not be empty"
)

type service struct {
	client *greenapiclient.Client
}

func NewService(client *greenapiclient.Client) Service {
	return &service{client: client}
}

func (s *service) GetSettings(ctx context.Context, input domain.GreenAPIGetInput) (*greenapiclient.Settings, srverr.ServerError) {
	if sErr := s.validateCredentials(input.Credentials); sErr != nil {
		return nil, sErr
	}
	result, err := s.client.GetSettings(ctx, input.Credentials.IDInstance, input.Credentials.APITokenInstance)
	if err != nil {
		return nil, srverr.NewServerError(ServiceErrInternal, "greenapi.GetSettings").SetDetails(err.Error())
	}
	return result, nil
}

func (s *service) GetStateInstance(ctx context.Context, input domain.GreenAPIGetInput) (*greenapiclient.StateInstanceResponse, srverr.ServerError) {
	if sErr := s.validateCredentials(input.Credentials); sErr != nil {
		return nil, sErr
	}
	result, err := s.client.GetStateInstance(ctx, input.Credentials.IDInstance, input.Credentials.APITokenInstance)
	if err != nil {
		return nil, srverr.NewServerError(ServiceErrInternal, "greenapi.GetStateInstance").SetDetails(err.Error())
	}
	return result, nil
}

func (s *service) SendMessage(ctx context.Context, input domain.GreenAPISendMessageInput) (*greenapiclient.SendResponse, srverr.ServerError) {
	if sErr := s.validateCredentials(input.Credentials); sErr != nil {
		return nil, sErr
	}
	if input.ChatID == "" {
		return nil, srverr.NewServerError(ServiceErrBadRequest, "greenapi.SendMessage").SetDetails(msgChatIDRequired)
	}
	if input.Message == "" {
		return nil, srverr.NewServerError(ServiceErrBadRequest, "greenapi.SendMessage").SetDetails(msgMessageRequired)
	}
	result, err := s.client.SendMessage(ctx, input.Credentials.IDInstance, input.Credentials.APITokenInstance, greenapiclient.SendMessageRequest{
		ChatID:          input.ChatID,
		Message:         input.Message,
		QuotedMessageID: input.QuotedMessageID,
		LinkPreview:     input.LinkPreview,
	})
	if err != nil {
		return nil, srverr.NewServerError(ServiceErrInternal, "greenapi.SendMessage").SetDetails(err.Error())
	}
	return result, nil
}

func (s *service) SendFileByURL(ctx context.Context, input domain.GreenAPISendFileByURLInput) (*greenapiclient.SendResponse, srverr.ServerError) {
	if sErr := s.validateCredentials(input.Credentials); sErr != nil {
		return nil, sErr
	}
	if input.ChatID == "" {
		return nil, srverr.NewServerError(ServiceErrBadRequest, "greenapi.SendFileByURL").SetDetails(msgChatIDRequired)
	}
	if input.URLFile == "" {
		return nil, srverr.NewServerError(ServiceErrBadRequest, "greenapi.SendFileByURL").SetDetails(msgURLFileRequired)
	}
	if input.FileName == "" {
		return nil, srverr.NewServerError(ServiceErrBadRequest, "greenapi.SendFileByURL").SetDetails(msgFileNameRequired)
	}
	result, err := s.client.SendFileByURL(ctx, input.Credentials.IDInstance, input.Credentials.APITokenInstance, greenapiclient.SendFileByURLRequest{
		ChatID:          input.ChatID,
		URLFile:         input.URLFile,
		FileName:        input.FileName,
		Caption:         input.Caption,
		QuotedMessageID: input.QuotedMessageID,
	})
	if err != nil {
		return nil, srverr.NewServerError(ServiceErrInternal, "greenapi.SendFileByURL").SetDetails(err.Error())
	}
	return result, nil
}

func (s *service) validateCredentials(creds domain.GreenAPICredentials) srverr.ServerError {
	if creds.IDInstance == "" {
		return srverr.NewServerError(ServiceErrBadRequest, "greenapi.validateCredentials").SetDetails(msgIDInstanceRequired)
	}
	if creds.APITokenInstance == "" {
		return srverr.NewServerError(ServiceErrBadRequest, "greenapi.validateCredentials").SetDetails(msgAPITokenRequired)
	}
	return nil
}
