package amcrest

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// PTZService handles PTZ (pan-tilt-zoom) related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 284-305 (Section 8.1)
type PTZService struct {
	client *Client
}

// PTZ action code constants.
const (
	PTZUp        = "Up"
	PTZDown      = "Down"
	PTZLeft      = "Left"
	PTZRight     = "Right"
	PTZZoomTele  = "ZoomTele"
	PTZZoomWide  = "ZoomWide"
	PTZFocusNear = "FocusNear"
	PTZFocusFar  = "FocusFar"
	PTZLeftUp    = "LeftUp"
	PTZRightUp   = "RightUp"
	PTZLeftDown     = "LeftDown"
	PTZRightDown    = "RightDown"
	PTZIrisLarge    = "IrisLarge"
	PTZIrisSmall    = "IrisSmall"
	PTZContinuously = "Continuously"
)

// Control sends a PTZ start command to the camera.
// ptz.cgi?action=start&channel=N&code=X&arg1=A&arg2=B&arg3=C
func (s *PTZService) Control(ctx context.Context, channel int, code string, arg1, arg2, arg3 int) error {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
		"code":    {code},
		"arg1":    {fmt.Sprintf("%d", arg1)},
		"arg2":    {fmt.Sprintf("%d", arg2)},
		"arg3":    {fmt.Sprintf("%d", arg3)},
	}
	return s.client.cgiAction(ctx, "ptz.cgi", "start", params)
}

// Stop sends a PTZ stop command to the camera.
// ptz.cgi?action=stop&channel=N&code=X&arg1=0&arg2=0&arg3=0
func (s *PTZService) Stop(ctx context.Context, channel int, code string) error {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
		"code":    {code},
		"arg1":    {"0"},
		"arg2":    {"0"},
		"arg3":    {"0"},
	}
	return s.client.cgiAction(ctx, "ptz.cgi", "stop", params)
}

// GotoPreset moves the camera to a saved preset position.
func (s *PTZService) GotoPreset(ctx context.Context, channel, preset int) error {
	return s.Control(ctx, channel, "GotoPreset", 0, preset, 0)
}

// SetPreset saves the current camera position as a preset.
func (s *PTZService) SetPreset(ctx context.Context, channel, preset int) error {
	return s.Control(ctx, channel, "SetPreset", 0, preset, 0)
}

// ClearPreset removes a saved preset position.
func (s *PTZService) ClearPreset(ctx context.Context, channel, preset int) error {
	return s.Control(ctx, channel, "ClearPreset", 0, preset, 0)
}

// GetPresets retrieves all configured presets for the given channel.
// Returns the raw response body as a string.
func (s *PTZService) GetPresets(ctx context.Context, channel int) (string, error) {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	return s.client.cgiGet(ctx, "ptz.cgi", "getPresets", params)
}

// GetStatus retrieves the current PTZ status for the given channel.
func (s *PTZService) GetStatus(ctx context.Context, channel int) (map[string]string, error) {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	body, err := s.client.cgiGet(ctx, "ptz.cgi", "getStatus", params)
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}

// GetConfig retrieves the PTZ configuration via configManager.
func (s *PTZService) GetConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getConfig(ctx, "Ptz")
}

// SetConfig sets PTZ configuration values via configManager.
// CGI: configManager.cgi?action=setConfig
func (s *PTZService) SetConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetProtocolList returns the list of supported PTZ protocols for a channel.
// CGI: ptz.cgi?action=getProtocolList&channel=N
func (s *PTZService) GetProtocolList(ctx context.Context, channel int) (string, error) {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	return s.client.cgiGet(ctx, "ptz.cgi", "getProtocolList", params)
}

// GetCaps returns the current PTZ protocol capabilities for a channel.
// CGI: ptz.cgi?action=getCurrentProtocolCaps&channel=N
func (s *PTZService) GetCaps(ctx context.Context, channel int) (string, error) {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	return s.client.cgiGet(ctx, "ptz.cgi", "getCurrentProtocolCaps", params)
}

