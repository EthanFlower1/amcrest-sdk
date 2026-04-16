package amcrest

import (
	"context"
	"errors"
	"testing"
)

func TestStorage(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	t.Run("GetAllDeviceInfo", func(t *testing.T) {
		if !hasStorageDevInfo {
			t.Skip("camera does not support storageDevice getDeviceAllInfo")
		}
		body, err := c.Storage.GetAllDeviceInfo(ctx)
		if err != nil {
			t.Fatalf("GetAllDeviceInfo: %v", err)
		}
		if body == "" {
			t.Fatal("expected non-empty device info response")
		}
		t.Logf("AllDeviceInfo response:\n%s", body)
	})

	t.Run("GetDeviceNames", func(t *testing.T) {
		if !hasStorageCollect {
			t.Skip("camera does not support storageDevice factory.getCollect")
		}
		body, err := c.Storage.GetDeviceNames(ctx)
		if err != nil {
			t.Fatalf("GetDeviceNames: %v", err)
		}
		if body == "" {
			t.Fatal("expected non-empty device names response")
		}
		t.Logf("DeviceNames response:\n%s", body)
	})

	t.Run("GetCaps", func(t *testing.T) {
		caps, err := c.Storage.GetCaps(ctx)
		if err != nil {
			t.Fatalf("GetCaps: %v", err)
		}
		if len(caps) == 0 {
			t.Fatal("expected non-empty storage caps")
		}
		for k, v := range caps {
			t.Logf("Caps.%s = %s", k, v)
		}
	})

	t.Run("GetDiskInfo", func(t *testing.T) {
		info, err := c.Storage.GetDiskInfo(ctx)
		if err != nil {
			var apiErr *APIError
			if errors.As(err, &apiErr) && apiErr.StatusCode == 400 {
				t.Skip("factory.getPortInfo not supported on this device")
			}
			t.Fatalf("GetDiskInfo: %v", err)
		}
		if len(info) == 0 {
			t.Fatal("expected non-empty disk info")
		}
		for k, v := range info {
			t.Logf("Disk.%s = %s", k, v)
		}
	})

	t.Run("GetNASConfig", func(t *testing.T) {
		cfg, err := c.Storage.GetNASConfig(ctx)
		if err != nil {
			t.Fatalf("GetNASConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty NAS config")
		}
		for k, v := range cfg {
			t.Logf("NAS.%s = %s", k, v)
		}
	})

	t.Run("GetStorageGroupConfig", func(t *testing.T) {
		cfg, err := c.Storage.GetStorageGroupConfig(ctx)
		if err != nil {
			t.Fatalf("GetStorageGroupConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty StorageGroup config")
		}
		for k, v := range cfg {
			t.Logf("StorageGroup.%s = %s", k, v)
		}
	})

	t.Run("GetStorageHealthAlarm", func(t *testing.T) {
		cfg, err := c.Storage.GetStorageHealthAlarm(ctx)
		if err != nil {
			t.Fatalf("GetStorageHealthAlarm: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty StorageHealthAlarm config")
		}
		for k, v := range cfg {
			t.Logf("HealthAlarm.%s = %s", k, v)
		}
	})

	t.Run("SetStorageHealthAlarm", func(t *testing.T) {
		original, err := c.Storage.GetStorageHealthAlarm(ctx)
		if err != nil {
			t.Fatalf("GetStorageHealthAlarm (save): %v", err)
		}

		// Find LowerLimit key.
		limitKey := ""
		origVal := ""
		for k, v := range original {
			if contains(k, "LowerLimit") {
				limitKey = k
				origVal = v
				break
			}
		}
		if limitKey == "" {
			t.Skip("no LowerLimit key found in StorageHealthAlarm config")
		}
		t.Logf("Original %s = %s", limitKey, origVal)

		setKey := limitKey
		if len(setKey) > 6 && setKey[:6] == "table." {
			setKey = setKey[6:]
		}

		defer func() {
			_ = c.Storage.SetStorageHealthAlarm(ctx, map[string]string{
				setKey: origVal,
			})
		}()

		newVal := "10"
		if origVal == "10" {
			newVal = "15"
		}
		err = c.Storage.SetStorageHealthAlarm(ctx, map[string]string{
			setKey: newVal,
		})
		skipOnSetError(t, err, "SetStorageHealthAlarm")

		updated, err := c.Storage.GetStorageHealthAlarm(ctx)
		if err != nil {
			t.Fatalf("GetStorageHealthAlarm (verify): %v", err)
		}
		if updated[limitKey] != newVal {
			t.Fatalf("expected %s=%q, got %q", limitKey, newVal, updated[limitKey])
		}
		t.Logf("Verified LowerLimit changed to %q", newVal)
	})
}
