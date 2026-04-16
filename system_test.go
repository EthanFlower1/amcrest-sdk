package amcrest

import (
	"context"
	"testing"
	"time"
)

func TestSystem(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	t.Run("GetDeviceType", func(t *testing.T) {
		v, err := c.System.GetDeviceType(ctx)
		if err != nil {
			t.Fatalf("GetDeviceType: %v", err)
		}
		if v == "" {
			t.Fatal("expected non-empty device type")
		}
		t.Logf("DeviceType: %s", v)
	})

	t.Run("GetHardwareVersion", func(t *testing.T) {
		v, err := c.System.GetHardwareVersion(ctx)
		if err != nil {
			t.Fatalf("GetHardwareVersion: %v", err)
		}
		if v == "" {
			t.Fatal("expected non-empty hardware version")
		}
		t.Logf("HardwareVersion: %s", v)
	})

	t.Run("GetSerialNumber", func(t *testing.T) {
		v, err := c.System.GetSerialNumber(ctx)
		if err != nil {
			t.Fatalf("GetSerialNumber: %v", err)
		}
		if v == "" {
			t.Fatal("expected non-empty serial number")
		}
		t.Logf("SerialNumber: %s", v)
	})

	t.Run("GetMachineName", func(t *testing.T) {
		v, err := c.System.GetMachineName(ctx)
		if err != nil {
			t.Fatalf("GetMachineName: %v", err)
		}
		if v == "" {
			t.Fatal("expected non-empty machine name")
		}
		t.Logf("MachineName: %s", v)
	})

	t.Run("GetSoftwareVersion", func(t *testing.T) {
		v, err := c.System.GetSoftwareVersion(ctx)
		if err != nil {
			t.Fatalf("GetSoftwareVersion: %v", err)
		}
		if v == "" {
			t.Fatal("expected non-empty software version")
		}
		t.Logf("SoftwareVersion: %s", v)
	})

	t.Run("GetVendor", func(t *testing.T) {
		v, err := c.System.GetVendor(ctx)
		if err != nil {
			t.Fatalf("GetVendor: %v", err)
		}
		if v == "" {
			t.Fatal("expected non-empty vendor")
		}
		t.Logf("Vendor: %s", v)
	})

	t.Run("GetDeviceClass", func(t *testing.T) {
		v, err := c.System.GetDeviceClass(ctx)
		if err != nil {
			t.Fatalf("GetDeviceClass: %v", err)
		}
		if v == "" {
			t.Fatal("expected non-empty device class")
		}
		t.Logf("DeviceClass: %s", v)
	})

	t.Run("GetCurrentTime", func(t *testing.T) {
		v, err := c.System.GetCurrentTime(ctx)
		if err != nil {
			t.Fatalf("GetCurrentTime: %v", err)
		}
		if v == "" {
			t.Fatal("expected non-empty current time")
		}
		t.Logf("CurrentTime: %s", v)
	})

	t.Run("SetCurrentTime", func(t *testing.T) {
		original, err := c.System.GetCurrentTime(ctx)
		if err != nil {
			t.Fatalf("GetCurrentTime (save): %v", err)
		}
		t.Logf("Original time: %s", original)

		newTime := time.Now().Add(1 * time.Second).Format("2006-1-2 15:04:05")
		if err := c.System.SetCurrentTime(ctx, newTime); err != nil {
			t.Fatalf("SetCurrentTime: %v", err)
		}
		t.Logf("Set time to: %s", newTime)

		updated, err := c.System.GetCurrentTime(ctx)
		if err != nil {
			t.Fatalf("GetCurrentTime (verify): %v", err)
		}
		t.Logf("Updated time: %s", updated)
		if updated == "" {
			t.Fatal("expected non-empty updated time")
		}

		restoreTime := time.Now().Format("2006-1-2 15:04:05")
		if err := c.System.SetCurrentTime(ctx, restoreTime); err != nil {
			t.Fatalf("SetCurrentTime (restore): %v", err)
		}
		t.Logf("Restored time to: %s", restoreTime)
	})

	t.Run("GetGeneralConfig", func(t *testing.T) {
		cfg, err := c.System.GetGeneralConfig(ctx)
		if err != nil {
			t.Fatalf("GetGeneralConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty General config")
		}
		for k, v := range cfg {
			t.Logf("General.%s = %s", k, v)
		}
	})

	t.Run("GetAutoMaintainConfig", func(t *testing.T) {
		cfg, err := c.System.GetAutoMaintainConfig(ctx)
		if err != nil {
			t.Fatalf("GetAutoMaintainConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty AutoMaintain config")
		}
		for k, v := range cfg {
			t.Logf("AutoMaintain.%s = %s", k, v)
		}
	})

	t.Run("GetLocalesConfig", func(t *testing.T) {
		cfg, err := c.System.GetLocalesConfig(ctx)
		if err != nil {
			t.Fatalf("GetLocalesConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty Locales config")
		}
		for k, v := range cfg {
			t.Logf("Locales.%s = %s", k, v)
		}
	})

	t.Run("GetHolidayConfig", func(t *testing.T) {
		cfg, err := c.System.GetHolidayConfig(ctx)
		if err != nil {
			t.Fatalf("GetHolidayConfig: %v", err)
		}
		t.Logf("HolidayConfig entries: %d", len(cfg))
		for k, v := range cfg {
			t.Logf("Holiday.%s = %s", k, v)
		}
	})

	t.Run("GetLanguageCaps", func(t *testing.T) {
		v, err := c.System.GetLanguageCaps(ctx)
		if err != nil {
			t.Fatalf("GetLanguageCaps: %v", err)
		}
		if v == "" {
			t.Fatal("expected non-empty language caps")
		}
		t.Logf("LanguageCaps: %s", v)
	})

	t.Run("GetLanguage", func(t *testing.T) {
		v, err := c.System.GetLanguage(ctx)
		if err != nil {
			t.Logf("GetLanguage not available: %v", err)
			return
		}
		t.Logf("Language: %q", v)
	})

	t.Run("GetSystemInfo", func(t *testing.T) {
		info, err := c.System.GetSystemInfo(ctx)
		if err != nil {
			t.Fatalf("GetSystemInfo: %v", err)
		}
		if len(info) == 0 {
			t.Fatal("expected non-empty system info")
		}
		for k, v := range info {
			t.Logf("SystemInfo.%s = %s", k, v)
		}
	})

	t.Run("GetOnvifVersion", func(t *testing.T) {
		v, err := c.System.GetOnvifVersion(ctx)
		if err != nil {
			t.Logf("GetOnvifVersion not available: %v", err)
			return
		}
		t.Logf("OnvifVersion: %s", v)
	})

	t.Run("GetHTTPAPIVersion", func(t *testing.T) {
		v, err := c.System.GetHTTPAPIVersion(ctx)
		if err != nil {
			t.Logf("GetHTTPAPIVersion not available: %v", err)
			return
		}
		t.Logf("HTTPAPIVersion: %s", v)
	})
}
