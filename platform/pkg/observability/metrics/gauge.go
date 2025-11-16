// platform/pkg/observability/metrics/gauge.go
package metrics

import (
	"sync"
)

// Gauge is a metric that can go up and down.
// Used for tracking memory usage, queue size, active connections, etc.
type Gauge interface {
	// Set sets the gauge to the given value
	Set(value float64)

	// Inc increments the gauge by 1
	Inc()

	// Dec decrements the gauge by 1
	Dec()

	// Add adds the given value to the gauge
	Add(delta float64)

	// Sub subtracts the given value from the gauge
	Sub(delta float64)

	// Value returns the current gauge value
	Value() float64

	// Name returns the gauge name
	Name() string

	// Labels returns the gauge labels
	Labels() map[string]string
}

// gauge implements Gauge.
type gauge struct {
	name   string
	help   string
	labels map[string]string
	value  float64
	mu     sync.RWMutex
}

// NewGauge creates a new gauge.
func NewGauge(name, help string, labels map[string]string) Gauge {
	if labels == nil {
		labels = make(map[string]string)
	}
	return &gauge{
		name:   name,
		help:   help,
		labels: labels,
		value:  0,
	}
}

// Set sets the gauge to the given value.
func (g *gauge) Set(value float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value = value
}

// Inc increments the gauge by 1.
func (g *gauge) Inc() {
	g.Add(1)
}

// Dec decrements the gauge by 1.
func (g *gauge) Dec() {
	g.Sub(1)
}

// Add adds the given value to the gauge.
func (g *gauge) Add(delta float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value += delta
}

// Sub subtracts the given value from the gauge.
func (g *gauge) Sub(delta float64) {
	g.Add(-delta)
}

// Value returns the current gauge value.
func (g *gauge) Value() float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.value
}

// Name returns the gauge name.
func (g *gauge) Name() string {
	return g.name
}

// Labels returns the gauge labels.
func (g *gauge) Labels() map[string]string {
	return g.labels
}