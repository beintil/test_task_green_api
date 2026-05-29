package transperr

import "github.com/google/uuid"

type transportError struct {
	error   string
	message string
	details string
	code    int
	txID    uuid.UUID
}

func NewTransportError(errStr string, code int) TransportError {
	return &transportError{
		error: errStr,
		code:  code,
		txID:  uuid.New(),
	}
}

func (t *transportError) Error() string                        { return t.error }
func (t *transportError) SetMessage(msg string) TransportError { t.message = msg; return t }
func (t *transportError) SetDetails(d string) TransportError   { t.details = d; return t }
func (t *transportError) GetMessage() string                   { return t.message }
func (t *transportError) GetDetails() string                   { return t.details }
func (t *transportError) GetCode() int                         { return t.code }
func (t *transportError) GetTransactionID() uuid.UUID          { return t.txID }
