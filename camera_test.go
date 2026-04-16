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

	t.Run("SetImageConfig", func(t *testing.T) {
		original, err := c.Camera.GetImageConfig(ctx)
		if err != nil {
			t.Fatalf("GetImageConfig (save): %v", err)
		}

		// Find Brightness key (raw config has keys like "table.VideoColor[0][0].Brightness").
		brightnessKey := ""
		origBrightness := ""
		for k, v := range original {
			if contains(k, "Brightness") {
				brightnessKey = k
				origBrightness = v
				break
			}
		}
		if brightnessKey == "" {
			t.Skip("no Brightness key found in ImageConfig")
		}
		t.Logf("Original %s = %s", brightnessKey, origBrightness)

		// Strip "table." prefix for set key.
		setKey := brightnessKey
		if len(setKey) > 6 && setKey[:6] == "table." {
			setKey = setKey[6:]
		}

		defer func() {
			_ = c.Camera.SetImageConfig(ctx, map[string]string{
				setKey: origBrightness,
			})
		}()

		// Nudge brightness by 1 (keep in safe range 1-99).
		newBrightness := "51"
		if origBrightness == "51" {
			newBrightness = "50"
		}
		err = c.Camera.SetImageConfig(ctx, map[string]string{
			setKey: newBrightness,
		})
		skipOnSetError(t, err, "SetImageConfig")

		updated, err := c.Camera.GetImageConfig(ctx)
		if err != nil {
			t.Fatalf("GetImageConfig (verify): %v", err)
		}
		if updated[brightnessKey] != newBrightness {
			t.Fatalf("expected %s=%q, got %q", brightnessKey, newBrightness, updated[brightnessKey])
		}
		t.Logf("Verified Brightness changed to %q", newBrightness)
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

	t.Run("SetExposureConfig", func(t *testing.T) {
		if capsInt(videoInputCaps, "caps.Exposure") == 0 && videoInputCaps != nil {
			t.Skip("camera does not support exposure control")
		}
		original, err := c.Camera.GetExposureConfig(ctx)
		if err != nil {
			t.Fatalf("GetExposureConfig (save): %v", err)
		}

		// Find a safe value to toggle - AntiFlicker is typically "50" or "60".
		antiFlickerKey := ""
		origVal := ""
		for k, v := range original {
			if contains(k, "AntiFlicker") {
				antiFlickerKey = k
				origVal = v
				break
			}
		}
		if antiFlickerKey == "" {
			t.Skip("no AntiFlicker key found in ExposureConfig")
		}
		t.Logf("Original %s = %s", antiFlickerKey, origVal)

		setKey := antiFlickerKey
		if len(setKey) > 6 && setKey[:6] == "table." {
			setKey = setKey[6:]
		}

		defer func() {
			_ = c.Camera.SetExposureConfig(ctx, map[string]string{
				setKey: origVal,
			})
		}()

		newVal := "60"
		if origVal == "60" {
			newVal = "50"
		}
		err = c.Camera.SetExposureConfig(ctx, map[string]string{
			setKey: newVal,
		})
		skipOnSetError(t, err, "SetExposureConfig")

		updated, err := c.Camera.GetExposureConfig(ctx)
		if err != nil {
			t.Fatalf("GetExposureConfig (verify): %v", err)
		}
		if updated[antiFlickerKey] != newVal {
			t.Fatalf("expected %s=%q, got %q", antiFlickerKey, newVal, updated[antiFlickerKey])
		}
		t.Logf("Verified AntiFlicker changed to %q", newVal)
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

	t.Run("SetDayNightConfig", func(t *testing.T) {
		if capsInt(videoInputCaps, "caps.DayNightColor") == 0 && videoInputCaps != nil {
			t.Skip("camera does not support DayNightColor")
		}
		original, err := c.Camera.GetDayNightConfig(ctx)
		if err != nil {
			t.Fatalf("GetDayNightConfig (save): %v", err)
		}

		// Find a Mode key to toggle.
		modeKey := ""
		origMode := ""
		for k, v := range original {
			if contains(k, "Mode") && !contains(k, "SwitchMode") {
				modeKey = k
				origMode = v
				break
			}
		}
		if modeKey == "" {
			t.Skip("no Mode key found in DayNightConfig")
		}
		t.Logf("Original %s = %s", modeKey, origMode)

		setKey := modeKey
		if len(setKey) > 6 && setKey[:6] == "table." {
			setKey = setKey[6:]
		}

		defer func() {
			_ = c.Camera.SetDayNightConfig(ctx, map[string]string{
				setKey: origMode,
			})
		}()

		// Toggle between known modes. Avoid setting invalid values.
		// Common modes: 0=Color, 1=BlackWhite, 2=Auto
		newMode := "2" // Auto
		if origMode == "2" {
			newMode = "0" // Color
		}
		err = c.Camera.SetDayNightConfig(ctx, map[string]string{
			setKey: newMode,
		})
		skipOnSetError(t, err, "SetDayNightConfig")

		updated, err := c.Camera.GetDayNightConfig(ctx)
		if err != nil {
			t.Fatalf("GetDayNightConfig (verify): %v", err)
		}
		if updated[modeKey] != newMode {
			t.Fatalf("expected %s=%q, got %q", modeKey, newMode, updated[modeKey])
		}
		t.Logf("Verified DayNight Mode changed to %q", newMode)
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
