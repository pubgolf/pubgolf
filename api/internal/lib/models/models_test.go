package models

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/pubgolf/pubgolf/api/internal/lib/dbtest"
)

var (
	_sharedEmptyDB        *sql.DB
	_sharedEmptyDBCleanup func()
	_sharedDB             *sql.DB
	_sharedDBCleanup      func()
)

func TestMain(m *testing.M) {
	os.Exit(executeTests(m))
}

func executeTests(m *testing.M) int {
	_sharedEmptyDB, _sharedEmptyDBCleanup = dbtest.New("models-test", false)
	defer _sharedEmptyDBCleanup()

	_sharedDB, _sharedDBCleanup = dbtest.New("models-test-migrated", false)
	defer _sharedDBCleanup()
	dbtest.Migrate(_sharedDB, dbtest.MigrationDir())

	return m.Run()
}

func initDB(t *testing.T) (context.Context, *sql.Tx, func()) {
	t.Helper()

	return dbtest.NewTestTx(t, _sharedEmptyDB)
}

func initMigratedDB(t *testing.T) (context.Context, *sql.Tx, func()) {
	t.Helper()

	return dbtest.NewTestTx(t, _sharedDB)
}
