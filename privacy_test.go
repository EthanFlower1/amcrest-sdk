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
		if !hasPrivacyConfig {
			t.Skip("camera does not support PrivacyMasking config")
		}
		cfg, err := c.Privacy.GetConfig(ctx)
		if err != nil {
			t.Fatalf("Privacy.GetConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty config")
		}
		for k, v := range cfg {
			t.Logf("Privacy.%s = %s", k, v)
		}
	})

	t.Run("GetEnable", func(t *testing.T) {
		if !hasPrivacyMaskCGI {
			t.Skip("camera does not support PrivacyMasking.cgi")
		}
		enabled, err := c.Privacy.GetEnable(ctx, 0)
		if err != nil {
			t.Fatalf("Privacy.GetEnable: %v", err)
		}
		t.Logf("PrivacyMasking enabled: %v", enabled)
	})

	t.Run("GetMasking", func(t *testing.T) {
		if !hasPrivacyMaskCGI {
			t.Skip("camera does not support PrivacyMasking.cgi")
		}
		if !hasPrivacyMask {
			t.Skip("camera reports CoverCount=0 or no privacy masking support")
		}
		body, err := c.Privacy.GetMasking(ctx, 0, 0, 10)
		if err != nil {
			t.Fatalf("Privacy.GetMasking: %v", err)
		}
		t.Logf("GetMasking response:\n%s", body)
	})
}
