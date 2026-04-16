package amcrest

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
)

// RecordingService handles recording-related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 177-208 (Section 4.10)
type RecordingService struct {
	client *Client
}

// FindFilesOpts specifies search criteria for media file queries.
type FindFilesOpts struct {
	Channel   int
	StartTime string // "YYYY-MM-DD HH:MM:SS"
	EndTime   string // "YYYY-MM-DD HH:MM:SS"
	Type      string // "dav", "jpg", "mp4"
}

// MediaFile represents a single recorded media file returned by the camera.
type MediaFile struct {
	Channel   int
	StartTime string
	EndTime   string
	Type      string
	FilePath  string
	Length    int
	Duration  int
}

// GetCaps returns recording manager capabilities.
// CGI: recordManager.cgi?action=getCaps
func (s *RecordingService) GetCaps(ctx context.Context) (map[string]string, error) {
	body, err := s.client.cgiGet(ctx, "recordManager.cgi", "getCaps", nil)
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}

// GetRecordConfig returns the Record configuration table.
// CGI: configManager.cgi?action=getConfig&name=Record
func (s *RecordingService) GetRecordConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "Record")
}

// GetRecordMode returns the RecordMode configuration table.
// CGI: configManager.cgi?action=getConfig&name=RecordMode
func (s *RecordingService) GetRecordMode(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "RecordMode")
}

// GetMediaGlobal returns the MediaGlobal configuration table with the
// "table.MediaGlobal." prefix stripped from keys.
// CGI: configManager.cgi?action=getConfig&name=MediaGlobal
func (s *RecordingService) GetMediaGlobal(ctx context.Context) (map[string]string, error) {
	return s.client.getConfig(ctx, "MediaGlobal")
}

// FindFiles searches for recorded media files matching the given criteria.
// It abstracts the stateful mediaFileFind flow: create -> find -> nextFile loop -> close -> destroy.
// CGI: mediaFileFind.cgi
func (s *RecordingService) FindFiles(ctx context.Context, opts FindFilesOpts) ([]MediaFile, error) {
	// Step 1: factory.create - obtain an object ID
	body, err := s.client.cgiGet(ctx, "mediaFileFind.cgi", "factory.create", nil)
	if err != nil {
		return nil, fmt.Errorf("amcrest: mediaFileFind factory.create: %w", err)
	}
	kv := parseKV(body)
	objectID := kv["result"]
	if objectID == "" {
		return nil, fmt.Errorf("amcrest: mediaFileFind factory.create returned no object ID")
	}

	// Ensure cleanup runs even on error.
	defer func() {
		// Close the finder.
		_ = s.mediaFileFindRaw(ctx, objectID, "close", "")
		// Destroy the factory object.
		_ = s.mediaFileFindRaw(ctx, objectID, "destroy", "")
	}()

	// Step 2: findFile - set search conditions.
	findExtra := fmt.Sprintf(
		"condition.Channel=%d"+
			"&condition.StartTime=%s"+
			"&condition.EndTime=%s",
		opts.Channel,
		amcrestEscape(opts.StartTime),
		amcrestEscape(opts.EndTime),
	)
	if opts.Type != "" {
		findExtra += "&condition.Types[0]=" + amcrestEscape(opts.Type)
	}

	body, err = s.mediaFileFindRawBody(ctx, objectID, "findFile", findExtra)
	if err != nil {
		return nil, fmt.Errorf("amcrest: mediaFileFind findFile: %w", err)
	}
	// Camera returns "OK" on success, or "Error" on failure.
	// If we got here without error from readBody, the search is set up.

	// Step 3: findNextFile in a loop, fetching up to 100 at a time.
	var files []MediaFile
	for {
		body, err = s.mediaFileFindRawBody(ctx, objectID, "findNextFile", "count=100")
		if err != nil {
			return nil, fmt.Errorf("amcrest: mediaFileFind findNextFile: %w", err)
		}

		batch := parseMediaFiles(body)
		if len(batch) == 0 {
			break
		}
		files = append(files, batch...)

		// Check if the camera indicated no more results.
		kv := parseKV(body)
		if kv["found"] == "0" {
			break
		}
	}

	return files, nil
}

// mediaFileFindRaw performs a mediaFileFind action using a raw query string
// to preserve literal brackets in parameter names.
func (s *RecordingService) mediaFileFindRaw(ctx context.Context, objectID, action, extra string) error {
	_, err := s.mediaFileFindRawBody(ctx, objectID, action, extra)
	return err
}

