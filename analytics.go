package amcrest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

// SetGlobalConfig updates the global video analytics configuration.
// PDF 9.6.2
func (s *AnalyticsService) SetGlobalConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// SetRuleConfig updates the video analytics rule configuration.
// PDF 9.6.3
func (s *AnalyticsService) SetRuleConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetLastEventInfo returns the last video analytics event info for the given channel.
// PDF 9.6.4: devVideoAnalyse.cgi?action=getLastEventInfo&channel=N
func (s *AnalyticsService) GetLastEventInfo(ctx context.Context, channel int) (string, error) {
	body, err := s.client.cgiGet(ctx, "devVideoAnalyse.cgi", "getLastEventInfo", url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	})
	if err != nil {
		return "", err
	}
	return body, nil
}

// GetGlobalDeviceParam returns the global device parameter configuration.
// PDF 9.6.5: getRawConfig GlobalDeviceParam
func (s *AnalyticsService) GetGlobalDeviceParam(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "GlobalDeviceParam")
}

// SetGlobalDeviceParam updates the global device parameter configuration.
// PDF 9.6.5
func (s *AnalyticsService) SetGlobalDeviceParam(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetTemplateRule returns the template rule for video analytics on the given channel.
// PDF 9.6.6: VideoInAnalyse.cgi?action=getTemplateRule&channel=N
func (s *AnalyticsService) GetTemplateRule(ctx context.Context, channel int) (string, error) {
	body, err := s.client.cgiGet(ctx, "VideoInAnalyse.cgi", "getTemplateRule", url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	})
	if err != nil {
		return "", err
	}
	return body, nil
}

// GetIntelliTourConfig returns the intelligent scheme tour configuration.
// PDF 9.6.7-8: getRawConfig IntelliSchemeTour
func (s *AnalyticsService) GetIntelliTourConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "IntelliSchemeTour")
}

// SetIntelliTourConfig updates the intelligent scheme tour configuration.
// PDF 9.6.7-8
func (s *AnalyticsService) SetIntelliTourConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetIntelliCaps returns the intelligence capabilities for the given channel.
// PDF 9.6.11: POST /cgi-bin/api/intelli/getCaps or cgiGet intelli.cgi
func (s *AnalyticsService) GetIntelliCaps(ctx context.Context, channel int) (string, error) {
	reqBody := map[string]int{"channel": channel}
	body, err := s.client.postRaw(ctx, "/cgi-bin/api/intelli/getCaps", reqBody)
	if err != nil {
		return "", fmt.Errorf("AnalyticsService.GetIntelliCaps: %w", err)
	}
	return body, nil
}

// EnableScene enables video analytics scene types for the given channel.
// PDF 9.6.16: devVideoAnalyse.cgi?action=enableScene&typeList=X&channel=N
func (s *AnalyticsService) EnableScene(ctx context.Context, channel int, types string) error {
	params := url.Values{
		"typeList": {types},
		"channel":  {fmt.Sprintf("%d", channel)},
	}
	return s.client.cgiAction(ctx, "devVideoAnalyse.cgi", "enableScene", params)
}

// DisableScene disables video analytics scene types for the given channel.
// PDF 9.6.16: devVideoAnalyse.cgi?action=disableScene&typeList=X&channel=N
func (s *AnalyticsService) DisableScene(ctx context.Context, channel int, types string) error {
	params := url.Values{
		"typeList": {types},
		"channel":  {fmt.Sprintf("%d", channel)},
	}
	return s.client.cgiAction(ctx, "devVideoAnalyse.cgi", "disableScene", params)
}

// ExportData exports security data of the given type, identified by key.
// PDF 9.6.9/13: POST /cgi-bin/api/SecurityImExport/exportData with JSON.
func (s *AnalyticsService) ExportData(ctx context.Context, key string, dataType int) ([]byte, error) {
	reqBody := map[string]interface{}{
		"key":      key,
		"dataType": dataType,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("AnalyticsService.ExportData: marshaling JSON: %w", err)
	}

	u := s.client.baseURL + "/cgi-bin/api/SecurityImExport/exportData"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("AnalyticsService.ExportData: creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("AnalyticsService.ExportData: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, &APIError{StatusCode: resp.StatusCode}
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("AnalyticsService.ExportData: reading body: %w", err)
	}
	return data, nil
}

// ImportData imports security data of the given type.
// PDF 9.6.10: POST /cgi-bin/api/SecurityImExport/importData
func (s *AnalyticsService) ImportData(ctx context.Context, dataType int, data []byte) error {
	reqBody := map[string]interface{}{
		"dataType": dataType,
		"data":     data,
	}
	return s.client.postJSON(ctx, "/cgi-bin/api/SecurityImExport/importData", reqBody, nil)
}

// SubscribeResourceUsage subscribes to intelligence resource usage events.
// PDF 9.6.12: cgiGet intelli.cgi attachResource
func (s *AnalyticsService) SubscribeResourceUsage(ctx context.Context) (string, error) {
	body, err := s.client.cgiGet(ctx, "intelli.cgi", "attachResource", nil)
	if err != nil {
		return "", fmt.Errorf("AnalyticsService.SubscribeResourceUsage: %w", err)
	}
	return body, nil
}

// SetPollingConfig sets the platform intelligent control polling configuration.
// PDF 9.6.14: POST /cgi-bin/api/intelli/setPollingConfig
func (s *AnalyticsService) SetPollingConfig(ctx context.Context, body interface{}) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/intelli/setPollingConfig", body, nil)
}
