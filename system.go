package amcrest

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
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

// SetGeneralConfig sets General configuration values. Keys should be
// prefixed with "General." (e.g., "General.MachineName").
func (s *SystemService) SetGeneralConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetLocalesConfig returns the Locales configuration table with the
// "table.Locales." prefix stripped from keys.
func (s *SystemService) GetLocalesConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getConfig(ctx, "Locales")
}

// SetLocalesConfig sets Locales configuration values. Keys should be
// prefixed with "Locales." (e.g., "Locales.DSTEnable").
func (s *SystemService) SetLocalesConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetHolidayConfig returns the Holiday configuration table without stripping
// any prefix, since Holiday entries use an indexed format.
func (s *SystemService) GetHolidayConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "Holiday")
}

// SetHolidayConfig sets Holiday configuration values.
func (s *SystemService) SetHolidayConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetLanguage returns the current UI language setting (e.g., "English").
func (s *SystemService) GetLanguage(ctx context.Context) (string, error) {
	cfg, err := s.client.getConfig(ctx, "Language")
	if err != nil {
		return "", err
	}
	return cfg["CurrentLanguage"], nil
}

// SetLanguage sets the current UI language (e.g., "English", "SimpChinese").
func (s *SystemService) SetLanguage(ctx context.Context, lang string) error {
	return s.client.setConfig(ctx, map[string]string{
		"Language.CurrentLanguage": lang,
	})
}

// GetSystemInfo returns system information (e.g., serialNumber, hardwareVersion).
// CGI: magicBox.cgi?action=getSystemInfo
func (s *SystemService) GetSystemInfo(ctx context.Context) (map[string]string, error) {
	body, err := s.client.cgiGet(ctx, "magicBox.cgi", "getSystemInfo", nil)
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}

// GetSystemInfoNew returns extended system information.
// CGI: magicBox.cgi?action=getSystemInfoNew
func (s *SystemService) GetSystemInfoNew(ctx context.Context) (map[string]string, error) {
	body, err := s.client.cgiGet(ctx, "magicBox.cgi", "getSystemInfoNew", nil)
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}

// GetTracingCode returns the device tracing code.
// CGI: magicBox.cgi?action=getTracingCode
func (s *SystemService) GetTracingCode(ctx context.Context) (string, error) {
	body, err := s.client.cgiGet(ctx, "magicBox.cgi", "getTracingCode", nil)
	if err != nil {
		return "", err
	}
	return parseKV(body)["tc"], nil
}

// GetCompleteMachineVersion returns the complete machine version string.
// POST /cgi-bin/api/MagicBox/getCompleteMachineVersion
func (s *SystemService) GetCompleteMachineVersion(ctx context.Context) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/MagicBox/getCompleteMachineVersion", nil)
}

// TCPTest tests TCP connectivity to the given IP and port.
// POST /cgi-bin/api/tcpConnect/tcpTest with JSON {Ip, Port}.
// Returns true if the connection succeeds.
func (s *SystemService) TCPTest(ctx context.Context, ip string, port int) (bool, error) {
	reqBody := struct {
		Ip   string `json:"Ip"`
		Port int    `json:"Port"`
	}{Ip: ip, Port: port}

	raw, err := s.client.postRaw(ctx, "/cgi-bin/api/tcpConnect/tcpTest", reqBody)
	if err != nil {
		return false, err
	}

	// Response contains "Connect" boolean field
	var result struct {
		Connect bool `json:"Connect"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &result); err != nil {
		// Fallback: parse as key-value
		kv := parseKV(raw)
		return kv["Connect"] == "true", nil
	}
	return result.Connect, nil
}

// GetRecordingStateAll returns the recording state for all channels.
// POST /cgi-bin/api/recordManager/getStateAll
func (s *SystemService) GetRecordingStateAll(ctx context.Context) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/recordManager/getStateAll", nil)
}

// AddCamera adds a camera (or cameras) by group.
// POST /cgi-bin/LogicDeviceManager.cgi?action=addCameraByGroup
// The body should be the request payload as defined by the API.
func (s *SystemService) AddCamera(ctx context.Context, body interface{}) (string, error) {
	params := url.Values{"action": {"addCameraByGroup"}}
	path := "/cgi-bin/LogicDeviceManager.cgi?" + params.Encode()
	return s.client.postRaw(ctx, path, body)
}

// GetCameraAll returns all camera information.
// POST /cgi-bin/api/LogicDeviceManager/getCameraAll
func (s *SystemService) GetCameraAll(ctx context.Context) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/LogicDeviceManager/getCameraAll", nil)
}

// GetCameraState returns the connection state for the given channels.
// POST /cgi-bin/api/LogicDeviceManager/getCameraState with JSON {channel: [...]}.
func (s *SystemService) GetCameraState(ctx context.Context, channels []int) (string, error) {
	strs := make([]string, len(channels))
	for i, ch := range channels {
		strs[i] = strconv.Itoa(ch)
	}
	reqBody := struct {
		Channel []string `json:"channel"`
	}{Channel: strs}
	return s.client.postRaw(ctx, "/cgi-bin/api/LogicDeviceManager/getCameraState", reqBody)
}
