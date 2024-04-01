// Package webapi implements a "pure HTTP" (i.e. not gRPC) server to allow the web app to securely persist auth tokens as cookies, while itself remaining statically generated.
package webapi

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/pubgolf/pubgolf/api/internal/lib/config"
)

// Router returns a Chi router with all `webapi` routes and middleware.
func Router(cfg *config.App) func(r chi.Router) {
	return func(r chi.Router) {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{cfg.HostOrigin},
			AllowedMethods:   []string{"POST"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))

		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", logIn(cfg))
			r.Post("/logout", logOut(cfg))
			r.Post("/generate-api-token", getAPIToken(cfg))
		})
	}
}
