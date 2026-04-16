package amcrest

import (
	"context"
	"net/url"
)

// SystemService handles system-related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 108-128 (Section 4.6)
type SystemService struct {
	client *Client
}

// GetDeviceType returns the device type (e.g., "IP4M-1041B").
// CGI: magicBox.cgi?action=getDeviceType
func (s *SystemService) GetDeviceType(ctx context.Context) (string, error) {
	body, err := s.client.cgiGet(ctx, "magicBox.cgi", "getDeviceType", nil)
	if err != nil {
		return "", err
	}
	return parseKV(body)["type"], nil
}

// GetHardwareVersion returns the hardware version string.
// CGI: magicBox.cgi?action=getHardwareVersion
func (s *SystemService) GetHardwareVersion(ctx context.Context) (string, error) {
	body, err := s.client.cgiGet(ctx, "magicBox.cgi", "getHardwareVersion", nil)
	if err != nil {
		return "", err
	}
	return parseKV(body)["version"], nil
}

// GetSerialNumber returns the device serial number.
// CGI: magicBox.cgi?action=getSerialNo
func (s *SystemService) GetSerialNumber(ctx context.Context) (string, error) {
	body, err := s.client.cgiGet(ctx, "magicBox.cgi", "getSerialNo", nil)
	if err != nil {
		return "", err
	}
	return parseKV(body)["sn"], nil
}

// GetMachineName returns the device machine name.
// CGI: magicBox.cgi?action=getMachineName
func (s *SystemService) GetMachineName(ctx context.Context) (string, error) {
	body, err := s.client.cgiGet(ctx, "magicBox.cgi", "getMachineName", nil)
	if err != nil {
		return "", err
	}
	return parseKV(body)["name"], nil
}

// GetSoftwareVersion returns the firmware/software version string.
// CGI: magicBox.cgi?action=getSoftwareVersion
func (s *SystemService) GetSoftwareVersion(ctx context.Context) (string, error) {
	body, err := s.client.cgiGet(ctx, "magicBox.cgi", "getSoftwareVersion", nil)
	if err != nil {
		return "", err
	}
	return parseKV(body)["version"], nil
}

// GetVendor returns the device vendor name (e.g., "Amcrest").
// CGI: magicBox.cgi?action=getVendor
func (s *SystemService) GetVendor(ctx context.Context) (string, error) {
	body, err := s.client.cgiGet(ctx, "magicBox.cgi", "getVendor", nil)
	if err != nil {
		return "", err
	}
	return parseKV(body)["vendor"], nil
}

// GetDeviceClass returns the device class (e.g., "IPC").
// CGI: magicBox.cgi?action=getDeviceClass
func (s *SystemService) GetDeviceClass(ctx context.Context) (string, error) {
	body, err := s.client.cgiGet(ctx, "magicBox.cgi", "getDeviceClass", nil)
	if err != nil {
		return "", err
	}
	return parseKV(body)["class"], nil
}

// GetCurrentTime returns the device's current date/time as a string (e.g., "2024-1-15 10:30:00").
// CGI: global.cgi?action=getCurrentTime
func (s *SystemService) GetCurrentTime(ctx context.Context) (string, error) {
	body, err := s.client.cgiGet(ctx, "global.cgi", "getCurrentTime", nil)
	if err != nil {
		return "", err
	}
	return parseKV(body)["result"], nil
}

// SetCurrentTime sets the device's date/time. The time parameter should be in
// "YYYY-M-D HH:MM:SS" format.
// CGI: global.cgi?action=setCurrentTime&time=<url-encoded-time>
func (s *SystemService) SetCurrentTime(ctx context.Context, timeStr string) error {
	return s.client.cgiAction(ctx, "global.cgi", "setCurrentTime", url.Values{
		"time": {timeStr},
	})
}

// Reboot reboots the device.
// CGI: magicBox.cgi?action=reboot
func (s *SystemService) Reboot(ctx context.Context) error {
	return s.client.cgiAction(ctx, "magicBox.cgi", "reboot", nil)
}

// Shutdown shuts down the device.
// CGI: magicBox.cgi?action=shutdown
func (s *SystemService) Shutdown(ctx context.Context) error {
	return s.client.cgiAction(ctx, "magicBox.cgi", "shutdown", nil)
}

// FactoryReset performs a factory reset. If keepNetwork is true, network settings
// are preserved (type=1); otherwise all settings are reset (type=0).
// CGI: magicBox.cgi?action=resetSystemEx&type=0|1
func (s *SystemService) FactoryReset(ctx context.Context, keepNetwork bool) error {
	resetType := "0"
	if keepNetwork {
		resetType = "1"
	}
	return s.client.cgiAction(ctx, "magicBox.cgi", "resetSystemEx", url.Values{
		"type": {resetType},
	})
}

// GetGeneralConfig returns the General configuration table with the
// "table.General." prefix stripped from keys.
func (s *SystemService) GetGeneralConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getConfig(ctx, "General")
}

// GetAutoMaintainConfig returns the AutoMaintain configuration table with the
// "table.AutoMaintain." prefix stripped from keys.
func (s *SystemService) GetAutoMaintainConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getConfig(ctx, "AutoMaintain")
}

// SetAutoMaintainConfig sets AutoMaintain configuration values. Keys should be
// prefixed with "AutoMaintain." (e.g., "AutoMaintain.AutoRebootDay").
func (s *SystemService) SetAutoMaintainConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetLanguageCaps returns a comma-separated list of supported languages.
// CGI: magicBox.cgi?action=getLanguageCaps
func (s *SystemService) GetLanguageCaps(ctx context.Context) (string, error) {
	body, err := s.client.cgiGet(ctx, "magicBox.cgi", "getLanguageCaps", nil)
	if err != nil {
		return "", err
	}
	return parseKV(body)["Languages"], nil
}

// GetOnvifVersion returns the ONVIF protocol version supported by the device.
// CGI: IntervideoManager.cgi?action=getVersion&Name=Onvif
func (s *SystemService) GetOnvifVersion(ctx context.Context) (string, error) {
	body, err := s.client.cgiGet(ctx, "IntervideoManager.cgi", "getVersion", url.Values{
		"Name": {"Onvif"},
	})
	if err != nil {
		return "", err
	}
	return parseKV(body)["version"], nil
}

// GetHTTPAPIVersion returns the HTTP/CGI API version supported by the device.
// CGI: IntervideoManager.cgi?action=getVersion&Name=CGI
func (s *SystemService) GetHTTPAPIVersion(ctx context.Context) (string, error) {
	body, err := s.client.cgiGet(ctx, "IntervideoManager.cgi", "getVersion", url.Values{
		"Name": {"CGI"},
	})
	if err != nil {
		return "", err
	}
	return parseKV(body)["version"], nil
}
