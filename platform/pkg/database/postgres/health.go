package postgres

import (
	"context"
	"time"

	pkgerrors "github.com/0xsj/scout/platform/pkg/errors"
)

// HealthCheck represents the result of a database health check.
type HealthCheck struct {
    // Healthy indicates if the database is healthy and ready to serve requests.
    Healthy bool `json:"healthy"`
    
    // Message provides additional context about the health status.
    Message string `json:"message,omitempty"`
    
    // Latency is the time taken to ping the database.
    Latency time.Duration `json:"latency"`
    
    // Stats contains connection pool statistics.
    Stats *PoolStats `json:"stats,omitempty"`
    
    // Error contains error details if the check failed.
    Error string `json:"error,omitempty"`
}

// Check performs a health check on the database client.
// This is suitable for Kubernetes liveness and readiness probes.
//
// The check:
// - Pings the database to verify connectivity
// - Measures response latency
// - Retrieves pool statistics
// - Determines overall health status
//
// Example:
//   health := postgres.Check(ctx, client)
//   if !health.Healthy {
//       log.Errorf("Database unhealthy: %s", health.Message)
//   }
func Check(ctx context.Context, client Client) *HealthCheck {
    start := time.Now()
    
    // Ping the database
    err := client.Ping(ctx)
    latency := time.Since(start)
    
    // Get pool statistics
    stats := client.Stats()
    
    // Build health check result
    check := &HealthCheck{
        Latency: latency,
        Stats:   stats,
    }
    
    // Check for ping errors
    if err != nil {
        check.Healthy = false
        check.Message = "database ping failed"
        check.Error = err.Error()
        return check
    }
    
    // Check pool health
    poolHealth := stats.HealthStatus()
    if poolHealth != "healthy" {
        check.Healthy = false
        check.Message = poolHealth
        return check
    }
    
    // All checks passed
    check.Healthy = true
    check.Message = "database is healthy"
    
    return check
}

// CheckWithTimeout performs a health check with a specified timeout.
// This is useful for ensuring health checks don't block indefinitely.
//
// Example:
//   health := postgres.CheckWithTimeout(client, 5*time.Second)
//   if !health.Healthy {
//       return http.StatusServiceUnavailable
//   }
func CheckWithTimeout(client Client, timeout time.Duration) *HealthCheck {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    
    return Check(ctx, client)
}

// IsReady determines if the database is ready to serve requests.
// This is suitable for Kubernetes readiness probes.
//
// A database is considered ready if:
// - It responds to pings
// - The connection pool has available connections
// - Response latency is acceptable
//
// Example:
//   if !postgres.IsReady(ctx, client) {
//       return http.StatusServiceUnavailable
//   }
func IsReady(ctx context.Context, client Client) bool {
    check := Check(ctx, client)
    
    // Not ready if unhealthy
    if !check.Healthy {
        return false
    }
    
    // Not ready if latency is too high (>1 second indicates issues)
    if check.Latency > time.Second {
        return false
    }
    
    // Not ready if pool is exhausted
    if check.Stats.AcquiredConns >= check.Stats.MaxConns {
        return false
    }
    
    return true
}

// IsAlive determines if the database connection is alive.
// This is suitable for Kubernetes liveness probes.
//
// A database is considered alive if it responds to pings,
// even if the pool is under stress.
//
// Example:
//   if !postgres.IsAlive(ctx, client) {
//       log.Fatal("Database connection lost, restarting...")
//   }
func IsAlive(ctx context.Context, client Client) bool {
    err := client.Ping(ctx)
    return err == nil
}

// WaitForReady blocks until the database is ready or the context times out.
// This is useful during application startup to wait for database availability.
//
// Example:
//   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//   defer cancel()
//   
//   if err := postgres.WaitForReady(ctx, client, 1*time.Second); err != nil {
//       log.Fatalf("Database never became ready: %v", err)
//   }
func WaitForReady(ctx context.Context, client Client, checkInterval time.Duration) error {
    ticker := time.NewTicker(checkInterval)
    defer ticker.Stop()
    
    // Try immediately first
    if IsReady(ctx, client) {
        return nil
    }
    
    for {
        select {
        case <-ctx.Done():
            return pkgerrors.New(
                pkgerrors.ErrorTypeTimeout,
                "DB_READY_TIMEOUT",
                "database did not become ready within timeout",
            ).WithCause(ctx.Err())
        
        case <-ticker.C:
            if IsReady(ctx, client) {
                return nil
            }
        }
    }
}

// WaitForAlive blocks until the database is alive or the context times out.
// This is useful during application startup to wait for initial database connectivity.
//
// Example:
//   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//   defer cancel()
//   
//   if err := postgres.WaitForAlive(ctx, client, 1*time.Second); err != nil {
//       log.Fatalf("Database never became alive: %v", err)
//   }
func WaitForAlive(ctx context.Context, client Client, checkInterval time.Duration) error {
    ticker := time.NewTicker(checkInterval)
    defer ticker.Stop()
    
    // Try immediately first
    if IsAlive(ctx, client) {
        return nil
    }
    
    for {
        select {
        case <-ctx.Done():
            return pkgerrors.New(
                pkgerrors.ErrorTypeTimeout,
                "DB_ALIVE_TIMEOUT",
                "database did not become alive within timeout",
            ).WithCause(ctx.Err())
        
        case <-ticker.C:
            if IsAlive(ctx, client) {
                return nil
            }
        }
    }
}

// HealthCheckFunc is a function that performs a health check.
// This allows for custom health check implementations.
type HealthCheckFunc func(context.Context, Client) *HealthCheck

// MonitorHealth periodically performs health checks and calls a callback.
// This is useful for background health monitoring and alerting.
//
// The monitor runs until the context is canceled.
//
// Example:
//   ctx, cancel := context.WithCancel(context.Background())
//   defer cancel()
//   
//   go postgres.MonitorHealth(ctx, client, 30*time.Second, func(health *HealthCheck) {
//       if !health.Healthy {
//           alerting.Send("Database unhealthy: " + health.Message)
//       }
//       metrics.RecordLatency("db.health.latency", health.Latency)
//   })
func MonitorHealth(
    ctx context.Context,
    client Client,
    interval time.Duration,
    callback func(*HealthCheck),
) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    
    // Run initial check immediately
    callback(Check(ctx, client))
    
    for {
        select {
        case <-ctx.Done():
            return
        
        case <-ticker.C:
            callback(Check(ctx, client))
        }
    }
}