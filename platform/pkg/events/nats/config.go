// platform/pkg/events/nats/config.go
package nats

import (
	"time"

	"github.com/0xsj/scout/platform/pkg/events"
)

// Config implements events.Config for NATS.
type Config struct {
	url           string
	maxReconnects int
	reconnectWait time.Duration
}

// NewConfig creates a new NATS config.
func NewConfig(url string, maxReconnects int, reconnectWait time.Duration) *Config {
	return &Config{
		url:           url,
		maxReconnects: maxReconnects,
		reconnectWait: reconnectWait,
	}
}

// URL returns the NATS connection URL.
func (c *Config) URL() string {
	return c.url
}

// MaxReconnects returns the maximum number of reconnection attempts.
func (c *Config) MaxReconnects() int {
	return c.maxReconnects
}

// ReconnectWait returns the wait time between reconnection attempts.
func (c *Config) ReconnectWait() time.Duration {
	return c.reconnectWait
}

// Ensure Config implements events.Config interface
var _ events.Config = (*Config)(nil)