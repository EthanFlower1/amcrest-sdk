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
		state, err := c.Upgrade.GetState(ctx)
		if err != nil {
			t.Logf("GetState not available: %v", err)
			return
		}
		if len(state) == 0 {
			t.Log("upgrade state returned empty map (no upgrade in progress)")
		}
		for k, v := range state {
			t.Logf("Upgrade.%s = %s", k, v)
		}
	})

	t.Run("CheckCloudUpdate", func(t *testing.T) {
		result, err := c.Upgrade.CheckCloudUpdate(ctx)
		if err != nil {
			t.Logf("CheckCloudUpdate not available: %v", err)
			return
		}
		for k, v := range result {
			t.Logf("CloudUpdate.%s = %s", k, v)
		}
	})

	t.Run("GetAutoUpgradeConfig", func(t *testing.T) {
		cfg, err := c.Upgrade.GetAutoUpgradeConfig(ctx)
		if err != nil {
			t.Logf("GetAutoUpgradeConfig not available: %v", err)
			return
		}
		for k, v := range cfg {
			t.Logf("AutoUpgrade.%s = %s", k, v)
		}
	})

	t.Run("GetCloudUpgradeMode", func(t *testing.T) {
		cfg, err := c.Upgrade.GetCloudUpgradeMode(ctx)
		if err != nil {
			t.Logf("GetCloudUpgradeMode not available: %v", err)
			return
		}
		for k, v := range cfg {
			t.Logf("CloudUpgrade.%s = %s", k, v)
		}
	})
}
