package amcrest

import (
	"context"
	"testing"
)

func TestCamera(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	t.Run("GetImageConfig", func(t *testing.T) {
		cfg, err := c.Camera.GetImageConfig(ctx)
		if err != nil {
			t.Fatalf("GetImageConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty config")
		}
		for k, v := range cfg {
			t.Logf("Image.%s = %s", k, v)
		}
	})

	t.Run("GetExposureConfig", func(t *testing.T) {
		if capsInt(videoInputCaps, "caps.Exposure") == 0 && videoInputCaps != nil {
			t.Skip("camera does not support exposure control")
		}
		cfg, err := c.Camera.GetExposureConfig(ctx)
		if err != nil {
			t.Fatalf("GetExposureConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty config")
		}
		for k, v := range cfg {
			t.Logf("Exposure.%s = %s", k, v)
		}
	})

	t.Run("GetBacklightConfig", func(t *testing.T) {
		if capsInt(videoInputCaps, "caps.Backlight") == 0 && videoInputCaps != nil {
			t.Skip("camera does not support backlight control")
		}
		cfg, err := c.Camera.GetBacklightConfig(ctx)
		if err != nil {
			t.Fatalf("GetBacklightConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty config")
		}
		for k, v := range cfg {
			t.Logf("Backlight.%s = %s", k, v)
		}
	})

	t.Run("GetWhiteBalanceConfig", func(t *testing.T) {
		cfg, err := c.Camera.GetWhiteBalanceConfig(ctx)
		if err != nil {
			t.Fatalf("GetWhiteBalanceConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty config")
		}
		for k, v := range cfg {
			t.Logf("WhiteBalance.%s = %s", k, v)
		}
	})

	t.Run("GetDayNightConfig", func(t *testing.T) {
		if capsInt(videoInputCaps, "caps.DayNightColor") == 0 && videoInputCaps != nil {
			t.Skip("camera does not support DayNightColor")
		}
		cfg, err := c.Camera.GetDayNightConfig(ctx)
		if err != nil {
			t.Fatalf("GetDayNightConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty config")
		}
		for k, v := range cfg {
			t.Logf("DayNight.%s = %s", k, v)
		}
	})

	t.Run("GetFocusStatus", func(t *testing.T) {
		if !hasFocusControl {
			t.Skip("camera does not support electric focus (ElectricFocus cap)")
		}
		cfg, err := c.Camera.GetFocusStatus(ctx, 0)
		if err != nil {
			t.Fatalf("GetFocusStatus: %v", err)
		}
		for k, v := range cfg {
			t.Logf("Focus.%s = %s", k, v)
		}
	})

	t.Run("GetLightingConfig", func(t *testing.T) {
		if !hasLighting {
			t.Skip("camera does not support lighting (InfraRed cap)")
		}
		cfg, err := c.Camera.GetLightingConfig(ctx)
		if err != nil {
			t.Fatalf("GetLightingConfig: %v", err)
		}
		for k, v := range cfg {
			t.Logf("Lighting.%s = %s", k, v)
		}
	})

	t.Run("GetVideoInOptions", func(t *testing.T) {
		if !hasVideoInOptions {
			t.Skip("camera does not support VideoInOptions config")
		}
		cfg, err := c.Camera.GetVideoInOptions(ctx)
		if err != nil {
			t.Fatalf("GetVideoInOptions: %v", err)
		}
		for k, v := range cfg {
			t.Logf("VideoInOptions.%s = %s", k, v)
		}
	})
}
