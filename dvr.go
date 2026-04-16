package amcrest

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
)

// DVRService handles DVR-related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 609-617 (Section 14)
type DVRService struct {
	client *Client
}

// StartFind begins a file search and returns the find ID.
// FileFindHelper.cgi?action=startFind&channel=N&startTime=X&endTime=Y
func (s *DVRService) StartFind(ctx context.Context, channel int, startTime, endTime string) (string, error) {
	params := url.Values{
		"channel":   {fmt.Sprintf("%d", channel)},
		"startTime": {startTime},
		"endTime":   {endTime},
	}
	body, err := s.client.cgiGet(ctx, "FileFindHelper.cgi", "startFind", params)
	if err != nil {
		return "", err
	}
	// Response format: findId=N
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "findId=") {
			return strings.TrimPrefix(line, "findId="), nil
		}
	}
	return "", fmt.Errorf("amcrest: findId not found in response: %s", strings.TrimSpace(body))
}

// FindNext retrieves the next batch of results for a file search.
// FileFindHelper.cgi?action=findNext&findId=X&count=N
func (s *DVRService) FindNext(ctx context.Context, findId string, count int) (string, error) {
	params := url.Values{
		"findId": {findId},
		"count":  {fmt.Sprintf("%d", count)},
	}
	return s.client.cgiGet(ctx, "FileFindHelper.cgi", "findNext", params)
}

// StopFind stops a file search.
// FileFindHelper.cgi?action=stopFind&findId=X
func (s *DVRService) StopFind(ctx context.Context, findId string) error {
	params := url.Values{
		"findId": {findId},
	}
	return s.client.cgiAction(ctx, "FileFindHelper.cgi", "stopFind", params)
}

// GetBandwidthLimit retrieves the bandwidth limit state.
// BandLimit.cgi?action=getLimitState
func (s *DVRService) GetBandwidthLimit(ctx context.Context) (map[string]string, error) {
	body, err := s.client.cgiGet(ctx, "BandLimit.cgi", "getLimitState", nil)
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}

// StartMotionFind begins a motion-based file search with optional extra parameters.
// FileFindHelper.cgi?action=startFind&channel=N&startTime=X&endTime=Y&...
func (s *DVRService) StartMotionFind(ctx context.Context, channel int, startTime, endTime string, params map[string]string) (string, error) {
	v := url.Values{
		"channel":   {fmt.Sprintf("%d", channel)},
		"startTime": {startTime},
		"endTime":   {endTime},
	}
	for k, val := range params {
		v.Set(k, val)
	}
	return s.client.cgiGet(ctx, "FileFindHelper.cgi", "startFind", v)
}

// GetBoundFiles retrieves bound files for the given channel and time range.
// FileFindHelper.cgi?action=getBoundFile&channel=N&startTime=X&endTime=Y
func (s *DVRService) GetBoundFiles(ctx context.Context, channel int, startTime, endTime string) (string, error) {
	params := url.Values{
		"channel":   {fmt.Sprintf("%d", channel)},
		"startTime": {startTime},
		"endTime":   {endTime},
	}
	return s.client.cgiGet(ctx, "FileFindHelper.cgi", "getBoundFile", params)
}

// AddProtection adds a protection condition to files.
// FileManager.cgi?action=addConditionList&...
func (s *DVRService) AddProtection(ctx context.Context, params map[string]string) error {
	v := url.Values{}
	for k, val := range params {
		v.Set(k, val)
	}
	return s.client.cgiAction(ctx, "FileManager.cgi", "addConditionList", v)
}

// CancelProtection cancels a protection condition on files.
// FileManager.cgi?action=cancelConditionList&...
func (s *DVRService) CancelProtection(ctx context.Context, params map[string]string) error {
	v := url.Values{}
	for k, val := range params {
		v.Set(k, val)
	}
	return s.client.cgiAction(ctx, "FileManager.cgi", "cancelConditionList", v)
}

// RemoveProtection removes a protection condition from files.
// FileManager.cgi?action=removeConditionList&...
func (s *DVRService) RemoveProtection(ctx context.Context, params map[string]string) error {
	v := url.Values{}
	for k, val := range params {
		v.Set(k, val)
	}
	return s.client.cgiAction(ctx, "FileManager.cgi", "removeConditionList", v)
}

// DownloadFile downloads a file from the DVR by name.
// GET /cgi-bin/FileManager.cgi?action=downloadFile&fileName=X
func (s *DVRService) DownloadFile(ctx context.Context, fileName string) ([]byte, error) {
	params := url.Values{
		"action":   {"downloadFile"},
		"fileName": {fileName},
	}
	resp, err := s.client.get(ctx, "/cgi-bin/FileManager.cgi", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("amcrest: download failed with status %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

// ListDirectory lists the contents of a directory on the DVR.
// POST /cgi-bin/api/FileManager/list
func (s *DVRService) ListDirectory(ctx context.Context, path string) (string, error) {
	body := map[string]interface{}{
		"path": path,
	}
	return s.client.postRaw(ctx, "/cgi-bin/api/FileManager/list", body)
}

// GetDaylight retrieves the daylight saving time (DST) configuration.
// global.cgi?action=getDST
func (s *DVRService) GetDaylight(ctx context.Context) (string, error) {
	return s.client.cgiGet(ctx, "global.cgi", "getDST", nil)
}
