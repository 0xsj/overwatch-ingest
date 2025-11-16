"""PostgreSQL error mapping to platform errors."""

from typing import Optional

import psycopg
from psycopg import errors as pg_errors

from scout_common.errors import (
    Error,
    ErrorType,
    new_error,
    not_found,
    already_exists,
    validation_error,
    internal_error,
)


# Error codes for database operations
CODE_CONNECTION_FAILED = "DB_CONNECTION_FAILED"
CODE_QUERY_FAILED = "DB_QUERY_FAILED"
CODE_TRANSACTION_FAILED = "DB_TRANSACTION_FAILED"
CODE_UNIQUE_VIOLATION = "DB_UNIQUE_VIOLATION"
CODE_FOREIGN_KEY_VIOLATION = "DB_FOREIGN_KEY_VIOLATION"
CODE_NOT_NULL_VIOLATION = "DB_NOT_NULL_VIOLATION"
CODE_CHECK_VIOLATION = "DB_CHECK_VIOLATION"
CODE_NO_ROWS = "DB_NO_ROWS"
CODE_DEADLOCK = "DB_DEADLOCK"
CODE_SERIALIZATION_FAILURE = "DB_SERIALIZATION_FAILURE"


# PostgreSQL error codes (SQLSTATE codes)
# See: https://www.postgresql.org/docs/current/errcodes-appendix.html
PG_ERR_CODE_UNIQUE_VIOLATION = "23505"
PG_ERR_CODE_FOREIGN_KEY_VIOLATION = "23503"
PG_ERR_CODE_NOT_NULL_VIOLATION = "23502"
PG_ERR_CODE_CHECK_VIOLATION = "23514"
PG_ERR_CODE_DEADLOCK = "40P01"
PG_ERR_CODE_SERIALIZATION_FAILURE = "40001"


def map_error(err: Exception, operation: str) -> Error:
    """
    Map a psycopg error to a platform error.
    
    This is the primary function for converting database errors to our error types.
    
    Error mapping strategy:
    - No rows found -> NotFound
    - Constraint violations -> AlreadyExists or Validation
    - Connection errors -> Unavailable
    - Transaction errors -> Database (retryable)
    - Other errors -> Internal
    
    Args:
        err: The original exception
        operation: Description of the operation that failed
        
    Returns:
        Platform Error with appropriate type and details
        
    Example:
        try:
            cursor.execute("SELECT * FROM users WHERE id = %s", (user_id,))
        except Exception as e:
            raise map_error(e, "find user by id")
    """
    if err is None:
        return None
    
    # Check for no rows found (psycopg doesn't have a specific exception for this)
    # This is typically handled by checking cursor.rowcount == 0
    
    # Check for PostgreSQL-specific errors
    if isinstance(err, psycopg.Error):
        return _map_psycopg_error(err, operation)
    
    # Check for connection errors
    if _is_connection_error(err):
        return new_error(
            error_type=ErrorType.UNAVAILABLE,
            code=CODE_CONNECTION_FAILED,
            message=f"{operation}: database connection failed",
            cause=err,
        )
    
    # Default: internal error
    return new_error(
        error_type=ErrorType.INTERNAL,
        code=CODE_QUERY_FAILED,
        message=f"{operation}: database operation failed",
        cause=err,
    )


