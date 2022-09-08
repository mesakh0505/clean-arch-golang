// This is a helper file containing methods that can be used to manage an sql's cross-cutting
// transactional concerns.
package sqlx

import (
	"context"
	"database/sql"
	"errors"
)

type key string

const (
	transactionKey key = "sqlx//tx"
)

var (
	ErrMissingTx = errors.New("missing transaction in context")
)

// Wraps the parent context in a new one, and stores a started transaction. If a transaction
// already exists in context, this method simply returns the same context. If the context is
// already done, this method does nothing.
//
// Example:
//	ctx, tx, err := sqlx.WithTx(ctx, db, &sql.TxOptions{
//		Isolation: sql.LevelDefault,
//		ReadOnly: false
//	})
//	if err != nil {
//		// complain that a transaction could not be started
//	}
//	err := someRepo.doOperation(ctx, params)
//	if err != nil {
//		tx.Rollback()
//		return
//	}
//	tx.Commit()
func WithTx(parent context.Context, db *sql.DB, options *sql.TxOptions) (context.Context, *sql.Tx, error) {
	tx, ok := TxFrom(parent)
	if ok {
		return parent, tx, nil // a transaction is already present
	}

	select {
	default:
	case <-parent.Done():
		return parent, tx, nil // nothing else to do here
	}

	tx, err := db.BeginTx(parent, options)
	if err != nil {
		return parent, tx, err
	}
	return context.WithValue(parent, transactionKey, tx), tx, nil
}

// Returns the transaction stored in the provided context, if any.
//
//	tx, hasTx := sqlx.TxFrom(ctx)
//  if !hasTx {
//		// either complain, or manually start one
//	}
//	// use tx
//
func TxFrom(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(transactionKey).(*sql.Tx)
	if !ok {
		return nil, false
	}
	return tx, true
}

// Commits the transaction injected to the context
func Commit(ctx context.Context) error {
	tx, ok := TxFrom(ctx)
	if !ok {
		return nil
	}
	return tx.Commit()
}

// Rolls back the transaction injected to the context
func Rollback(ctx context.Context) error {
	tx, ok := TxFrom(ctx)
	if !ok {
		return nil
	}
	return tx.Rollback()
}
