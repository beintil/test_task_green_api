package middleware

import (
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/darkness/green_api/internal/shared/dto"
	"github.com/darkness/green_api/internal/shared/response"
	srverr "github.com/darkness/green_api/internal/shared/server_error"
	transperr "github.com/darkness/green_api/internal/shared/transport_error"
	"github.com/darkness/green_api/pkg/logger"

	"github.com/gorilla/mux"
)

type Middleware interface {
	PanicRecovery(next http.Handler) http.Handler
	CORS(next http.Handler) http.Handler
}

type middleware struct {
	log                logger.Logger
	httpResponse       response.HttpResponse
	converter          transperr.ErrorConverter
	allowedCORSOrigins string
}

func NewMiddleware(
	log logger.Logger,
	httpResponse response.HttpResponse,
	converter transperr.ErrorConverter,
	allowedCORSOrigins string,
) Middleware {
	return &middleware{
		log:                log,
		httpResponse:       httpResponse,
		converter:          converter,
		allowedCORSOrigins: allowedCORSOrigins,
	}
}

func (m *middleware) PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				m.log.Errorf("panic recovered: %v\nstack:\n%s", rec, debug.Stack())
				m.httpResponse.ErrorResponse(w, r,
					dto.TransportErrorToModel(
						m.converter.ToHTTP(srverr.NewServerError(srverr.ErrInternalServerError, "middleware.PanicRecovery")),
					),
				)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (m *middleware) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if m.allowedCORSOrigins != "" && (m.allowedCORSOrigins == "*" || originAllowed(m.allowedCORSOrigins, origin)) {
			allowOrigin := m.allowedCORSOrigins
			if allowOrigin != "*" {
				allowOrigin = origin
			}
			w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "86400")
			if allowOrigin != "*" {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func originAllowed(allowed, origin string) bool {
	if origin == "" {
		return false
	}
	for _, o := range strings.Split(strings.TrimSpace(allowed), ",") {
		if strings.TrimSpace(o) == origin {
			return true
		}
	}
	return false
}

func SetupCORS(router *mux.Router, allowedCORSOrigins string, mid Middleware) {
	m, ok := mid.(*middleware)
	if ok {
		m.allowedCORSOrigins = allowedCORSOrigins
	}
	router.Use(mid.CORS)
	router.Methods(http.MethodOptions).PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}
