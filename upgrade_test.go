package amcrest

import (
	"context"
	"testing"
)

func TestUpgrade(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	t.Run("GetState", func(t *testing.T) {
		if !hasUpgrade {
			t.Skip("camera does not support upgrader.cgi getState")
		}
		state, err := c.Upgrade.GetState(ctx)
		if err != nil {
			t.Fatalf("GetState: %v", err)
		}
		if len(state) == 0 {
			t.Log("upgrade state returned empty map (no upgrade in progress)")
		}
		for k, v := range state {
			t.Logf("Upgrade.%s = %s", k, v)
		}
	})

	t.Run("CheckCloudUpdate", func(t *testing.T) {
		if !hasCloudUpgrade {
			t.Skip("camera does not support CloudUpgrader/check")
		}
		result, err := c.Upgrade.CheckCloudUpdate(ctx)
		if err != nil {
			t.Fatalf("CheckCloudUpdate: %v", err)
		}
		for k, v := range result {
			t.Logf("CloudUpdate.%s = %s", k, v)
		}
	})

	t.Run("GetAutoUpgradeConfig", func(t *testing.T) {
		if !hasAutoUpgrade {
			t.Skip("camera does not support AutoUpgrade config")
		}
		cfg, err := c.Upgrade.GetAutoUpgradeConfig(ctx)
		if err != nil {
			t.Fatalf("GetAutoUpgradeConfig: %v", err)
		}
		for k, v := range cfg {
			t.Logf("AutoUpgrade.%s = %s", k, v)
		}
	})

	t.Run("GetCloudUpgradeMode", func(t *testing.T) {
		if !hasCloudMode {
			t.Skip("camera does not support CloudUpgrade config")
		}
		cfg, err := c.Upgrade.GetCloudUpgradeMode(ctx)
		if err != nil {
			t.Fatalf("GetCloudUpgradeMode: %v", err)
		}
		for k, v := range cfg {
			t.Logf("CloudUpgrade.%s = %s", k, v)
		}
	})
}
