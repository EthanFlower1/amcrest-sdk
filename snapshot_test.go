package amcrest

import (
	"context"
	"testing"
)

func TestSnapshotGet(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	data, err := c.Snapshot.Get(ctx, 1)
	if err != nil {
		t.Fatalf("Snapshot.Get: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("expected non-empty snapshot data")
	}

	// Verify JPEG magic bytes (SOI marker: 0xFF 0xD8)
	if len(data) < 2 || data[0] != 0xFF || data[1] != 0xD8 {
		t.Fatalf("expected JPEG magic bytes (0xFF 0xD8), got (0x%02X 0x%02X)", data[0], data[1])
	}

	t.Logf("Snapshot size: %d bytes", len(data))
}

func TestGetSnapConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	cfg, err := c.Snapshot.GetSnapConfig(ctx)
	if err != nil {
		t.Fatalf("Snapshot.GetSnapConfig: %v", err)
	}

	if len(cfg) == 0 {
		t.Fatal("expected non-empty snap config map")
	}

	t.Logf("Snap config keys: %d", len(cfg))
	for k, v := range cfg {
		t.Logf("  %s = %s", k, v)
	}
}