// mediaFileFindRawBody performs a mediaFileFind action using a raw query string
// and returns the response body.
func (s *RecordingService) mediaFileFindRawBody(ctx context.Context, objectID, action, extra string) (string, error) {
	query := "action=" + url.QueryEscape(action) + "&object=" + url.QueryEscape(objectID)
	if extra != "" {
		query += "&" + extra
	}
	path := "/cgi-bin/mediaFileFind.cgi?" + query

	resp, err := s.client.get(ctx, path, nil)
	if err != nil {
		return "", err
	}
	return readBody(resp)
}

// amcrestEscape encodes a query value for Amcrest cameras. Unlike url.QueryEscape,
// it uses %20 for spaces (not +) and leaves colons unencoded, which is what the
// camera firmware expects.
func amcrestEscape(s string) string {
	return strings.ReplaceAll(s, " ", "%20")
}

// parseMediaFiles parses the multi-item response from findNextFile.
// Lines look like:
//
//	items[0].Channel=0
//	items[0].StartTime=2024-01-15 00:00:00
//	items[0].EndTime=2024-01-15 01:00:00
//	items[0].Type=dav
//	items[0].FilePath=/mnt/sd/2024-01-15/001/dav/00/00.00.00-01.00.00.dav
//	items[0].Length=12345
//	items[0].Duration=3600
func parseMediaFiles(body string) []MediaFile {
	// Group fields by index.
	groups := make(map[int]map[string]string)

	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "items[") {
			continue
		}
		// Parse "items[N].Field=Value"
		closeBracket := strings.Index(line, "]")
		if closeBracket < 0 {
			continue
		}
		idxStr := line[len("items["):closeBracket]
		idx, err := strconv.Atoi(idxStr)
		if err != nil {
			continue
		}

		rest := line[closeBracket+1:] // ".Field=Value"
		if len(rest) < 2 || rest[0] != '.' {
			continue
		}
		eqIdx := strings.Index(rest, "=")
		if eqIdx < 0 {
			continue
		}
		field := rest[1:eqIdx]
		value := strings.TrimSpace(rest[eqIdx+1:])

		if groups[idx] == nil {
			groups[idx] = make(map[string]string)
		}
		groups[idx][field] = value
	}

	files := make([]MediaFile, 0, len(groups))
	for i := 0; i < len(groups); i++ {
		g, ok := groups[i]
		if !ok {
			break
		}
		ch, _ := strconv.Atoi(g["Channel"])
		length, _ := strconv.Atoi(g["Length"])
		duration, _ := strconv.Atoi(g["Duration"])
		files = append(files, MediaFile{
			Channel:   ch,
			StartTime: g["StartTime"],
			EndTime:   g["EndTime"],
			Type:      g["Type"],
			FilePath:  g["FilePath"],
			Length:    length,
			Duration:  duration,
		})
	}

	return files
}

