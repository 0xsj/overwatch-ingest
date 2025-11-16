// platform/pkg/observability/metrics/registry.go
package metrics

import (
	"sort"
	"strings"
	"sync"
)

// Registry stores and manages metrics.
type Registry struct {
	counters   map[string]Counter
	histograms map[string]Histogram
	gauges     map[string]Gauge
	mu         sync.RWMutex
}

// NewRegistry creates a new metric registry.
func NewRegistry() *Registry {
	return &Registry{
		counters:   make(map[string]Counter),
		histograms: make(map[string]Histogram),
		gauges:     make(map[string]Gauge),
	}
}

// DefaultRegistry is the global default registry.
var DefaultRegistry = NewRegistry()

// RegisterCounter registers a new counter.
func (r *Registry) RegisterCounter(name, help string, labels map[string]string) Counter {
	key := metricKey(name, labels)
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if existing, ok := r.counters[key]; ok {
		return existing
	}
	
	counter := NewCounter(name, help, labels)
	r.counters[key] = counter
	return counter
}

// RegisterHistogram registers a new histogram with default buckets.
func (r *Registry) RegisterHistogram(name, help string, labels map[string]string) Histogram {
	return r.RegisterHistogramWithBuckets(name, help, labels, DefaultBuckets)
}

// RegisterHistogramWithBuckets registers a new histogram with custom buckets.
func (r *Registry) RegisterHistogramWithBuckets(name, help string, labels map[string]string, buckets []float64) Histogram {
	key := metricKey(name, labels)
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if existing, ok := r.histograms[key]; ok {
		return existing
	}
	
	histogram := NewHistogramWithBuckets(name, help, labels, buckets)
	r.histograms[key] = histogram
	return histogram
}

// RegisterGauge registers a new gauge.
func (r *Registry) RegisterGauge(name, help string, labels map[string]string) Gauge {
	key := metricKey(name, labels)
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if existing, ok := r.gauges[key]; ok {
		return existing
	}
	
	gauge := NewGauge(name, help, labels)
	r.gauges[key] = gauge
	return gauge
}

// GetCounter retrieves a counter by name and labels.
func (r *Registry) GetCounter(name string, labels map[string]string) (Counter, bool) {
	key := metricKey(name, labels)
	r.mu.RLock()
	defer r.mu.RUnlock()
	counter, ok := r.counters[key]
	return counter, ok
}

// GetHistogram retrieves a histogram by name and labels.
func (r *Registry) GetHistogram(name string, labels map[string]string) (Histogram, bool) {
	key := metricKey(name, labels)
	r.mu.RLock()
	defer r.mu.RUnlock()
	histogram, ok := r.histograms[key]
	return histogram, ok
}

// GetGauge retrieves a gauge by name and labels.
func (r *Registry) GetGauge(name string, labels map[string]string) (Gauge, bool) {
	key := metricKey(name, labels)
	r.mu.RLock()
	defer r.mu.RUnlock()
	gauge, ok := r.gauges[key]
	return gauge, ok
}

// Counters returns all registered counters.
func (r *Registry) Counters() []Counter {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	result := make([]Counter, 0, len(r.counters))
	for _, c := range r.counters {
		result = append(result, c)
	}
	return result
}

// Histograms returns all registered histograms.
func (r *Registry) Histograms() []Histogram {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	result := make([]Histogram, 0, len(r.histograms))
	for _, h := range r.histograms {
		result = append(result, h)
	}
	return result
}

// Gauges returns all registered gauges.
func (r *Registry) Gauges() []Gauge {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	result := make([]Gauge, 0, len(r.gauges))
	for _, g := range r.gauges {
		result = append(result, g)
	}
	return result
}

// metricKey generates a unique key for a metric based on name and labels.
func metricKey(name string, labels map[string]string) string {
	if len(labels) == 0 {
		return name
	}
	
	// Sort labels for consistent keys
	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	
	var sb strings.Builder
	sb.WriteString(name)
	sb.WriteString("{")
	for i, k := range keys {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(labels[k])
	}
	sb.WriteString("}")
	
	return sb.String()
}

// Package-level convenience functions using DefaultRegistry

// RegisterCounter registers a counter in the default registry.
func RegisterCounter(name, help string, labels map[string]string) Counter {
	return DefaultRegistry.RegisterCounter(name, help, labels)
}

// RegisterHistogram registers a histogram in the default registry.
func RegisterHistogram(name, help string, labels map[string]string) Histogram {
	return DefaultRegistry.RegisterHistogram(name, help, labels)
}

// RegisterGauge registers a gauge in the default registry.
func RegisterGauge(name, help string, labels map[string]string) Gauge {
	return DefaultRegistry.RegisterGauge(name, help, labels)
}