package amcrest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
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

// FindTrafficFlowHistory searches for traffic flow history records with the
// given condition parameters.
// CGI: recordFinder.cgi?action=find&name=TrafficFlow&...
func (s *TrafficService) FindTrafficFlowHistory(ctx context.Context, params map[string]string) (string, error) {
	qv := url.Values{
		"name": {"TrafficFlow"},
	}
	for k, v := range params {
		qv.Set(k, v)
	}
	return s.client.cgiGet(ctx, "recordFinder.cgi", "find", qv)
}

// DownloadExportFile downloads an exported file of the given type and returns
// the raw binary data.
// CGI: trafficRecord.cgi?action=downloadFile&Type=X
func (s *TrafficService) DownloadExportFile(ctx context.Context, fileType string) ([]byte, error) {
	params := url.Values{
		"action": {"downloadFile"},
		"Type":   {fileType},
	}
	resp, err := s.client.get(ctx, "/cgi-bin/trafficRecord.cgi", params)
	if err != nil {
		return nil, fmt.Errorf("amcrest: downloading export file: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("failed to download export file type %s", fileType),
		}
	}
	return io.ReadAll(resp.Body)
}

// ExportTrafficFlow starts an asynchronous export of traffic flow data.
// CGI: trafficRecord.cgi?action=exportAsyncFile&name=TrafficFlow
func (s *TrafficService) ExportTrafficFlow(ctx context.Context) error {
	return s.client.cgiAction(ctx, "trafficRecord.cgi", "exportAsyncFile", url.Values{
		"name": {"TrafficFlow"},
	})
}

// GetTrafficFlowExportState returns the export state for traffic flow data.
// CGI: trafficRecord.cgi?action=getFileExportState&name=TrafficFlow
func (s *TrafficService) GetTrafficFlowExportState(ctx context.Context) (map[string]string, error) {
	body, err := s.client.cgiGet(ctx, "trafficRecord.cgi", "getFileExportState", url.Values{
		"name": {"TrafficFlow"},
	})
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}

// DownloadTrafficFlow downloads the exported traffic flow file.
// CGI: trafficRecord.cgi?action=downloadFile&Type=TrafficFlow
func (s *TrafficService) DownloadTrafficFlow(ctx context.Context) ([]byte, error) {
	return s.DownloadExportFile(ctx, "TrafficFlow")
}

// ExportSnapEventInfo starts an asynchronous export of traffic snap event info
// with the given condition parameters.
// CGI: trafficRecord.cgi?action=exportAsyncFileByConditon&name=TrafficSnapEventInfo&...
func (s *TrafficService) ExportSnapEventInfo(ctx context.Context, params map[string]string) error {
	qv := url.Values{
		"name": {"TrafficSnapEventInfo"},
	}
	for k, v := range params {
		qv.Set(k, v)
	}
	return s.client.cgiAction(ctx, "trafficRecord.cgi", "exportAsyncFileByConditon", qv)
}

// GetSnapEventExportState returns the export state for traffic snap event info.
// CGI: trafficRecord.cgi?action=getFileExportState&name=TrafficSnapEventInfo
func (s *TrafficService) GetSnapEventExportState(ctx context.Context) (map[string]string, error) {
	body, err := s.client.cgiGet(ctx, "trafficRecord.cgi", "getFileExportState", url.Values{
		"name": {"TrafficSnapEventInfo"},
	})
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}

// DownloadSnapEventInfo downloads the exported traffic snap event info file.
// CGI: trafficRecord.cgi?action=downloadFile&Type=TrafficSnapEventInfo
func (s *TrafficService) DownloadSnapEventInfo(ctx context.Context) ([]byte, error) {
	return s.DownloadExportFile(ctx, "TrafficSnapEventInfo")
}

// OpenUnlicensedDetection opens unlicensed vehicle detection on the given channel.
// CGI: trafficSnap.cgi?action=open&name=UnlicensedVehicle&channel=N
func (s *TrafficService) OpenUnlicensedDetection(ctx context.Context, channel int) error {
	return s.client.cgiAction(ctx, "trafficSnap.cgi", "open", url.Values{
		"name":    {"UnlicensedVehicle"},
		"channel": {strconv.Itoa(channel)},
	})
}

// CloseUnlicensedDetection closes unlicensed vehicle detection on the given channel.
// CGI: trafficSnap.cgi?action=close&name=UnlicensedVehicle&channel=N
func (s *TrafficService) CloseUnlicensedDetection(ctx context.Context, channel int) error {
	return s.client.cgiAction(ctx, "trafficSnap.cgi", "close", url.Values{
		"name":    {"UnlicensedVehicle"},
		"channel": {strconv.Itoa(channel)},
	})
}

// SubscribeVehiclesDistribution subscribes to vehicles distribution data on the
// given channel. This is a long-lived stream; the raw response body is returned.
// CGI: vehiclesDistribution.cgi?action=attach&channel=N
func (s *TrafficService) SubscribeVehiclesDistribution(ctx context.Context, channel int) (string, error) {
	return s.client.cgiGet(ctx, "vehiclesDistribution.cgi", "attach", url.Values{
		"channel": {strconv.Itoa(channel)},
	})
}
