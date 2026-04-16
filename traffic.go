package amcrest

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// TrafficService handles traffic-related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 440-466 (Section 10)
type TrafficService struct {
	client *Client
}

// InsertRecord inserts a record into a traffic list (e.g., TrafficBlackList or
// TrafficRedList) via recordUpdater.cgi. The params map may include keys such
// as PlateNumber, MasterOfCar, etc.
func (s *TrafficService) InsertRecord(ctx context.Context, name string, params map[string]string) error {
	qv := url.Values{}
	qv.Set("name", name)
	for k, v := range params {
		qv.Set(k, v)
	}
	return s.client.cgiAction(ctx, "recordUpdater.cgi", "insert", qv)
}

// RemoveRecord removes a record from a traffic list by record number.
func (s *TrafficService) RemoveRecord(ctx context.Context, name string, recno int) error {
	return s.client.cgiAction(ctx, "recordUpdater.cgi", "remove", url.Values{
		"name":  {name},
		"recno": {fmt.Sprintf("%d", recno)},
	})
}

// FindRecord searches for records in the given traffic list and returns the
// raw response body.
func (s *TrafficService) FindRecord(ctx context.Context, name string) (string, error) {
	return s.client.cgiGet(ctx, "recordFinder.cgi", "find", url.Values{
		"name": {name},
	})
}

// UpdateRecord updates an existing record in a traffic list by record number.
// CGI: recordUpdater.cgi?action=update&name=X&recno=N&...
func (s *TrafficService) UpdateRecord(ctx context.Context, name string, recno int, params map[string]string) error {
	qv := url.Values{}
	qv.Set("name", name)
	qv.Set("recno", fmt.Sprintf("%d", recno))
	for k, v := range params {
		qv.Set(k, v)
	}
	return s.client.cgiAction(ctx, "recordUpdater.cgi", "update", qv)
}

// RemoveRecordEx removes records from a traffic list using extended parameters.
// CGI: recordUpdater.cgi?action=removeEx&name=X&...
func (s *TrafficService) RemoveRecordEx(ctx context.Context, name string, params map[string]string) error {
	qv := url.Values{}
	qv.Set("name", name)
	for k, v := range params {
		qv.Set(k, v)
	}
	return s.client.cgiAction(ctx, "recordUpdater.cgi", "removeEx", qv)
}

// StartFlowSearch starts a traffic flow statistics search.
// API: POST /api/trafficFlowStat/startFind
func (s *TrafficService) StartFlowSearch(ctx context.Context, body interface{}) (string, error) {
	return s.client.postRaw(ctx, "/api/trafficFlowStat/startFind", body)
}

// DoFlowSearch performs a traffic flow statistics search query.
// API: POST /api/trafficFlowStat/doFind
func (s *TrafficService) DoFlowSearch(ctx context.Context, body interface{}) (string, error) {
	return s.client.postRaw(ctx, "/api/trafficFlowStat/doFind", body)
}

// StopFlowSearch stops a traffic flow statistics search by token.
// API: POST /api/trafficFlowStat/stopFind
func (s *TrafficService) StopFlowSearch(ctx context.Context, token int) error {
	body := map[string]interface{}{
		"token": token,
	}
	_, err := s.client.postRaw(ctx, "/api/trafficFlowStat/stopFind", body)
	if err != nil {
		return fmt.Errorf("traffic StopFlowSearch: %w", err)
	}
	return nil
}

// ImportRecords uploads a binary file of records into the given traffic list.
// CGI: POST /cgi-bin/trafficRecord.cgi?action=uploadFile&name=X
func (s *TrafficService) ImportRecords(ctx context.Context, name string, data []byte) error {
	u := s.client.baseURL + "/cgi-bin/trafficRecord.cgi?" + url.Values{
		"action": {"uploadFile"},
		"name":   {name},
	}.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("amcrest: creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("amcrest: executing request: %w", err)
	}
	return checkOK(resp)
}

