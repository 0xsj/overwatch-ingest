-- name: CreateQuarantinedRecord :exec
INSERT INTO quarantined_records (
    id, tenant_id, source_id, source_type, raw_data_id, ingest_record_id,
    raw_data,
    reason, reason_detail, anomalies,
    confidence_overall, confidence_source_reliability, confidence_data_completeness,
    confidence_temporal_freshness, confidence_signature_trust, confidence_factors,
    resolution, resolved_by, resolved_by_did, resolution_notes, modified_data,
    ingest_signer, resolver_signer,
    quarantined_at, expires_at, resolved_at
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7,
    $8, $9, $10,
    $11, $12, $13, $14, $15, $16,
    $17, $18, $19, $20, $21,
    $22, $23,
    $24, $25, $26
);

-- name: UpdateQuarantinedRecord :exec
UPDATE quarantined_records SET
    resolution = $2,
    resolved_by = $3,
    resolved_by_did = $4,
    resolution_notes = $5,
    modified_data = $6,
    resolver_signer = $7,
    resolved_at = $8
WHERE id = $1;

-- name: FindQuarantinedRecordByID :one
SELECT * FROM quarantined_records WHERE id = $1;

-- name: FindQuarantinedRecordByIngestRecordID :one
SELECT * FROM quarantined_records WHERE ingest_record_id = $1;

-- name: ListQuarantinedRecords :many
SELECT * FROM quarantined_records
WHERE 
    (sqlc.narg('tenant_id')::varchar IS NULL OR tenant_id = sqlc.narg('tenant_id'))
    AND (sqlc.narg('source_id')::varchar IS NULL OR source_id = sqlc.narg('source_id'))
    AND (sqlc.narg('source_type')::varchar IS NULL OR source_type = sqlc.narg('source_type'))
    AND (sqlc.narg('reason')::quarantine_reason IS NULL OR reason = sqlc.narg('reason'))
    AND (sqlc.narg('resolution')::quarantine_resolution IS NULL OR resolution = sqlc.narg('resolution'))
    AND (sqlc.narg('start_time')::timestamptz IS NULL OR quarantined_at >= sqlc.narg('start_time'))
    AND (sqlc.narg('end_time')::timestamptz IS NULL OR quarantined_at <= sqlc.narg('end_time'))
ORDER BY quarantined_at DESC
LIMIT $1 OFFSET $2;

-- name: CountQuarantinedRecords :one
SELECT COUNT(*) FROM quarantined_records
WHERE 
    (sqlc.narg('tenant_id')::varchar IS NULL OR tenant_id = sqlc.narg('tenant_id'))
    AND (sqlc.narg('source_id')::varchar IS NULL OR source_id = sqlc.narg('source_id'))
    AND (sqlc.narg('source_type')::varchar IS NULL OR source_type = sqlc.narg('source_type'))
    AND (sqlc.narg('reason')::quarantine_reason IS NULL OR reason = sqlc.narg('reason'))
    AND (sqlc.narg('resolution')::quarantine_resolution IS NULL OR resolution = sqlc.narg('resolution'))
    AND (sqlc.narg('start_time')::timestamptz IS NULL OR quarantined_at >= sqlc.narg('start_time'))
    AND (sqlc.narg('end_time')::timestamptz IS NULL OR quarantined_at <= sqlc.narg('end_time'));

-- name: ListPendingQuarantinedRecords :many
SELECT * FROM quarantined_records
WHERE 
    resolution = 'pending'
    AND (sqlc.narg('tenant_id')::varchar IS NULL OR tenant_id = sqlc.narg('tenant_id'))
ORDER BY 
    CASE reason
        WHEN 'signature_invalid' THEN 1
        WHEN 'validation_failed' THEN 2
        WHEN 'anomaly_detected' THEN 3
        WHEN 'low_confidence' THEN 4
        WHEN 'duplicate_suspected' THEN 5
        WHEN 'manual_review' THEN 6
    END,
    quarantined_at ASC
LIMIT $1 OFFSET $2;

-- name: ListExpiredQuarantinedRecords :many
SELECT * FROM quarantined_records
WHERE 
    resolution = 'pending'
    AND expires_at IS NOT NULL
    AND expires_at <= $1
LIMIT $2;

-- name: BulkUpdateResolution :exec
UPDATE quarantined_records SET
    resolution = $2,
    resolved_by = $3,
    resolved_by_did = $4,
    resolution_notes = $5,
    resolved_at = $6
WHERE id = ANY($1::varchar[]);

-- name: DeleteQuarantinedRecord :exec
DELETE FROM quarantined_records WHERE id = $1;

-- name: GetQuarantineByReason :many
SELECT reason::text, COUNT(*) AS count
FROM quarantined_records
WHERE 
    (sqlc.narg('tenant_id')::varchar IS NULL OR tenant_id = sqlc.narg('tenant_id'))
    AND (sqlc.narg('start_time')::timestamptz IS NULL OR quarantined_at >= sqlc.narg('start_time'))
    AND (sqlc.narg('end_time')::timestamptz IS NULL OR quarantined_at <= sqlc.narg('end_time'))
GROUP BY reason;

-- name: CountPendingByTenant :one
SELECT COUNT(*) FROM quarantined_records
WHERE 
    resolution = 'pending'
    AND (sqlc.narg('tenant_id')::varchar IS NULL OR tenant_id = sqlc.narg('tenant_id'));