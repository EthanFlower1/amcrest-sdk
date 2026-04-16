package amcrest

import (
	"context"
	"testing"
)

func TestUpload(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	t.Run("GetPictureUploadConfig", func(t *testing.T) {
		if !hasUploadPicture {
			t.Skip("camera does not support PictureHttpUpload config")
		}
		cfg, err := c.Upload.GetPictureUploadConfig(ctx)
		if err != nil {
			t.Fatalf("GetPictureUploadConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Log("PictureHttpUpload config returned empty map")
		}
		for k, v := range cfg {
			t.Logf("PictureUpload.%s = %s", k, v)
		}
	})

	t.Run("SetPictureUploadConfig", func(t *testing.T) {
		if !hasUploadPicture {
			t.Skip("camera does not support PictureHttpUpload config")
		}
		original, err := c.Upload.GetPictureUploadConfig(ctx)
		if err != nil {
			t.Fatalf("GetPictureUploadConfig (save): %v", err)
		}

		// Find Enable key.
		enableKey := ""
		origVal := ""
		for k, v := range original {
			if contains(k, "Enable") {
				enableKey = k
				origVal = v
				break
			}
		}
		if enableKey == "" {
			t.Skip("no Enable key found in PictureUploadConfig")
		}
		t.Logf("Original %s = %s", enableKey, origVal)

		setKey := enableKey
		if len(setKey) > 6 && setKey[:6] == "table." {
			setKey = setKey[6:]
		}

		defer func() {
			_ = c.Upload.SetPictureUploadConfig(ctx, map[string]string{
				setKey: origVal,
			})
		}()

		newVal := "true"
		if origVal == "true" {
			newVal = "false"
		}
		err = c.Upload.SetPictureUploadConfig(ctx, map[string]string{
			setKey: newVal,
		})
		skipOnSetError(t, err, "SetPictureUploadConfig")

		updated, err := c.Upload.GetPictureUploadConfig(ctx)
		if err != nil {
			t.Fatalf("GetPictureUploadConfig (verify): %v", err)
		}
		if updated[enableKey] != newVal {
			t.Fatalf("expected %s=%q, got %q", enableKey, newVal, updated[enableKey])
		}
		t.Logf("Verified PictureUpload Enable changed to %q", newVal)
	})

	t.Run("GetEventUploadConfig", func(t *testing.T) {
		if !hasUploadEvent {
			t.Skip("camera does not support EventHttpUpload config")
		}
		cfg, err := c.Upload.GetEventUploadConfig(ctx)
		if err != nil {
			t.Fatalf("GetEventUploadConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Log("EventHttpUpload config returned empty map")
		}
		for k, v := range cfg {
			t.Logf("EventUpload.%s = %s", k, v)
		}
	})

	t.Run("GetReportUploadConfig", func(t *testing.T) {
		if !hasUploadReport {
			t.Skip("camera does not support ReportHttpUpload config")
		}
		cfg, err := c.Upload.GetReportUploadConfig(ctx)
		if err != nil {
			t.Fatalf("GetReportUploadConfig: %v", err)
		}
		for k, v := range cfg {
			t.Logf("ReportUpload.%s = %s", k, v)
		}
	})
}
