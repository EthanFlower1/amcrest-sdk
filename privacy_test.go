package amcrest

import (
	"context"
	"testing"
)

func TestPrivacy(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	t.Run("GetConfig", func(t *testing.T) {
		cfg, err := c.Privacy.GetConfig(ctx)
		if err != nil {
			t.Logf("Privacy.GetConfig not available: %v", err)
			return
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty config")
		}
		for k, v := range cfg {
			t.Logf("Privacy.%s = %s", k, v)
		}
	})

	t.Run("GetEnable", func(t *testing.T) {
		enabled, err := c.Privacy.GetEnable(ctx, 0)
		if err != nil {
			t.Logf("Privacy.GetEnable not available: %v", err)
			return
		}
		t.Logf("PrivacyMasking enabled: %v", enabled)
	})

	t.Run("GetMasking", func(t *testing.T) {
		coverCount := capsInt(videoInputCaps, "CoverCount")
		if coverCount <= 0 && videoInputCaps != nil {
			t.Skip("camera reports CoverCount=0, no privacy masking regions")
		}
		body, err := c.Privacy.GetMasking(ctx, 0, 0, 10)
		if err != nil {
			t.Logf("Privacy.GetMasking not available: %v", err)
			return
		}
		t.Logf("GetMasking response:\n%s", body)
	})
}
