CREATE TYPE quarantine_reason AS ENUM (
    'validation_failed',
    'low_confidence',
    'anomaly_detected',
    'signature_invalid',
    'duplicate_suspected',
    'manual_review'
);

CREATE TYPE quarantine_resolution AS ENUM (
    'pending',
    'approved',
    'modified',
    'rejected',
    'expired'
);

CREATE TABLE quarantined_records (
    id                  VARCHAR(36) PRIMARY KEY,
    tenant_id           VARCHAR(36),
    source_id           VARCHAR(36) NOT NULL,
    source_type         VARCHAR(100) NOT NULL,
    raw_data_id         VARCHAR(255) NOT NULL,
    ingest_record_id    VARCHAR(36) NOT NULL REFERENCES ingest_records(id),
    
    -- Raw data preserved for review
    raw_data            JSONB NOT NULL,
    
    -- Quarantine reason
    reason              quarantine_reason NOT NULL,
    reason_detail       TEXT NOT NULL DEFAULT '',
    anomalies           JSONB NOT NULL DEFAULT '[]',
    
    -- Confidence at time of quarantine
    confidence_overall            REAL NOT NULL DEFAULT 0.0,
    confidence_source_reliability REAL NOT NULL DEFAULT 0.0,
    confidence_data_completeness  REAL NOT NULL DEFAULT 0.0,
    confidence_temporal_freshness REAL NOT NULL DEFAULT 0.0,
    confidence_signature_trust    REAL NOT NULL DEFAULT 0.0,
    confidence_factors            JSONB NOT NULL DEFAULT '[]',
    
    -- Resolution
    resolution          quarantine_resolution NOT NULL DEFAULT 'pending',
    resolved_by         VARCHAR(255),
    resolved_by_did     VARCHAR(255),
    resolution_notes    TEXT,
    modified_data       JSONB,
    
    -- Provenance
    ingest_signer       JSONB,
    resolver_signer     JSONB,
    
    -- Timestamps
    quarantined_at      TIMESTAMPTZ NOT NULL,
    expires_at          TIMESTAMPTZ,
    resolved_at         TIMESTAMPTZ,
    
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for common query patterns
CREATE INDEX idx_quarantined_records_tenant_id ON quarantined_records(tenant_id) WHERE tenant_id IS NOT NULL;
CREATE INDEX idx_quarantined_records_source_id ON quarantined_records(source_id);
CREATE INDEX idx_quarantined_records_source_type ON quarantined_records(source_type);
CREATE INDEX idx_quarantined_records_ingest_record_id ON quarantined_records(ingest_record_id);
CREATE INDEX idx_quarantined_records_reason ON quarantined_records(reason);
CREATE INDEX idx_quarantined_records_resolution ON quarantined_records(resolution);
CREATE INDEX idx_quarantined_records_quarantined_at ON quarantined_records(quarantined_at);
CREATE INDEX idx_quarantined_records_expires_at ON quarantined_records(expires_at) WHERE expires_at IS NOT NULL;

-- Composite indexes
CREATE INDEX idx_quarantined_records_tenant_resolution ON quarantined_records(tenant_id, resolution);
CREATE INDEX idx_quarantined_records_pending ON quarantined_records(resolution, quarantined_at) WHERE resolution = 'pending';