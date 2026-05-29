package greenapi

import (
	"net/http"

	"github.com/darkness/green_api/internal/runner"
	"github.com/darkness/green_api/internal/shared/middleware"
	"github.com/darkness/green_api/internal/shared/response"
	transperr "github.com/darkness/green_api/internal/shared/transport_error"

	"github.com/gorilla/mux"
)

type runnerV1 struct {
	router  *mux.Router
	handler Handler
}

func NewRunnerHandlerV1(
	router *mux.Router,
	service Service,
	httpResp response.HttpResponse,
	converter transperr.ErrorConverter,
) runner.Handler {
	return &runnerV1{
		router:  router.PathPrefix("/v1").Subrouter(),
		handler: NewHandler(service, httpResp, converter),
	}
}

func (r *runnerV1) Init() []runner.Runner {
	return []runner.Runner{r}
}

func (r *runnerV1) RouterWithVersion() *mux.Router {
	return r.router
}

func (r *runnerV1) Run(router *mux.Router, _ middleware.Middleware) {
	api := router.PathPrefix("/green-api").Subrouter()

	api.HandleFunc("/settings", r.handler.handleGetSettings).Methods(http.MethodPost)
	api.HandleFunc("/state", r.handler.handleGetStateInstance).Methods(http.MethodPost)
	api.HandleFunc("/send-message", r.handler.handleSendMessage).Methods(http.MethodPost)
	api.HandleFunc("/send-file-by-url", r.handler.handleSendFileByURL).Methods(http.MethodPost)
}
