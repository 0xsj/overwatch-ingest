package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	pkgerrors "github.com/0xsj/scout/platform/pkg/errors"
	"github.com/0xsj/scout/platform/pkg/observability/logger"
)

// Client provides PostgreSQL database operations with connection pooling.
// It wraps pgxpool.Pool and provides a clean interface for database operations.
type Client interface {
    // Pool returns the underlying connection pool for advanced usage.
    // Use this when you need direct access to pgxpool features.
    Pool() *pgxpool.Pool
    
    // Ping verifies the connection to the database is still alive.
    Ping(ctx context.Context) error
    
    // Close closes all connections in the pool.
    // After calling Close, the Client should not be used.
    Close()
    
    // Stats returns connection pool statistics.
    Stats() *PoolStats
    
    // Acquire gets a connection from the pool.
    // The connection must be released back to the pool when done.
    Acquire(ctx context.Context) (*pgxpool.Conn, error)
    
    // BeginTx starts a new transaction.
    // The transaction must be committed or rolled back when done.
    BeginTx(ctx context.Context) (pgx.Tx, error)
    
    // QueryRow executes a query that returns at most one row.
    QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
    
    // Query executes a query that returns rows.
    Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
    
    // Exec executes a query that doesn't return rows (INSERT, UPDATE, DELETE).
    Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

// client is the concrete implementation of Client.
type client struct {
    pool   *pgxpool.Pool
    config *Config
    logger logger.Logger
}

// NewClient creates a new PostgreSQL client with connection pooling.
//
// Example:
//   cfg := postgres.DefaultConfig().
//       WithDSN("postgres://user:pass@localhost:5432/db").
//       WithLogger(appLogger)
//   
//   client, err := postgres.NewClient(ctx, cfg)
//   if err != nil {
//       return err
//   }
//   defer client.Close()
func NewClient(ctx context.Context, config *Config) (Client, error) {
    if config == nil {
        return nil, pkgerrors.RequiredField("config")
    }
    
    // Validate configuration
    if err := config.Validate(); err != nil {
        return nil, err
    }
    
    // Convert to pgxpool config
    poolConfig, err := config.ToPgxConfig(ctx)
    if err != nil {
        return nil, err
    }
    
    // Create connection pool
    pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
    if err != nil {
        return nil, ConnectionError("create pool", err)
    }
    
    // Verify connection
    if err := pool.Ping(ctx); err != nil {
        pool.Close()
        return nil, ConnectionError("ping database", err)
    }
    
    c := &client{
        pool:   pool,
        config: config,
        logger: config.Logger(),
    }
    
    // Log successful connection
    if c.logger != nil {
        c.logger.Info("postgres client connected",
            "max_connections", config.MaxConnections(),
            "min_connections", config.MinConnections(),
        )
    }
    
    return c, nil
}

// Pool returns the underlying connection pool.
func (c *client) Pool() *pgxpool.Pool {
    return c.pool
}

// Ping verifies the connection to the database.
func (c *client) Ping(ctx context.Context) error {
    if err := c.pool.Ping(ctx); err != nil {
        return MapError(err, "ping")
    }
    return nil
}

// Close closes all connections in the pool.
func (c *client) Close() {
    if c.logger != nil {
        c.logger.Info("postgres client closing")
    }
    
    c.pool.Close()
    
    if c.logger != nil {
        c.logger.Info("postgres client closed")
    }
}

// Stats returns connection pool statistics.
func (c *client) Stats() *PoolStats {
    stat := c.pool.Stat()
    
    return &PoolStats{
        AcquireCount:            stat.AcquireCount(),
        AcquireDuration:         stat.AcquireDuration(),
        AcquiredConns:           stat.AcquiredConns(),
        CanceledAcquireCount:    stat.CanceledAcquireCount(),
        ConstructingConns:       stat.ConstructingConns(),
        EmptyAcquireCount:       stat.EmptyAcquireCount(),
        IdleConns:               stat.IdleConns(),
        MaxConns:                stat.MaxConns(),
        TotalConns:              stat.TotalConns(),
        NewConnsCount:           stat.NewConnsCount(),
        MaxLifetimeDestroyCount: stat.MaxLifetimeDestroyCount(),
        MaxIdleDestroyCount:     stat.MaxIdleDestroyCount(),
    }
}

// Acquire gets a connection from the pool.
// The caller is responsible for releasing the connection.
//
// Example:
//   conn, err := client.Acquire(ctx)
//   if err != nil {
//       return err
//   }
//   defer conn.Release()
//   
//   // Use conn...
func (c *client) Acquire(ctx context.Context) (*pgxpool.Conn, error) {
    conn, err := c.pool.Acquire(ctx)
    if err != nil {
        return nil, MapError(err, "acquire connection")
    }
    return conn, nil
}

// BeginTx starts a new transaction.
// The caller is responsible for committing or rolling back the transaction.
//
// Example:
//   tx, err := client.BeginTx(ctx)
//   if err != nil {
//       return err
//   }
//   defer tx.Rollback(ctx) // Safe to call even after Commit
//   
//   // Do work...
//   
//   if err := tx.Commit(ctx); err != nil {
//       return err
//   }
func (c *client) BeginTx(ctx context.Context) (pgx.Tx, error) {
    tx, err := c.pool.Begin(ctx)
    if err != nil {
        return nil, MapError(err, "begin transaction")
    }
    return tx, nil
}

// QueryRow executes a query that returns at most one row.
// Errors are deferred until Row's Scan method is called.
//
// Example:
//   var name string
//   var age int
//   err := client.QueryRow(ctx, "SELECT name, age FROM users WHERE id = $1", userID).
//       Scan(&name, &age)
//   if err != nil {
//       if errors.HasCode(err, postgres.CodeNoRows) {
//           return ErrUserNotFound
//       }
//       return err
//   }
func (c *client) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
    return c.pool.QueryRow(ctx, sql, args...)
}

