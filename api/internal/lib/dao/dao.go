package dao

import (
	"database/sql"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
)

// Queries holds references to all data stores and provides query methods.
type Queries struct {
	dbc dbc.Querier
}

// New returns a concrete implementation of `QueryProvider`.
func New(db *sql.DB) *Queries {
	return &Queries{
		dbc: dbc.New(db),
	}
}
