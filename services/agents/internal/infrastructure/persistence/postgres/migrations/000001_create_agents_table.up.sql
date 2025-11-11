-- services/agents/internal/infrastructure/persistence/postgres/migrations/000001_create_agents_table.up.sql
-- Agents table for storing agent state

CREATE TABLE IF NOT EXISTS agents (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    model VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'initialized',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for common queries
CREATE INDEX idx_agents_status ON agents(status);
CREATE INDEX idx_agents_provider ON agents(provider);
CREATE INDEX idx_agents_created_at ON agents(created_at DESC);

-- Comments for documentation
COMMENT ON TABLE agents IS 'AI/ML agents that process tasks';
COMMENT ON COLUMN agents.id IS 'Unique agent identifier';
COMMENT ON COLUMN agents.name IS 'Agent name';
COMMENT ON COLUMN agents.provider IS 'LLM provider (anthropic, openai, etc)';
COMMENT ON COLUMN agents.model IS 'Specific model identifier';
COMMENT ON COLUMN agents.status IS 'Agent operational status';
COMMENT ON COLUMN agents.created_at IS 'When the agent was created';
COMMENT ON COLUMN agents.updated_at IS 'Last update timestamp';