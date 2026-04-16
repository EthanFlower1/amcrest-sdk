//go:build dangerous

package amcrest

import (
	"context"
	"testing"
)

func TestSystemReboot(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	if err := c.System.Reboot(ctx); err != nil {
		t.Fatalf("Reboot: %v", err)
	}
	t.Log("Reboot command sent successfully")
}

func TestSystemShutdown(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	if err := c.System.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}
	t.Log("Shutdown command sent successfully")
}

func TestSystemFactoryReset(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	// Keep network settings so the camera remains reachable.
	if err := c.System.FactoryReset(ctx, true); err != nil {
		t.Fatalf("FactoryReset: %v", err)
	}
	t.Log("FactoryReset (keepNetwork=true) command sent successfully")
}
