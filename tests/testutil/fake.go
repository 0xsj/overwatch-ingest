// Package testutil provides testing utilities for the ingest service.
package testutil

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"
)

// Fake provides generators for fake test data.
var Fake = &fakeGenerator{}

type fakeGenerator struct {
	counter int64
}

// String generates a random string with the given prefix.
func (f *fakeGenerator) String(prefix string) string {
	f.counter++
	return fmt.Sprintf("%s_%d_%s", prefix, f.counter, f.randomHex(4))
}

// Hex generates a random hex string of the given byte length.
func (f *fakeGenerator) Hex(byteLength int) string {
	return f.randomHex(byteLength)
}

// ID generates a fake ULID-like string.
func (f *fakeGenerator) ID() string {
	return strings.ToUpper(f.randomHex(13))
}

// Nonce generates a cryptographic nonce.
func (f *fakeGenerator) Nonce(length int) string {
	return f.randomHex(length)
}

// DIDKey generates a fake did:key string.
func (f *fakeGenerator) DIDKey() string {
	return fmt.Sprintf("did:key:z6Mk%s", f.randomHex(32))
}

// RawDataID generates a fake raw data identifier.
func (f *fakeGenerator) RawDataID() string {
	f.counter++
	return fmt.Sprintf("raw_%d_%s", f.counter, f.randomHex(8))
}

// SourceType returns a random ingest source type string.
func (f *fakeGenerator) SourceType() string {
	types := []string{
		"ais",
		"adsb",
		"satellite",
		"radar",
		"sigint",
		"osint",
		"humint",
		"cyber",
	}
	return f.randomChoice(types)
}

// Payload generates a fake data payload with common ingest fields.
func (f *fakeGenerator) Payload() map[string]any {
	f.counter++
	return map[string]any{
		"id":        f.counter,
		"data":      f.randomHex(16),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
}

// VesselPayload generates a fake AIS vessel payload.
func (f *fakeGenerator) VesselPayload() map[string]any {
	f.counter++
	return map[string]any{
		"mmsi":      fmt.Sprintf("%d", 200000000+f.randomInt(0, 799999999)),
		"lat":       float64(f.randomInt(-90, 90)) + float64(f.randomInt(0, 999999))/1000000.0,
		"lon":       float64(f.randomInt(-180, 180)) + float64(f.randomInt(0, 999999))/1000000.0,
		"speed":     float64(f.randomInt(0, 30)),
		"course":    float64(f.randomInt(0, 360)),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
}

// AircraftPayload generates a fake ADS-B aircraft payload.
func (f *fakeGenerator) AircraftPayload() map[string]any {
	f.counter++
	return map[string]any{
		"icao":      fmt.Sprintf("%06X", f.randomInt(0, 16777215)),
		"callsign":  fmt.Sprintf("TST%04d", f.randomInt(0, 9999)),
		"lat":       float64(f.randomInt(-90, 90)) + float64(f.randomInt(0, 999999))/1000000.0,
		"lon":       float64(f.randomInt(-180, 180)) + float64(f.randomInt(0, 999999))/1000000.0,
		"altitude":  float64(f.randomInt(0, 45000)),
		"speed":     float64(f.randomInt(0, 600)),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
}

// LocationPayload generates a generic location payload.
func (f *fakeGenerator) LocationPayload() map[string]any {
	f.counter++
	return map[string]any{
		"lat":       float64(f.randomInt(-90, 90)) + float64(f.randomInt(0, 999999))/1000000.0,
		"lon":       float64(f.randomInt(-180, 180)) + float64(f.randomInt(0, 999999))/1000000.0,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
}

// Metadata generates fake metadata.
func (f *fakeGenerator) Metadata() map[string]string {
	return map[string]string{
		"content_type": "application/json",
		"source_url":   fmt.Sprintf("https://api-%d.example.com/v1/data", f.counter),
	}
}

// EntityType returns a random entity type.
func (f *fakeGenerator) EntityType() string {
	types := []string{"vessel", "aircraft", "location", "unknown"}
	return f.randomChoice(types)
}

// EntityID generates a fake entity identifier.
func (f *fakeGenerator) EntityID() string {
	f.counter++
	return fmt.Sprintf("entity_%d_%s", f.counter, f.randomHex(4))
}

// AnomalyMessage generates a fake anomaly description.
func (f *fakeGenerator) AnomalyMessage() string {
	messages := []string{
		"value out of expected range",
		"invalid format detected",
		"required field missing",
		"unexpected value encountered",
		"temporal anomaly: timestamp in the future",
		"statistical outlier detected",
		"duplicate record suspected",
		"suspicious data pattern",
	}
	return f.randomChoice(messages)
}

// FieldName generates a fake field name.
func (f *fakeGenerator) FieldName() string {
	fields := []string{
		"lat", "lon", "speed", "course", "altitude",
		"timestamp", "mmsi", "icao", "callsign",
		"heading", "status", "destination",
	}
	return f.randomChoice(fields)
}

// ConfidenceValue generates a random confidence score between 0 and 1.
func (f *fakeGenerator) ConfidenceValue() float64 {
	return float64(f.randomInt(0, 100)) / 100.0
}

// ReliabilityScore generates a random reliability score.
func (f *fakeGenerator) ReliabilityScore() float64 {
	return float64(f.randomInt(30, 100)) / 100.0
}

// ErrorMessage generates a fake error message.
func (f *fakeGenerator) ErrorMessage() string {
	messages := []string{
		"connection timeout",
		"validation failed: schema mismatch",
		"anomaly threshold exceeded",
		"signature verification failed",
		"confidence below threshold",
		"duplicate record detected",
		"source reliability too low",
	}
	return f.randomChoice(messages)
}

// FutureTime generates a time in the future.
func (f *fakeGenerator) FutureTime(maxOffset time.Duration) time.Time {
	offset := f.Duration(time.Minute, maxOffset)
	return time.Now().Add(offset)
}

// PastTime generates a time in the past.
func (f *fakeGenerator) PastTime(maxOffset time.Duration) time.Time {
	offset := f.Duration(time.Minute, maxOffset)
	return time.Now().Add(-offset)
}

// Duration generates a random duration between min and max.
func (f *fakeGenerator) Duration(min, max time.Duration) time.Duration {
	minNanos := min.Nanoseconds()
	maxNanos := max.Nanoseconds()
	deltaNanos := f.randomInt64(0, maxNanos-minNanos)
	return time.Duration(minNanos + deltaNanos)
}

// --- Helpers ---

func (f *fakeGenerator) randomHex(byteLength int) string {
	bytes := make([]byte, byteLength)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (f *fakeGenerator) randomChoice(choices []string) string {
	idx := f.randomInt(0, len(choices))
	return choices[idx]
}

func (f *fakeGenerator) randomInt(min, max int) int {
	if max <= min {
		return min
	}
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	return min + int(n.Int64())
}

func (f *fakeGenerator) randomInt64(min, max int64) int64 {
	if max <= min {
		return min
	}
	n, _ := rand.Int(rand.Reader, big.NewInt(max-min))
	return min + n.Int64()
}
