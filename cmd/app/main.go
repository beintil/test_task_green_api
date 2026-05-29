package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/darkness/green_api/internal/config"
	greenapi "github.com/darkness/green_api/internal/modules/green_api"
	"github.com/darkness/green_api/internal/runner"
	http2 "github.com/darkness/green_api/internal/server/http"
	"github.com/darkness/green_api/internal/shared/middleware"
	"github.com/darkness/green_api/internal/shared/response"
	transperr "github.com/darkness/green_api/internal/shared/transport_error"
	greenapiclient "github.com/darkness/green_api/pkg/client/green_api"
	"github.com/darkness/green_api/pkg/logger"

	"github.com/gorilla/mux"
)

//go:embed web/index.html
var indexHTML []byte

func main() {
	log, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	cfg := config.MustConfig(log)

	var (
		httpResp = response.NewHTTPResponse(log, true)
		convert  = transperr.NewErrorConverter()
		mid      = middleware.NewMiddleware(log, httpResp, convert, cfg.Handler.AllowedCORSOrigins)
		router   = mux.NewRouter()
	)

	middleware.SetupCORS(router, cfg.Handler.AllowedCORSOrigins, mid)
	router.Use(middleware.AccessLog(log))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(indexHTML)
	}).Methods(http.MethodGet)

	router.HandleFunc("/health", http2.HealthHandler).Methods(http.MethodGet)

	initBusinessLogic(router, mid, httpResp, convert)

	httpServer := http2.NewServer(&cfg.Server, router)

	logRegisteredEndpoints(log, cfg.Server.Port, router)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("HTTP server failed: %v", err)
		}
	}()

	log.Infof("server listening on port [%d] | Env %s", cfg.Server.Port, cfg.Env)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Errorf("error during shutdown: %s", err)
	}
	log.Info("server shutdown")
}

func initBusinessLogic(
	router *mux.Router,
	mid middleware.Middleware,
	httpResp response.HttpResponse,
	convert transperr.ErrorConverter,
) {
	greenAPIClient := greenapiclient.NewClient(greenapiclient.Config{})
	greenAPIService := greenapi.NewService(greenAPIClient)

	runner.InitHandlers(router, mid,
		greenapi.NewRunnerHandlerV1(router, greenAPIService, httpResp, convert),
	)
}

func logRegisteredEndpoints(log logger.Logger, port int, router *mux.Router) {
	endpoints := http2.RegisteredEndpoints(router)
	if len(endpoints) == 0 {
		return
	}
	baseURL := fmt.Sprintf("http://localhost:%d", port)
	log.Infof("HTTP endpoints (base: %s):", baseURL)
	for _, ep := range endpoints {
		log.Infof("  %-8s %s", ep.Method, ep.Path)
	}
}
