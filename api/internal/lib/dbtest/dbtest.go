// Package dbtest provides helpers for running queries against a live database for testing purposes.
package dbtest

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	epg "github.com/fergusstrange/embedded-postgres"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/require"
)

// New returns an empty, ephemeral database. Namespace must be a unique identifier for the instance.
func New(namespace string, enableLogging bool) (*sql.DB, func()) {
	port, err := freeport.GetFreePort()
	guardSetup(err, "find open port")

	dbURL := fmt.Sprintf("host=localhost port=%d user=postgres password=postgres dbname=postgres sslmode=disable", port)

	tempDir, err := os.MkdirTemp(os.TempDir(), "pubgolf-api-"+namespace+"-")
	guardSetup(err, "make temp dir")

	pgConfig := epg.DefaultConfig().
		Version(epg.V16).
		Port(uint32(port)).
		RuntimePath(tempDir)
	if !enableLogging {
		pgConfig = pgConfig.Logger(io.Discard)
	}

	db := epg.NewDatabase(pgConfig)
	guardSetup(db.Start(), "start test DB")

	cfg, err := pgx.ParseConfig(dbURL)
	guardSetup(err, "prase test DB config")

	conn := stdlib.OpenDB(*cfg)

	return conn, func() {
		guardSetup(db.Stop(), "stop test DB")
	}
}

// Migrate migrates a test DB to the current schema.
func Migrate(conn *sql.DB, migrationDir string) {
	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	guardSetup(err, "create DB driver")
	migrator, err := migrate.NewWithDatabaseInstance("file://"+migrationDir, "postgres", driver)
	guardSetup(err, "create migrator")
	guardSetup(migrator.Up(), "migrate test DB")
}

// MigrationDir returns the relative path to the migrations directory. Must be called directly from within the package being tested, since that's what sets the current working directory of the test process.
func MigrationDir() string {
	_, filename, _, _ := runtime.Caller(1)
	testDir := filepath.Dir(filename)
	parts := strings.Split(testDir, filepath.FromSlash("/pubgolf/api/"))
	numDirs := strings.Count(parts[1], string(filepath.Separator))

	return filepath.FromSlash(strings.Repeat("../", numDirs) + "db/migrations")
}

// NewTestTx returns a context, transaction and cleanup function to isolate a test run into a rolled-back transaction.
func NewTestTx(t *testing.T, db *sql.DB) (context.Context, *sql.Tx, func()) {
	t.Helper()

	ctx, cancel := newTestContext(t)
	tx, rollback := newTx(ctx, t, db)

	return ctx, tx, func() {
		rollback()
		cancel()
	}
}

// newTestContext returns a context.Context which respects the test's timeout.
func newTestContext(t *testing.T) (context.Context, func()) {
	t.Helper()

	ctx := context.Background()
	cancel := func() {}

	if dl, ok := t.Deadline(); ok {
		ctx, cancel = context.WithDeadline(ctx, dl)
	}

	return ctx, cancel
}

// newTx returns a transaction and a cleanup function to roll back changes from aa test case.
func newTx(ctx context.Context, t *testing.T, db *sql.DB) (*sql.Tx, func()) {
	t.Helper()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	require.NoError(t, err)

	return tx, func() {
		err = tx.Rollback()
		require.NoError(t, err)
	}
}

// guard logs and exits on error.
func guardSetup(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %v", msg, err.Error())
	}
}
