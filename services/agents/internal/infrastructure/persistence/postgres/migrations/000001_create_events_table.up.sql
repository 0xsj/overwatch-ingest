-- services/agents/internal/infrastructure/persistence/postgres/migrations/000002_create_events_table.up.sql
-- Event Store Schema for Agents Service

-- Events table stores all domain events (source of truth)
CREATE TABLE IF NOT EXISTS agent_events (
    id BIGSERIAL PRIMARY KEY,
    
    -- Aggregate identification
    aggregate_id UUID NOT NULL,
    aggregate_type VARCHAR(50) NOT NULL DEFAULT 'agent',
    
    -- Event metadata
    event_type VARCHAR(100) NOT NULL,
    event_version INT NOT NULL DEFAULT 1,
    sequence_number BIGINT NOT NULL,
    
    -- Event data (JSONB for flexibility and querying)
    event_data JSONB NOT NULL,
    
    -- Additional metadata
    metadata JSONB,
    
    -- Timestamps
    occurred_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Constraints
    UNIQUE(aggregate_id, sequence_number)
);

-- Indexes for performance
CREATE INDEX idx_agent_events_aggregate_id ON agent_events(aggregate_id);
CREATE INDEX idx_agent_events_aggregate_id_sequence ON agent_events(aggregate_id, sequence_number);
CREATE INDEX idx_agent_events_event_type ON agent_events(event_type);
CREATE INDEX idx_agent_events_occurred_at ON agent_events(occurred_at DESC);
CREATE INDEX idx_agent_events_created_at ON agent_events(created_at DESC);
CREATE INDEX idx_agent_events_data ON agent_events USING gin(event_data);

-- Comments for documentation
COMMENT ON TABLE agent_events IS 'Event store for Agent aggregate events';
COMMENT ON COLUMN agent_events.aggregate_id IS 'UUID of the agent aggregate';
COMMENT ON COLUMN agent_events.event_type IS 'Type of event (AgentInitialized, AgentTaskReceived, etc)';
COMMENT ON COLUMN agent_events.sequence_number IS 'Position in the aggregate event stream';
COMMENT ON COLUMN agent_events.event_data IS 'Event payload as JSON';
COMMENT ON COLUMN agent_events.metadata IS 'Additional context (user_id, correlation_id, etc)';