// SetRecordConfig sets Record configuration values. Keys should be prefixed
// with "Record." (e.g., "Record[0].TimeSection[0][0]").
func (s *RecordingService) SetRecordConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// SetRecordMode sets RecordMode configuration values. Keys should be prefixed
// with "RecordMode." (e.g., "RecordMode[0].Mode").
func (s *RecordingService) SetRecordMode(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// SetMediaGlobal sets MediaGlobal configuration values. Keys should be prefixed
// with "MediaGlobal." (e.g., "MediaGlobal.PacketLength").
func (s *RecordingService) SetMediaGlobal(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetRecordStateAll returns the recording status of all channels.
// POST: /cgi-bin/api/recordManager/getStateAll
func (s *RecordingService) GetRecordStateAll(ctx context.Context) (string, error) {
	body, err := s.client.postRaw(ctx, "/cgi-bin/api/recordManager/getStateAll", struct{}{})
	if err != nil {
		return "", fmt.Errorf("RecordingService.GetRecordStateAll: %w", err)
	}
	return body, nil
}

// FindFilesWithFilter searches for recorded media files with additional database filter conditions.
// The dbFilter map keys are condition.DB.* suffixes (e.g., "FaceDetectionRecordFilter.Sex" -> "Man").
func (s *RecordingService) FindFilesWithFilter(ctx context.Context, opts FindFilesOpts, dbFilter map[string]string) ([]MediaFile, error) {
	// Step 1: factory.create - obtain an object ID
	body, err := s.client.cgiGet(ctx, "mediaFileFind.cgi", "factory.create", nil)
	if err != nil {
		return nil, fmt.Errorf("amcrest: mediaFileFind factory.create: %w", err)
	}
	kv := parseKV(body)
	objectID := kv["result"]
	if objectID == "" {
		return nil, fmt.Errorf("amcrest: mediaFileFind factory.create returned no object ID")
	}

	// Ensure cleanup runs even on error.
	defer func() {
		_ = s.mediaFileFindRaw(ctx, objectID, "close", "")
		_ = s.mediaFileFindRaw(ctx, objectID, "destroy", "")
	}()

	// Step 2: findFile - set search conditions.
	findExtra := fmt.Sprintf(
		"condition.Channel=%d"+
			"&condition.StartTime=%s"+
			"&condition.EndTime=%s",
		opts.Channel,
		amcrestEscape(opts.StartTime),
		amcrestEscape(opts.EndTime),
	)
	if opts.Type != "" {
		findExtra += "&condition.Types[0]=" + amcrestEscape(opts.Type)
	}
	// Add DB filter conditions.
	for k, v := range dbFilter {
		findExtra += "&condition.DB." + amcrestEscape(k) + "=" + amcrestEscape(v)
	}

	body, err = s.mediaFileFindRawBody(ctx, objectID, "findFile", findExtra)
	if err != nil {
		return nil, fmt.Errorf("amcrest: mediaFileFind findFile: %w", err)
	}

	// Step 3: findNextFile in a loop, fetching up to 100 at a time.
	var files []MediaFile
	for {
		body, err = s.mediaFileFindRawBody(ctx, objectID, "findNextFile", "count=100")
		if err != nil {
			return nil, fmt.Errorf("amcrest: mediaFileFind findNextFile: %w", err)
		}

		batch := parseMediaFiles(body)
		if len(batch) == 0 {
			break
		}
		files = append(files, batch...)

		kv := parseKV(body)
		if kv["found"] == "0" {
			break
		}
	}

	return files, nil
}

// findFilesWithEvents searches for media files with event conditions and DB filter parameters.
// It extends the FindFiles flow by adding Flags, Events, and DB filter conditions.
func (s *RecordingService) findFilesWithEvents(ctx context.Context, opts FindFilesOpts, eventName, filterPrefix string, filter map[string]string) ([]MediaFile, error) {
	if opts.Type == "" {
		opts.Type = "jpg"
	}

	// Step 1: factory.create
	body, err := s.client.cgiGet(ctx, "mediaFileFind.cgi", "factory.create", nil)
	if err != nil {
		return nil, fmt.Errorf("amcrest: mediaFileFind factory.create: %w", err)
	}
	kv := parseKV(body)
	objectID := kv["result"]
	if objectID == "" {
		return nil, fmt.Errorf("amcrest: mediaFileFind factory.create returned no object ID")
	}

	defer func() {
		_ = s.mediaFileFindRaw(ctx, objectID, "close", "")
		_ = s.mediaFileFindRaw(ctx, objectID, "destroy", "")
	}()

	// Step 2: findFile with event conditions
	findExtra := fmt.Sprintf(
		"condition.Channel=%d"+
			"&condition.StartTime=%s"+
			"&condition.EndTime=%s",
		opts.Channel,
		amcrestEscape(opts.StartTime),
		amcrestEscape(opts.EndTime),
	)
	if opts.Type != "" {
		findExtra += "&condition.Types[0]=" + amcrestEscape(opts.Type)
	}
	findExtra += "&condition.Flags[0]=Event"
	findExtra += "&condition.Events[0]=" + amcrestEscape(eventName)

	// Add DB filter conditions with the appropriate prefix.
	for k, v := range filter {
		findExtra += "&condition.DB." + filterPrefix + amcrestEscape(k) + "=" + amcrestEscape(v)
	}

	body, err = s.mediaFileFindRawBody(ctx, objectID, "findFile", findExtra)
	if err != nil {
		return nil, fmt.Errorf("amcrest: mediaFileFind findFile: %w", err)
	}

	// Step 3: findNextFile loop
	var files []MediaFile
	for {
		body, err = s.mediaFileFindRawBody(ctx, objectID, "findNextFile", "count=100")
		if err != nil {
			return nil, fmt.Errorf("amcrest: mediaFileFind findNextFile: %w", err)
		}

		batch := parseMediaFiles(body)
		if len(batch) == 0 {
			break
		}
		files = append(files, batch...)

		kv := parseKV(body)
		if kv["found"] == "0" {
			break
		}
	}

	return files, nil
}

// FindFaceDetectionFiles searches for face detection event recordings.
// Filter keys are prefixed with "FaceDetectionRecordFilter." automatically.
func (s *RecordingService) FindFaceDetectionFiles(ctx context.Context, opts FindFilesOpts, filter map[string]string) ([]MediaFile, error) {
	return s.findFilesWithEvents(ctx, opts, "FaceDetection", "FaceDetectionRecordFilter.", filter)
}

// FindFaceRecognitionFiles searches for face recognition event recordings.
// Filter keys are prefixed with "FaceRecognitionRecordFilter." automatically.
func (s *RecordingService) FindFaceRecognitionFiles(ctx context.Context, opts FindFilesOpts, filter map[string]string) ([]MediaFile, error) {
	return s.findFilesWithEvents(ctx, opts, "FaceRecognition", "FaceRecognitionRecordFilter.", filter)
}

// FindHumanTraitFiles searches for human trait event recordings.
// Filter keys are prefixed with "HumanTraitRecordFilter." automatically.
func (s *RecordingService) FindHumanTraitFiles(ctx context.Context, opts FindFilesOpts, filter map[string]string) ([]MediaFile, error) {
	return s.findFilesWithEvents(ctx, opts, "HumanTrait", "HumanTraitRecordFilter.", filter)
}

// FindTrafficCarFiles searches for traffic car event recordings.
// Filter keys are prefixed with "TrafficCar." automatically.
func (s *RecordingService) FindTrafficCarFiles(ctx context.Context, opts FindFilesOpts, filter map[string]string) ([]MediaFile, error) {
	return s.findFilesWithEvents(ctx, opts, "TrafficCar", "TrafficCar.", filter)
}

// FindIVSFiles searches for IVS (intelligent video surveillance) event recordings.
// Filter keys are prefixed with "IVS." automatically.
func (s *RecordingService) FindIVSFiles(ctx context.Context, opts FindFilesOpts, filter map[string]string) ([]MediaFile, error) {
	return s.findFilesWithEvents(ctx, opts, "IVS", "IVS.", filter)
}

// FindNonMotorFiles searches for non-motor vehicle event recordings.
// Filter keys are prefixed with "NonMotorRecordFilter." automatically.
func (s *RecordingService) FindNonMotorFiles(ctx context.Context, opts FindFilesOpts, filter map[string]string) ([]MediaFile, error) {
	return s.findFilesWithEvents(ctx, opts, "NonMotor", "NonMotorRecordFilter.", filter)
}

// DownloadByTime downloads a recording by time range and returns the raw bytes.
// CGI: loadfile.cgi?action=startLoad&channel=N&startTime=X&endTime=Y&subtype=Z
func (s *RecordingService) DownloadByTime(ctx context.Context, channel int, startTime, endTime string, subtype int) ([]byte, error) {
	params := url.Values{
		"channel":   {strconv.Itoa(channel)},
		"startTime": {startTime},
		"endTime":   {endTime},
		"subtype":   {strconv.Itoa(subtype)},
	}

	path := "/cgi-bin/loadfile.cgi"
	resp, err := s.client.get(ctx, path, url.Values{
		"action":    {"startLoad"},
		"channel":   params["channel"],
		"startTime": params["startTime"],
		"endTime":   params["endTime"],
		"subtype":   params["subtype"],
	})
	if err != nil {
		return nil, fmt.Errorf("amcrest: downloading by time: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    "failed to download by time",
		}
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("amcrest: reading download body: %w", err)
	}
	return data, nil
}

// DownloadEncrypted downloads an encrypted recording file and returns the raw bytes.
// CGI: RecordStreamInterleaved.cgi?action=attachStream&path=X&password=Y
func (s *RecordingService) DownloadEncrypted(ctx context.Context, filePath, password string) ([]byte, error) {
	params := url.Values{
		"action":   {"attachStream"},
		"path":     {filePath},
		"password": {password},
	}

	resp, err := s.client.get(ctx, "/cgi-bin/RecordStreamInterleaved.cgi", params)
	if err != nil {
		return nil, fmt.Errorf("amcrest: downloading encrypted file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("failed to download encrypted %s", filePath),
		}
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("amcrest: reading encrypted file body: %w", err)
	}
	return data, nil
}

// GetAlarmCount returns the count of alarm recordings.
// POST: /cgi-bin/api/mediaFileFind/getCount
func (s *RecordingService) GetAlarmCount(ctx context.Context) (int, error) {
	var result struct {
		Count int `json:"count"`
	}
	err := s.client.postJSON(ctx, "/cgi-bin/api/mediaFileFind/getCount", struct{}{}, &result)
	if err != nil {
		return 0, err
	}
	return result.Count, nil
}

// DownloadFile downloads a recorded file from the camera and returns the raw bytes.
// The filePath should be the path returned by FindFiles (e.g., "/mnt/sd/...").
// CGI: GET /cgi-bin/RPC_Loadfile/<filePath>
func (s *RecordingService) DownloadFile(ctx context.Context, filePath string) ([]byte, error) {
	// Strip leading slash from filePath to avoid double slash.
	filePath = strings.TrimPrefix(filePath, "/")
	path := "/cgi-bin/RPC_Loadfile/" + filePath

	resp, err := s.client.get(ctx, path, nil)
	if err != nil {
		return nil, fmt.Errorf("amcrest: downloading file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("failed to download %s", filePath),
		}
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("amcrest: reading file body: %w", err)
	}
	return data, nil
}
