package dao

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/dbtest"
	"github.com/pubgolf/pubgolf/api/internal/lib/testguard"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	_sharedDB        *sql.DB
	_sharedDBCleanup func()
)

type mockDBCCall struct {
	ShouldCall bool
	Args       []any
	Return     []any
}

type fixedClock time.Time

func (c fixedClock) Now() time.Time { return time.Time(c) }

func (c mockDBCCall) Bind(m *dbc.MockQuerier, name string) {
	if c.ShouldCall {
		m.On(name, c.Args...).Return(c.Return...).Once()
	}
}

func TestMain(m *testing.M) {
	testguard.UnitTest()

	_sharedDB, _sharedDBCleanup = dbtest.NewConn("dao")
	dbtest.Migrate(_sharedDB, dbtest.MigrationDir())

	code := m.Run()

	_sharedDBCleanup()

	// Check for leaks after cleanup.
	if code == 0 {
		err := goleak.Find(
			// database/sql spawns a persistent goroutine to open connections on demand; it exits
			// when the DB is closed, but may still be winding down at check time.
			goleak.IgnoreTopFunction("database/sql.(*DB).connectionOpener"),
			// dao package initializes expirable LRU caches at package level; their eviction
			// goroutines run for the lifetime of the process and can't be stopped.
			goleak.IgnoreTopFunction("github.com/hashicorp/golang-lru/v2/expirable.NewLRU[...].func1"),
			// HTTP/2 client keep-alive reader from test infrastructure (dbtest) HTTP calls.
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

func Test_txQuerier(t *testing.T) {
	t.Parallel()

	t.Run("Succeeds when DAO is constructed with New()", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()

		dao, err := New(ctx, _sharedDB, Options{})
		require.NoError(t, err)

		tx, err := dao.db.BeginTx(ctx, &sql.TxOptions{})
		require.NoError(t, err)

		_, err = dao.txQuerier(tx)
		require.NoError(t, err)
	})

	t.Run("Succeeds with mock DAO", func(t *testing.T) {
		t.Parallel()

		db, dbMock, err := sqlmock.New()
		require.NoError(t, err)

		defer db.Close()

		dbMock.ExpectBegin()

		m := new(dbc.MockQuerier)
		dao, err := New(t.Context(), db, Options{Querier: m})
		require.NoError(t, err)

		tx, err := dao.db.BeginTx(t.Context(), &sql.TxOptions{})
		require.NoError(t, err)

		_, err = dao.txQuerier(tx)
		require.NoError(t, err)
	})
}

func Test_clockInjection(t *testing.T) {
	t.Parallel()

	t.Run("uses injected clock", func(t *testing.T) {
		t.Parallel()

		fixedTime := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
		dao, err := New(t.Context(), nil, Options{
			Clock:   fixedClock(fixedTime),
			Querier: new(dbc.MockQuerier),
		})
		require.NoError(t, err)
		assert.Equal(t, fixedTime, dao.clock.Now())
	})

	t.Run("defaults to real time when clock is nil", func(t *testing.T) {
		t.Parallel()

		before := time.Now()
		dao, err := New(t.Context(), nil, Options{Querier: new(dbc.MockQuerier)})
		require.NoError(t, err)

		got := dao.clock.Now()
		assert.False(t, got.Before(before))
		assert.False(t, got.After(time.Now()))
	})
}
