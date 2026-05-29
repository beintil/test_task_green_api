package response

import (
	"encoding/json"
	"github.com/darkness/green_api/models"
	"github.com/darkness/green_api/pkg/logger"
	"net/http"
)

type httpResponse struct {
	log     logger.Logger
	isDebug bool
}

func NewHTTPResponse(log logger.Logger, isDebug bool) HttpResponse {
	return &httpResponse{log: log, isDebug: isDebug}
}

func (m *httpResponse) ErrorResponse(w http.ResponseWriter, r *http.Request, err *models.TransportError) {
	m.WriteResponse(w, r, int(*err.Code), err)
}

func (m *httpResponse) WriteResponse(w http.ResponseWriter, r *http.Request, code int, resp any) {
	if code == http.StatusNoContent || (code == http.StatusCreated && resp == nil) {
		w.WriteHeader(code)
		return
	}
	raw, err := json.Marshal(resp)
	if err != nil {
		m.log.Errorf("error marshalling response: %v", err)
		return
	}
	if m.isDebug || code >= 500 {
		m.log.Debugf("[%s] '%s' [%d]", r.Method, r.URL.String(), code)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_, _ = w.Write(raw)
}
