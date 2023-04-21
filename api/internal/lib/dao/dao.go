package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/telemetry"
)

var (
	// ErrAlreadyCreated indicates that a create operation has failed due to a uniqueness violation.
	ErrAlreadyCreated = errors.New("entity already exists")
	// ErrTransactedQuerier indicates that the underlying DBC.Querier could not be used to create a transaction-compatible version of itself.
	ErrTransactedQuerier = errors.New("cannot construct transacted querier")
)

const fallbackVenueImage = "https://assets.pubgolf.co/images/venues/348x348/server-fallback.jpg"

// Queries holds references to all data stores and provides query methods.
type Queries struct {
	db  *sql.DB
	tx  *sql.Tx
	dbc dbc.Querier
}

// New returns a concrete implementation of `QueryProvider`.
func New(ctx context.Context, db *sql.DB) (*Queries, error) {
	q, err := dbc.Prepare(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("prepare dbc queries: %w", err)
	}

	return &Queries{
		dbc: q,
	}, nil
}

	}
}

// daoSpan annotates a DAO method with a span for tracing.
var daoSpan = telemetry.AutoSpan("dao")
