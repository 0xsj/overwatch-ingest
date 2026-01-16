package validation

import (
	"context"
	"math"
	"time"

	"github.com/0xsj/overwatch-ingest/internal/domain/model"
	"github.com/0xsj/overwatch-ingest/internal/port/outbound/validation"
)

const anomalyDetectorVersion = "1.0.0"

type anomalyDetector struct {
	registry *SourceTypeRegistry
}

func NewAnomalyDetector(registry *SourceTypeRegistry) validation.AnomalyDetector {
	if registry == nil {
		registry = NewSourceTypeRegistry()
	}
	return &anomalyDetector{
		registry: registry,
	}
}

func (d *anomalyDetector) Detect(ctx context.Context, sourceType string, payload map[string]any, metadata validation.DetectionMetadata) validation.AnomalyDetectionResult {
	var anomalies []model.Anomaly

	anomalies = append(anomalies, d.detectTemporalAnomalies(payload, metadata)...)
	anomalies = append(anomalies, d.detectRangeAnomalies(sourceType, payload)...)
	anomalies = append(anomalies, d.detectVelocityAnomalies(sourceType, payload, metadata)...)
	anomalies = append(anomalies, d.detectSuspiciousPatterns(sourceType, payload, metadata)...)

	return validation.NewAnomalyDetectionResult(anomalies, anomalyDetectorVersion)
}

func (d *anomalyDetector) SupportsSourceType(sourceType string) bool {
	return d.registry.Supports(sourceType)
}

func (d *anomalyDetector) Version() string {
	return anomalyDetectorVersion
}

func (d *anomalyDetector) detectTemporalAnomalies(payload map[string]any, metadata validation.DetectionMetadata) []model.Anomaly {
	var anomalies []model.Anomaly

	ts, ok := extractTimestamp(payload, "timestamp")
	if !ok {
		return anomalies
	}

	now := time.Now()

	if ts.After(now.Add(5 * time.Minute)) {
		anomalies = append(anomalies, model.TemporalAnomaly(
			"timestamp",
			model.AnomalySeverityWarning,
			"timestamp is in the future",
			ts,
		))
	}

	if ts.Before(now.Add(-24 * time.Hour)) {
		anomalies = append(anomalies, model.TemporalAnomaly(
			"timestamp",
			model.AnomalySeverityInfo,
			"timestamp is more than 24 hours old",
			ts,
		))
	}

	if ts.Before(now.Add(-30 * 24 * time.Hour)) {
		anomalies = append(anomalies, model.TemporalAnomaly(
			"timestamp",
			model.AnomalySeverityWarning,
			"timestamp is more than 30 days old",
			ts,
		))
	}

	if metadata.SourceTimestamp.IsPresent() {
		sourceTs := metadata.SourceTimestamp.MustGet().Time()
		if !ts.Equal(sourceTs) {
			diff := ts.Sub(sourceTs).Abs()
			if diff > 1*time.Hour {
				anomalies = append(anomalies, model.TemporalAnomaly(
					"timestamp",
					model.AnomalySeverityWarning,
					"significant time drift between source and payload timestamp",
					ts,
				).WithContext("drift_seconds", diff.Seconds()))
			}
		}
	}

	return anomalies
}

func (d *anomalyDetector) detectRangeAnomalies(sourceType string, payload map[string]any) []model.Anomaly {
	var anomalies []model.Anomaly

	switch sourceType {
	case "ais":
		anomalies = append(anomalies, d.detectAISRangeAnomalies(payload)...)
	case "satellite":
		anomalies = append(anomalies, d.detectSatelliteRangeAnomalies(payload)...)
	case "sigint":
		anomalies = append(anomalies, d.detectSigintRangeAnomalies(payload)...)
	}

	if lat, ok := getFloat64(payload, "latitude"); ok {
		if lat < -90 || lat > 90 {
			anomalies = append(anomalies, model.OutOfRangeAnomaly(
				"latitude",
				model.AnomalySeverityError,
				-90.0, 90.0, lat,
			))
		}
	}

	if lon, ok := getFloat64(payload, "longitude"); ok {
		if lon < -180 || lon > 180 {
			anomalies = append(anomalies, model.OutOfRangeAnomaly(
				"longitude",
				model.AnomalySeverityError,
				-180.0, 180.0, lon,
			))
		}
	}

	return anomalies
}

func (d *anomalyDetector) detectAISRangeAnomalies(payload map[string]any) []model.Anomaly {
	var anomalies []model.Anomaly

	if speed, ok := getFloat64(payload, "speed"); ok {
		if speed < 0 {
			anomalies = append(anomalies, model.OutOfRangeAnomaly(
				"speed",
				model.AnomalySeverityError,
				0.0, 102.2, speed,
			))
		} else if speed > 50 {
			anomalies = append(anomalies, model.SuspiciousAnomaly(
				"speed",
				model.AnomalySeverityWarning,
				"unusually high speed for vessel",
				map[string]any{"speed": speed},
			))
		}
	}

	if course, ok := getFloat64(payload, "course"); ok {
		if course < 0 || course >= 360 {
			anomalies = append(anomalies, model.OutOfRangeAnomaly(
				"course",
				model.AnomalySeverityError,
				0.0, 359.9, course,
			))
		}
	}

	if heading, ok := getFloat64(payload, "heading"); ok {
		if heading < 0 || heading >= 360 {
			anomalies = append(anomalies, model.OutOfRangeAnomaly(
				"heading",
				model.AnomalySeverityError,
				0.0, 359.9, heading,
			))
		}
	}

	return anomalies
}