// ExportRecords exports records from the given traffic list asynchronously and
// returns the raw response body.
// CGI: trafficRecord.cgi?action=exportAsyncFile&name=X
func (s *TrafficService) ExportRecords(ctx context.Context, name string) (string, error) {
	return s.client.cgiGet(ctx, "trafficRecord.cgi", "exportAsyncFile", url.Values{
		"name": {name},
	})
}

// OpenStrobe opens the traffic strobe light on the specified channel.
// CGI: trafficSnap.cgi?action=openStrobe&channel=N
func (s *TrafficService) OpenStrobe(ctx context.Context, channel int) error {
	return s.client.cgiAction(ctx, "trafficSnap.cgi", "openStrobe", url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	})
}

// ManualSnap triggers a manual traffic snapshot on the specified channel.
// CGI: trafficSnap.cgi?action=manSnap&channel=N
func (s *TrafficService) ManualSnap(ctx context.Context, channel int) error {
	return s.client.cgiAction(ctx, "trafficSnap.cgi", "manSnap", url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	})
}

// StartRecordFind starts a paginated record search for the given list name.
// Returns the raw response containing the search token.
// CGI: recordFinder.cgi?action=startFind&name=X
func (s *TrafficService) StartRecordFind(ctx context.Context, name string, params map[string]string) (string, error) {
	qv := url.Values{
		"name": {name},
	}
	for k, v := range params {
		qv.Set(k, v)
	}
	return s.client.cgiGet(ctx, "recordFinder.cgi", "startFind", qv)
}

// DoRecordFind retrieves a page of record search results using the given token.
// CGI: recordFinder.cgi?action=doFind&token=T&count=N
func (s *TrafficService) DoRecordFind(ctx context.Context, token string, count int) (string, error) {
	return s.client.cgiGet(ctx, "recordFinder.cgi", "doFind", url.Values{
		"token": {token},
		"count": {fmt.Sprintf("%d", count)},
	})
}

// StopRecordFind stops a paginated record search and releases the token.
// CGI: recordFinder.cgi?action=stopFind&token=T
func (s *TrafficService) StopRecordFind(ctx context.Context, token string) error {
	return s.client.cgiAction(ctx, "recordFinder.cgi", "stopFind", url.Values{
		"token": {token},
	})
}

// GetRecordCount returns the total number of records for the given list name.
// CGI: recordFinder.cgi?action=getQuerySize&name=X
func (s *TrafficService) GetRecordCount(ctx context.Context, name string) (int, error) {
	body, err := s.client.cgiGet(ctx, "recordFinder.cgi", "getQuerySize", url.Values{
		"name": {name},
	})
	if err != nil {
		return 0, fmt.Errorf("traffic GetRecordCount: %w", err)
	}
	kv := parseKV(body)
	count := 0
	if v, ok := kv["count"]; ok {
		fmt.Sscanf(v, "%d", &count)
	}
	return count, nil
}

// CloseStrobe closes the traffic strobe light on the specified channel.
// CGI: trafficSnap.cgi?action=closeStrobe&channel=N
func (s *TrafficService) CloseStrobe(ctx context.Context, channel int) error {
	return s.client.cgiAction(ctx, "trafficSnap.cgi", "closeStrobe", url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	})
}

// GetExportPercent returns the current file export progress.
// CGI: trafficRecord.cgi?action=getFileExportState
func (s *TrafficService) GetExportPercent(ctx context.Context) (map[string]string, error) {
	body, err := s.client.cgiGet(ctx, "trafficRecord.cgi", "getFileExportState", nil)
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}

// GetTrafficSnapConfig retrieves the TrafficSnap configuration table without
// stripping key prefixes.
// CGI: configManager.cgi?action=getConfig&name=TrafficSnap
func (s *TrafficService) GetTrafficSnapConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "TrafficSnap")
}

// SetTrafficSnapConfig updates TrafficSnap configuration values. Keys should be
// prefixed with "TrafficSnap." (e.g., "TrafficSnap.SnapMode").
// CGI: configManager.cgi?action=setConfig
func (s *TrafficService) SetTrafficSnapConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}
