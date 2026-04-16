package amcrest

import (
	"context"
	"testing"
)

func TestPeripheral(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	t.Run("GetCoaxialIOStatus", func(t *testing.T) {
		if !hasCoaxialIO {
			t.Skip("camera does not support CoaxialIO")
		}
		status, err := c.Peripheral.GetCoaxialIOStatus(ctx, 0)
		if err != nil {
			t.Fatalf("GetCoaxialIOStatus: %v", err)
		}
		for k, v := range status {
			t.Logf("CoaxialIO.%s = %s", k, v)
		}
	})

	t.Run("GetFlashlightConfig", func(t *testing.T) {
		if !hasFlashlight {
			t.Skip("camera does not support FlashLight config")
		}
		cfg, err := c.Peripheral.GetFlashlightConfig(ctx)
		if err != nil {
			t.Fatalf("GetFlashlightConfig: %v", err)
		}
		for k, v := range cfg {
			t.Logf("Flashlight.%s = %s", k, v)
		}
	})

	t.Run("GetGPSConfig", func(t *testing.T) {
		requireCapability(t, hasGPS, "GPS")
		cfg, err := c.Peripheral.GetGPSConfig(ctx)
		if err != nil {
			t.Fatalf("GetGPSConfig: %v", err)
		}
		for k, v := range cfg {
			t.Logf("GPS.%s = %s", k, v)
		}
	})

	t.Run("GetFishEyeConfig", func(t *testing.T) {
		if !hasFishEye {
			t.Skip("camera does not support FishEye config")
		}
		cfg, err := c.Peripheral.GetFishEyeConfig(ctx)
		if err != nil {
			t.Fatalf("GetFishEyeConfig: %v", err)
		}
		for k, v := range cfg {
			t.Logf("FishEye.%s = %s", k, v)
		}
	})
}
