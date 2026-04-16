package amcrest

import (
	"context"
	"errors"
	"testing"
)

func TestCameraGetImageConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Camera.GetImageConfig(ctx)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			t.Skipf("GetImageConfig not supported (HTTP %d)", apiErr.StatusCode)
		}
		t.Fatalf("GetImageConfig: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config map")
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}

func TestCameraGetExposureConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Camera.GetExposureConfig(ctx)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			t.Skipf("GetExposureConfig not supported (HTTP %d)", apiErr.StatusCode)
		}
		t.Fatalf("GetExposureConfig: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config map")
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}

func TestCameraGetBacklightConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Camera.GetBacklightConfig(ctx)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			t.Skipf("GetBacklightConfig not supported (HTTP %d)", apiErr.StatusCode)
		}
		t.Fatalf("GetBacklightConfig: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config map")
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}

func TestCameraGetWhiteBalanceConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Camera.GetWhiteBalanceConfig(ctx)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			t.Skipf("GetWhiteBalanceConfig not supported (HTTP %d)", apiErr.StatusCode)
		}
		t.Fatalf("GetWhiteBalanceConfig: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config map")
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}

func TestCameraGetDayNightConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Camera.GetDayNightConfig(ctx)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			t.Skipf("GetDayNightConfig not supported (HTTP %d)", apiErr.StatusCode)
		}
		t.Fatalf("GetDayNightConfig: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config map")
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}

func TestCameraGetFocusStatus(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Camera.GetFocusStatus(ctx, 0)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			t.Skipf("GetFocusStatus not supported (HTTP %d)", apiErr.StatusCode)
		}
		t.Fatalf("GetFocusStatus: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil status map")
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}

func TestCameraGetLightingConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Camera.GetLightingConfig(ctx)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			t.Skipf("GetLightingConfig not supported (HTTP %d)", apiErr.StatusCode)
		}
		t.Fatalf("GetLightingConfig: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config map")
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}

func TestCameraGetVideoInOptions(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Camera.GetVideoInOptions(ctx)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			t.Skipf("GetVideoInOptions not supported (HTTP %d)", apiErr.StatusCode)
		}
		t.Fatalf("GetVideoInOptions: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config map")
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}
