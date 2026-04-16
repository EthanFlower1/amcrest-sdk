package amcrest

import (
	"context"
	"errors"
	"testing"
)

func TestThermalGetCaps(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	caps, err := c.Thermal.GetCaps(ctx)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			t.Skipf("RadiometryManager getCaps not supported on this device (HTTP %d)", apiErr.StatusCode)
		}
		t.Fatalf("GetCaps: %v", err)
	}
	if len(caps) == 0 {
		t.Fatal("expected non-empty caps")
	}
	for k, v := range caps {
		t.Logf("Caps.%s = %s", k, v)
	}
}

func TestThermalGetThermographyOptions(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Thermal.GetThermographyOptions(ctx)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			t.Skipf("ThermographyOptions config not supported on this device (HTTP %d)", apiErr.StatusCode)
		}
		t.Fatalf("GetThermographyOptions: %v", err)
	}
	if len(cfg) == 0 {
		t.Fatal("expected non-empty ThermographyOptions config")
	}
	for k, v := range cfg {
		t.Logf("ThermographyOptions.%s = %s", k, v)
	}
}

func TestThermalGetRadiometryCaps(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	v, err := c.Thermal.GetRadiometryCaps(ctx)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			t.Skipf("RadiometryManager getRadiometryCaps not supported on this device (HTTP %d)", apiErr.StatusCode)
		}
		t.Fatalf("GetRadiometryCaps: %v", err)
	}
	if v == "" {
		t.Fatal("expected non-empty radiometry caps")
	}
	t.Logf("RadiometryCaps:\n%s", v)
}

func TestThermalGetFireWarningConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Thermal.GetFireWarningConfig(ctx)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			t.Skipf("FireWarning config not supported on this device (HTTP %d)", apiErr.StatusCode)
		}
		t.Fatalf("GetFireWarningConfig: %v", err)
	}
	if len(cfg) == 0 {
		t.Fatal("expected non-empty FireWarning config")
	}
	for k, v := range cfg {
		t.Logf("FireWarning.%s = %s", k, v)
	}
}
