package amcrest

import (
	"context"
	"fmt"
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
	PTZLeftDown  = "LeftDown"
	PTZRightDown = "RightDown"
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
