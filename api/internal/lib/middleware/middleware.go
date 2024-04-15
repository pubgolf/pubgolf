// Package middleware contains logic for pre- or post-processing of requests at the HTTP or Connect (gRPC) level, such as auth, logging and panic recovery.
package middleware

import (
	"fmt"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/riandyrn/otelchi"

	"github.com/pubgolf/pubgolf/api/internal/lib/telemetry"
)

// ConnectInterceptors returns the standard set of middleware for the gRPC servers.
func ConnectInterceptors() ([]connect.Interceptor, error) {
	otel, err := otelconnect.NewInterceptor()
	if err != nil {
		return nil, fmt.Errorf("construct OTel interceptor: %w", err)
	}

	return []connect.Interceptor{
		otel,
		NewLoggingInterceptor(),
		NewRecoveringInterceptor(),
	}, nil
}

// ChiMiddleware returns the standard set of middleware for the HTTP handlers.
func ChiMiddleware(r chi.Router) chi.Middlewares {
	return chi.Middlewares{
		otelchi.Middleware(telemetry.ServiceName, otelchi.WithChiRoutes(r)),
		middleware.RealIP,
		middleware.Logger,
		Recoverer,
	}
}
