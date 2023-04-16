package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/bufbuild/connect-go"
	chim "github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Recoverer is a direct copy of `github.com/go-chi/chi/v5/middleware.Recoverer`, but adds the recovered stack trace to the OTel span (in addition to logging and responding to the client with a 500 error).
func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler {
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

				span := trace.SpanFromContext(r.Context())
				span.SetAttributes(attribute.String("error.stack_trace", string(debug.Stack())))

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// NewRecoveringInterceptor handles panics that occur within a gRPC handler by recovering, adding the stack trace to the logs and OTel span, and returning a connect error message to the client.
func NewRecoveringInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (res connect.AnyResponse, err error) {
			defer func() {
				if rvr := recover(); rvr != nil {
					span := trace.SpanFromContext(ctx)
					span.SetAttributes(attribute.String("error.stack_trace", string(debug.Stack())))

					chim.PrintPrettyStack(rvr)

					res = nil
					err = connect.NewError(connect.CodeInternal, fmt.Errorf("connect middleware recovered panic: %v", rvr))
				}
			}()

			res, err = next(ctx, req)

			return // Implicit return due to named return vars
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
