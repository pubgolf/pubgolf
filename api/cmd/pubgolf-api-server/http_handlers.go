package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/pubgolf/pubgolf/api/internal/lib/config"
)

// status reports DB reachability and live pool stats for dashboards and debugging.
// This is NOT a liveness probe — do not wire it to any deploy health check signal.
func status(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		dbStatus := "ok"
		code := http.StatusOK

		pingErr := db.PingContext(ctx)
		if pingErr != nil {
			dbStatus = pingErr.Error()
			code = http.StatusServiceUnavailable
		}

		stats := db.Stats()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		fmt.Fprintf(w,
			`{"db":%q,"pool":{"open":%d,"in_use":%d,"idle":%d,"wait_count":%d}}`,
			dbStatus,
			stats.OpenConnections,
			stats.InUse,
			stats.Idle,
			stats.WaitCount,
		)
	}
}

// healthCheck returns a 200 if the app is online and able to process requests.
func healthCheck(cfg *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "Saluton mundo, from `%s`!", cfg.EnvName)
	}
}

var robotsTxtProd = `
User-agent: *
Disallow: /admin/
Disallow: /rpc/
`[1:]

var robotsTxt = `
User-agent: *
Disallow: /
`[1:]

// robots returns a permissive robots.txt for production and disallows all indexing in other envs.
func robots(cfg *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if cfg.EnvName == config.DeployEnvProd {
			fmt.Fprint(w, robotsTxtProd)

			return
		}

		fmt.Fprint(w, robotsTxt)
	}
}