func (d *anomalyDetector) detectSatelliteRangeAnomalies(payload map[string]any) []model.Anomaly {
	var anomalies []model.Anomaly

	if cloudCover, ok := getFloat64(payload, "cloud_cover"); ok {
		if cloudCover < 0 || cloudCover > 100 {
			anomalies = append(anomalies, model.OutOfRangeAnomaly(
				"cloud_cover",
				model.AnomalySeverityError,
				0.0, 100.0, cloudCover,
			))
		}
	}

	if resolution, ok := getFloat64(payload, "resolution"); ok {
		if resolution <= 0 {
			anomalies = append(anomalies, model.OutOfRangeAnomaly(
				"resolution",
				model.AnomalySeverityError,
				0.0, "positive", resolution,
			))
		}
	}

	return anomalies
}

func (d *anomalyDetector) detectSigintRangeAnomalies(payload map[string]any) []model.Anomaly {
	var anomalies []model.Anomaly

	if freq, ok := getFloat64(payload, "frequency"); ok {
		if freq <= 0 {
			anomalies = append(anomalies, model.OutOfRangeAnomaly(
				"frequency",
				model.AnomalySeverityError,
				0.0, "positive", freq,
			))
		}
	}

	if strength, ok := getFloat64(payload, "strength"); ok {
		if strength < -200 || strength > 50 {
			anomalies = append(anomalies, model.OutOfRangeAnomaly(
				"strength",
				model.AnomalySeverityCritical,
				-200.0, 50.0, strength,
			))
		}
	}

	return anomalies
}

func (d *anomalyDetector) detectVelocityAnomalies(sourceType string, payload map[string]any, metadata validation.DetectionMetadata) []model.Anomaly {
	var anomalies []model.Anomaly

	if metadata.PreviousPayload == nil {
		return anomalies
	}

	if sourceType != "ais" {
		return anomalies
	}

	currLat, currLatOk := getFloat64(payload, "latitude")
	currLon, currLonOk := getFloat64(payload, "longitude")
	prevLat, prevLatOk := getFloat64(metadata.PreviousPayload, "latitude")
	prevLon, prevLonOk := getFloat64(metadata.PreviousPayload, "longitude")

	if !currLatOk || !currLonOk || !prevLatOk || !prevLonOk {
		return anomalies
	}

	currTs, currTsOk := extractTimestamp(payload, "timestamp")
	prevTs, prevTsOk := extractTimestamp(metadata.PreviousPayload, "timestamp")

	if !currTsOk || !prevTsOk {
		return anomalies
	}

	timeDiff := currTs.Sub(prevTs)
	if timeDiff <= 0 {
		return anomalies
	}

	distanceKm := haversineDistance(prevLat, prevLon, currLat, currLon)
	speedKnots := (distanceKm / timeDiff.Hours()) * 0.539957

	if speedKnots > 60 {
		anomalies = append(anomalies, model.SuspiciousAnomaly(
			"position",
			model.AnomalySeverityWarning,
			"implied speed between positions exceeds plausible vessel speed",
			map[string]any{
				"implied_speed_knots": speedKnots,
				"distance_km":         distanceKm,
				"time_diff_seconds":   timeDiff.Seconds(),
			},
		))
	}

	return anomalies
}

func (d *anomalyDetector) detectSuspiciousPatterns(sourceType string, payload map[string]any, metadata validation.DetectionMetadata) []model.Anomaly {
	var anomalies []model.Anomaly

	if sourceType == "ais" {
		if lat, latOk := getFloat64(payload, "latitude"); latOk {
			if lon, lonOk := getFloat64(payload, "longitude"); lonOk {
				if lat == 0 && lon == 0 {
					anomalies = append(anomalies, model.SuspiciousAnomaly(
						"position",
						model.AnomalySeverityWarning,
						"position at null island (0,0) is suspicious",
						nil,
					))
				}
			}
		}

		if mmsi, ok := getString(payload, "mmsi"); ok {
			if len(mmsi) == 9 && (mmsi[0] == '0' || mmsi == "123456789") {
				anomalies = append(anomalies, model.SuspiciousAnomaly(
					"mmsi",
					model.AnomalySeverityWarning,
					"MMSI appears to be invalid or placeholder",
					map[string]any{"mmsi": mmsi},
				))
			}
		}
	}

	return anomalies
}

func extractTimestamp(payload map[string]any, field string) (time.Time, bool) {
	val, ok := payload[field]
	if !ok {
		return time.Time{}, false
	}

	switch v := val.(type) {
	case time.Time:
		return v, true
	case string:
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			return t, true
		}
		if t, err := time.Parse(time.RFC3339Nano, v); err == nil {
			return t, true
		}
	case int64:
		return time.Unix(v, 0), true
	case float64:
		return time.Unix(int64(v), 0), true
	}

	return time.Time{}, false
}

func getFloat64(payload map[string]any, field string) (float64, bool) {
	val, ok := payload[field]
	if !ok {
		return 0, false
	}
	return toFloat64(val)
}

func getString(payload map[string]any, field string) (string, bool) {
	val, ok := payload[field]
	if !ok {
		return "", false
	}
	return toString(val)
}

func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}
