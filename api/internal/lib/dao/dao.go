// Package dao contains methods for accessing the database.
package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/telemetry"
)

var (
	// ErrAlreadyCreated indicates that a create operation has failed due to a uniqueness violation.
	ErrAlreadyCreated = errors.New("entity already exists")
	// ErrTransactedQuerier indicates that the underlying DBC.Querier could not be used to create a transaction-compatible version of itself.
	ErrTransactedQuerier = errors.New("cannot construct transacted querier")
	// ErrInvariantViolation indicates an unexpected modeling condition within the database.
	ErrInvariantViolation = errors.New("data model invariant violated")
)

const fallbackVenueImage = "https://assets.pubgolf.co/images/venues/348x348/server-fallback.jpg"

// Queries holds references to all data stores and provides query methods.
type Queries struct {
	db  *sql.DB
	tx  *sql.Tx
	dbc dbc.Querier
}

// New returns a concrete implementation of `QueryProvider`.
func New(ctx context.Context, db *sql.DB, forcePreparedQueries bool) (*Queries, error) {
	q, err := dbc.Prepare(ctx, db)
	if err != nil {
		if forcePreparedQueries {
			return nil, fmt.Errorf("prepare dbc queries: %w", err)
		}

		log.Printf("Failed to prepare queries, initializing DAO with lazy query parsing: %+v", err)
	}

	return &Queries{
		db:  db,
		dbc: q,
	}, nil
}

func (q *Queries) useTx(ctx context.Context, query func(ctx context.Context, q *Queries) error) error {
	defer telemetry.FnSpan(&ctx)()

	// Avoid nested transactions.
	if q.tx != nil {
		return query(ctx, q)
	}

	tx, err := q.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}

	tDBC, err := q.txQuerier(tx)
	if err != nil {
		return fmt.Errorf("acquire transacted DAO: %w", err)
	}

	err = query(ctx, &Queries{tx: tx, dbc: tDBC})
	if err != nil {
		tx.Rollback() //nolint:errcheck // Already recovering from query error.

		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (q *Queries) txQuerier(tx *sql.Tx) (dbc.Querier, //nolint:ireturn // Needed to wrap type-specific logic for mocking. Needs to be *inside* the function signature to disable ireturn but not trigger nolint
	error,
) {
	if q.tx != nil {
		return q.dbc, nil
	}

	var tDBC dbc.Querier

	switch dbc := q.dbc.(type) {
	case *dbc.Queries:
		return dbc.WithTx(tx), nil
	case *dbc.MockQuerier:
		return q.dbc, nil
	default:
		return tDBC, fmt.Errorf("type %T does not implement WithTx(tx *sql.Tx) dbc.Querier: %w", dbc, ErrTransactedQuerier)
	}
}

// daoSpan annotates a DAO method with a span for tracing.
var daoSpan = telemetry.AutoSpan("dao")
