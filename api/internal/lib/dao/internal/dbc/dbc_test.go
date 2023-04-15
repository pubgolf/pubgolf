package dbc_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/dbtest"
)

var (
	_sharedDB        *sql.DB
	_sharedDBCleanup func()
	_sharedDBC       *dbc.Queries
)

func TestMain(m *testing.M) {
	os.Exit(executeTests(m))
}

func executeTests(m *testing.M) int {
	_sharedDB, _sharedDBCleanup = dbtest.New("dbc-test", false)
	defer _sharedDBCleanup()

	dbtest.Migrate(_sharedDB, dbtest.MigrationDir())
	_sharedDBC = dbc.New(_sharedDB)

	return m.Run()
}

func initDB(t *testing.T) (context.Context, *sql.Tx, func()) {
	return dbtest.NewTestTx(t, _sharedDB)
}
