package dto

import (
	"github.com/darkness/green_api/internal/domain"
	"github.com/darkness/green_api/models"
	greenapi "github.com/darkness/green_api/pkg/client/green_api"
)

func GreenAPIGetRequestToDomain(req *models.GreenAPIGetRequest) domain.GreenAPIGetInput {
	return domain.GreenAPIGetInput{
		Credentials: domain.GreenAPICredentials{
			IDInstance:       *req.IDInstance,
			APITokenInstance: *req.APITokenInstance,
		},
	}
}

func GreenAPISendMessageRequestToDomain(req *models.GreenAPISendMessageRequest) domain.GreenAPISendMessageInput {
	return domain.GreenAPISendMessageInput{
		Credentials: domain.GreenAPICredentials{
			IDInstance:       *req.IDInstance,
			APITokenInstance: *req.APITokenInstance,
		},
		ChatID:          *req.ChatID,
		Message:         *req.Message,
		QuotedMessageID: req.QuotedMessageID,
		LinkPreview:     &req.LinkPreview,
	}
}

func GreenAPISendFileByURLRequestToDomain(req *models.GreenAPISendFileByURLRequest) domain.GreenAPISendFileByURLInput {
	return domain.GreenAPISendFileByURLInput{
		Credentials: domain.GreenAPICredentials{
			IDInstance:       *req.IDInstance,
			APITokenInstance: *req.APITokenInstance,
		},
		ChatID:          *req.ChatID,
		URLFile:         *req.URLFile,
		FileName:        *req.FileName,
		Caption:         req.Caption,
		QuotedMessageID: req.QuotedMessageID,
	}
}

func GreenAPISettingsToResponse(s *greenapi.Settings) *models.GreenAPISettingsResponse {
	return &models.GreenAPISettingsResponse{
		Wid:                               s.WID,
		WebhookURL:                        s.WebhookURL,
		WebhookURLToken:                   s.WebhookURLToken,
		DelaySendMessagesMilliseconds:     int32(s.DelaySendMessagesMilliseconds),
		LinkPreview:                       s.LinkPreview,
		AutoTyping:                        int32(s.AutoTyping),
		OutgoingWebhook:                   string(s.OutgoingWebhook),
		OutgoingMessageWebhook:            string(s.OutgoingMessageWebhook),
		OutgoingAPIMessageWebhook:         string(s.OutgoingAPIMessageWebhook),
		IncomingWebhook:                   string(s.IncomingWebhook),
		StateWebhook:                      string(s.StateWebhook),
		StatusInstanceWebhook:             string(s.StatusInstanceWebhook),
		DeviceWebhook:                     string(s.DeviceWebhook),
		PollMessageWebhook:                string(s.PollMessageWebhook),
		IncomingCallWebhook:               string(s.IncomingCallWebhook),
		IncomingBlockWebhook:              string(s.IncomingBlockWebhook),
		EditedMessageWebhook:              string(s.EditedMessageWebhook),
		DeletedMessageWebhook:             string(s.DeletedMessageWebhook),
		MarkIncomingMessagesReaded:        string(s.MarkIncomingMessagesReaded),
		MarkIncomingMessagesReadedOnReply: string(s.MarkIncomingMessagesReadedOnReply),
		EnableMessagesHistory:             string(s.EnableMessagesHistory),
		KeepOnlineStatus:                  string(s.KeepOnlineStatus),
		SharedSession:                     string(s.SharedSession),
	}
}

func GreenAPIStateToResponse(state *greenapi.StateInstanceResponse) *models.GreenAPIStateResponse {
	return &models.GreenAPIStateResponse{
		StateInstance: string(state.State),
	}
}

func GreenAPISendResponseToModel(resp *greenapi.SendResponse) *models.GreenAPISendResponse {
	return &models.GreenAPISendResponse{
		IDMessage: &resp.IDMessage,
	}
}
