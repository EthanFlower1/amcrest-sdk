package amcrest

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestSystemGetDeviceType(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	v, err := c.System.GetDeviceType(ctx)
	if err != nil {
		t.Fatalf("GetDeviceType: %v", err)
	}
	if v == "" {
		t.Fatal("expected non-empty device type")
	}
	t.Logf("DeviceType: %s", v)
}

func TestSystemGetHardwareVersion(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	v, err := c.System.GetHardwareVersion(ctx)
	if err != nil {
		t.Fatalf("GetHardwareVersion: %v", err)
	}
	if v == "" {
		t.Fatal("expected non-empty hardware version")
	}
	t.Logf("HardwareVersion: %s", v)
}

func TestSystemGetSerialNumber(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	v, err := c.System.GetSerialNumber(ctx)
	if err != nil {
		t.Fatalf("GetSerialNumber: %v", err)
	}
	if v == "" {
		t.Fatal("expected non-empty serial number")
	}
	t.Logf("SerialNumber: %s", v)
}

func TestSystemGetMachineName(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	v, err := c.System.GetMachineName(ctx)
	if err != nil {
		t.Fatalf("GetMachineName: %v", err)
	}
	if v == "" {
		t.Fatal("expected non-empty machine name")
	}
	t.Logf("MachineName: %s", v)
}

func TestSystemGetSoftwareVersion(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	v, err := c.System.GetSoftwareVersion(ctx)
	if err != nil {
		t.Fatalf("GetSoftwareVersion: %v", err)
	}
	if v == "" {
		t.Fatal("expected non-empty software version")
	}
	t.Logf("SoftwareVersion: %s", v)
}

func TestSystemGetVendor(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	v, err := c.System.GetVendor(ctx)
	if err != nil {
		t.Fatalf("GetVendor: %v", err)
	}
	if v == "" {
		t.Fatal("expected non-empty vendor")
	}
	t.Logf("Vendor: %s", v)
}

func TestSystemGetDeviceClass(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	v, err := c.System.GetDeviceClass(ctx)
	if err != nil {
		t.Fatalf("GetDeviceClass: %v", err)
	}
	if v == "" {
		t.Fatal("expected non-empty device class")
	}
	t.Logf("DeviceClass: %s", v)
}

func TestSystemGetCurrentTime(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	v, err := c.System.GetCurrentTime(ctx)
	if err != nil {
		t.Fatalf("GetCurrentTime: %v", err)
	}
	if v == "" {
		t.Fatal("expected non-empty current time")
	}
	t.Logf("CurrentTime: %s", v)
}

func TestSystemSetCurrentTime(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	// Save the current time.
	original, err := c.System.GetCurrentTime(ctx)
	if err != nil {
		t.Fatalf("GetCurrentTime (save): %v", err)
	}
	t.Logf("Original time: %s", original)

	// Set a new time (1 second from now to ensure it's different).
	newTime := time.Now().Add(1 * time.Second).Format("2006-1-2 15:04:05")
	if err := c.System.SetCurrentTime(ctx, newTime); err != nil {
		t.Fatalf("SetCurrentTime: %v", err)
	}
	t.Logf("Set time to: %s", newTime)

	// Verify the time was updated.
	updated, err := c.System.GetCurrentTime(ctx)
	if err != nil {
		t.Fatalf("GetCurrentTime (verify): %v", err)
	}
	t.Logf("Updated time: %s", updated)
	if updated == "" {
		t.Fatal("expected non-empty updated time")
	}

	// Restore the original time (use current wall clock since the camera clock drifted).
	restoreTime := time.Now().Format("2006-1-2 15:04:05")
	if err := c.System.SetCurrentTime(ctx, restoreTime); err != nil {
		t.Fatalf("SetCurrentTime (restore): %v", err)
	}
	t.Logf("Restored time to: %s", restoreTime)
}

func TestSystemGetGeneralConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
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
}

func TestSystemGetAutoMaintainConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
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
}

func TestSystemGetLanguageCaps(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	v, err := c.System.GetLanguageCaps(ctx)
	if err != nil {
		t.Fatalf("GetLanguageCaps: %v", err)
	}
	if v == "" {
		t.Fatal("expected non-empty language caps")
	}
	t.Logf("LanguageCaps: %s", v)
}

func TestSystemGetOnvifVersion(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	v, err := c.System.GetOnvifVersion(ctx)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 501 {
			t.Skip("IntervideoManager not supported on this device")
		}
		t.Fatalf("GetOnvifVersion: %v", err)
	}
	if v == "" {
		t.Fatal("expected non-empty ONVIF version")
	}
	t.Logf("OnvifVersion: %s", v)
}

func TestSystemGetHTTPAPIVersion(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	v, err := c.System.GetHTTPAPIVersion(ctx)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 501 {
			t.Skip("IntervideoManager not supported on this device")
		}
		t.Fatalf("GetHTTPAPIVersion: %v", err)
	}
	if v == "" {
		t.Fatal("expected non-empty HTTP API version")
	}
	t.Logf("HTTPAPIVersion: %s", v)
}
