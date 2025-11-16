// platform/pkg/events/nats/publisher_test.go
package nats_test

import (
	"context"
	"testing"
	"time"

	"github.com/0xsj/scout/platform/pkg/events/nats"
)

func TestNATSPublisher(t *testing.T) {
	// Skip if NATS is not running
	cfg := nats.NewConfig("nats://localhost:4224", 3, 2*time.Second)
	
	publisher, err := nats.NewPublisher(cfg)
	if err != nil {
		t.Skip("NATS not available, skipping test")
		return
	}
	defer publisher.Close()

	ctx := context.Background()
	
	// Test publish
	err = publisher.PublishJSON(ctx, "test.subject", map[string]string{
		"message": "hello from go",
	})
	if err != nil {
		t.Fatalf("Failed to publish: %v", err)
	}

	t.Log("✅ Successfully published to NATS")
}