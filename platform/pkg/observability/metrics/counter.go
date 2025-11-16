// platform/pkg/observability/metrics/counter.go
package metrics

import (
	"sync"
)

// Counter is a monotonically increasing metric.
// Used for counting events (requests, errors, messages published, etc.)
type Counter interface {
	// Inc increments the counter by 1
	Inc()

	// Add adds the given value to the counter
	Add(delta float64)

	// Value returns the current counter value
	Value() float64

	// Name returns the counter name
	Name() string

	// Labels returns the counter labels
	Labels() map[string]string
}

// counter implements Counter.
type counter struct {
	name   string
	help   string
	labels map[string]string
	value  float64
	mu     sync.RWMutex
}

// NewCounter creates a new counter.
func NewCounter(name, help string, labels map[string]string) Counter {
	if labels == nil {
		labels = make(map[string]string)
	}
	return &counter{
		name:   name,
		help:   help,
		labels: labels,
		value:  0,
	}
}

// Inc increments the counter by 1.
func (c *counter) Inc() {
	c.Add(1)
}

// Add adds the given value to the counter.
func (c *counter) Add(delta float64) {
	if delta < 0 {
		// Counters can only increase
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value += delta
}

// Value returns the current counter value.
func (c *counter) Value() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value
}

// Name returns the counter name.
func (c *counter) Name() string {
	return c.name
}

// Labels returns the counter labels.
func (c *counter) Labels() map[string]string {
	return c.labels
}