package greenapi

import (
	"encoding/json"
	"net/http"

	"github.com/darkness/green_api/internal/shared/dto"
	"github.com/darkness/green_api/internal/shared/response"
	srverr "github.com/darkness/green_api/internal/shared/server_error"
	transperr "github.com/darkness/green_api/internal/shared/transport_error"
	"github.com/darkness/green_api/models"
)

type handler struct {
	service   Service
	httpResp  response.HttpResponse
	converter transperr.ErrorConverter
}

func NewHandler(
	service Service,
	httpResp response.HttpResponse,
	converter transperr.ErrorConverter,
) Handler {
	return &handler{
		service:   service,
		httpResp:  httpResp,
		converter: converter,
	}
}

func (h *handler) handleGetSettings(w http.ResponseWriter, r *http.Request) {
	var req models.GreenAPIGetRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}
	settings, sErr := h.service.GetSettings(r.Context(), dto.GreenAPIGetRequestToDomain(&req))
	if sErr != nil {
		h.httpResp.ErrorResponse(w, r, dto.TransportErrorToModel(h.converter.ToHTTP(sErr)))
		return
	}
	h.httpResp.WriteResponse(w, r, http.StatusOK, dto.GreenAPISettingsToResponse(settings))
}

func (h *handler) handleGetStateInstance(w http.ResponseWriter, r *http.Request) {
	var req models.GreenAPIGetRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}
	state, sErr := h.service.GetStateInstance(r.Context(), dto.GreenAPIGetRequestToDomain(&req))
	if sErr != nil {
		h.httpResp.ErrorResponse(w, r, dto.TransportErrorToModel(h.converter.ToHTTP(sErr)))
		return
	}
	h.httpResp.WriteResponse(w, r, http.StatusOK, dto.GreenAPIStateToResponse(state))
}

func (h *handler) handleSendMessage(w http.ResponseWriter, r *http.Request) {
	var req models.GreenAPISendMessageRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}
	result, sErr := h.service.SendMessage(r.Context(), dto.GreenAPISendMessageRequestToDomain(&req))
	if sErr != nil {
		h.httpResp.ErrorResponse(w, r, dto.TransportErrorToModel(h.converter.ToHTTP(sErr)))
		return
	}
	h.httpResp.WriteResponse(w, r, http.StatusOK, dto.GreenAPISendResponseToModel(result))
}

func (h *handler) handleSendFileByURL(w http.ResponseWriter, r *http.Request) {
	var req models.GreenAPISendFileByURLRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}
	result, sErr := h.service.SendFileByURL(r.Context(), dto.GreenAPISendFileByURLRequestToDomain(&req))
	if sErr != nil {
		h.httpResp.ErrorResponse(w, r, dto.TransportErrorToModel(h.converter.ToHTTP(sErr)))
		return
	}
	h.httpResp.WriteResponse(w, r, http.StatusOK, dto.GreenAPISendResponseToModel(result))
}

func (h *handler) decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		h.httpResp.ErrorResponse(w, r,
			dto.TransportErrorToModel(h.converter.ToHTTP(
				srverr.NewServerError(ServiceErrBadRequest, "greenapi.decodeJSON").SetDetails(err.Error()),
			)),
		)
		return false
	}
	return true
}
