"""PostgreSQL transaction helpers."""

from contextlib import contextmanager
from enum import Enum
from typing import Callable, Generator, Optional, TypeVar

from psycopg import Connection, IsolationLevel, Transaction

from scout_common.database.postgres.client import Client
from scout_common.database.postgres.errors import (
    transaction_error,
    CODE_DEADLOCK,
    CODE_SERIALIZATION_FAILURE,
)
from scout_common.errors import Error, has_code


class IsoLevel(Enum):
    """Transaction isolation levels."""
    
    READ_UNCOMMITTED = IsolationLevel.READ_UNCOMMITTED
    READ_COMMITTED = IsolationLevel.READ_COMMITTED
    REPEATABLE_READ = IsolationLevel.REPEATABLE_READ
    SERIALIZABLE = IsolationLevel.SERIALIZABLE


class AccessMode(Enum):
    """Transaction access modes."""
    
    READ_WRITE = "read_write"
    READ_ONLY = "read_only"


T = TypeVar("T")


@contextmanager
def with_transaction(
    client: Client,
    isolation_level: IsoLevel = IsoLevel.READ_COMMITTED,
    read_only: bool = False,
) -> Generator[Connection, None, None]:
    """
    Execute a function within a database transaction.
    
    Automatically handles transaction begin, commit, and rollback.
    - If the context exits normally, the transaction is committed
    - If an exception occurs, the transaction is rolled back
    - Panics/exceptions are re-raised after rollback
    
    Args:
        client: Database client
        isolation_level: Transaction isolation level
        read_only: Whether transaction is read-only
        
    Yields:
        Database connection within a transaction
        
    Raises:
        Error: If transaction fails
        
    Example:
        with with_transaction(client) as conn:
            cur = conn.cursor()
            cur.execute("INSERT INTO users (name) VALUES (%s)", ("Alice",))
            cur.execute("INSERT INTO profiles (user_id) VALUES (%s)", (user_id,))
            # Automatic commit on success, rollback on error
    """
    with client.connection() as conn:
        # Set transaction properties
        conn.isolation_level = isolation_level.value
        conn.read_only = read_only
        
        try:
            with conn.transaction():
                yield conn
        except Exception as e:
            # Transaction is automatically rolled back by psycopg
            raise transaction_error("transaction failed", e)


@contextmanager
def with_read_only_transaction(
    client: Client,
) -> Generator[Connection, None, None]:
    """
    Execute a function within a read-only transaction.
    
    This is useful for queries that need a consistent snapshot of the database.
    
    Args:
        client: Database client
        
    Yields:
        Database connection within a read-only transaction
        
    Example:
        users = []
        with with_read_only_transaction(client) as conn:
            cur = conn.cursor()
            cur.execute("SELECT * FROM users")
            users = cur.fetchall()
    """
    with with_transaction(client, read_only=True) as conn:
        yield conn


@contextmanager
def with_serializable_transaction(
    client: Client,
) -> Generator[Connection, None, None]:
    """
    Execute a function within a serializable transaction.
    
    This provides the highest isolation level but may result in
    serialization failures that require retry.
    
    Args:
        client: Database client
        
    Yields:
        Database connection within a serializable transaction
        
    Example:
        with with_serializable_transaction(client) as conn:
            cur = conn.cursor()
            # Critical operations that need strict isolation
            cur.execute("SELECT balance FROM accounts WHERE id = %s FOR UPDATE", (id,))
            balance = cur.fetchone()[0]
            cur.execute("UPDATE accounts SET balance = %s WHERE id = %s", (balance - amount, id))
    """
    with with_transaction(client, isolation_level=IsoLevel.SERIALIZABLE) as conn:
        yield conn


TxFunc = Callable[[Connection], T]


