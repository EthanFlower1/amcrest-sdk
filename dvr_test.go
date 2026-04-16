package amcrest

import (
	"context"
	"testing"
)

func TestDVR(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	t.Run("GetBandwidthLimit", func(t *testing.T) {
		if !hasBandwidthLimit {
			t.Skip("camera does not support BandwidthLimit config")
		}
		cfg, err := c.DVR.GetBandwidthLimit(ctx)
		if err != nil {
			t.Fatalf("GetBandwidthLimit: %v", err)
		}
		for k, v := range cfg {
			t.Logf("Bandwidth.%s = %s", k, v)
		}
	})

	t.Run("StartFindAndStop", func(t *testing.T) {
		if !hasDVR {
			t.Skip("camera does not support DVR media find")
		}
		findId, err := c.DVR.StartFind(ctx, 0, "2024-01-01 00:00:00", "2024-01-02 00:00:00")
		if err != nil {
			t.Fatalf("DVR StartFind: %v", err)
		}
		t.Logf("findId: %s", findId)

		body, err := c.DVR.FindNext(ctx, findId, 10)
		if err != nil {
			t.Fatalf("FindNext: %v", err)
		}
		t.Logf("FindNext response:\n%s", body)

		if err := c.DVR.StopFind(ctx, findId); err != nil {
			t.Fatalf("StopFind: %v", err)
		}
		t.Log("StopFind succeeded")
	})
}
