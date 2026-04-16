package amcrest

import (
	"context"
	"testing"
)

func TestSnapshot(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	t.Run("Get", func(t *testing.T) {
		data, err := c.Snapshot.Get(ctx, 1)
		if err != nil {
			t.Fatalf("Snapshot.Get: %v", err)
		}
		if len(data) == 0 {
			t.Fatal("expected non-empty snapshot data")
		}
		if len(data) < 2 || data[0] != 0xFF || data[1] != 0xD8 {
			t.Fatalf("expected JPEG magic bytes (0xFF 0xD8), got (0x%02X 0x%02X)", data[0], data[1])
		}
		t.Logf("Snapshot size: %d bytes", len(data))
	})

	t.Run("GetSnapConfig", func(t *testing.T) {
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
	})

	t.Run("GetWithType", func(t *testing.T) {
		if !hasSnapWithType {
			t.Skip("camera does not support snapshot GetWithType")
		}
		data, err := c.Snapshot.GetWithType(ctx, 1, 0)
		if err != nil {
			t.Fatalf("GetWithType: %v", err)
		}
		if len(data) == 0 {
			t.Fatal("expected non-empty snapshot data")
		}
		if len(data) >= 2 && data[0] == 0xFF && data[1] == 0xD8 {
			t.Logf("GetWithType returned JPEG: %d bytes", len(data))
		} else {
			t.Logf("GetWithType returned %d bytes", len(data))
		}
	})
}
