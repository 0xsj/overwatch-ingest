-- name: CreateIngestRecord :exec
INSERT INTO ingest_records (
    id, tenant_id, source_id, source_type, raw_data_id,
    status,
    validation_valid, validation_schema_valid, validation_fields_present, 
    validation_fields_missing, validation_anomalies, validation_validator_version,
    confidence_overall, confidence_source_reliability, confidence_data_completeness,
    confidence_temporal_freshness, confidence_signature_trust, confidence_factors,
    entity_type, entity_id, event_ids,
    rejection_reason, quarantine_id,
    source_signer, collector_signer, ingest_signer,
    source_signature_verified, collector_signature_verified,
    received_at, processed_at
) VALUES (
    $1, $2, $3, $4, $5,
    $6,
    $7, $8, $9, $10, $11, $12,
    $13, $14, $15, $16, $17, $18,
    $19, $20, $21,
    $22, $23,
    $24, $25, $26,
    $27, $28,
    $29, $30
);

-- name: UpdateIngestRecord :exec
UPDATE ingest_records SET
    status = $2,
    validation_valid = $3,
    validation_schema_valid = $4,
    validation_fields_present = $5,
    validation_fields_missing = $6,
    validation_anomalies = $7,
    validation_validator_version = $8,
    confidence_overall = $9,
    confidence_source_reliability = $10,
    confidence_data_completeness = $11,
    confidence_temporal_freshness = $12,
    confidence_signature_trust = $13,
    confidence_factors = $14,
    entity_type = $15,
    entity_id = $16,
    event_ids = $17,
    rejection_reason = $18,
    quarantine_id = $19,
    source_signer = $20,
    collector_signer = $21,
    ingest_signer = $22,
    source_signature_verified = $23,
    collector_signature_verified = $24,
    processed_at = $25
WHERE id = $1;

-- name: FindIngestRecordByID :one
SELECT * FROM ingest_records WHERE id = $1;

-- name: FindIngestRecordByRawDataID :one
SELECT * FROM ingest_records WHERE raw_data_id = $1;

-- name: ListIngestRecords :many
SELECT * FROM ingest_records
WHERE 
    (sqlc.narg('tenant_id')::varchar IS NULL OR tenant_id = sqlc.narg('tenant_id'))
    AND (sqlc.narg('source_id')::varchar IS NULL OR source_id = sqlc.narg('source_id'))
    AND (sqlc.narg('source_type')::varchar IS NULL OR source_type = sqlc.narg('source_type'))
    AND (sqlc.narg('status')::ingest_status IS NULL OR status = sqlc.narg('status'))
    AND (sqlc.narg('entity_type')::varchar IS NULL OR entity_type = sqlc.narg('entity_type'))
    AND (sqlc.narg('start_time')::timestamptz IS NULL OR received_at >= sqlc.narg('start_time'))
    AND (sqlc.narg('end_time')::timestamptz IS NULL OR received_at <= sqlc.narg('end_time'))
ORDER BY received_at DESC
LIMIT $1 OFFSET $2;

-- name: CountIngestRecords :one
SELECT COUNT(*) FROM ingest_records
WHERE 
    (sqlc.narg('tenant_id')::varchar IS NULL OR tenant_id = sqlc.narg('tenant_id'))
    AND (sqlc.narg('source_id')::varchar IS NULL OR source_id = sqlc.narg('source_id'))
    AND (sqlc.narg('source_type')::varchar IS NULL OR source_type = sqlc.narg('source_type'))
    AND (sqlc.narg('status')::ingest_status IS NULL OR status = sqlc.narg('status'))
    AND (sqlc.narg('entity_type')::varchar IS NULL OR entity_type = sqlc.narg('entity_type'))
    AND (sqlc.narg('start_time')::timestamptz IS NULL OR received_at >= sqlc.narg('start_time'))
    AND (sqlc.narg('end_time')::timestamptz IS NULL OR received_at <= sqlc.narg('end_time'));

