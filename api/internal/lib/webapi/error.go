package webapi

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/trace"
)

type errorCode string

const (
	errorCodeGenericRetryable    errorCode = "ERR_GENERIC_SERVER_RETRYABLE"
	errorCodeGenericNonRetryable errorCode = "ERR_GENERIC_SERVER_NON_RETRYABLE"
	errorCodeMalformedRequest    errorCode = "ERR_MALFORMED_REQUEST"
	errorCodeNotAuthorized       errorCode = "ERR_NOT_AUTHORIZED"
)

var errorCodeHTTPStatuses = map[errorCode]int{
	errorCodeGenericRetryable:    http.StatusInternalServerError,
	errorCodeGenericNonRetryable: http.StatusInternalServerError,
	errorCodeMalformedRequest:    http.StatusBadRequest,
	errorCodeNotAuthorized:       http.StatusUnauthorized,
}

func allErrorCodes() []errorCode {
	all := make([]errorCode, 0, len(errorCodeHTTPStatuses))
	for e := range errorCodeHTTPStatuses {
		all = append(all, e)
	}

	return all
}

func (e errorCode) HTTPStatus() int {
	return errorCodeHTTPStatuses[e]
}

type errorResponse struct {
	PrivateErr error     `json:"-"`
	RootSpanID string    `json:"root_span_id"`
	Code       errorCode `json:"error_code"`
	Msg        string    `json:"error_message"`
}

func newErrorResponse(ctx context.Context, code errorCode, msg string, err error) errorResponse {
	return errorResponse{
		PrivateErr: err,
		RootSpanID: trace.SpanFromContext(ctx).SpanContext().TraceID().String(),
		Code:       code,
		Msg:        msg,
	}
}

func (e errorResponse) Render(w http.ResponseWriter, r *http.Request) {
	log.Printf("handler error: %+v\n", e.PrivateErr)
	render.Status(r, e.Code.HTTPStatus())
	render.JSON(w, r, e)
}
