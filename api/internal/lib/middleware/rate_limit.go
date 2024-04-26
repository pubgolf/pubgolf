package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"

	"github.com/pubgolf/pubgolf/api/internal/lib/config"
)

func rateLimitKey(r *http.Request) (string, error) {
	path := r.URL.Path
	clientID := r.Header.Get(tokenHeader)

	if nonAuthPath(r) || clientID == "" {
		ip, err := httprate.KeyByRealIP(r)
		if err != nil {
			return "", fmt.Errorf("parse true IP: %w", err)
		}

		clientID = ip
	}

	return fmt.Sprintf("%s:%s", clientID, path), nil
}

// rateLimitApply determines whether the config for the player-specific rate or the anonymous (IP-based) rate limit should apply.
func rateLimitApply(cfg *config.App, playerSpecificRate bool) func(r *http.Request) bool {
	return func(r *http.Request) bool {
		if cfg.EnvName == config.DeployEnvE2ETest {
			return false
		}

		limitByIP := nonAuthPath(r) && r.Header.Get(deviceIDHeader) == ""

		if playerSpecificRate {
			return !limitByIP
		}

		return limitByIP
	}
}

// RateLimiter applies rate limiting to the RPC routes.
func RateLimiter(cfg *config.App) chi.Middlewares {
	return chi.Middlewares{
		middleware.Maybe(
			httprate.Limit(30, 15*time.Second, httprate.WithKeyFuncs(rateLimitKey)),
			rateLimitApply(cfg, false),
		),
		middleware.Maybe(
			httprate.Limit(2, 1*time.Second, httprate.WithKeyFuncs(rateLimitKey)),
			rateLimitApply(cfg, true),
		),
	}
}
