package amcrest

import (
	"context"
	"testing"
)

func TestAccessControl(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	requireCapability(t, hasAccessCtrl, "Access Control")

	t.Run("GetDoorStatus", func(t *testing.T) {
		status, err := c.AccessControl.GetDoorStatus(ctx, 0)
		if err != nil {
			t.Fatalf("GetDoorStatus: %v", err)
		}
		t.Logf("Door status: %s", status)
	})

	t.Run("GetGeneralConfig", func(t *testing.T) {
		cfg, err := c.AccessControl.GetGeneralConfig(ctx)
		if err != nil {
			t.Fatalf("GetGeneralConfig: %v", err)
		}
		for k, v := range cfg {
			t.Logf("General.%s = %s", k, v)
		}
	})

	t.Run("GetControlConfig", func(t *testing.T) {
		cfg, err := c.AccessControl.GetControlConfig(ctx)
		if err != nil {
			t.Fatalf("GetControlConfig: %v", err)
		}
		for k, v := range cfg {
			t.Logf("Control.%s = %s", k, v)
		}
	})
}
