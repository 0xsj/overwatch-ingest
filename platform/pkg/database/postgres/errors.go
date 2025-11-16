package postgres

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	pkgerrors "github.com/0xsj/scout/platform/pkg/errors"
)

// Error codes for database operations
const (
    CodeConnectionFailed   pkgerrors.Code = "DB_CONNECTION_FAILED"
    CodeQueryFailed        pkgerrors.Code = "DB_QUERY_FAILED"
    CodeTransactionFailed  pkgerrors.Code = "DB_TRANSACTION_FAILED"
    CodeUniqueViolation    pkgerrors.Code = "DB_UNIQUE_VIOLATION"
    CodeForeignKeyViolation pkgerrors.Code = "DB_FOREIGN_KEY_VIOLATION"
    CodeNotNullViolation   pkgerrors.Code = "DB_NOT_NULL_VIOLATION"
    CodeCheckViolation     pkgerrors.Code = "DB_CHECK_VIOLATION"
    CodeNoRows             pkgerrors.Code = "DB_NO_ROWS"
    CodeDeadlock           pkgerrors.Code = "DB_DEADLOCK"
    CodeSerializationFailure pkgerrors.Code = "DB_SERIALIZATION_FAILURE"
)

// PostgreSQL error codes (from pgconn)
// See: https://www.postgresql.org/docs/current/errcodes-appendix.html
const (
    pgErrCodeUniqueViolation      = "23505"
    pgErrCodeForeignKeyViolation  = "23503"
    pgErrCodeNotNullViolation     = "23502"
    pgErrCodeCheckViolation       = "23514"
    pgErrCodeDeadlock             = "40P01"
    pgErrCodeSerializationFailure = "40001"
)

// MapError maps a pgx error to a platform error.
// This is the primary function for converting database errors to our error types.
//
// Error mapping strategy:
// - pgx.ErrNoRows -> NotFound
// - Constraint violations -> AlreadyExists or Validation
// - Connection errors -> Unavailable (automatically retryable)
// - Deadlocks/Serialization failures -> Database (automatically retryable)
// - Other errors -> Internal
func MapError(err error, operation string) *pkgerrors.Error {
    if err == nil {
        return nil
    }
    
    // Check for pgx.ErrNoRows (common case)
    if errors.Is(err, pgx.ErrNoRows) {
        return pkgerrors.New(
            pkgerrors.ErrorTypeNotFound,
            CodeNoRows,
            fmt.Sprintf("%s: no rows found", operation),
        ).WithCause(err)
    }
    
    // Check for PostgreSQL-specific errors
    var pgErr *pgconn.PgError
    if errors.As(err, &pgErr) {
        return mapPgError(pgErr, operation)
    }
    
    // Check for connection errors
    if isConnectionError(err) {
        return pkgerrors.New(
            pkgerrors.ErrorTypeUnavailable,
            CodeConnectionFailed,
            fmt.Sprintf("%s: database connection failed", operation),
        ).WithCause(err)
        // Note: ErrorTypeUnavailable is automatically retryable
    }
    
    // Default: internal error
    return pkgerrors.New(
        pkgerrors.ErrorTypeInternal,
        CodeQueryFailed,
        fmt.Sprintf("%s: database operation failed", operation),
    ).WithCause(err)
}