def _map_psycopg_error(err: psycopg.Error, operation: str) -> Error:
    """Map PostgreSQL-specific errors to platform errors."""
    
    # Get SQLSTATE code if available
    sqlstate = getattr(err, "sqlstate", None)
    
    if not sqlstate:
        # No SQLSTATE, return generic error
        return new_error(
            error_type=ErrorType.INTERNAL,
            code=CODE_QUERY_FAILED,
            message=f"{operation}: database error",
            cause=err,
        )
    
    # Map based on SQLSTATE code
    if sqlstate == PG_ERR_CODE_UNIQUE_VIOLATION:
        details = {}
        if hasattr(err, "diag"):
            if err.diag.constraint_name:
                details["constraint"] = err.diag.constraint_name
            if err.diag.message_detail:
                details["detail"] = err.diag.message_detail
        
        return new_error(
            error_type=ErrorType.ALREADY_EXISTS,
            code=CODE_UNIQUE_VIOLATION,
            message=f"{operation}: unique constraint violation",
            cause=err,
            details=details,
        )
    
    elif sqlstate == PG_ERR_CODE_FOREIGN_KEY_VIOLATION:
        details = {}
        if hasattr(err, "diag"):
            if err.diag.constraint_name:
                details["constraint"] = err.diag.constraint_name
            if err.diag.message_detail:
                details["detail"] = err.diag.message_detail
        
        return new_error(
            error_type=ErrorType.VALIDATION,
            code=CODE_FOREIGN_KEY_VIOLATION,
            message=f"{operation}: foreign key constraint violation",
            cause=err,
            details=details,
        )
    
    elif sqlstate == PG_ERR_CODE_NOT_NULL_VIOLATION:
        details = {}
        if hasattr(err, "diag"):
            if err.diag.column_name:
                details["column"] = err.diag.column_name
            if err.diag.table_name:
                details["table"] = err.diag.table_name
        
        return new_error(
            error_type=ErrorType.VALIDATION,
            code=CODE_NOT_NULL_VIOLATION,
            message=f"{operation}: not null constraint violation",
            cause=err,
            details=details,
        )
    
    elif sqlstate == PG_ERR_CODE_CHECK_VIOLATION:
        details = {}
        if hasattr(err, "diag"):
            if err.diag.constraint_name:
                details["constraint"] = err.diag.constraint_name
            if err.diag.message_detail:
                details["detail"] = err.diag.message_detail
        
        return new_error(
            error_type=ErrorType.VALIDATION,
            code=CODE_CHECK_VIOLATION,
            message=f"{operation}: check constraint violation",
            cause=err,
            details=details,
        )
    
    elif sqlstate == PG_ERR_CODE_DEADLOCK:
        details = {}
        if hasattr(err, "diag") and err.diag.message_detail:
            details["detail"] = err.diag.message_detail
        
        return new_error(
            error_type=ErrorType.DATABASE,
            code=CODE_DEADLOCK,
            message=f"{operation}: deadlock detected",
            cause=err,
            details=details,
        )
    
    elif sqlstate == PG_ERR_CODE_SERIALIZATION_FAILURE:
        details = {}
        if hasattr(err, "diag") and err.diag.message_detail:
            details["detail"] = err.diag.message_detail
        
        return new_error(
            error_type=ErrorType.DATABASE,
            code=CODE_SERIALIZATION_FAILURE,
            message=f"{operation}: serialization failure",
            cause=err,
            details=details,
        )
    
    else:
        # Unknown PostgreSQL error
        details = {"pg_code": sqlstate}
        if hasattr(err, "diag"):
            if err.diag.message_primary:
                details["pg_message"] = err.diag.message_primary
            if err.diag.message_detail:
                details["detail"] = err.diag.message_detail
        
        return new_error(
            error_type=ErrorType.INTERNAL,
            code=CODE_QUERY_FAILED,
            message=f"{operation}: database error (code: {sqlstate})",
            cause=err,
            details=details,
        )


def _is_connection_error(err: Exception) -> bool:
    """Check if an error is a connection-related error."""
    if err is None:
        return False
    
    # Check for psycopg connection errors
    if isinstance(err, (psycopg.OperationalError, psycopg.InterfaceError)):
        return True
    
    # Check error message for common connection errors
    err_msg = str(err).lower()
    
    connection_indicators = [
        "connection refused",
        "connection reset",
        "connection closed",
        "no such host",
        "network is unreachable",
        "timeout",
        "failed to connect",
        "could not connect",
    ]
    
    return any(indicator in err_msg for indicator in connection_indicators)


# Helper constructors for common database errors


def connection_error(operation: str, cause: Optional[Exception] = None) -> Error:
    """
    Create a connection failed error.
    
    Args:
        operation: Description of the operation
        cause: The underlying exception
        
    Returns:
        Error with UNAVAILABLE type (automatically retryable)
    """
    return new_error(
        error_type=ErrorType.UNAVAILABLE,
        code=CODE_CONNECTION_FAILED,
        message=f"{operation}: failed to connect to database",
        cause=cause,
    )


def query_error(operation: str, cause: Optional[Exception] = None) -> Error:
    """
    Create a generic query failed error.
    
    Args:
        operation: Description of the operation
        cause: The underlying exception
        
    Returns:
        Error with INTERNAL type
    """
    return new_error(
        error_type=ErrorType.INTERNAL,
        code=CODE_QUERY_FAILED,
        message=f"{operation}: query execution failed",
        cause=cause,
    )


def transaction_error(operation: str, cause: Optional[Exception] = None) -> Error:
    """
    Create a transaction failed error.
    
    Args:
        operation: Description of the operation
        cause: The underlying exception
        
    Returns:
        Error with DATABASE type (automatically retryable)
    """
    return new_error(
        error_type=ErrorType.DATABASE,
        code=CODE_TRANSACTION_FAILED,
        message=f"{operation}: transaction failed",
        cause=cause,
    )


def not_found_error(resource_type: str, resource_id: str) -> Error:
    """
    Create a not found error for database queries.
    
    Args:
        resource_type: Type of resource (e.g., "user", "incident")
        resource_id: ID of the resource
        
    Returns:
        Error with NOT_FOUND type
    """
    return not_found(resource_type, resource_id, details={"source": "database"})