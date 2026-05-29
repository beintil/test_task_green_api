package response

import (
	"github.com/darkness/green_api/models"
	"net/http"
)

type HttpResponse interface {
	ErrorResponse(w http.ResponseWriter, r *http.Request, err *models.TransportError)
	WriteResponse(w http.ResponseWriter, r *http.Request, code int, resp any)
}