-- name: ListIngestRecordsBySource :many
SELECT * FROM ingest_records
WHERE source_id = $1
    AND (sqlc.narg('status')::ingest_status IS NULL OR status = sqlc.narg('status'))
    AND (sqlc.narg('start_time')::timestamptz IS NULL OR received_at >= sqlc.narg('start_time'))
    AND (sqlc.narg('end_time')::timestamptz IS NULL OR received_at <= sqlc.narg('end_time'))
ORDER BY received_at DESC
LIMIT $2 OFFSET $3;

-- name: DeleteIngestRecord :exec
DELETE FROM ingest_records WHERE id = $1;

-- name: GetIngestStats :one
SELECT
    COUNT(*) AS total_records,
    COUNT(*) FILTER (WHERE status = 'accepted') AS accepted_records,
    COUNT(*) FILTER (WHERE status = 'rejected') AS rejected_records,
    COUNT(*) FILTER (WHERE status = 'quarantined') AS quarantined_records,
    COUNT(*) FILTER (WHERE status = 'pending') AS pending_records,
    COALESCE(AVG(confidence_overall), 0)::real AS average_confidence,
    COUNT(*) FILTER (WHERE source_signature_verified = true) AS source_signatures_verified,
    COUNT(*) FILTER (WHERE source_signature_verified = false) AS source_signatures_failed,
    COUNT(*) FILTER (WHERE collector_signature_verified = true) AS collector_signatures_verified,
    COUNT(*) FILTER (WHERE collector_signature_verified = false) AS collector_signatures_failed
FROM ingest_records
WHERE 
    (sqlc.narg('tenant_id')::varchar IS NULL OR tenant_id = sqlc.narg('tenant_id'))
    AND (sqlc.narg('source_id')::varchar IS NULL OR source_id = sqlc.narg('source_id'))
    AND (sqlc.narg('source_type')::varchar IS NULL OR source_type = sqlc.narg('source_type'))
    AND (sqlc.narg('start_time')::timestamptz IS NULL OR received_at >= sqlc.narg('start_time'))
    AND (sqlc.narg('end_time')::timestamptz IS NULL OR received_at <= sqlc.narg('end_time'));

-- name: GetRecordsBySource :many
SELECT source_id, COUNT(*) AS count
FROM ingest_records
WHERE 
    (sqlc.narg('tenant_id')::varchar IS NULL OR tenant_id = sqlc.narg('tenant_id'))
    AND (sqlc.narg('start_time')::timestamptz IS NULL OR received_at >= sqlc.narg('start_time'))
    AND (sqlc.narg('end_time')::timestamptz IS NULL OR received_at <= sqlc.narg('end_time'))
GROUP BY source_id;

-- name: GetRecordsBySourceType :many
SELECT source_type, COUNT(*) AS count
FROM ingest_records
WHERE 
    (sqlc.narg('tenant_id')::varchar IS NULL OR tenant_id = sqlc.narg('tenant_id'))
    AND (sqlc.narg('start_time')::timestamptz IS NULL OR received_at >= sqlc.narg('start_time'))
    AND (sqlc.narg('end_time')::timestamptz IS NULL OR received_at <= sqlc.narg('end_time'))
GROUP BY source_type;

-- name: GetRecordsByStatus :many
SELECT status::text, COUNT(*) AS count
FROM ingest_records
WHERE 
    (sqlc.narg('tenant_id')::varchar IS NULL OR tenant_id = sqlc.narg('tenant_id'))
    AND (sqlc.narg('start_time')::timestamptz IS NULL OR received_at >= sqlc.narg('start_time'))
    AND (sqlc.narg('end_time')::timestamptz IS NULL OR received_at <= sqlc.narg('end_time'))
GROUP BY status;

-- name: GetRecordsByEntityType :many
SELECT entity_type, COUNT(*) AS count
FROM ingest_records
WHERE 
    entity_type IS NOT NULL
    AND (sqlc.narg('tenant_id')::varchar IS NULL OR tenant_id = sqlc.narg('tenant_id'))
    AND (sqlc.narg('start_time')::timestamptz IS NULL OR received_at >= sqlc.narg('start_time'))
    AND (sqlc.narg('end_time')::timestamptz IS NULL OR received_at <= sqlc.narg('end_time'))
GROUP BY entity_type;