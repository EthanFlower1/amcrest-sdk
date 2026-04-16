package amcrest

import (
	"context"
	"testing"
)

func TestThermal(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	requireCapability(t, hasThermal, "Thermal Imaging")

	t.Run("GetCaps", func(t *testing.T) {
		caps, err := c.Thermal.GetCaps(ctx)
		if err != nil {
			t.Fatalf("GetCaps: %v", err)
		}
		if len(caps) == 0 {
			t.Fatal("expected non-empty caps")
		}
		for k, v := range caps {
			t.Logf("Caps.%s = %s", k, v)
		}
	})

	t.Run("GetThermographyOptions", func(t *testing.T) {
		cfg, err := c.Thermal.GetThermographyOptions(ctx)
		if err != nil {
			t.Fatalf("GetThermographyOptions: %v", err)
		}
		for k, v := range cfg {
			t.Logf("ThermographyOptions.%s = %s", k, v)
		}
	})

	t.Run("GetRadiometryCaps", func(t *testing.T) {
		v, err := c.Thermal.GetRadiometryCaps(ctx)
		if err != nil {
			t.Fatalf("GetRadiometryCaps: %v", err)
		}
		t.Logf("RadiometryCaps:\n%s", v)
	})

	t.Run("GetFireWarningConfig", func(t *testing.T) {
		cfg, err := c.Thermal.GetFireWarningConfig(ctx)
		if err != nil {
			t.Fatalf("GetFireWarningConfig: %v", err)
		}
		for k, v := range cfg {
			t.Logf("FireWarning.%s = %s", k, v)
		}
	})
}
