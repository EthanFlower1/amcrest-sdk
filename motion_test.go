package amcrest

import (
	"context"
	"testing"
)

func TestMotion(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	t.Run("GetConfig", func(t *testing.T) {
		cfg, err := c.Motion.GetConfig(ctx)
		if err != nil {
			t.Fatalf("Motion.GetConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty config")
		}
		for k, v := range cfg {
			t.Logf("Motion.%s = %s", k, v)
		}
	})

	t.Run("SetConfig", func(t *testing.T) {
		original, err := c.Motion.GetConfig(ctx)
		if err != nil {
			t.Fatalf("Motion.GetConfig (save): %v", err)
		}

		// Find Enable key (e.g., "table.MotionDetect[0].Enable").
		enableKey := ""
		origVal := ""
		for k, v := range original {
			if contains(k, "Enable") && contains(k, "MotionDetect[0]") {
				enableKey = k
				origVal = v
				break
			}
		}
		if enableKey == "" {
			t.Skip("no MotionDetect[0].Enable key found")
		}
		t.Logf("Original %s = %s", enableKey, origVal)

		setKey := enableKey
		if len(setKey) > 6 && setKey[:6] == "table." {
			setKey = setKey[6:]
		}

		defer func() {
			_ = c.Motion.SetConfig(ctx, map[string]string{
				setKey: origVal,
			})
		}()

		newVal := "true"
		if origVal == "true" {
			newVal = "false"
		}
		err = c.Motion.SetConfig(ctx, map[string]string{
			setKey: newVal,
		})
		skipOnSetError(t, err, "Motion.SetConfig")

		updated, err := c.Motion.GetConfig(ctx)
		if err != nil {
			t.Fatalf("Motion.GetConfig (verify): %v", err)
		}
		if updated[enableKey] != newVal {
			t.Fatalf("expected %s=%q, got %q", enableKey, newVal, updated[enableKey])
		}
		t.Logf("Verified MotionDetect Enable changed to %q", newVal)
	})

	t.Run("GetSmartMotionConfig", func(t *testing.T) {
		if !hasSmartMotion {
			t.Skip("camera does not support SmartMotionDetect config")
		}
		cfg, err := c.Motion.GetSmartMotionConfig(ctx)
		if err != nil {
			t.Fatalf("GetSmartMotionConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty SmartMotion config")
		}
		for k, v := range cfg {
			t.Logf("SmartMotion.%s = %s", k, v)
		}
	})

	t.Run("SetSmartMotionConfig", func(t *testing.T) {
		if !hasSmartMotion {
			t.Skip("camera does not support SmartMotionDetect config")
		}
		original, err := c.Motion.GetSmartMotionConfig(ctx)
		if err != nil {
			t.Fatalf("GetSmartMotionConfig (save): %v", err)
		}

		// Find a Sensitivity key to modify.
		sensKey := ""
		origVal := ""
		for k, v := range original {
			if contains(k, "Sensitivity") {
				sensKey = k
				origVal = v
				break
			}
		}
		if sensKey == "" {
			// Fall back to Enable.
			for k, v := range original {
				if contains(k, "Enable") {
					sensKey = k
					origVal = v
					break
				}
			}
		}
		if sensKey == "" {
			t.Skip("no Sensitivity or Enable key found in SmartMotionConfig")
		}
		t.Logf("Original %s = %s", sensKey, origVal)

		setKey := sensKey
		if len(setKey) > 6 && setKey[:6] == "table." {
			setKey = setKey[6:]
		}

		defer func() {
			_ = c.Motion.SetSmartMotionConfig(ctx, map[string]string{
				setKey: origVal,
			})
		}()

		// Toggle value.
		newVal := origVal
		if contains(sensKey, "Sensitivity") {
			newVal = "50"
			if origVal == "50" {
				newVal = "60"
			}
		} else {
			newVal = "true"
			if origVal == "true" {
				newVal = "false"
			}
		}
		err = c.Motion.SetSmartMotionConfig(ctx, map[string]string{
			setKey: newVal,
		})
		skipOnSetError(t, err, "SetSmartMotionConfig")

		updated, err := c.Motion.GetSmartMotionConfig(ctx)
		if err != nil {
			t.Fatalf("GetSmartMotionConfig (verify): %v", err)
		}
		if updated[sensKey] != newVal {
			t.Fatalf("expected %s=%q, got %q", sensKey, newVal, updated[sensKey])
		}
		t.Logf("Verified SmartMotion value changed to %q", newVal)
	})

	t.Run("GetLAEConfig", func(t *testing.T) {
		if !hasLAEConfig {
			t.Skip("camera does not support LAEConfig config")
		}
		cfg, err := c.Motion.GetLAEConfig(ctx)
		if err != nil {
			t.Fatalf("GetLAEConfig: %v", err)
		}
		for k, v := range cfg {
			t.Logf("LAE.%s = %s", k, v)
		}
	})
}
