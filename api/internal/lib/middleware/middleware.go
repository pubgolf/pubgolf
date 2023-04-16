package middleware

import (
	"github.com/bufbuild/connect-go"
	otelconnect "github.com/bufbuild/connect-opentelemetry-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// ConnectInterceptors returns the standard set of middleware for the gRPC servers.
func ConnectInterceptors() []connect.Interceptor {
	return []connect.Interceptor{
		otelconnect.NewInterceptor(),
		NewLoggingInterceptor(),
		NewRecoveringInterceptor(),
	}
}

// ChiMiddleware returns the standard set of middleware for the HTTP handlers.
func ChiMiddleware() chi.Middlewares {
	return chi.Middlewares{
		middleware.RealIP,
		middleware.Logger,
		Recoverer,
	}
}
