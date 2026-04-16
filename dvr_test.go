package amcrest

import (
	"context"
	"errors"
	"testing"
)

func skipIfNoDVR(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		return
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		t.Skipf("camera does not support DVR feature (HTTP %d), skipping", apiErr.StatusCode)
	}
}

func TestDVRStartFindAndStop(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	findId, err := c.DVR.StartFind(ctx, 0, "2024-01-01 00:00:00", "2024-01-02 00:00:00")
	skipIfNoDVR(t, err)
	if err != nil {
		t.Fatalf("StartFind: %v", err)
	}
	t.Logf("findId: %s", findId)

	body, err := c.DVR.FindNext(ctx, findId, 10)
	if err != nil {
		t.Fatalf("FindNext: %v", err)
	}
	t.Logf("FindNext response:\n%s", body)

	err = c.DVR.StopFind(ctx, findId)
	if err != nil {
		t.Fatalf("StopFind: %v", err)
	}
}

func TestDVRGetBandwidthLimit(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	cfg, err := c.DVR.GetBandwidthLimit(ctx)
	skipIfNoDVR(t, err)
	if err != nil {
		t.Fatalf("GetBandwidthLimit: %v", err)
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}
