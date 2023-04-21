package dao

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/assert"

	"github.com/pubgolf/pubgolf/api/internal/lib/dao/internal/dbc"
	"github.com/pubgolf/pubgolf/api/internal/lib/dbtest"
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

func Test_txQuerier(t *testing.T) {
	t.Run("Succeeds when DAO is constructed with New()", func(t *testing.T) {
		tempDB, tempDBCleanup := dbtest.New("dbc-test", false)
		defer tempDBCleanup()
		dbtest.Migrate(tempDB, dbtest.MigrationDir())

		ctx := context.Background()

		dao, err := New(ctx, tempDB)
		assert.NoError(t, err)

		tx, err := dao.db.BeginTx(ctx, &sql.TxOptions{})
		assert.NoError(t, err)

		_, err = dao.txQuerier(tx)
		assert.NoError(t, err)
	})

	t.Run("Fails with mock DAO", func(t *testing.T) {
		db, dbMock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		dbMock.ExpectBegin()

		m := new(dbc.MockQuerier)
		dao := Queries{dbc: m, db: db}

		tx, err := dao.db.BeginTx(context.Background(), &sql.TxOptions{})
		assert.NoError(t, err)

		_, err = dao.txQuerier(tx)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrTransactedQuerier)
	})
}
