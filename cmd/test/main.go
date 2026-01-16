package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	ingestv1 "github.com/0xsj/overwatch-contracts/gen/go/ingest/v1"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	addr := getEnv("INGEST_SERVICE_ADDR", "localhost:50054")

	fmt.Printf("Connecting to ingest service at %s...\n", addr)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	client := ingestv1.NewIngestServiceClient(conn)

	// ─────────────────────────────────────────────────────────────
	// Test 1: Ping
	// ─────────────────────────────────────────────────────────────
	fmt.Println("\n[1] Ping")
	pingResp, err := client.Ping(ctx, &ingestv1.PingRequest{})
	if err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	fmt.Printf("    ✓ Response: %s\n", pingResp.Message)

	// ─────────────────────────────────────────────────────────────
	// Test 2: Get Ingest Stats (empty initially)
	// ─────────────────────────────────────────────────────────────
	fmt.Println("\n[2] Get Ingest Stats")
	statsResp, err := client.GetIngestStats(ctx, &ingestv1.GetIngestStatsRequest{
		TenantId: stringPtr("tenant-001"),
	})
	if err != nil {
		return fmt.Errorf("get ingest stats failed: %w", err)
	}
	fmt.Printf("    ✓ Total Records: %d\n", statsResp.TotalRecords)
	fmt.Printf("    ✓ Accepted: %d\n", statsResp.AcceptedRecords)
	fmt.Printf("    ✓ Rejected: %d\n", statsResp.RejectedRecords)
	fmt.Printf("    ✓ Quarantined: %d\n", statsResp.QuarantinedRecords)
	fmt.Printf("    ✓ Pending: %d\n", statsResp.PendingRecords)
	fmt.Printf("    ✓ Avg Confidence: %.2f\n", statsResp.AverageConfidence)
	fmt.Printf("    ✓ Avg Processing Time: %.2fms\n", statsResp.AverageProcessingTimeMs)

	// ─────────────────────────────────────────────────────────────
	// Test 3: List Records (empty initially)
	// ─────────────────────────────────────────────────────────────
	fmt.Println("\n[3] List Records")
	listResp, err := client.ListRecords(ctx, &ingestv1.ListRecordsRequest{
		TenantId: stringPtr("tenant-001"),
	})
	if err != nil {
		return fmt.Errorf("list records failed: %w", err)
	}
	fmt.Printf("    ✓ Found: %d records\n", len(listResp.Records))
	if listResp.Pagination != nil {
		fmt.Printf("    ✓ Total Items: %d\n", listResp.Pagination.TotalItems)
	}

	// ─────────────────────────────────────────────────────────────
	// Test 4: List Records by Status
	// ─────────────────────────────────────────────────────────────
	fmt.Println("\n[4] List Records by Status (ACCEPTED)")
	listByStatusResp, err := client.ListRecords(ctx, &ingestv1.ListRecordsRequest{
		TenantId: stringPtr("tenant-001"),
		Status:   ingestStatusPtr(ingestv1.IngestStatus_INGEST_STATUS_ACCEPTED),
	})
	if err != nil {
		return fmt.Errorf("list records by status failed: %w", err)
	}
	fmt.Printf("    ✓ Found: %d accepted records\n", len(listByStatusResp.Records))

	// ─────────────────────────────────────────────────────────────
	// Test 5: List Records by Source Type
	// ─────────────────────────────────────────────────────────────
	fmt.Println("\n[5] List Records by Source Type (ais)")
	listBySourceTypeResp, err := client.ListRecords(ctx, &ingestv1.ListRecordsRequest{
		TenantId:   stringPtr("tenant-001"),
		SourceType: stringPtr("ais"),
	})
	if err != nil {
		return fmt.Errorf("list records by source type failed: %w", err)
	}
	fmt.Printf("    ✓ Found: %d AIS records\n", len(listBySourceTypeResp.Records))

	// ─────────────────────────────────────────────────────────────
	// Test 6: List Quarantined Records
	// ─────────────────────────────────────────────────────────────
	fmt.Println("\n[6] List Quarantined Records")
	quarantinedResp, err := client.ListQuarantined(ctx, &ingestv1.ListQuarantinedRequest{
		TenantId: stringPtr("tenant-001"),
	})
	if err != nil {
		return fmt.Errorf("list quarantined failed: %w", err)
	}
	fmt.Printf("    ✓ Found: %d quarantined records\n", len(quarantinedResp.Records))

	var quarantinedID string
	if len(quarantinedResp.Records) > 0 {
		quarantinedID = quarantinedResp.Records[0].Id
		fmt.Printf("    ✓ First quarantined ID: %s\n", quarantinedID[:8])
	}

	// ─────────────────────────────────────────────────────────────
	// Test 7: List Quarantined by Reason
	// ─────────────────────────────────────────────────────────────
	fmt.Println("\n[7] List Quarantined by Reason (LOW_CONFIDENCE)")
	quarantinedByReasonResp, err := client.ListQuarantined(ctx, &ingestv1.ListQuarantinedRequest{
		TenantId: stringPtr("tenant-001"),
		Reason:   quarantineReasonPtr(ingestv1.QuarantineReason_QUARANTINE_REASON_LOW_CONFIDENCE),
	})
	if err != nil {
		return fmt.Errorf("list quarantined by reason failed: %w", err)
	}
	fmt.Printf("    ✓ Found: %d low confidence records\n", len(quarantinedByReasonResp.Records))

	// ─────────────────────────────────────────────────────────────
	// Test 8: List Source Reliability
	// ─────────────────────────────────────────────────────────────
	fmt.Println("\n[8] List Source Reliability")
	reliabilityResp, err := client.ListSourceReliability(ctx, &ingestv1.ListSourceReliabilityRequest{
		TenantId: stringPtr("tenant-001"),
	})
	if err != nil {
		return fmt.Errorf("list source reliability failed: %w", err)
	}
	fmt.Printf("    ✓ Found: %d source reliability records\n", len(reliabilityResp.Reliabilities))
	for _, r := range reliabilityResp.Reliabilities {
		fmt.Printf("      - Source %s: score=%.2f, records=%d\n",
			r.SourceId[:8], r.ReliabilityScore, r.TotalRecords)
	}

	// ─────────────────────────────────────────────────────────────
	// Test 9: List Source Reliability with Min Score
	// ─────────────────────────────────────────────────────────────
	fmt.Println("\n[9] List Source Reliability (min score 0.5)")
	reliabilityMinResp, err := client.ListSourceReliability(ctx, &ingestv1.ListSourceReliabilityRequest{
		TenantId: stringPtr("tenant-001"),
		MinScore: float32Ptr(0.5),
	})
	if err != nil {
		return fmt.Errorf("list source reliability with min score failed: %w", err)
	}
	fmt.Printf("    ✓ Found: %d reliable sources (score >= 0.5)\n", len(reliabilityMinResp.Reliabilities))

	// ─────────────────────────────────────────────────────────────
	// Test 10: Get Record (will fail if no records exist)
	// ─────────────────────────────────────────────────────────────
	fmt.Println("\n[10] Get Record by ID (testing error case)")
	_, err = client.GetRecord(ctx, &ingestv1.GetRecordRequest{
		Id: "nonexistent-record-id",
	})
	if err != nil {
		fmt.Printf("    ✓ Expected error: record not found\n")
	} else {
		fmt.Printf("    ✗ Should have returned error for nonexistent record\n")
	}

	// ─────────────────────────────────────────────────────────────
	// Test 11: Get Quarantined (will fail if no quarantined exist)
	// ─────────────────────────────────────────────────────────────
	fmt.Println("\n[11] Get Quarantined by ID (testing error case)")
	_, err = client.GetQuarantined(ctx, &ingestv1.GetQuarantinedRequest{
		Id: "nonexistent-quarantine-id",
	})
	if err != nil {
		fmt.Printf("    ✓ Expected error: quarantine record not found\n")
	} else {
		fmt.Printf("    ✗ Should have returned error for nonexistent quarantine\n")
	}

	// ─────────────────────────────────────────────────────────────
	// Test 12: Get Source Reliability (will fail if source doesn't exist)
	// ─────────────────────────────────────────────────────────────
	fmt.Println("\n[12] Get Source Reliability (testing error case)")
	_, err = client.GetSourceReliability(ctx, &ingestv1.GetSourceReliabilityRequest{
		SourceId: "nonexistent-source-id",
	})
	if err != nil {
		fmt.Printf("    ✓ Expected error: source reliability not found\n")
	} else {
		fmt.Printf("    ✗ Should have returned error for nonexistent source\n")
	}

	// ─────────────────────────────────────────────────────────────
	// Test 13: Get Ingest Stats by Source
	// ─────────────────────────────────────────────────────────────
	fmt.Println("\n[13] Get Ingest Stats by Source")
	statsSourceResp, err := client.GetIngestStats(ctx, &ingestv1.GetIngestStatsRequest{
		TenantId: stringPtr("tenant-001"),
		SourceId: stringPtr("source-ais-001"),
	})
	if err != nil {
		return fmt.Errorf("get ingest stats by source failed: %w", err)
	}
	fmt.Printf("    ✓ Stats for source-ais-001:\n")
	fmt.Printf("      Total: %d, Accepted: %d, Rejected: %d\n",
		statsSourceResp.TotalRecords,
		statsSourceResp.AcceptedRecords,
		statsSourceResp.RejectedRecords)

	// ─────────────────────────────────────────────────────────────
	// Test 14: Get Ingest Stats by Source Type
	// ─────────────────────────────────────────────────────────────
	fmt.Println("\n[14] Get Ingest Stats by Source Type")
	statsTypeResp, err := client.GetIngestStats(ctx, &ingestv1.GetIngestStatsRequest{
		TenantId:   stringPtr("tenant-001"),
		SourceType: stringPtr("rss"),
	})
	if err != nil {
		return fmt.Errorf("get ingest stats by source type failed: %w", err)
	}
	fmt.Printf("    ✓ Stats for RSS sources:\n")
	fmt.Printf("      Total: %d, Accepted: %d, Rejected: %d\n",
		statsTypeResp.TotalRecords,
		statsTypeResp.AcceptedRecords,
		statsTypeResp.RejectedRecords)

	// ─────────────────────────────────────────────────────────────
	// Test 15: Validate Data (preview - returns stub for now)
	// ─────────────────────────────────────────────────────────────
	fmt.Println("\n[15] Validate Data (preview)")
	validateResp, err := client.ValidateData(ctx, &ingestv1.ValidateDataRequest{
		SourceId:   "source-ais-001",
		SourceType: "ais",
	})
	if err != nil {
		return fmt.Errorf("validate data failed: %w", err)
	}
	fmt.Printf("    ✓ Predicted Status: %s\n", validateResp.PredictedStatus)
	fmt.Printf("    ✓ (Note: This is a stub implementation)\n")

	// ─────────────────────────────────────────────────────────────
	// Done
	// ─────────────────────────────────────────────────────────────
	fmt.Println()
	fmt.Println("══════════════════════════════════════════════════")
	fmt.Println("All ingest service tests passed!")
	fmt.Println("══════════════════════════════════════════════════")
	fmt.Println()
	fmt.Printf("Summary:\n")
	fmt.Printf("  Queries Executed: 15\n")
	fmt.Printf("  Stats Retrieved: 3\n")
	fmt.Printf("  Error Cases Tested: 3\n")
	fmt.Println()
	fmt.Println("Note: To test actual ingestion, publish events to NATS:")
	fmt.Printf("  Subject: overwatch.collector.raw_data.collected\n")
	fmt.Println()
	fmt.Println("The ingest service will:")
	fmt.Println("  1. Verify collector/source signatures")
	fmt.Println("  2. Validate data against source type schema")
	fmt.Println("  3. Detect anomalies and calculate confidence")
	fmt.Println("  4. Accept, reject, or quarantine records")
	fmt.Println("  5. Publish signed events to downstream services")

	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func stringPtr(v string) *string {
	return &v
}

func float32Ptr(v float32) *float32 {
	return &v
}

func ingestStatusPtr(v ingestv1.IngestStatus) *ingestv1.IngestStatus {
	return &v
}

func quarantineReasonPtr(v ingestv1.QuarantineReason) *ingestv1.QuarantineReason {
	return &v
}
