package sqlx

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestWithTxIgnoresDoneContext(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	parent, cancel := context.WithDeadline(context.TODO(), time.Now().Add(-1*time.Millisecond))
	defer cancel()
	ctx, tx, err := WithTx(parent, db, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  true,
	})
	require.Nil(t, err)
	require.Nil(t, tx)
	require.Equal(t, parent, ctx)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestWithTxReusesExistingTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	mock.ExpectBegin()

	expectedTx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	parent := context.WithValue(context.TODO(), transactionKey, expectedTx)
	_ = db.Close() // prematurely close, so that opening a new tx will fail

	ctx, tx, err := WithTx(parent, db, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  true,
	})
	require.Equal(t, expectedTx, tx)
	require.Equal(t, parent, ctx)
	require.Nil(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestWithTxBeginsTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	mock.ExpectBegin()

	ctx, tx, err := WithTx(context.TODO(), db, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	require.NotNil(t, ctx)
	require.NotNil(t, tx)
	require.Nil(t, err)

	actual := ctx.Value(transactionKey)
	require.Equal(t, tx, actual)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
