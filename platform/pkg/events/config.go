// platform/pkg/events/config.go
package events

import (
	"time"
)

// Config defines common configuration for event bus providers.
type Config interface {
	// URL returns the connection URL
	URL() string

	// MaxReconnects returns the maximum number of reconnection attempts
	MaxReconnects() int

	// ReconnectWait returns the wait time between reconnection attempts
	ReconnectWait() time.Duration
}