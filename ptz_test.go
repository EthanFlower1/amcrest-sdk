package amcrest

import (
	"context"
	"errors"
	"testing"
)

// skipIfNoPTZ skips the test if the camera does not support PTZ.
func skipIfNoPTZ(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		return
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) && apiErr.StatusCode == 400 {
		t.Skip("camera does not support PTZ, skipping")
	}
}

func TestPTZGetStatus(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	status, err := c.PTZ.GetStatus(ctx, 0)
	skipIfNoPTZ(t, err)
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	if len(status) == 0 {
		t.Fatal("expected non-empty status map")
	}
	t.Logf("PTZ status: %v", status)
}

func TestPTZGetConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	config, err := c.PTZ.GetConfig(ctx)
	skipIfNoPTZ(t, err)
	if err != nil {
		t.Fatalf("GetConfig failed: %v", err)
	}
	if len(config) == 0 {
		t.Fatal("expected non-empty config map")
	}
	t.Logf("PTZ config: %v", config)
}

func TestPTZGetPresets(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	presets, err := c.PTZ.GetPresets(ctx, 0)
	skipIfNoPTZ(t, err)
	if err != nil {
		t.Fatalf("GetPresets failed: %v", err)
	}
	// Some cameras may return an empty presets list, which is valid.
	t.Logf("PTZ presets: %s", presets)
}
