package amcrest

import (
	"context"
	"testing"
)

func TestPTZ(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	requireCapability(t, hasPTZ, "PTZ")

	t.Run("GetStatus", func(t *testing.T) {
		status, err := c.PTZ.GetStatus(ctx, 0)
		if err != nil {
			t.Fatalf("GetStatus: %v", err)
		}
		if len(status) == 0 {
			t.Fatal("expected non-empty status map")
		}
		for k, v := range status {
			t.Logf("Status.%s = %s", k, v)
		}
	})

	t.Run("GetConfig", func(t *testing.T) {
		config, err := c.PTZ.GetConfig(ctx)
		if err != nil {
			t.Fatalf("GetConfig: %v", err)
		}
		if len(config) == 0 {
			t.Fatal("expected non-empty config map")
		}
		for k, v := range config {
			t.Logf("Config.%s = %s", k, v)
		}
	})

	t.Run("GetPresets", func(t *testing.T) {
		presets, err := c.PTZ.GetPresets(ctx, 0)
		if err != nil {
			t.Fatalf("GetPresets: %v", err)
		}
		t.Logf("PTZ presets: %s", presets)
	})

	t.Run("GetCaps", func(t *testing.T) {
		caps, err := c.PTZ.GetCaps(ctx, 0)
		if err != nil {
			t.Fatalf("GetCaps: %v", err)
		}
		if caps == "" {
			t.Fatal("expected non-empty caps")
		}
		t.Logf("PTZ caps (first 500 chars): %.500s", caps)
	})

	t.Run("GetViewRangeStatus", func(t *testing.T) {
		status, err := c.PTZ.GetViewRangeStatus(ctx, 0)
		if err != nil {
			t.Logf("GetViewRangeStatus not available: %v", err)
			return
		}
		t.Logf("ViewRangeStatus: %s", status)
	})

	t.Run("GetEPTZConfig", func(t *testing.T) {
		cfg, err := c.PTZ.GetEPTZConfig(ctx)
		if err != nil {
			t.Logf("GetEPTZConfig not available: %v", err)
			return
		}
		for k, v := range cfg {
			t.Logf("EPTZ.%s = %s", k, v)
		}
	})

	t.Run("GetAutoMovementConfig", func(t *testing.T) {
		cfg, err := c.PTZ.GetAutoMovementConfig(ctx)
		if err != nil {
			t.Logf("GetAutoMovementConfig not available: %v", err)
			return
		}
		for k, v := range cfg {
			t.Logf("AutoMovement.%s = %s", k, v)
		}
	})
}
