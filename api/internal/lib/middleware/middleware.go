package middleware

import (
	"net/http"
	"time"

	"github.com/bufbuild/connect-go"
	otelconnect "github.com/bufbuild/connect-opentelemetry-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/riandyrn/otelchi"

	"github.com/pubgolf/pubgolf/api/internal/lib/telemetry"
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
func ChiMiddleware(r chi.Router) chi.Middlewares {
	return chi.Middlewares{
		otelchi.Middleware(telemetry.ServiceName, otelchi.WithChiRoutes(r)),
		middleware.RealIP,
		middleware.Logger,
		httprate.Limit(10, 1*time.Second, httprate.WithKeyFuncs(func(r *http.Request) (string, error) {
			return r.Header.Get("X-PubGolf-User-ID"), nil
		})),
		Recoverer,
	}
}
