package main

import (
	"fmt"
	"net/http"

	"github.com/pubgolf/pubgolf/api/internal/lib/config"
)

// healthCheck returns a 200 if the app is online and able to process requests.
func healthCheck(cfg *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Saluton mundo, from `%s`!", cfg.EnvName)
	}
}
