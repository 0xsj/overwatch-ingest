package service

import (
	"fmt"
	"strings"
)

type EntityInferenceService interface {
	InferEntityType(sourceType string, payload map[string]any) string
	InferEntityID(sourceType string, payload map[string]any) string
}

type entityInferenceService struct {
	typeRules map[string][]EntityTypeRule
	idRules   map[string][]EntityIDRule
}

type EntityTypeRule struct {
	Field      string
	EntityType string
}

type EntityIDRule struct {
	Field  string
	Prefix string
}

func NewEntityInferenceService() EntityInferenceService {
	return &entityInferenceService{
		typeRules: defaultTypeRules(),
		idRules:   defaultIDRules(),
	}
}

func defaultTypeRules() map[string][]EntityTypeRule {
	return map[string][]EntityTypeRule{
		"ais": {
			{Field: "mmsi", EntityType: "vessel"},
		},
		"adsb": {
			{Field: "icao", EntityType: "aircraft"},
			{Field: "icao24", EntityType: "aircraft"},
		},
		"gps": {
			{Field: "device_id", EntityType: "device"},
		},
		"social": {
			{Field: "user_id", EntityType: "account"},
			{Field: "username", EntityType: "account"},
		},
		"*": {
			{Field: "mmsi", EntityType: "vessel"},
			{Field: "icao", EntityType: "aircraft"},
			{Field: "icao24", EntityType: "aircraft"},
			{Field: "device_id", EntityType: "device"},
			{Field: "user_id", EntityType: "account"},
			{Field: "lat", EntityType: "location"},
		},
	}
}

func defaultIDRules() map[string][]EntityIDRule {
	return map[string][]EntityIDRule{
		"ais": {
			{Field: "mmsi", Prefix: "mmsi"},
		},
		"adsb": {
			{Field: "icao", Prefix: "icao"},
			{Field: "icao24", Prefix: "icao"},
		},
		"gps": {
			{Field: "device_id", Prefix: "device"},
		},
		"social": {
			{Field: "user_id", Prefix: "user"},
			{Field: "username", Prefix: "user"},
		},
		"*": {
			{Field: "mmsi", Prefix: "mmsi"},
			{Field: "icao", Prefix: "icao"},
			{Field: "icao24", Prefix: "icao"},
			{Field: "device_id", Prefix: "device"},
			{Field: "user_id", Prefix: "user"},
		},
	}
}

func (s *entityInferenceService) InferEntityType(sourceType string, payload map[string]any) string {
	normalizedType := strings.ToLower(sourceType)

	if rules, ok := s.typeRules[normalizedType]; ok {
		for _, rule := range rules {
			if _, exists := payload[rule.Field]; exists {
				return rule.EntityType
			}
		}
	}

	if rules, ok := s.typeRules["*"]; ok {
		for _, rule := range rules {
			if _, exists := payload[rule.Field]; exists {
				return rule.EntityType
			}
		}
	}

	return "unknown"
}

func (s *entityInferenceService) InferEntityID(sourceType string, payload map[string]any) string {
	normalizedType := strings.ToLower(sourceType)

	if rules, ok := s.idRules[normalizedType]; ok {
		for _, rule := range rules {
			if value, exists := payload[rule.Field]; exists {
				return formatEntityID(rule.Prefix, value)
			}
		}
	}

	if rules, ok := s.idRules["*"]; ok {
		for _, rule := range rules {
			if value, exists := payload[rule.Field]; exists {
				return formatEntityID(rule.Prefix, value)
			}
		}
	}

	return "unknown"
}

func formatEntityID(prefix string, value any) string {
	return fmt.Sprintf("%s:%s", prefix, toString(value))
}

func toString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%.0f", val)
		}
		return fmt.Sprintf("%g", val)
	case int:
		return fmt.Sprintf("%d", val)
	case int64:
		return fmt.Sprintf("%d", val)
	case int32:
		return fmt.Sprintf("%d", val)
	default:
		return fmt.Sprintf("%v", v)
	}
}
