package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	pkgerrors "github.com/0xsj/scout/platform/pkg/errors"
)

// TxFunc is a function that performs work within a transaction.
// If the function returns an error, the transaction is rolled back.
// If the function returns nil, the transaction is committed.
type TxFunc func(pgx.Tx) error

// WithTx executes a function within a database transaction.
// It automatically handles transaction begin, commit, and rollback.
//
// If fn returns an error, the transaction is rolled back.
// If fn returns nil, the transaction is committed.
// Panics are recovered and cause a rollback.
//
// Example:
//   err := postgres.WithTx(ctx, client, func(tx pgx.Tx) error {
//       // Do transactional work
//       _, err := tx.Exec(ctx, "INSERT INTO users (name) VALUES ($1)", "Alice")
//       if err != nil {
//           return err
//       }
//       
//       _, err = tx.Exec(ctx, "INSERT INTO profiles (user_id) VALUES ($1)", userID)
//       return err
//   })
func WithTx(ctx context.Context, client Client, fn TxFunc) (err error) {
    // Begin transaction
    tx, err := client.BeginTx(ctx)
    if err != nil {
        return TransactionError("begin transaction", err)
    }
    
    // Ensure rollback on panic or error
    defer func() {
        if p := recover(); p != nil {
            // Panic occurred, rollback
            _ = tx.Rollback(ctx)
            err = pkgerrors.New(
                pkgerrors.ErrorTypeInternal,
                CodeTransactionFailed,
                fmt.Sprintf("transaction panicked: %v", p),
            )
        } else if err != nil {
            // Error occurred, rollback
            if rbErr := tx.Rollback(ctx); rbErr != nil {
                // Log rollback error but return original error
                err = pkgerrors.New(
                    pkgerrors.ErrorTypeInternal,
                    CodeTransactionFailed,
                    "transaction failed and rollback failed",
                ).WithCause(err).
                    WithDetail("rollback_error", rbErr.Error())
            }
        } else {
            // Success, commit
            if commitErr := tx.Commit(ctx); commitErr != nil {
                err = TransactionError("commit transaction", commitErr)
            }
        }
    }()
    
    // Execute the function
    err = fn(tx)
    return err
}

// WithTxFromPool executes a function within a transaction from a pool.
// This is a convenience wrapper when you have a *pgxpool.Pool directly.
//
// Example:
//   err := postgres.WithTxFromPool(ctx, pool, func(tx pgx.Tx) error {
//       // Do transactional work
//       return nil
//   })
func WithTxFromPool(ctx context.Context, pool *pgxpool.Pool, fn TxFunc) error {
    tx, err := pool.Begin(ctx)
    if err != nil {
        return TransactionError("begin transaction", err)
    }
    
    defer func() {
        if p := recover(); p != nil {
            _ = tx.Rollback(ctx)
            panic(p) // Re-panic after rollback
        } else if err != nil {
            _ = tx.Rollback(ctx)
        } else {
            err = tx.Commit(ctx)
        }
    }()
    
    err = fn(tx)
    return err
}

// TxOptions contains options for transaction execution.
type TxOptions struct {
    // IsoLevel is the transaction isolation level.
    IsoLevel pgx.TxIsoLevel
    
    // AccessMode is the transaction access mode (read-write or read-only).
    AccessMode pgx.TxAccessMode
    
    // DeferrableMode indicates if the transaction is deferrable.
    DeferrableMode pgx.TxDeferrableMode
}

// DefaultTxOptions returns the default transaction options.
// Uses READ COMMITTED isolation level and read-write mode.
func DefaultTxOptions() *TxOptions {
    return &TxOptions{
        IsoLevel:       pgx.ReadCommitted,
        AccessMode:     pgx.ReadWrite,
        DeferrableMode: pgx.NotDeferrable,
    }
}

// ReadOnlyTxOptions returns transaction options for read-only transactions.
func ReadOnlyTxOptions() *TxOptions {
    return &TxOptions{
        IsoLevel:       pgx.ReadCommitted,
        AccessMode:     pgx.ReadOnly,
        DeferrableMode: pgx.NotDeferrable,
    }
}

// SerializableTxOptions returns transaction options for serializable transactions.
// Use this for transactions that need the highest isolation level.
func SerializableTxOptions() *TxOptions {
    return &TxOptions{
        IsoLevel:       pgx.Serializable,
        AccessMode:     pgx.ReadWrite,
        DeferrableMode: pgx.NotDeferrable,
    }
}

