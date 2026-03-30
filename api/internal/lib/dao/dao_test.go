package dao

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func (c mockDBCCall) Bind(m *dbc.MockQuerier, name string) {
	if c.ShouldCall {
		m.On(name, c.Args...).Return(c.Return...).Once()
	}
}

func TestMain(m *testing.M) {
	testguard.UnitTest()
	os.Exit(executeTests(m))
}

func executeTests(m *testing.M) int {
	_sharedDB, _sharedDBCleanup = dbtest.NewConn("dao")
	defer _sharedDBCleanup()

	dbtest.Migrate(_sharedDB, dbtest.MigrationDir())

	return m.Run()
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
			Clock:   func() time.Time { return fixedTime },
			Querier: new(dbc.MockQuerier),
		})
		require.NoError(t, err)
		assert.Equal(t, fixedTime, dao.now())
	})

	t.Run("defaults to real time when clock is nil", func(t *testing.T) {
		t.Parallel()

		before := time.Now()
		dao, err := New(t.Context(), nil, Options{Querier: new(dbc.MockQuerier)})
		require.NoError(t, err)

		got := dao.now()
		assert.False(t, got.Before(before))
		assert.False(t, got.After(time.Now()))
	})
}
