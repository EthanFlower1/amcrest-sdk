package amcrest

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
)

// LogService handles log-related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 208-212 (Section 4.11)
type LogService struct {
	client *Client
}

// LogEntry represents a single log record returned by the device.
type LogEntry struct {
	Time   string
	Type   string
	User   string
	Detail string
}

// Find searches the device log for entries between startTime and endTime.
// Times should be in "YYYY-MM-DD HH:MM:SS" format. logType is optional; pass
// an empty string to retrieve all types.
//
// The method abstracts the three-step CGI search flow:
//  1. startFind – returns a search token and total count
//  2. doFind   – fetches entries in batches of 100
//  3. stopFind – releases the search token
func (s *LogService) Find(ctx context.Context, startTime, endTime string, logType string) ([]LogEntry, error) {
	// Step 1: startFind
	params := url.Values{
		"condition.StartTime": {startTime},
		"condition.EndTime":   {endTime},
	}
	if logType != "" {
		params.Set("condition.Type", logType)
	}

	body, err := s.client.cgiGet(ctx, "log.cgi", "startFind", params)
	if err != nil {
		return nil, fmt.Errorf("log startFind: %w", err)
	}

	kv := parseKV(body)
	token := kv["token"]
	if token == "" {
		return nil, fmt.Errorf("log startFind: missing token in response: %s", strings.TrimSpace(body))
	}

	totalCount, _ := strconv.Atoi(kv["count"])

	// Ensure we always call stopFind when done.
	defer func() {
		_ = s.client.cgiAction(ctx, "log.cgi", "stopFind", url.Values{
			"token": {token},
		})
	}()

	if totalCount == 0 {
		return nil, nil
	}

	// Step 2: doFind – fetch in batches of 100
	var entries []LogEntry
	fetched := 0
	for fetched < totalCount {
		batch := 100
		if totalCount-fetched < batch {
			batch = totalCount - fetched
		}

		doBody, err := s.client.cgiGet(ctx, "log.cgi", "doFind", url.Values{
			"token": {token},
			"count": {strconv.Itoa(batch)},
		})
		if err != nil {
			return entries, fmt.Errorf("log doFind: %w", err)
		}

		parsed := parseLogEntries(doBody)
		if len(parsed) == 0 {
			break // no more results
		}
		entries = append(entries, parsed...)
		fetched += len(parsed)
	}

	return entries, nil
}

// parseLogEntries parses the items[N].Key=Value format returned by doFind.
func parseLogEntries(body string) []LogEntry {
	kv := parseKV(body)

	// Determine the maximum index present.
	maxIdx := -1
	for k := range kv {
		if !strings.HasPrefix(k, "items[") {
			continue
		}
		end := strings.Index(k, "]")
		if end < 0 {
			continue
		}
		idx, err := strconv.Atoi(k[len("items["):end])
		if err != nil {
			continue
		}
		if idx > maxIdx {
			maxIdx = idx
		}
	}

	if maxIdx < 0 {
		return nil
	}

	entries := make([]LogEntry, 0, maxIdx+1)
	for i := 0; i <= maxIdx; i++ {
		prefix := fmt.Sprintf("items[%d].", i)
		entries = append(entries, LogEntry{
			Time:   kv[prefix+"Time"],
			Type:   kv[prefix+"Type"],
			User:   kv[prefix+"User"],
			Detail: kv[prefix+"Detail"],
		})
	}
	return entries
}

// Clear deletes all log entries on the device.
// CGI: log.cgi?action=clear
func (s *LogService) Clear(ctx context.Context) error {
	return s.client.cgiAction(ctx, "log.cgi", "clear", nil)
}

// SeekFind retrieves log entries by seeking to a specific offset within a
// search token obtained from Find's startFind step. This is useful for
// paginated access to log results.
// CGI: Log.cgi?action=doSeekFind&token=T&offset=O&count=C
func (s *LogService) SeekFind(ctx context.Context, token, offset, count int) ([]LogEntry, error) {
	body, err := s.client.cgiGet(ctx, "Log.cgi", "doSeekFind", url.Values{
		"token":  {strconv.Itoa(token)},
		"offset": {strconv.Itoa(offset)},
		"count":  {strconv.Itoa(count)},
	})
	if err != nil {
		return nil, fmt.Errorf("log doSeekFind: %w", err)
	}
	return parseLogEntries(body), nil
}

// ExportEncrypted exports an encrypted log archive (ZIP) for the given time
// range, protected by the supplied password. Times should be in
// "YYYY-MM-DD HH:MM:SS" format.
// CGI: GET /cgi-bin/Log.exportEncrypedLog?action=All&condition.StartTime=X&condition.EndTime=Y&password=P
func (s *LogService) ExportEncrypted(ctx context.Context, startTime, endTime, password string) ([]byte, error) {
	params := url.Values{
		"action":              {"All"},
		"condition.StartTime": {startTime},
		"condition.EndTime":   {endTime},
		"password":            {password},
	}

	resp, err := s.client.get(ctx, "/cgi-bin/Log.exportEncrypedLog", params)
	if err != nil {
		return nil, fmt.Errorf("log exportEncrypted: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, &APIError{StatusCode: resp.StatusCode}
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("log exportEncrypted: reading body: %w", err)
	}
	return data, nil
}

// GetDebugInfoRedirConfig returns the DebugInfoRedir configuration table with
// the "table.DebugInfoRedir." prefix stripped from keys.
func (s *LogService) GetDebugInfoRedirConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getConfig(ctx, "DebugInfoRedir")
}

// SetDebugInfoRedirConfig sets DebugInfoRedir configuration values. Keys should
// be prefixed with "DebugInfoRedir." (e.g., "DebugInfoRedir.Enable").
func (s *LogService) SetDebugInfoRedirConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// Backup downloads a binary log backup covering the given time range.
// Times should be in "YYYY-MM-DD HH:MM:SS" format.
// CGI: GET /cgi-bin/Log.backup?action=All&condition.StartTime=X&condition.EndTime=Y
func (s *LogService) Backup(ctx context.Context, startTime, endTime string) ([]byte, error) {
	params := url.Values{
		"action":              {"All"},
		"condition.StartTime": {startTime},
		"condition.EndTime":   {endTime},
	}

	resp, err := s.client.get(ctx, "/cgi-bin/Log.backup", params)
	if err != nil {
		return nil, fmt.Errorf("log backup: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, &APIError{StatusCode: resp.StatusCode}
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("log backup: reading body: %w", err)
	}
	return data, nil
}
