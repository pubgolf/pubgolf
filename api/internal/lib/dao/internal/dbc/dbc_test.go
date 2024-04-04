package dbc_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/dbtest"

	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	t.Helper()

	return dbtest.NewTestTx(t, _sharedDB)
}

func TestPrepare(t *testing.T) {
	t.Parallel()

	ctx, tx, cleanup := initDB(t)
	defer cleanup()

	_, err := dbc.Prepare(ctx, tx)
	assert.NoError(t, err, "Preparation of queries failed")
}
