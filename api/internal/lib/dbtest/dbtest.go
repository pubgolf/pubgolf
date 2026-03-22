// Package dbtest provides helpers for running queries against a live database for testing purposes.
package dbtest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/url"
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

	"github.com/pubgolf/pubgolf/api/internal/lib/testguard"
)

// ErrMissingRequiredArg is returned when a required argument is missing.
var ErrMissingRequiredArg = errors.New("missing required arg")

// ErrPortOutOfRange is returned when a port number is outside the valid range.
var ErrPortOutOfRange = errors.New("port out of range")

// namespacePrefix is prepended to ephemeral files or databases.
const namespacePrefix = "pubgolf-api-dbtest-"

// createEmbeddedURL initializes an embedded instance of Postgres, returning a URL and cleanup function.
func createEmbeddedURL(namespace string, enableLogging bool) (string, func()) {
	port, err := freeport.GetFreePort()
	guardSetup(err, "find open port")

	if port < 0 || port > math.MaxUint32 {
		guardSetup(ErrPortOutOfRange, "validate port range")
	}

	dbURL := fmt.Sprintf("host=localhost port=%d user=postgres password=postgres dbname=postgres sslmode=disable", port)

	tempDir, err := os.MkdirTemp(os.TempDir(), namespacePrefix+namespace+"-")
	guardSetup(err, "make temp dir")

	pgConfig := epg.DefaultConfig().
		Version(epg.V16).
		Port(uint32(port)). //nolint:gosec // Port is validated above to be in [0, MaxUint32].
		RuntimePath(tempDir)
	if !enableLogging {
		pgConfig = pgConfig.Logger(io.Discard)
	}

	db := epg.NewDatabase(pgConfig)
	guardSetup(db.Start(), "start test DB")

	return dbURL, func() {
		guardSetup(db.Stop(), "stop test DB")
	}
}

// createSharedURL initializes a new database in a running instance of Postgres, provided via the PUBGOLF_SHARED_DB_URL env var, returning a URL and cleanup function.
func createSharedURL(namespace string) (string, func()) {
	url, err := url.Parse(os.Getenv("PUBGOLF_SHARED_DB_URL"))
	guardSetup(err, "parse PUBGOLF_SHARED_DB_URL env var as url")

	pw, ok := url.User.Password()
	if !ok {
		guardSetup(ErrMissingRequiredArg, "parse password from PUBGOLF_SHARED_DB_URL")
	}

	cfg, err := pgx.ParseConfig(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", url.Hostname(), url.Port(), url.User.Username(), pw, strings.Trim(url.Path, "/")))
	guardSetup(err, "prase test DB config")

	conn := stdlib.OpenDB(*cfg)
	dbName := sharedDBName(namespace)
	_, err = conn.ExecContext(context.Background(), "CREATE DATABASE "+dbName)
	guardSetup(err, "create shared DB "+dbName)

	dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", url.Hostname(), url.Port(), url.User.Username(), pw, dbName)

	return dbURL, func() {
		_, err = conn.ExecContext(context.Background(), fmt.Sprintf("DROP DATABASE %s WITH (FORCE)", dbName))
		guardSetup(err, "clean up shared DB "+dbName)
	}
}

// sharedDBName returns the database name for shared-postgres mode,
// including a worktree slug suffix when PUBGOLF_WORKTREE_SLUG is set.
func sharedDBName(namespace string) string {
	base := strings.ReplaceAll(namespacePrefix+namespace, "-", "_")
	if slug := os.Getenv("PUBGOLF_WORKTREE_SLUG"); slug != "" {
		base += "_" + strings.ReplaceAll(slug, "-", "_")
	}

	return base
}

// NewURL returns a URL for an empty, ephemeral database, as well as a cleanup function. Namespace must be a unique identifier for the instance.
func NewURL(namespace string) (string, func()) {
	f := testguard.GetDBFlags()

	if f.UseSharedPostgres {
		return createSharedURL(namespace)
	}

	return createEmbeddedURL(namespace, f.EnablePostgresLog)
}

// NewConn returns a connection to an empty, ephemeral database, as well as a cleanup function. Namespace must be a unique identifier for the instance.
func NewConn(namespace string) (*sql.DB, func()) {
	dbURL, cleanup := NewURL(namespace)

	cfg, err := pgx.ParseConfig(dbURL)
	guardSetup(err, "prase test DB config")

	conn := stdlib.OpenDB(*cfg)

	return conn, cleanup
}

// Migrate migrates a test DB to the current schema.
func Migrate(conn *sql.DB, migrationDir string) {
	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	guardSetup(err, "create DB driver")
	migrator, err := migrate.NewWithDatabaseInstance("file://"+migrationDir, "postgres", driver)
	guardSetup(err, "create migrator")
	guardSetup(migrator.Up(), "migrate test DB")
}

// MigrationDir returns the absolute path to the migrations directory.
func MigrationDir() string {
	_, filename, _, _ := runtime.Caller(0)

	dir := filepath.Dir(filename)

	for range 20 { // bounded: prevents infinite walk in misconfigured environments
		_, statErr := os.Stat(filepath.Join(dir, "go.mod"))
		if statErr == nil {
			return filepath.Join(dir, "api", "internal", "db", "migrations")
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}

		dir = parent
	}

	panic("dbtest.MigrationDir: could not locate repo root (go.mod not found)")
}

// NewTestTx returns a context, transaction and cleanup function to isolate a test run into a rolled-back transaction.
func NewTestTx(t *testing.T, db *sql.DB) (context.Context, *sql.Tx, func()) {
	t.Helper()

	ctx := t.Context()
	tx, rollback := newTx(ctx, t, db)

	return ctx, tx, func() {
		rollback()
	}
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
		log.Panicf("%s: %v", msg, err.Error()) //nolint:gosec // Log injection is not a concern in test helpers.
	}
}
