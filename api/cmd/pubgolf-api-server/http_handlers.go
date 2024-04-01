package main

import (
	"fmt"
	"net/http"

	"github.com/pubgolf/pubgolf/api/internal/lib/config"
)

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
