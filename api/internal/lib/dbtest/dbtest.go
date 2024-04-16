// Package dbtest provides helpers for running queries against a live database for testing purposes.
package dbtest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
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

// namespacePrefix is prepended to ephemeral files or databases.
const namespacePrefix = "pubgolf-api-dbtest-"

// createEmbeddedURL initializes an embedded instance of Postgres, returning a URL and cleanup function.
func createEmbeddedURL(namespace string, enableLogging bool) (string, func()) {
	port, err := freeport.GetFreePort()
	guardSetup(err, "find open port")

	dbURL := fmt.Sprintf("host=localhost port=%d user=postgres password=postgres dbname=postgres sslmode=disable", port)

	tempDir, err := os.MkdirTemp(os.TempDir(), namespacePrefix+namespace+"-")
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
		guardSetup(errors.New("missing required arg"), "parse password from PUBGOLF_SHARED_DB_URL") //nolint:goerr113
	}

	cfg, err := pgx.ParseConfig(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", url.Hostname(), url.Port(), url.User.Username(), pw, strings.Trim(url.Path, "/")))
	guardSetup(err, "prase test DB config")

	conn := stdlib.OpenDB(*cfg)
	dbName := strings.ReplaceAll(namespacePrefix+namespace, "-", "_")
	_, err = conn.Exec("CREATE DATABASE " + dbName)
	guardSetup(err, "create shared DB "+dbName)

	dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", url.Hostname(), url.Port(), url.User.Username(), pw, dbName)

	return dbURL, func() {
		_, err = conn.Exec(fmt.Sprintf("DROP DATABASE %s WITH (FORCE)", dbName))
		guardSetup(err, "clean up shared DB "+dbName)
	}
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
