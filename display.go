package amcrest

import (
	"context"
	"fmt"
	"net/url"
)

// DisplayService handles display-related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 279-283 (Section 7)
type DisplayService struct {
	client *Client
}

// GetGUIConfig returns the GUISet configuration as a raw key-value map.
// CGI: configManager.cgi?action=getConfig&name=GUISet
func (s *DisplayService) GetGUIConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "GUISet")
}

// GetSplitMode returns the current split mode for the given channel.
// CGI: split.cgi?action=getMode&channel=N
func (s *DisplayService) GetSplitMode(ctx context.Context, channel int) (map[string]string, error) {
	body, err := s.client.cgiGet(ctx, "split.cgi", "getMode", url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	})
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}

// GetMonitorTour returns the MonitorTour configuration as a raw key-value map.
// CGI: configManager.cgi?action=getConfig&name=MonitorTour
func (s *DisplayService) GetMonitorTour(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "MonitorTour")
}

// SetGUIConfig updates GUISet configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *DisplayService) SetGUIConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// SetSplitMode sets the split mode for the given channel.
// CGI: split.cgi?action=setMode&channel=N&mode=X&group=G
func (s *DisplayService) SetSplitMode(ctx context.Context, channel int, mode string, group int) error {
	return s.client.cgiAction(ctx, "split.cgi", "setMode", url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
		"mode":    {mode},
		"group":   {fmt.Sprintf("%d", group)},
	})
}

// SetMonitorTour updates MonitorTour configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *DisplayService) SetMonitorTour(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// EnableTour enables or disables the monitor tour for the given channel.
// CGI: split.cgi?action=enableTour&channel=N&enable=true/false
func (s *DisplayService) EnableTour(ctx context.Context, channel int, enable bool) error {
	return s.client.cgiAction(ctx, "split.cgi", "enableTour", url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
		"enable":  {fmt.Sprintf("%t", enable)},
	})
}

// GetMonitorCollection returns the MonitorCollection configuration as a raw key-value map.
// CGI: configManager.cgi?action=getConfig&name=MonitorCollection
func (s *DisplayService) GetMonitorCollection(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "MonitorCollection")
}

// SetMonitorCollection updates MonitorCollection configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *DisplayService) SetMonitorCollection(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}
