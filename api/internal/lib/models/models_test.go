package models

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"go.uber.org/goleak"

	"github.com/pubgolf/pubgolf/api/internal/lib/dbtest"
	"github.com/pubgolf/pubgolf/api/internal/lib/testguard"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	_sharedEmptyDB        *sql.DB
	_sharedEmptyDBCleanup func()
	_sharedDB             *sql.DB
	_sharedDBCleanup      func()
)

func TestMain(m *testing.M) {
	testguard.UnitTest()

	_sharedEmptyDB, _sharedEmptyDBCleanup = dbtest.NewConn("models-test")
	_sharedDB, _sharedDBCleanup = dbtest.NewConn("models-test-migrated")
	dbtest.Migrate(_sharedDB, dbtest.MigrationDir())

	code := m.Run()

	_sharedEmptyDBCleanup()
	_sharedDBCleanup()

	if code == 0 {
		err := goleak.Find(
			// database/sql spawns a persistent goroutine to open connections on demand; it exits
			// when the DB is closed, but may still be winding down at check time.
			goleak.IgnoreTopFunction("database/sql.(*DB).connectionOpener"),
			// Background cache eviction goroutine from expirable LRU cache used by config package.
			goleak.IgnoreTopFunction("github.com/hashicorp/golang-lru/v2/expirable.NewLRU[...].func1"),
			// HTTP/2 client keep-alive reader spawned by Go's net/http transport during DB connections over TLS.
			goleak.IgnoreTopFunction("net/http.(*http2clientConnReadLoop).run"),
			// TLS read goroutine — the same HTTP/2 keep-alive reader, but on CI the stack unwinds
			// into crypto/tls rather than net/http depending on timing.
			goleak.IgnoreTopFunction("crypto/tls.(*Conn).Read"),
			// Low-level poll wait backing the TLS/HTTP goroutines above; appears on CI when the
			// goroutine is parked in the network poller rather than in user-space Read.
			goleak.IgnoreTopFunction("internal/poll.runtime_pollWait"),
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

	return dbtest.NewTestTx(t, _sharedEmptyDB)
}

func initMigratedDB(t *testing.T) (context.Context, *sql.Tx, func()) {
	t.Helper()

	return dbtest.NewTestTx(t, _sharedDB)
}
