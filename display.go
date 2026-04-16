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
