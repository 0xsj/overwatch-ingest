# Overwatch Ingest Service

The Ingest service is responsible for **validating, scoring, and routing data**. It receives raw data from the Collector, validates against schemas, detects anomalies, scores confidence, and routes to Entity/Event services or quarantine.

## Responsibilities

| What it does | What it doesn't do |
|--------------|-------------------|
| Consume raw data from Collector | Fetch data from sources |
| Validate against schemas | Maintain source connections |
| Detect anomalies | Store entities/events |
| Score confidence | Perform correlation |
| Route: accept / quarantine / reject | Make intelligence assessments |
| Verify provenance chain | |
| Track source reliability | |

## Architecture
```
Collector Service                        Ingest Service
     │                                        │
     │ collector.raw.*                        │
     ▼                                        │
┌──────────────────────────────────────────────────────────────┐
│                                                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │   Schema    │  │  Anomaly    │  │ Confidence  │          │
│  │  Validator  │─▶│  Detector   │─▶│   Scorer    │          │
│  └─────────────┘  └─────────────┘  └─────────────┘          │
│                          │                                   │
│         ┌────────────────┼────────────────┐                 │
│         ▼                ▼                ▼                 │
│    [ACCEPTED]       [QUARANTINE]     [REJECTED]             │
│    conf > 0.7       0.3 < conf < 0.7  conf < 0.3            │
│         │                │                │                 │
│         ▼                ▼                ▼                 │
│   Entity/Event      Human Review     Dead Letter            │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

## Event Flow

**Subscribes to:**
- `overwatch.collector.raw.api`
- `overwatch.collector.raw.rss`
- `overwatch.collector.raw.webhook`
- `overwatch.collector.raw.stream`

**Publishes to:**
- `overwatch.ingest.record.accepted`
- `overwatch.ingest.record.rejected`
- `overwatch.ingest.record.quarantined`
- `overwatch.ingest.entity.observed` → Entity Service
- `overwatch.ingest.event.recorded` → Event Service

## Provenance

Ingest verifies the full chain and adds its own signature:
```json
{
  "id": "01KF0...",
  "raw_data_id": "01KF0...",
  "status": "ACCEPTED",
  "confidence": { "overall": 0.92 },
  "entity_type": "vessel",
  "entity_id": "mmsi:123456789",
  
  "source_signer": {
    "did": "did:key:z6MkPartner...",
    "verified": true
  },
  "collector_signer": {
    "did": "did:key:z6MkCollector...",
    "verified": true
  },
  "ingest_signer": {
    "did": "did:key:z6MkIngest...",
    "signer_type": "SERVICE",
    "signature": "base64..."
  }
}
```

## Validation Pipeline
```
Raw Data
    │
    ▼
┌─────────────────────────────────────────┐
│ 1. Signature Verification               │
│    - Verify collector signature         │
│    - Verify source signature (if any)   │
└────────────────┬────────────────────────┘
                 ▼
┌─────────────────────────────────────────┐
│ 2. Schema Validation                    │
│    - Required fields present            │
│    - Types correct                       │
│    - Format validation                   │
└────────────────┬────────────────────────┘
                 ▼
┌─────────────────────────────────────────┐
│ 3. Anomaly Detection                    │
│    - Range checks (lat: -90 to 90)      │
│    - Temporal checks (not future)       │
│    - Statistical outliers               │
└────────────────┬────────────────────────┘
                 ▼
┌─────────────────────────────────────────┐
│ 4. Confidence Scoring                   │
│    - Source reliability (historical)    │
│    - Data completeness                  │
│    - Temporal freshness                 │
│    - Signature trust                    │
└────────────────┬────────────────────────┘
                 ▼
┌─────────────────────────────────────────┐
│ 5. Routing Decision                     │
│    - Accept → Entity/Event              │
│    - Quarantine → Human review          │
│    - Reject → Dead letter               │
└─────────────────────────────────────────┘
```

## Key Domain Models

| Model | Purpose |
|-------|---------|
| `IngestRecord` | Result of validation with status and scores |
| `ValidationResult` | Schema validity, anomalies found |
| `ConfidenceScore` | Overall and component confidence scores |
| `QuarantinedRecord` | Data awaiting human review |
| `SourceReliability` | Historical accuracy tracking per source |

## Configuration
```env
INGEST_SERVER_PORT=50055
INGEST_DATABASE_NAME=overwatch_ingest
INGEST_NATS_URL=nats://localhost:4230
INGEST_NATS_SUBJECT_PREFIX=overwatch.ingest
INGEST_SERVICE_IDENTITY_ID=ingest-service
INGEST_SERVICE_IDENTITY_NAME=ingest
INGEST_SERVICE_IDENTITY_GENERATE_IF_MISSING=true
INGEST_QUARANTINE_EXPIRY_HOURS=72
INGEST_CONFIDENCE_THRESHOLD_ACCEPT=0.7
INGEST_CONFIDENCE_THRESHOLD_REJECT=0.3
```

## gRPC API
```protobuf
service IngestService {
  rpc Ping(PingRequest) returns (PingResponse);
  
  // Record access
  rpc GetRecord(GetRecordRequest) returns (GetRecordResponse);
  rpc ListRecords(ListRecordsRequest) returns (ListRecordsResponse);
  
  // Quarantine management
  rpc GetQuarantined(GetQuarantinedRequest) returns (GetQuarantinedResponse);
  rpc ListQuarantined(ListQuarantinedRequest) returns (ListQuarantinedResponse);
  rpc ResolveQuarantined(ResolveQuarantinedRequest) returns (ResolveQuarantinedResponse);
  
  // Reprocessing
  rpc ReprocessRecord(ReprocessRecordRequest) returns (ReprocessRecordResponse);
  
  // Validation preview
  rpc ValidateData(ValidateDataRequest) returns (ValidateDataResponse);
  
  // Source reliability
  rpc GetSourceReliability(GetSourceReliabilityRequest) returns (GetSourceReliabilityResponse);
}
```

## References

- **Protos**: `overwatch-contracts/proto/ingest/v1/`
- **Architecture pattern**: `overwatch-source/` (hexagonal, CQRS)
- **Provenance verification**: `overwatch-pkg/provenance/`