// Query executes a query that returns rows.
// The caller is responsible for closing the rows.
//
// Example:
//   rows, err := client.Query(ctx, "SELECT id, name FROM users WHERE active = $1", true)
//   if err != nil {
//       return err
//   }
//   defer rows.Close()
//   
//   for rows.Next() {
//       var id string
//       var name string
//       if err := rows.Scan(&id, &name); err != nil {
//           return err
//       }
//       // Process row...
//   }
//   
//   if err := rows.Err(); err != nil {
//       return err
//   }
func (c *client) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
    rows, err := c.pool.Query(ctx, sql, args...)
    if err != nil {
        return nil, MapError(err, "query")
    }
    return rows, nil
}

// Exec executes a query that doesn't return rows.
//
// Example:
//   tag, err := client.Exec(ctx, "UPDATE users SET active = $1 WHERE id = $2", false, userID)
//   if err != nil {
//       return err
//   }
//   
//   if tag.RowsAffected() == 0 {
//       return ErrUserNotFound
//   }
func (c *client) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
    tag, err := c.pool.Exec(ctx, sql, args...)
    if err != nil {
        return tag, MapError(err, "exec")
    }
    return tag, nil
}

// PoolStats contains connection pool statistics.
// This is a simplified wrapper around pgxpool.Stat for easier consumption.
type PoolStats struct {
    // AcquireCount is the cumulative count of successful acquires from the pool.
    AcquireCount int64
    
    // AcquireDuration is the total duration of all successful acquires.
    AcquireDuration time.Duration
    
    // AcquiredConns is the number of currently acquired connections.
    AcquiredConns int32
    
    // CanceledAcquireCount is the cumulative count of acquires canceled by context.
    CanceledAcquireCount int64
    
    // ConstructingConns is the number of connections being constructed.
    ConstructingConns int32
    
    // EmptyAcquireCount is the count of successful acquires that waited for a connection.
    EmptyAcquireCount int64
    
    // IdleConns is the number of idle connections in the pool.
    IdleConns int32
    
    // MaxConns is the maximum size of the pool.
    MaxConns int32
    
    // TotalConns is the total number of connections in the pool.
    TotalConns int32
    
    // NewConnsCount is the cumulative count of new connections created.
    NewConnsCount int64
    
    // MaxLifetimeDestroyCount is the count of connections destroyed due to max lifetime.
    MaxLifetimeDestroyCount int64
    
    // MaxIdleDestroyCount is the count of connections destroyed due to max idle time.
    MaxIdleDestroyCount int64
}

// String returns a human-readable representation of pool statistics.
func (s *PoolStats) String() string {
    return fmt.Sprintf(
        "total=%d idle=%d acquired=%d constructing=%d max=%d",
        s.TotalConns,
        s.IdleConns,
        s.AcquiredConns,
        s.ConstructingConns,
        s.MaxConns,
    )
}

// HealthStatus returns whether the pool is healthy based on statistics.
func (s *PoolStats) HealthStatus() string {
    if s.TotalConns == 0 {
        return "unhealthy: no connections"
    }
    
    if s.AcquiredConns >= s.MaxConns {
        return "warning: pool exhausted"
    }
    
    if s.IdleConns == 0 && s.AcquiredConns > 0 {
        return "warning: no idle connections"
    }
    
    return "healthy"
}