// WithTxOptions executes a function within a transaction with custom options.
//
// Example:
//   opts := postgres.SerializableTxOptions()
//   err := postgres.WithTxOptions(ctx, client, opts, func(tx pgx.Tx) error {
//       // Do transactional work with serializable isolation
//       return nil
//   })
func WithTxOptions(ctx context.Context, client Client, opts *TxOptions, fn TxFunc) (err error) {
    // Begin transaction with options
    tx, err := client.Pool().BeginTx(ctx, pgx.TxOptions{
        IsoLevel:       opts.IsoLevel,
        AccessMode:     opts.AccessMode,
        DeferrableMode: opts.DeferrableMode,
    })
    if err != nil {
        return TransactionError("begin transaction with options", err)
    }
    
    // Same defer pattern as WithTx
    defer func() {
        if p := recover(); p != nil {
            _ = tx.Rollback(ctx)
            err = pkgerrors.New(
                pkgerrors.ErrorTypeInternal,
                CodeTransactionFailed,
                fmt.Sprintf("transaction panicked: %v", p),
            )
        } else if err != nil {
            if rbErr := tx.Rollback(ctx); rbErr != nil {
                err = pkgerrors.New(
                    pkgerrors.ErrorTypeInternal,
                    CodeTransactionFailed,
                    "transaction failed and rollback failed",
                ).WithCause(err).
                    WithDetail("rollback_error", rbErr.Error())
            }
        } else {
            if commitErr := tx.Commit(ctx); commitErr != nil {
                err = TransactionError("commit transaction", commitErr)
            }
        }
    }()
    
    err = fn(tx)
    return err
}

// WithReadOnlyTx executes a function within a read-only transaction.
// This is useful for queries that need a consistent snapshot of the database.
//
// Example:
//   var users []User
//   err := postgres.WithReadOnlyTx(ctx, client, func(tx pgx.Tx) error {
//       rows, err := tx.Query(ctx, "SELECT * FROM users")
//       if err != nil {
//           return err
//       }
//       defer rows.Close()
//       
//       // Scan results...
//       return nil
//   })
func WithReadOnlyTx(ctx context.Context, client Client, fn TxFunc) error {
    return WithTxOptions(ctx, client, ReadOnlyTxOptions(), fn)
}

// WithSerializableTx executes a function within a serializable transaction.
// This provides the highest isolation level but may result in serialization failures.
//
// Example:
//   err := postgres.WithSerializableTx(ctx, client, func(tx pgx.Tx) error {
//       // Critical operations that need strict isolation
//       return nil
//   })
func WithSerializableTx(ctx context.Context, client Client, fn TxFunc) error {
    return WithTxOptions(ctx, client, SerializableTxOptions(), fn)
}

// WithRetryableTx executes a function within a transaction with automatic retry.
// If the transaction fails due to serialization failure or deadlock,
// it will be retried up to maxRetries times.
//
// Example:
//   err := postgres.WithRetryableTx(ctx, client, 3, func(tx pgx.Tx) error {
//       // Operations that might encounter serialization conflicts
//       return nil
//   })
func WithRetryableTx(
    ctx context.Context,
    client Client,
    maxRetries int,
    fn TxFunc,
) error {
    var lastErr error
    
    for attempt := 0; attempt <= maxRetries; attempt++ {
        err := WithTx(ctx, client, fn)
        
        if err == nil {
            // Success
            return nil
        }
        
        lastErr = err
        
        // Check if error is retryable (serialization failure or deadlock)
        if !isRetryableTransactionError(err) {
            // Not retryable, fail immediately
            return err
        }
        
        // Check if we've exhausted retries
        if attempt == maxRetries {
            break
        }
        
        // Log retry attempt (if logger available)
        // Wait a bit before retrying (simple backoff)
        // In production, you might want exponential backoff
    }
    
    return pkgerrors.New(
        pkgerrors.ErrorTypeInternal,
        CodeTransactionFailed,
        fmt.Sprintf("transaction failed after %d retries", maxRetries),
    ).WithCause(lastErr)
}

// isRetryableTransactionError checks if a transaction error is retryable.
func isRetryableTransactionError(err error) bool {
    if err == nil {
        return false
    }
    
    // Check for serialization failure or deadlock codes
    if pkgerrors.HasCode(err, CodeSerializationFailure) {
        return true
    }
    
    if pkgerrors.HasCode(err, CodeDeadlock) {
        return true
    }
    
    return false
}

// TxContext wraps a transaction with context for repository usage.
// This allows repositories to accept either a Client or an existing transaction.
type TxContext interface {
    QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
    Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
    Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

// GetTxContext returns a TxContext from either a Client or pgx.Tx.
// This allows repositories to work with both.
//
// Example:
//   type UserRepository struct {
//       db postgres.Client
//   }
//   
//   func (r *UserRepository) Save(ctx context.Context, user *User, tx pgx.Tx) error {
//       txCtx := postgres.GetTxContext(r.db, tx)
//       _, err := txCtx.Exec(ctx, "INSERT INTO users ...", ...)
//       return err
//   }
func GetTxContext(client Client, tx pgx.Tx) TxContext {
    if tx != nil {
        return tx
    }
    return &clientTxContext{client: client}
}

// clientTxContext adapts Client to TxContext interface.
type clientTxContext struct {
    client Client
}

func (c *clientTxContext) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
    return c.client.QueryRow(ctx, sql, args...)
}

func (c *clientTxContext) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
    return c.client.Query(ctx, sql, args...)
}

func (c *clientTxContext) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
    return c.client.Exec(ctx, sql, args...)
}