// MoveDirectly moves the camera using screen coordinates.
// CGI: ptzBase.cgi?action=moveDirectly&channel=N&startpoint[0]=X&startpoint[1]=Y&endpoint[0]=X&endpoint[1]=Y
func (s *PTZService) MoveDirectly(ctx context.Context, channel int, startX, startY, endX, endY int) error {
	// Build raw query string to preserve brackets unencoded.
	raw := fmt.Sprintf("action=moveDirectly&channel=%d&startpoint[0]=%d&startpoint[1]=%d&endpoint[0]=%d&endpoint[1]=%d",
		channel, startX, startY, endX, endY)
	u := s.client.baseURL + "/cgi-bin/ptzBase.cgi?" + raw
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return fmt.Errorf("amcrest: creating request: %w", err)
	}
	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("amcrest: executing request: %w", err)
	}
	return checkOK(resp)
}

// MoveRelatively moves the camera by relative amounts.
// CGI: ptz.cgi?action=moveRelatively&channel=N&...
func (s *PTZService) MoveRelatively(ctx context.Context, channel int, params map[string]string) error {
	qv := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	for k, v := range params {
		qv.Set(k, v)
	}
	return s.client.cgiAction(ctx, "ptz.cgi", "moveRelatively", qv)
}

// MoveAbsolutely moves the camera to an absolute position.
// CGI: ptz.cgi?action=moveAbsolutely&channel=N&...
func (s *PTZService) MoveAbsolutely(ctx context.Context, channel int, params map[string]string) error {
	qv := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	for k, v := range params {
		qv.Set(k, v)
	}
	return s.client.cgiAction(ctx, "ptz.cgi", "moveAbsolutely", qv)
}

// SetPresetName saves the current position as a named preset.
// CGI: ptz.cgi?action=SetPreset&channel=N&PresetName=X&index=N
func (s *PTZService) SetPresetName(ctx context.Context, channel, preset int, name string) error {
	params := url.Values{
		"channel":    {fmt.Sprintf("%d", channel)},
		"index":      {fmt.Sprintf("%d", preset)},
		"PresetName": {name},
	}
	return s.client.cgiAction(ctx, "ptz.cgi", "SetPreset", params)
}

// StartTour starts a PTZ tour.
func (s *PTZService) StartTour(ctx context.Context, channel, tourID int) error {
	return s.Control(ctx, channel, "StartTour", 0, tourID, 0)
}

// StopTour stops a PTZ tour.
func (s *PTZService) StopTour(ctx context.Context, channel int) error {
	return s.Control(ctx, channel, "StopTour", 0, 0, 0)
}

// AddTourPreset adds a preset to a tour.
func (s *PTZService) AddTourPreset(ctx context.Context, channel, tourID, preset int) error {
	return s.Control(ctx, channel, "AddTour", 0, tourID, preset)
}

// DeleteTourPreset removes a preset from a tour.
func (s *PTZService) DeleteTourPreset(ctx context.Context, channel, tourID, preset int) error {
	return s.Control(ctx, channel, "DelTour", 0, tourID, preset)
}

// ClearTour removes all presets from a tour.
func (s *PTZService) ClearTour(ctx context.Context, channel, tourID int) error {
	return s.Control(ctx, channel, "ClearTour", 0, tourID, 0)
}

// SetLeftLimit sets the left scan limit at the current position.
func (s *PTZService) SetLeftLimit(ctx context.Context, channel int) error {
	return s.Control(ctx, channel, "SetLeftLimit", 0, 0, 0)
}

// SetRightLimit sets the right scan limit at the current position.
func (s *PTZService) SetRightLimit(ctx context.Context, channel int) error {
	return s.Control(ctx, channel, "SetRightLimit", 0, 0, 0)
}

// AutoScanOn starts automatic scanning between left and right limits.
func (s *PTZService) AutoScanOn(ctx context.Context, channel int) error {
	return s.Control(ctx, channel, "AutoScanOn", 0, 0, 0)
}

