package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// transactable is a struct that can be used to perform database transactions
type transactable struct {
	db *sqlx.DB
}

func (r *transactable) WithTransaction(ctx context.Context,
	txFunc func(context.Context, *sqlx.Tx) error,
) error {
	var err error

	dbTx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			dbTx.Rollback() //nolint:errcheck
			panic(p)
		} else if err != nil {
			dbTx.Rollback() //nolint:errcheck
		} else {
			err = dbTx.Commit()
		}
	}()

	err = txFunc(ctx, dbTx)

	return err
}
