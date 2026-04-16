package amcrest

import (
	"context"
	"errors"
	"testing"
)

func skipIfNoPeripheral(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		return
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		t.Skipf("camera does not support peripheral feature (HTTP %d), skipping", apiErr.StatusCode)
	}
}

func TestPeripheralGetCoaxialIOStatus(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	status, err := c.Peripheral.GetCoaxialIOStatus(ctx, 0)
	skipIfNoPeripheral(t, err)
	if err != nil {
		t.Fatalf("GetCoaxialIOStatus: %v", err)
	}
	for k, v := range status {
		t.Logf("%s = %s", k, v)
	}
}

func TestPeripheralGetFlashlightConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	cfg, err := c.Peripheral.GetFlashlightConfig(ctx)
	skipIfNoPeripheral(t, err)
	if err != nil {
		t.Fatalf("GetFlashlightConfig: %v", err)
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}

func TestPeripheralGetGPSConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	cfg, err := c.Peripheral.GetGPSConfig(ctx)
	skipIfNoPeripheral(t, err)
	if err != nil {
		t.Fatalf("GetGPSConfig: %v", err)
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}

func TestPeripheralGetFishEyeConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	cfg, err := c.Peripheral.GetFishEyeConfig(ctx)
	skipIfNoPeripheral(t, err)
	if err != nil {
		t.Fatalf("GetFishEyeConfig: %v", err)
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}
