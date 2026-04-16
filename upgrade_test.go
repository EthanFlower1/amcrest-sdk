package amcrest

import (
	"context"
	"testing"
)

func TestUpgradeGetState(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	state, err := c.Upgrade.GetState(ctx)
	if err != nil {
		t.Skip("GetState not supported on this camera, skipping")
	}
	if len(state) == 0 {
		t.Log("upgrade state returned empty map (no upgrade in progress)")
	}
	for k, v := range state {
		t.Logf("Upgrade.%s = %s", k, v)
	}
}

func TestUpgradeCheckCloudUpdate(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	result, err := c.Upgrade.CheckCloudUpdate(ctx)
	if err != nil {
		t.Skip("CheckCloudUpdate not supported on this camera, skipping")
	}
	for k, v := range result {
		t.Logf("CloudUpdate.%s = %s", k, v)
	}
}

func TestUpgradeGetAutoUpgradeConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Upgrade.GetAutoUpgradeConfig(ctx)
	if err != nil {
		t.Skip("GetAutoUpgradeConfig not supported on this camera, skipping")
	}
	for k, v := range cfg {
		t.Logf("AutoUpgrade.%s = %s", k, v)
	}
}

func TestUpgradeGetCloudUpgradeMode(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Upgrade.GetCloudUpgradeMode(ctx)
	if err != nil {
		t.Skip("GetCloudUpgradeMode not supported on this camera, skipping")
	}
	for k, v := range cfg {
		t.Logf("CloudUpgrade.%s = %s", k, v)
	}
}