// AutoScanOff stops automatic scanning.
func (s *PTZService) AutoScanOff(ctx context.Context, channel int) error {
	return s.Control(ctx, channel, "AutoScanOff", 0, 0, 0)
}

// SetPatternBegin starts recording a PTZ pattern.
func (s *PTZService) SetPatternBegin(ctx context.Context, channel int) error {
	return s.Control(ctx, channel, "SetPatternBegin", 0, 0, 0)
}

// SetPatternEnd stops recording a PTZ pattern.
func (s *PTZService) SetPatternEnd(ctx context.Context, channel int) error {
	return s.Control(ctx, channel, "SetPatternEnd", 0, 0, 0)
}

// StartPattern starts playing a recorded PTZ pattern.
func (s *PTZService) StartPattern(ctx context.Context, channel int) error {
	return s.Control(ctx, channel, "StartPattern", 0, 0, 0)
}

// StopPattern stops playing a PTZ pattern.
func (s *PTZService) StopPattern(ctx context.Context, channel int) error {
	return s.Control(ctx, channel, "StopPattern", 0, 0, 0)
}

// AutoPanOn starts automatic panning.
func (s *PTZService) AutoPanOn(ctx context.Context, channel int) error {
	return s.Control(ctx, channel, "AutoPanOn", 0, 0, 0)
}

// AutoPanOff stops automatic panning.
func (s *PTZService) AutoPanOff(ctx context.Context, channel int) error {
	return s.Control(ctx, channel, "AutoPanOff", 0, 0, 0)
}

// GetAutoMovementConfig returns the PtzAutoMovement configuration.
// CGI: configManager.cgi?action=getConfig&name=PtzAutoMovement
func (s *PTZService) GetAutoMovementConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "PtzAutoMovement")
}

// SetAutoMovementConfig sets PtzAutoMovement configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *PTZService) SetAutoMovementConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// Restart restarts the PTZ subsystem.
func (s *PTZService) Restart(ctx context.Context, channel int) error {
	return s.Control(ctx, channel, "Restart", 0, 0, 0)
}

// Reset resets the PTZ subsystem to factory defaults.
func (s *PTZService) Reset(ctx context.Context, channel int) error {
	return s.Control(ctx, channel, "Reset", 0, 0, 0)
}

// MenuOpen opens the PTZ OSD menu.
func (s *PTZService) MenuOpen(ctx context.Context, channel int) error {
	return s.Control(ctx, channel, "Menu", 0, 0, 0)
}

// MenuClose closes the PTZ OSD menu.
func (s *PTZService) MenuClose(ctx context.Context, channel int) error {
	return s.Control(ctx, channel, "Exit", 0, 0, 0)
}

// MenuEnter confirms the current selection in the PTZ OSD menu.
func (s *PTZService) MenuEnter(ctx context.Context, channel int) error {
	return s.Control(ctx, channel, "Enter", 0, 0, 0)
}

// MenuUp moves up in the PTZ OSD menu.
func (s *PTZService) MenuUp(ctx context.Context, channel int) error {
	return s.Control(ctx, channel, "MenuUp", 0, 0, 0)
}

// MenuDown moves down in the PTZ OSD menu.
func (s *PTZService) MenuDown(ctx context.Context, channel int) error {
	return s.Control(ctx, channel, "MenuDown", 0, 0, 0)
}

// MenuLeft moves left in the PTZ OSD menu.
func (s *PTZService) MenuLeft(ctx context.Context, channel int) error {
	return s.Control(ctx, channel, "MenuLeft", 0, 0, 0)
}

// MenuRight moves right in the PTZ OSD menu.
func (s *PTZService) MenuRight(ctx context.Context, channel int) error {
	return s.Control(ctx, channel, "MenuRight", 0, 0, 0)
}

// GetEPTZConfig returns the EptzLink (electronic PTZ) configuration.
// CGI: configManager.cgi?action=getConfig&name=EptzLink
func (s *PTZService) GetEPTZConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "EptzLink")
}

// SetEPTZConfig sets EptzLink (electronic PTZ) configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *PTZService) SetEPTZConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

