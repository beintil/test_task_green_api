package dto

import (
	transperr "github.com/darkness/green_api/internal/shared/transport_error"
	"github.com/darkness/green_api/models"
	"github.com/go-openapi/strfmt"
)

func TransportErrorToModel(err transperr.TransportError) *models.TransportError {
	code := int32(err.GetCode())
	msg := err.GetMessage()
	return &models.TransportError{
		Code:          &code,
		Message:       &msg,
		Error:         err.Error(),
		Details:       err.GetDetails(),
		TransactionID: strfmt.UUID(err.GetTransactionID().String()),
	}
}
