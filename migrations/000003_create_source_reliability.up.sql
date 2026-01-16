CREATE TABLE source_reliability (
    source_id             VARCHAR(36) PRIMARY KEY,
    tenant_id             VARCHAR(36),
    
    reliability_score     REAL NOT NULL DEFAULT 0.5,
    
    -- Counters
    total_records         BIGINT NOT NULL DEFAULT 0,
    accepted_records      BIGINT NOT NULL DEFAULT 0,
    rejected_records      BIGINT NOT NULL DEFAULT 0,
    quarantined_records   BIGINT NOT NULL DEFAULT 0,
    corroborated_records  BIGINT NOT NULL DEFAULT 0,
    disputed_records      BIGINT NOT NULL DEFAULT 0,
    
    -- Time window for stats
    calculated_at         TIMESTAMPTZ NOT NULL,
    window_start          TIMESTAMPTZ NOT NULL,
    window_end            TIMESTAMPTZ NOT NULL,
    
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_source_reliability_tenant_id ON source_reliability(tenant_id) WHERE tenant_id IS NOT NULL;
CREATE INDEX idx_source_reliability_score ON source_reliability(reliability_score);
CREATE INDEX idx_source_reliability_tenant_score ON source_reliability(tenant_id, reliability_score);