def with_retryable_transaction(
    client: Client,
    max_retries: int,
    fn: TxFunc[T],
    isolation_level: IsoLevel = IsoLevel.READ_COMMITTED,
) -> T:
    """
    Execute a function within a transaction with automatic retry.
    
    If the transaction fails due to serialization failure or deadlock,
    it will be retried up to max_retries times.
    
    Args:
        client: Database client
        max_retries: Maximum number of retry attempts
        fn: Function to execute within transaction
        isolation_level: Transaction isolation level
        
    Returns:
        Result from the function
        
    Raises:
        Error: If transaction fails after all retries
        
    Example:
        def transfer_funds(conn: Connection) -> None:
            cur = conn.cursor()
            cur.execute("SELECT balance FROM accounts WHERE id = %s FOR UPDATE", (from_id,))
            balance = cur.fetchone()[0]
            
            if balance < amount:
                raise ValueError("Insufficient funds")
            
            cur.execute("UPDATE accounts SET balance = balance - %s WHERE id = %s", (amount, from_id))
            cur.execute("UPDATE accounts SET balance = balance + %s WHERE id = %s", (amount, to_id))
        
        with_retryable_transaction(client, 3, transfer_funds)
    """
    last_error: Optional[Error] = None
    
    for attempt in range(max_retries + 1):
        try:
            with with_transaction(client, isolation_level=isolation_level) as conn:
                return fn(conn)
        
        except Error as e:
            last_error = e
            
            # Check if error is retryable
            if not _is_retryable_transaction_error(e):
                # Not retryable, fail immediately
                raise
            
            # Check if we've exhausted retries
            if attempt == max_retries:
                break
            
            # Log retry attempt (optional - could add logger parameter)
            # time.sleep() for backoff (optional - could add backoff strategy)
    
    # All retries exhausted
    raise Error(
        error_type="INTERNAL",
        code="TRANSACTION_FAILED",
        message=f"transaction failed after {max_retries} retries",
        cause=last_error,
    )


def _is_retryable_transaction_error(err: Error) -> bool:
    """Check if a transaction error is retryable."""
    if err is None:
        return False
    
    # Check for serialization failure or deadlock codes
    if has_code(err, CODE_SERIALIZATION_FAILURE):
        return True
    
    if has_code(err, CODE_DEADLOCK):
        return True
    
    return False


# Optional transaction context for repositories


class TxContext:
    """
    Transaction context that can wrap either a Client or Connection.
    
    This allows repositories to accept either a client (auto-transaction)
    or an existing connection (manual transaction control).
    """
    
    def __init__(self, client: Optional[Client] = None, conn: Optional[Connection] = None):
        """
        Initialize transaction context.
        
        Args:
            client: Database client (for auto-transaction)
            conn: Existing connection (for manual transaction)
        """
        if client is None and conn is None:
            raise ValueError("Either client or conn must be provided")
        
        self._client = client
        self._conn = conn
    
    @contextmanager
    def cursor(self, **kwargs):
        """Get a cursor, handling transaction context automatically."""
        if self._conn is not None:
            # Using existing connection
            with self._conn.cursor(**kwargs) as cur:
                yield cur
        else:
            # Using client - create connection
            with self._client.cursor(**kwargs) as cur:
                yield cur


def get_tx_context(client: Client, conn: Optional[Connection] = None) -> TxContext:
    """
    Get a transaction context from either a Client or Connection.
    
    Args:
        client: Database client
        conn: Optional existing connection
        
    Returns:
        TxContext that works with both
        
    Example:
        class UserRepository:
            def __init__(self, db: Client):
                self.db = db
            
            def save(self, user: dict, conn: Optional[Connection] = None):
                tx = get_tx_context(self.db, conn)
                with tx.cursor() as cur:
                    cur.execute("INSERT INTO users ...", (...))
        
        # Without transaction
        repo.save(user)
        
        # With transaction
        with with_transaction(client) as conn:
            repo.save(user1, conn)
            repo.save(user2, conn)
    """
    return TxContext(client=client, conn=conn)