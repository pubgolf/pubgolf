package dbc_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/dbtest"
	"github.com/pubgolf/pubgolf/api/internal/lib/testguard"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	_sharedDB        *sql.DB
	_sharedDBCleanup func()
	_sharedDBC       *dbc.Queries
)

func TestMain(m *testing.M) {
	testguard.UnitTest()

	_sharedDB, _sharedDBCleanup = dbtest.NewConn("dbc")
	dbtest.Migrate(_sharedDB, dbtest.MigrationDir())
	_sharedDBC = dbc.New(_sharedDB)

	code := m.Run()

	_sharedDBCleanup()

	if code == 0 {
		err := goleak.Find(
			// database/sql spawns a persistent goroutine to open connections on demand; it exits
			// when the DB is closed, but may still be winding down at check time.
			goleak.IgnoreTopFunction("database/sql.(*DB).connectionOpener"),
			// Background cache eviction goroutine from expirable LRU cache used by config package.
			goleak.IgnoreTopFunction("github.com/hashicorp/golang-lru/v2/expirable.NewLRU[...].func1"),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "goleak: %v\n", err)
			os.Exit(1)
		}
	}

	os.Exit(code)
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
