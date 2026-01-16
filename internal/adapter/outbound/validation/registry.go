package validation

type SourceTypeConfig struct {
	SourceType      string
	RequiredFields  []string
	OptionalFields  []string
	SchemaID        string
	SchemaVersion   string
	BaseReliability float64
}

var defaultSourceTypes = map[string]SourceTypeConfig{
	"ais": {
		SourceType:      "ais",
		RequiredFields:  []string{"mmsi", "latitude", "longitude", "timestamp"},
		OptionalFields:  []string{"course", "speed", "heading", "name", "imo", "callsign", "vessel_type", "destination"},
		SchemaID:        "ais-v1",
		SchemaVersion:   "1.0.0",
		BaseReliability: 0.8,
	},
	"social_media": {
		SourceType:      "social_media",
		RequiredFields:  []string{"platform", "content", "author_id", "timestamp"},
		OptionalFields:  []string{"location", "hashtags", "mentions", "media_urls", "engagement"},
		SchemaID:        "social-media-v1",
		SchemaVersion:   "1.0.0",
		BaseReliability: 0.5,
	},
	"protest": {
		SourceType:      "protest",
		RequiredFields:  []string{"location", "timestamp", "source_url"},
		OptionalFields:  []string{"estimated_size", "description", "organizers", "demands", "images"},
		SchemaID:        "protest-v1",
		SchemaVersion:   "1.0.0",
		BaseReliability: 0.6,
	},
	"disaster": {
		SourceType:      "disaster",
		RequiredFields:  []string{"type", "location", "timestamp", "severity"},
		OptionalFields:  []string{"affected_area", "casualties", "damage_estimate", "source_url", "images"},
		SchemaID:        "disaster-v1",
		SchemaVersion:   "1.0.0",
		BaseReliability: 0.7,
	},
	"satellite": {
		SourceType:      "satellite",
		RequiredFields:  []string{"image_id", "capture_time", "bounding_box", "resolution"},
		OptionalFields:  []string{"cloud_cover", "sensor_type", "bands", "metadata"},
		SchemaID:        "satellite-v1",
		SchemaVersion:   "1.0.0",
		BaseReliability: 0.9,
	},
	"sigint": {
		SourceType:      "sigint",
		RequiredFields:  []string{"signal_type", "frequency", "timestamp", "location"},
		OptionalFields:  []string{"duration", "strength", "modulation", "callsign", "metadata"},
		SchemaID:        "sigint-v1",
		SchemaVersion:   "1.0.0",
		BaseReliability: 0.85,
	},
	"humint": {
		SourceType:      "humint",
		RequiredFields:  []string{"report_id", "timestamp", "classification"},
		OptionalFields:  []string{"location", "summary", "reliability_rating", "source_rating"},
		SchemaID:        "humint-v1",
		SchemaVersion:   "1.0.0",
		BaseReliability: 0.6,
	},
	"generic": {
		SourceType:      "generic",
		RequiredFields:  []string{"timestamp"},
		OptionalFields:  []string{"data"},
		SchemaID:        "generic-v1",
		SchemaVersion:   "1.0.0",
		BaseReliability: 0.5,
	},
}

type SourceTypeRegistry struct {
	configs map[string]SourceTypeConfig
}

func NewSourceTypeRegistry() *SourceTypeRegistry {
	configs := make(map[string]SourceTypeConfig, len(defaultSourceTypes))
	for k, v := range defaultSourceTypes {
		configs[k] = v
	}
	return &SourceTypeRegistry{configs: configs}
}

func (r *SourceTypeRegistry) Get(sourceType string) (SourceTypeConfig, bool) {
	cfg, ok := r.configs[sourceType]
	return cfg, ok
}

func (r *SourceTypeRegistry) GetOrDefault(sourceType string) SourceTypeConfig {
	if cfg, ok := r.configs[sourceType]; ok {
		return cfg
	}
	return r.configs["generic"]
}

func (r *SourceTypeRegistry) Register(cfg SourceTypeConfig) {
	r.configs[cfg.SourceType] = cfg
}

func (r *SourceTypeRegistry) Supports(sourceType string) bool {
	_, ok := r.configs[sourceType]
	return ok
}

func (r *SourceTypeRegistry) AllSourceTypes() []string {
	types := make([]string, 0, len(r.configs))
	for k := range r.configs {
		types = append(types, k)
	}
	return types
}
