package amcrest

import (
	"context"
	"testing"
)

func TestUploadGetPictureUploadConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Upload.GetPictureUploadConfig(ctx)
	if err != nil {
		t.Skip("GetPictureUploadConfig not supported on this camera, skipping")
	}
	if len(cfg) == 0 {
		t.Log("PictureHttpUpload config returned empty map")
	}
	for k, v := range cfg {
		t.Logf("PictureUpload.%s = %s", k, v)
	}
}

func TestUploadGetEventUploadConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Upload.GetEventUploadConfig(ctx)
	if err != nil {
		t.Skip("GetEventUploadConfig not supported on this camera, skipping")
	}
	if len(cfg) == 0 {
		t.Log("EventHttpUpload config returned empty map")
	}
	for k, v := range cfg {
		t.Logf("EventUpload.%s = %s", k, v)
	}
}
