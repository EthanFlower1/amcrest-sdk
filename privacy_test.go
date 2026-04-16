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

	t.Run("SetEnable", func(t *testing.T) {
		if !hasPrivacyMaskCGI {
			t.Skip("camera does not support PrivacyMasking.cgi")
		}
		original, err := c.Privacy.GetEnable(ctx, 0)
		if err != nil {
			t.Fatalf("Privacy.GetEnable (save): %v", err)
		}
		t.Logf("Original PrivacyMasking enabled: %v", original)

		defer func() {
			_ = c.Privacy.SetEnable(ctx, 0, original)
		}()

		newVal := !original
		err = c.Privacy.SetEnable(ctx, 0, newVal)
		skipOnSetError(t, err, "Privacy.SetEnable")

		updated, err := c.Privacy.GetEnable(ctx, 0)
		if err != nil {
			t.Fatalf("Privacy.GetEnable (verify): %v", err)
		}
		if updated != newVal {
			t.Fatalf("expected enabled=%v, got %v", newVal, updated)
		}
		t.Logf("Verified PrivacyMasking enabled changed to %v", newVal)
	})

	t.Run("MaskingCRUD", func(t *testing.T) {
		if !hasPrivacyMaskCGI {
			t.Skip("camera does not support PrivacyMasking.cgi")
		}
		if !hasPrivacyMask {
			t.Skip("camera reports CoverCount=0 or no privacy masking support")
		}

		// Create a mask region. Use a small rectangle in the corner.
		maskParams := map[string]string{
			"PrivacyMasking.Index":  "0",
			"PrivacyMasking.Name":   "sdk-test-mask",
			"PrivacyMasking.Left":   "0",
			"PrivacyMasking.Top":    "0",
			"PrivacyMasking.Right":  "1000",
			"PrivacyMasking.Bottom": "1000",
		}
		err := c.Privacy.SetMasking(ctx, 0, maskParams)
		skipOnSetError(t, err, "SetMasking")
		t.Log("Created privacy mask at index 0")

		// Ensure cleanup.
		defer func() {
			_ = c.Privacy.DeleteMasking(ctx, 0, 0)
		}()

		// Read it back.
		body, err := c.Privacy.GetMasking(ctx, 0, 0, 10)
		if err != nil {
			t.Fatalf("GetMasking (read): %v", err)
		}
		t.Logf("GetMasking response:\n%s", body)

		// Delete it.
		err = c.Privacy.DeleteMasking(ctx, 0, 0)
		skipOnSetError(t, err, "DeleteMasking")
		t.Log("Deleted privacy mask at index 0")
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
