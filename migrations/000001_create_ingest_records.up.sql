CREATE TYPE ingest_status AS ENUM (
    'pending',
    'accepted',
    'quarantined',
    'rejected'
);

CREATE TABLE ingest_records (
    id                            VARCHAR(36) PRIMARY KEY,
    tenant_id                     VARCHAR(36),
    source_id                     VARCHAR(36) NOT NULL,
    source_type                   VARCHAR(100) NOT NULL,
    raw_data_id                   VARCHAR(255) NOT NULL,
    
    status                        ingest_status NOT NULL DEFAULT 'pending',
    
    -- Validation result (stored as JSONB)
    validation_valid              BOOLEAN NOT NULL DEFAULT false,
    validation_schema_valid       BOOLEAN NOT NULL DEFAULT false,
    validation_fields_present     TEXT[] NOT NULL DEFAULT '{}',
    validation_fields_missing     TEXT[] NOT NULL DEFAULT '{}',
    validation_anomalies          JSONB NOT NULL DEFAULT '[]',
    validation_validator_version  VARCHAR(50),
    
    -- Confidence score
    confidence_overall            REAL NOT NULL DEFAULT 0.0,
    confidence_source_reliability REAL NOT NULL DEFAULT 0.0,
    confidence_data_completeness  REAL NOT NULL DEFAULT 0.0,
    confidence_temporal_freshness REAL NOT NULL DEFAULT 0.0,
    confidence_signature_trust    REAL NOT NULL DEFAULT 0.0,
    confidence_factors            JSONB NOT NULL DEFAULT '[]',
    
    -- Entity routing
    entity_type                   VARCHAR(100),
    entity_id                     VARCHAR(255),
    event_ids                     TEXT[] NOT NULL DEFAULT '{}',
    
    -- Rejection/Quarantine
    rejection_reason              TEXT,
    quarantine_id                 VARCHAR(36),
    
    -- Provenance (stored as JSONB for flexibility)
    source_signer                 JSONB,
    collector_signer              JSONB,
    ingest_signer                 JSONB,
    source_signature_verified     BOOLEAN,
    collector_signature_verified  BOOLEAN,
    
    -- Timestamps
    received_at                   TIMESTAMPTZ NOT NULL,
    processed_at                  TIMESTAMPTZ NOT NULL,
    
    -- Indexes will reference these
    created_at                    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for common query patterns
CREATE INDEX idx_ingest_records_tenant_id ON ingest_records(tenant_id) WHERE tenant_id IS NOT NULL;
CREATE INDEX idx_ingest_records_source_id ON ingest_records(source_id);
CREATE INDEX idx_ingest_records_source_type ON ingest_records(source_type);
CREATE INDEX idx_ingest_records_status ON ingest_records(status);
CREATE INDEX idx_ingest_records_raw_data_id ON ingest_records(raw_data_id);
CREATE INDEX idx_ingest_records_entity_type ON ingest_records(entity_type) WHERE entity_type IS NOT NULL;
CREATE INDEX idx_ingest_records_received_at ON ingest_records(received_at);
CREATE INDEX idx_ingest_records_processed_at ON ingest_records(processed_at);

-- Composite indexes for common filtered queries
CREATE INDEX idx_ingest_records_tenant_status ON ingest_records(tenant_id, status);
CREATE INDEX idx_ingest_records_source_status ON ingest_records(source_id, status);
CREATE INDEX idx_ingest_records_tenant_source ON ingest_records(tenant_id, source_id);