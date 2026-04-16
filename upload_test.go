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
		cfg, err := c.Upload.GetPictureUploadConfig(ctx)
		if err != nil {
			t.Logf("GetPictureUploadConfig not available: %v", err)
			return
		}
		if len(cfg) == 0 {
			t.Log("PictureHttpUpload config returned empty map")
		}
		for k, v := range cfg {
			t.Logf("PictureUpload.%s = %s", k, v)
		}
	})

	t.Run("GetEventUploadConfig", func(t *testing.T) {
		cfg, err := c.Upload.GetEventUploadConfig(ctx)
		if err != nil {
			t.Logf("GetEventUploadConfig not available: %v", err)
			return
		}
		if len(cfg) == 0 {
			t.Log("EventHttpUpload config returned empty map")
		}
		for k, v := range cfg {
			t.Logf("EventUpload.%s = %s", k, v)
		}
	})

	t.Run("GetReportUploadConfig", func(t *testing.T) {
		cfg, err := c.Upload.GetReportUploadConfig(ctx)
		if err != nil {
			t.Logf("GetReportUploadConfig not available: %v", err)
			return
		}
		for k, v := range cfg {
			t.Logf("ReportUpload.%s = %s", k, v)
		}
	})
}
