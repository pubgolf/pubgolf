package handlers

import (
	"context"
	"database/sql"

	log "github.com/sirupsen/logrus"

	pg "github.com/pubgolf/pubgolf/api/proto/pubgolf"
)

// RequestData is a struct for passing standard context from the middleware (logging and auth) to individual gRPC
// handler implementations.
type RequestData struct {
	// Handlers are required to preform all DB calls in this transaction.
	Tx *sql.Tx

	Ctx context.Context
	Log *log.Entry

	// Values inferred from the auth token, if present
	EventID    string
	PlayerID   string
	PlayerRole pg.PlayerRole
}
