package dbc_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	epg "github.com/fergusstrange/embedded-postgres"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/phayes/freeport"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
)

var _sharedDB *sql.DB
var _sharedDBCleanup func()

func TestMain(m *testing.M) {
	_sharedDB, _sharedDBCleanup = setupTestDB()
	defer _sharedDBCleanup()

	os.Exit(m.Run())
}

func initDBCTest(t *testing.T) (context.Context, *dbc.Queries, func()) {
	ctx, cancel := setupTestContext(t)
	return ctx, dbc.New(_sharedDB), cancel
}

func setupTestContext(t *testing.T) (context.Context, func()) {
	ctx := context.Background()
	cancel := func() {}

	if dl, ok := t.Deadline(); ok {
		ctx, cancel = context.WithDeadline(ctx, dl)
	}

	return ctx, cancel
}

func setupTestDB() (*sql.DB, func()) {
	port, err := freeport.GetFreePort()
	guardSetup(err, "find open port")

	var dbURL = fmt.Sprintf("host=localhost port=%d user=postgres password=postgres dbname=postgres sslmode=disable", port)
	const migrationSrc = "file://../../../../db/migrations"

	db := epg.NewDatabase(epg.DefaultConfig().Version(epg.V15).Port(uint32(port)))
	guardSetup(db.Start(), "start test DB")

	cfg, err := pgx.ParseConfig(dbURL)
	guardSetup(err, "prase test DB config")
	conn := stdlib.OpenDB(*cfg)

	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	guardSetup(err, "create DB driver")
	migrator, err := migrate.NewWithDatabaseInstance(migrationSrc, "postgres", driver)
	guardSetup(err, "create migrator")
	guardSetup(migrator.Up(), "migrate test DB")

	return conn, func() {
		guardSetup(db.Stop(), "stop test DB")
	}
}

// guard logs and exits on error.
func guardSetup(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err.Error())
	}
}
