package amcrest

import (
	"context"
	"fmt"
	"net/url"
)

// CameraService handles camera-related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 235-258 (Section 5.4)
type CameraService struct {
	client *Client
}

// GetImageConfig returns the VideoColor configuration as a raw key-value map.
// CGI: configManager.cgi?action=getConfig&name=VideoColor
func (s *CameraService) GetImageConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "VideoColor")
}

// SetImageConfig sets VideoColor configuration parameters.
// Each key in params should be a full VideoColor config key
// (e.g., "VideoColor[0].Brightness" = "50").
// CGI: configManager.cgi?action=setConfig&VideoColor[0].Brightness=50
func (s *CameraService) SetImageConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetExposureConfig returns the VideoInExposure configuration as a raw key-value map.
// CGI: configManager.cgi?action=getConfig&name=VideoInExposure
func (s *CameraService) GetExposureConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "VideoInExposure")
}

// GetBacklightConfig returns the VideoInBacklight configuration as a raw key-value map.
// CGI: configManager.cgi?action=getConfig&name=VideoInBacklight
func (s *CameraService) GetBacklightConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "VideoInBacklight")
}

// GetWhiteBalanceConfig returns the VideoInWhiteBalance configuration as a raw key-value map.
// CGI: configManager.cgi?action=getConfig&name=VideoInWhiteBalance
func (s *CameraService) GetWhiteBalanceConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "VideoInWhiteBalance")
}

// GetDayNightConfig returns the VideoInDayNight configuration as a raw key-value map.
// CGI: configManager.cgi?action=getConfig&name=VideoInDayNight
func (s *CameraService) GetDayNightConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "VideoInDayNight")
}

// AutoFocus triggers an auto-focus operation on the specified channel.
// CGI: devVideoInput.cgi?action=autoFocus&channel=N
func (s *CameraService) AutoFocus(ctx context.Context, channel int) error {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	return s.client.cgiAction(ctx, "devVideoInput.cgi", "autoFocus", params)
}

// GetFocusStatus returns the focus status for the specified channel.
// CGI: devVideoInput.cgi?action=getFocusStatus&channel=N
func (s *CameraService) GetFocusStatus(ctx context.Context, channel int) (map[string]string, error) {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	body, err := s.client.cgiGet(ctx, "devVideoInput.cgi", "getFocusStatus", params)
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}

// GetLightingConfig returns the Lighting configuration as a raw key-value map.
// CGI: configManager.cgi?action=getConfig&name=Lighting
func (s *CameraService) GetLightingConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "Lighting")
}

// GetVideoInOptions returns the VideoInOptions configuration as a raw key-value map.
// CGI: configManager.cgi?action=getConfig&name=VideoInOptions
func (s *CameraService) GetVideoInOptions(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "VideoInOptions")
}

// GetSharpnessConfig returns the VideoInSharpness configuration as a raw key-value map.
// CGI: configManager.cgi?action=getConfig&name=VideoInSharpness
func (s *CameraService) GetSharpnessConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "VideoInSharpness")
}

// SetSharpnessConfig updates VideoInSharpness configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *CameraService) SetSharpnessConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetFlipConfig returns the VideoImageControl configuration as a raw key-value map.
// CGI: configManager.cgi?action=getConfig&name=VideoImageControl
func (s *CameraService) GetFlipConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "VideoImageControl")
}

// SetFlipConfig updates VideoImageControl configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *CameraService) SetFlipConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// SetExposureConfig updates VideoInExposure configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *CameraService) SetExposureConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// SetBacklightConfig updates VideoInBacklight configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *CameraService) SetBacklightConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// SetWhiteBalanceConfig updates VideoInWhiteBalance configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *CameraService) SetWhiteBalanceConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// SetDayNightConfig updates VideoInDayNight configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *CameraService) SetDayNightConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// AdjustFocus adjusts focus and zoom for the specified channel.
// CGI: devVideoInput.cgi?action=adjustFocus&channel=N&focus=F&zoom=Z
func (s *CameraService) AdjustFocus(ctx context.Context, channel int, focus, zoom int) error {
	return s.client.cgiAction(ctx, "devVideoInput.cgi", "adjustFocus", url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
		"focus":   {fmt.Sprintf("%d", focus)},
		"zoom":    {fmt.Sprintf("%d", zoom)},
	})
}

// AdjustFocusContinuously adjusts focus and zoom continuously for the specified channel.
// CGI: devVideoInput.cgi?action=adjustFocusContinuously&channel=N&focus=F&zoom=Z
func (s *CameraService) AdjustFocusContinuously(ctx context.Context, channel int, focus, zoom int) error {
	return s.client.cgiAction(ctx, "devVideoInput.cgi", "adjustFocusContinuously", url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
		"focus":   {fmt.Sprintf("%d", focus)},
		"zoom":    {fmt.Sprintf("%d", zoom)},
	})
}

// GetZoomConfig returns the VideoInZoom configuration as a raw key-value map.
// CGI: configManager.cgi?action=getConfig&name=VideoInZoom
func (s *CameraService) GetZoomConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "VideoInZoom")
}

// SetZoomConfig updates VideoInZoom configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *CameraService) SetZoomConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetFocusConfig returns the VideoInFocus configuration as a raw key-value map.
// CGI: configManager.cgi?action=getConfig&name=VideoInFocus
func (s *CameraService) GetFocusConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "VideoInFocus")
}

// SetFocusConfig updates VideoInFocus configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *CameraService) SetFocusConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// SetLightingConfig updates Lighting configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *CameraService) SetLightingConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// SetVideoInOptions updates VideoInOptions configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *CameraService) SetVideoInOptions(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}
