"""PostgreSQL database client with connection pooling."""

# Configuration
from scout_common.database.postgres.config import (
    Config,
    default_config,
    from_service_config,
    mask_password,
)

# Client
from scout_common.database.postgres.client import (
    Client,
    PoolStats,
    create_client,
)

# Errors
from scout_common.database.postgres.errors import (
    # Error codes
    CODE_CONNECTION_FAILED,
    CODE_QUERY_FAILED,
    CODE_TRANSACTION_FAILED,
    CODE_UNIQUE_VIOLATION,
    CODE_FOREIGN_KEY_VIOLATION,
    CODE_NOT_NULL_VIOLATION,
    CODE_CHECK_VIOLATION,
    CODE_NO_ROWS,
    CODE_DEADLOCK,
    CODE_SERIALIZATION_FAILURE,
    # Error mapping
    map_error,
    # Error constructors
    connection_error,
    query_error,
    transaction_error,
    not_found_error,
)

# Health checks
from scout_common.database.postgres.health import (
    HealthCheck,
    check,
    check_with_timeout,
    is_ready,
    is_alive,
    wait_for_ready,
    wait_for_alive,
    monitor_health,
)

# Transactions
from scout_common.database.postgres.transaction import (
    IsoLevel,
    AccessMode,
    with_transaction,
    with_read_only_transaction,
    with_serializable_transaction,
    with_retryable_transaction,
    TxContext,
    get_tx_context,
)


__all__ = [
    # Configuration
    "Config",
    "default_config",
    "from_service_config",
    "mask_password",
    
    # Client
    "Client",
    "PoolStats",
    "create_client",
    
    # Error codes
    "CODE_CONNECTION_FAILED",
    "CODE_QUERY_FAILED",
    "CODE_TRANSACTION_FAILED",
    "CODE_UNIQUE_VIOLATION",
    "CODE_FOREIGN_KEY_VIOLATION",
    "CODE_NOT_NULL_VIOLATION",
    "CODE_CHECK_VIOLATION",
    "CODE_NO_ROWS",
    "CODE_DEADLOCK",
    "CODE_SERIALIZATION_FAILURE",
    
    # Error functions
    "map_error",
    "connection_error",
    "query_error",
    "transaction_error",
    "not_found_error",
    
    # Health checks
    "HealthCheck",
    "check",
    "check_with_timeout",
    "is_ready",
    "is_alive",
    "wait_for_ready",
    "wait_for_alive",
    "monitor_health",
    
    # Transactions
    "IsoLevel",
    "AccessMode",
    "with_transaction",
    "with_read_only_transaction",
    "with_serializable_transaction",
    "with_retryable_transaction",
    "TxContext",
    "get_tx_context",
]