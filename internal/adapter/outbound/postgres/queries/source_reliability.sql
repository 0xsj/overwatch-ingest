-- name: UpsertSourceReliability :exec
INSERT INTO source_reliability (
    source_id, tenant_id,
    reliability_score,
    total_records, accepted_records, rejected_records, quarantined_records,
    corroborated_records, disputed_records,
    calculated_at, window_start, window_end,
    created_at, updated_at
) VALUES (
    $1, $2,
    $3,
    $4, $5, $6, $7,
    $8, $9,
    $10, $11, $12,
    NOW(), NOW()
)
ON CONFLICT (source_id) DO UPDATE SET
    tenant_id = EXCLUDED.tenant_id,
    reliability_score = EXCLUDED.reliability_score,
    total_records = EXCLUDED.total_records,
    accepted_records = EXCLUDED.accepted_records,
    rejected_records = EXCLUDED.rejected_records,
    quarantined_records = EXCLUDED.quarantined_records,
    corroborated_records = EXCLUDED.corroborated_records,
    disputed_records = EXCLUDED.disputed_records,
    calculated_at = EXCLUDED.calculated_at,
    window_start = EXCLUDED.window_start,
    window_end = EXCLUDED.window_end,
    updated_at = NOW();

-- name: UpdateSourceReliability :exec
UPDATE source_reliability SET
    reliability_score = $2,
    total_records = $3,
    accepted_records = $4,
    rejected_records = $5,
    quarantined_records = $6,
    corroborated_records = $7,
    disputed_records = $8,
    calculated_at = $9,
    window_start = $10,
    window_end = $11,
    updated_at = NOW()
WHERE source_id = $1;

-- name: IncrementAccepted :exec
UPDATE source_reliability SET
    total_records = total_records + 1,
    accepted_records = accepted_records + 1,
    calculated_at = NOW(),
    window_end = NOW(),
    updated_at = NOW()
WHERE source_id = $1;

-- name: IncrementRejected :exec
UPDATE source_reliability SET
    total_records = total_records + 1,
    rejected_records = rejected_records + 1,
    calculated_at = NOW(),
    window_end = NOW(),
    updated_at = NOW()
WHERE source_id = $1;

-- name: IncrementQuarantined :exec
UPDATE source_reliability SET
    total_records = total_records + 1,
    quarantined_records = quarantined_records + 1,
    calculated_at = NOW(),
    window_end = NOW(),
    updated_at = NOW()
WHERE source_id = $1;

-- name: IncrementCorroborated :exec
UPDATE source_reliability SET
    corroborated_records = corroborated_records + 1,
    calculated_at = NOW(),
    updated_at = NOW()
WHERE source_id = $1;

-- name: IncrementDisputed :exec
UPDATE source_reliability SET
    disputed_records = disputed_records + 1,
    calculated_at = NOW(),
    updated_at = NOW()
WHERE source_id = $1;

-- name: FindSourceReliabilityByID :one
SELECT * FROM source_reliability WHERE source_id = $1;

-- name: ListSourceReliability :many
SELECT * FROM source_reliability
WHERE 
    (sqlc.narg('tenant_id')::varchar IS NULL OR tenant_id = sqlc.narg('tenant_id'))
    AND (sqlc.narg('min_score')::real IS NULL OR reliability_score >= sqlc.narg('min_score'))
    AND (sqlc.narg('max_score')::real IS NULL OR reliability_score <= sqlc.narg('max_score'))
ORDER BY reliability_score DESC
LIMIT $1 OFFSET $2;

-- name: CountSourceReliability :one
SELECT COUNT(*) FROM source_reliability
WHERE 
    (sqlc.narg('tenant_id')::varchar IS NULL OR tenant_id = sqlc.narg('tenant_id'))
    AND (sqlc.narg('min_score')::real IS NULL OR reliability_score >= sqlc.narg('min_score'))
    AND (sqlc.narg('max_score')::real IS NULL OR reliability_score <= sqlc.narg('max_score'));

-- name: ListUnreliableSources :many
SELECT * FROM source_reliability
WHERE 
    reliability_score < $1
    AND (sqlc.narg('tenant_id')::varchar IS NULL OR tenant_id = sqlc.narg('tenant_id'))
ORDER BY reliability_score ASC
LIMIT $2 OFFSET $3;

-- name: ResetSourceReliability :exec
UPDATE source_reliability SET
    reliability_score = 0.5,
    total_records = 0,
    accepted_records = 0,
    rejected_records = 0,
    quarantined_records = 0,
    corroborated_records = 0,
    disputed_records = 0,
    window_start = $2,
    window_end = $2,
    calculated_at = NOW(),
    updated_at = NOW()
WHERE source_id = $1;

-- name: DeleteSourceReliability :exec
DELETE FROM source_reliability WHERE source_id = $1;

-- name: ExistsSourceReliability :one
SELECT EXISTS(SELECT 1 FROM source_reliability WHERE source_id = $1);