// mapPgError maps PostgreSQL-specific errors to platform errors.
func mapPgError(pgErr *pgconn.PgError, operation string) *pkgerrors.Error {
    switch pgErr.Code {
    case pgErrCodeUniqueViolation:
        return pkgerrors.New(
            pkgerrors.ErrorTypeAlreadyExists,
            CodeUniqueViolation,
            fmt.Sprintf("%s: unique constraint violation", operation),
        ).WithCause(pgErr).
            WithDetail("constraint", pgErr.ConstraintName).
            WithDetail("detail", pgErr.Detail)
    
    case pgErrCodeForeignKeyViolation:
        return pkgerrors.New(
            pkgerrors.ErrorTypeValidation,
            CodeForeignKeyViolation,
            fmt.Sprintf("%s: foreign key constraint violation", operation),
        ).WithCause(pgErr).
            WithDetail("constraint", pgErr.ConstraintName).
            WithDetail("detail", pgErr.Detail)
    
    case pgErrCodeNotNullViolation:
        return pkgerrors.New(
            pkgerrors.ErrorTypeValidation,
            CodeNotNullViolation,
            fmt.Sprintf("%s: not null constraint violation", operation),
        ).WithCause(pgErr).
            WithDetail("column", pgErr.ColumnName).
            WithDetail("table", pgErr.TableName)
    
    case pgErrCodeCheckViolation:
        return pkgerrors.New(
            pkgerrors.ErrorTypeValidation,
            CodeCheckViolation,
            fmt.Sprintf("%s: check constraint violation", operation),
        ).WithCause(pgErr).
            WithDetail("constraint", pgErr.ConstraintName).
            WithDetail("detail", pgErr.Detail)
    
    case pgErrCodeDeadlock:
        // Use ErrorTypeDatabase which is automatically retryable
        return pkgerrors.New(
            pkgerrors.ErrorTypeDatabase,
            CodeDeadlock,
            fmt.Sprintf("%s: deadlock detected", operation),
        ).WithCause(pgErr).
            WithDetail("detail", pgErr.Detail)
        // Note: ErrorTypeDatabase is automatically retryable
    
    case pgErrCodeSerializationFailure:
        // Use ErrorTypeDatabase which is automatically retryable
        return pkgerrors.New(
            pkgerrors.ErrorTypeDatabase,
            CodeSerializationFailure,
            fmt.Sprintf("%s: serialization failure", operation),
        ).WithCause(pgErr).
            WithDetail("detail", pgErr.Detail)
        // Note: ErrorTypeDatabase is automatically retryable
    
    default:
        // Unknown PostgreSQL error
        return pkgerrors.New(
            pkgerrors.ErrorTypeInternal,
            CodeQueryFailed,
            fmt.Sprintf("%s: database error (code: %s)", operation, pgErr.Code),
        ).WithCause(pgErr).
            WithDetail("pg_code", pgErr.Code).
            WithDetail("pg_message", pgErr.Message).
            WithDetail("detail", pgErr.Detail)
    }
}

// isConnectionError checks if an error is a connection-related error.
func isConnectionError(err error) bool {
    if err == nil {
        return false
    }
    
    // Check error message for common connection errors
    // This is a heuristic - pgx doesn't provide typed connection errors
    errMsg := err.Error()
    
    connectionIndicators := []string{
        "connection refused",
        "connection reset",
        "connection closed",
        "no such host",
        "network is unreachable",
        "i/o timeout",
        "context deadline exceeded",
        "failed to connect",
        "dial tcp",
    }
    
    for _, indicator := range connectionIndicators {
        if contains(errMsg, indicator) {
            return true
        }
    }
    
    return false
}

// contains checks if a string contains a substring (case-insensitive helper).
func contains(s, substr string) bool {
    for i := 0; i <= len(s)-len(substr); i++ {
        if s[i:i+len(substr)] == substr {
            return true
        }
    }
    return false
}

// Helper constructors for common database errors

// ConnectionError creates a connection failed error.
func ConnectionError(operation string, cause error) *pkgerrors.Error {
    return pkgerrors.New(
        pkgerrors.ErrorTypeUnavailable,
        CodeConnectionFailed,
        fmt.Sprintf("%s: failed to connect to database", operation),
    ).WithCause(cause)
    // Note: ErrorTypeUnavailable is automatically retryable
}

// QueryError creates a generic query failed error.
func QueryError(operation string, cause error) *pkgerrors.Error {
    return pkgerrors.New(
        pkgerrors.ErrorTypeInternal,
        CodeQueryFailed,
        fmt.Sprintf("%s: query execution failed", operation),
    ).WithCause(cause)
}

// TransactionError creates a transaction failed error.
func TransactionError(operation string, cause error) *pkgerrors.Error {
    return pkgerrors.New(
        pkgerrors.ErrorTypeDatabase,
        CodeTransactionFailed,
        fmt.Sprintf("%s: transaction failed", operation),
    ).WithCause(cause)
    // Note: ErrorTypeDatabase is automatically retryable
}

// NotFoundError creates a not found error for database queries.
func NotFoundError(resourceType, resourceID string) *pkgerrors.Error {
    return pkgerrors.NotFound(resourceType, resourceID).
        WithDetail("source", "database")
}