package dao

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/telemetry"
)

// Queries holds references to all data stores and provides query methods.
type Queries struct {
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
