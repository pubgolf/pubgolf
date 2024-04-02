package dao

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/dbtest"
)

var (
	_sharedDB        *sql.DB
	_sharedDBCleanup func()
)

type mockDBCCall struct {
	ShouldCall bool
	Args       []interface{}
	Return     []interface{}
}

func (c mockDBCCall) Bind(m *dbc.MockQuerier, name string) {
	if c.ShouldCall {
		m.On(name, c.Args...).Return(c.Return...)
	}
}

func TestMain(m *testing.M) {
	os.Exit(executeTests(m))
}

func executeTests(m *testing.M) int {
	_sharedDB, _sharedDBCleanup = dbtest.New("dbc-test", false)
	defer _sharedDBCleanup()

	dbtest.Migrate(_sharedDB, dbtest.MigrationDir())

	return m.Run()
}

func Test_txQuerier(t *testing.T) {
	t.Parallel()

	t.Run("Succeeds when DAO is constructed with New()", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		dao, err := New(ctx, _sharedDB, false)
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
		dao := Queries{dbc: m, db: db}

		tx, err := dao.db.BeginTx(context.Background(), &sql.TxOptions{})
		require.NoError(t, err)

		_, err = dao.txQuerier(tx)
		require.NoError(t, err)
	})
}
