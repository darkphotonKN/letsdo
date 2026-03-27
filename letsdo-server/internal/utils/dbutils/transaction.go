package dbutils

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// ExecTx executes a function within a database transaction.
// It automatically handles transaction initialization, commit/rollback, and panic recovery.
//
// Usage:
//
//	err := dbutils.ExecTx(ctx, db, func(tx *sqlx.Tx) error {
//	    // Perform multiple database operations using tx
//	    if err := repo.CreateWithTx(ctx, tx, entity1); err != nil {
//	        return err // Transaction will be rolled back
//	    }
//	    if err := repo.UpdateWithTx(ctx, tx, entity2); err != nil {
//	        return err // Transaction will be rolled back
//	    }
//	    return nil // Transaction will be committed
//	})
func ExecTx(ctx context.Context, db *sqlx.DB, fn func(tx *sqlx.Tx) error) (err error) {
	tx, txBeginErr := db.BeginTxx(ctx, nil)
	if txBeginErr != nil {
		return fmt.Errorf("failed to begin transaction: %w", txBeginErr)
	}

	// Handle panics and errors, rolling back if they occur
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			// Re-throw panic after rollback
			panic(fmt.Sprintf("transaction rolled back due to panic: %v", p))
		}

		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("failed to rollback transaction: %w (original error: %w)", rbErr, err)
			}
		}
	}()

	// Execute the function with the transaction
	err = fn(tx)
	if err != nil {
		return err
	}

	// Commit the transaction
	if commitErr := tx.Commit(); commitErr != nil {
		return fmt.Errorf("failed to commit transaction: %w", commitErr)
	}

	return nil
}

// ExecTxWithOptions executes a function within a database transaction with custom options.
// This allows for more control over transaction isolation levels and other settings.
func ExecTxWithOptions(ctx context.Context, db *sqlx.DB, opts *sqlx.TxOptions, fn func(tx *sqlx.Tx) error) (err error) {
	tx, txBeginErr := db.BeginTxx(ctx, opts)
	if txBeginErr != nil {
		return fmt.Errorf("failed to begin transaction with options: %w", txBeginErr)
	}

	// Handle panics and errors, rolling back if they occur
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			// Re-throw panic after rollback
			panic(fmt.Sprintf("transaction rolled back due to panic: %v", p))
		}

		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("failed to rollback transaction: %w (original error: %w)", rbErr, err)
			}
		}
	}()

	// Execute the function with the transaction
	err = fn(tx)
	if err != nil {
		return err
	}

	// Commit the transaction
	if commitErr := tx.Commit(); commitErr != nil {
		return fmt.Errorf("failed to commit transaction: %w", commitErr)
	}

	return nil
}