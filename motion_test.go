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
