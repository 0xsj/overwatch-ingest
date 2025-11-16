// platform/pkg/observability/metrics/histogram.go
package metrics

import (
	"math"
	"sync"
)

// Histogram tracks the distribution of values.
// Used for tracking latencies, request sizes, etc.
type Histogram interface {
	// Observe records a value
	Observe(value float64)

	// Count returns the total number of observations
	Count() uint64

	// Sum returns the sum of all observed values
	Sum() float64

	// Buckets returns the histogram buckets with counts
	Buckets() map[float64]uint64

	// Name returns the histogram name
	Name() string

	// Labels returns the histogram labels
	Labels() map[string]string
}

// histogram implements Histogram.
type histogram struct {
	name    string
	help    string
	labels  map[string]string
	buckets []float64
	counts  []uint64
	sum     float64
	count   uint64
	mu      sync.RWMutex
}

// DefaultBuckets are the default histogram buckets (in seconds).
// Suitable for HTTP request latencies: 5ms, 10ms, 25ms, 50ms, 100ms, 250ms, 500ms, 1s, 2.5s, 5s, 10s
var DefaultBuckets = []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}

// NewHistogram creates a new histogram with default buckets.
func NewHistogram(name, help string, labels map[string]string) Histogram {
	return NewHistogramWithBuckets(name, help, labels, DefaultBuckets)
}

// NewHistogramWithBuckets creates a new histogram with custom buckets.
func NewHistogramWithBuckets(name, help string, labels map[string]string, buckets []float64) Histogram {
	if labels == nil {
		labels = make(map[string]string)
	}

	// Sort buckets and add +Inf
	sortedBuckets := make([]float64, len(buckets)+1)
	copy(sortedBuckets, buckets)
	sortedBuckets[len(buckets)] = math.Inf(1)

	return &histogram{
		name:    name,
		help:    help,
		labels:  labels,
		buckets: sortedBuckets,
		counts:  make([]uint64, len(sortedBuckets)),
		sum:     0,
		count:   0,
	}
}

// Observe records a value.
func (h *histogram) Observe(value float64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.sum += value
	h.count++

	// Find bucket and increment
	for i, bucket := range h.buckets {
		if value <= bucket {
			h.counts[i]++
		}
	}
}

// Count returns the total number of observations.
func (h *histogram) Count() uint64 {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.count
}

// Sum returns the sum of all observed values.
func (h *histogram) Sum() float64 {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.sum
}

// Buckets returns the histogram buckets with counts.
func (h *histogram) Buckets() map[float64]uint64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make(map[float64]uint64)
	for i, bucket := range h.buckets {
		result[bucket] = h.counts[i]
	}
	return result
}

// Name returns the histogram name.
func (h *histogram) Name() string {
	return h.name
}

// Labels returns the histogram labels.
func (h *histogram) Labels() map[string]string {
	return h.labels
}