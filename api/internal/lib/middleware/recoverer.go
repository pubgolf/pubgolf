package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"connectrpc.com/connect"
	chim "github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// errRecoveredPanic indicated a panic has been caught and wrapped in a regular error.
var errRecoveredPanic = errors.New("recovered panic")

// Recoverer is a direct copy of `github.com/go-chi/chi/v5/middleware.Recoverer`, but adds the recovered stack trace to the OTel span (in addition to logging and responding to the client with a 500 error).
func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func(ctx context.Context) {
			if rvr := recover(); rvr != nil {
				if err, ok := rvr.(error); ok && errors.Is(err, http.ErrAbortHandler) {
					// we don't recover http.ErrAbortHandler so the response
					// to the client is aborted, this should not be logged
					panic(rvr)
				}

				logEntry := chim.GetLogEntry(r)
				if logEntry != nil {
					logEntry.Panic(rvr, debug.Stack())
				} else {
					chim.PrintPrettyStack(rvr)
				}

				span := trace.SpanFromContext(ctx)
				span.SetAttributes(attribute.String("error.stack_trace", string(debug.Stack())))

				w.WriteHeader(http.StatusInternalServerError)
			}
		}(r.Context())

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// NewRecoveringInterceptor handles panics that occur within a gRPC handler by recovering, adding the stack trace to the logs and OTel span, and returning a connect error message to the client.
func NewRecoveringInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (res connect.AnyResponse, err error) { //nolint:nonamedreturns
			defer func() {
				if rvr := recover(); rvr != nil {
					span := trace.SpanFromContext(ctx)
					span.SetAttributes(attribute.String("error.stack_trace", string(debug.Stack())))

					chim.PrintPrettyStack(rvr)

					res = nil
					err = connect.NewError(connect.CodeInternal, fmt.Errorf("connect middleware recovered from panic %q: %w", rvr, errRecoveredPanic))
				}
			}()

			res, err = next(ctx, req)

			return // Implicit return due to named return vars
		})
	}

	return connect.UnaryInterceptorFunc(interceptor)
}
