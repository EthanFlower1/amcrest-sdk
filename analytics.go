package amcrest

import (
	"context"
	"fmt"
	"net/url"
)

// AnalyticsService handles video analytics related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 381-401 (Section 9.6)
type AnalyticsService struct {
	client *Client
}

// GetCaps returns the video analytics capabilities for the given channel.
// The response is returned as a raw string because the structure is complex
// and varies by device model.
// CGI: devVideoAnalyse.cgi?action=getCaps&channel=N
func (s *AnalyticsService) GetCaps(ctx context.Context, channel int) (string, error) {
	body, err := s.client.cgiGet(ctx, "devVideoAnalyse.cgi", "getCaps", url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	})
	if err != nil {
		return "", err
	}
	return body, nil
}

// GetGlobalConfig returns the global video analytics configuration.
// Uses configManager getRawConfig with name VideoAnalyseGlobal.
func (s *AnalyticsService) GetGlobalConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "VideoAnalyseGlobal")
}

// GetRuleConfig returns the video analytics rule configuration.
// Uses configManager getRawConfig with name VideoAnalyseRule.
func (s *AnalyticsService) GetRuleConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "VideoAnalyseRule")
}

// GetSceneList returns the video analytics scene list for the given channel.
// The response is returned as a raw string because the structure is complex
// and varies by device model.
// CGI: devVideoAnalyse.cgi?action=getSceneList&channel=N
func (s *AnalyticsService) GetSceneList(ctx context.Context, channel int) (string, error) {
	body, err := s.client.cgiGet(ctx, "devVideoAnalyse.cgi", "getSceneList", url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	})
	if err != nil {
		return "", err
	}
	return body, nil
}
