package greenapi

import (
	"context"
	"net/http"

	"github.com/darkness/green_api/internal/domain"
	srverr "github.com/darkness/green_api/internal/shared/server_error"
	greenapiclient "github.com/darkness/green_api/pkg/client/green_api"
)

type Service interface {
	GetSettings(ctx context.Context, input domain.GreenAPIGetInput) (*greenapiclient.Settings, srverr.ServerError)
	GetStateInstance(ctx context.Context, input domain.GreenAPIGetInput) (*greenapiclient.StateInstanceResponse, srverr.ServerError)
	SendMessage(ctx context.Context, input domain.GreenAPISendMessageInput) (*greenapiclient.SendResponse, srverr.ServerError)
	SendFileByURL(ctx context.Context, input domain.GreenAPISendFileByURLInput) (*greenapiclient.SendResponse, srverr.ServerError)
}

type Handler interface {
	handleGetSettings(w http.ResponseWriter, r *http.Request)
	handleGetStateInstance(w http.ResponseWriter, r *http.Request)
	handleSendMessage(w http.ResponseWriter, r *http.Request)
	handleSendFileByURL(w http.ResponseWriter, r *http.Request)
}
