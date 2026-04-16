package amcrest

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// SnapshotService handles snapshot-related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 66-70 (Section 4.4)
type SnapshotService struct {
	client *Client
}

// Get captures a snapshot from the given channel and returns the raw JPEG bytes.
// Channel is 1-based (e.g., 1 for the first video channel).
func (s *SnapshotService) Get(ctx context.Context, channel int) ([]byte, error) {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	resp, err := s.client.get(ctx, "/cgi-bin/snapshot.cgi", params)
	if err != nil {
		return nil, fmt.Errorf("amcrest: snapshot get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, &APIError{StatusCode: resp.StatusCode}
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("amcrest: reading snapshot body: %w", err)
	}

	return data, nil
}

// GetSnapConfig retrieves the snapshot configuration table.
func (s *SnapshotService) GetSnapConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getConfig(ctx, "Snap")
}

// SetSnapConfig sets snapshot configuration values. Keys should be
// prefixed with "Snap." (e.g., "Snap.Channel[0].SnapMode").
func (s *SnapshotService) SetSnapConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetWithType captures a snapshot from the given channel with a specific type.
// Type 0 is a real-time snapshot from the frontend encoder; type 1 is from
// the local decode channel. Returns the raw JPEG bytes.
func (s *SnapshotService) GetWithType(ctx context.Context, channel, snapType int) ([]byte, error) {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
		"type":    {fmt.Sprintf("%d", snapType)},
	}
	resp, err := s.client.get(ctx, "/cgi-bin/snapshot.cgi", params)
	if err != nil {
		return nil, fmt.Errorf("amcrest: snapshot getWithType: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, &APIError{StatusCode: resp.StatusCode}
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("amcrest: reading snapshot body: %w", err)
	}
	return data, nil
}

// SnapshotEvent represents a single event from the snapshot attach stream.
// Code and Action describe the event type; ImageData contains raw JPEG bytes
// if the event includes an image (may be nil for text-only events).
type SnapshotEvent struct {
	Code      string
	Action    string
	ImageData []byte
	Raw       string
}

// SnapshotStream reads snapshot events from a long-lived multipart stream.
type SnapshotStream struct {
	resp    *http.Response
	scanner *bufio.Scanner
	err     error
}

// Next returns the next snapshot event from the stream. Returns false when the
// stream ends or an error occurs.
func (ss *SnapshotStream) Next() (*SnapshotEvent, bool) {
	var block []string
	inBlock := false

	for ss.scanner.Scan() {
		line := ss.scanner.Text()

		if strings.HasPrefix(line, "--") && !inBlock {
			inBlock = true
			block = nil
			continue
		}

		if inBlock {
			if strings.HasPrefix(line, "--") {
				evt := parseSnapshotEvent(block)
				if evt != nil {
					block = nil
					inBlock = true
					return evt, true
				}
				block = nil
				continue
			}
			block = append(block, line)
		}
	}

	if err := ss.scanner.Err(); err != nil {
		ss.err = err
	}

	if len(block) > 0 {
		evt := parseSnapshotEvent(block)
		if evt != nil {
			return evt, true
		}
	}

	return nil, false
}

// Err returns any error encountered during stream reading.
func (ss *SnapshotStream) Err() error {
	return ss.err
}

// Close closes the underlying HTTP response body.
func (ss *SnapshotStream) Close() error {
	if ss.resp != nil && ss.resp.Body != nil {
		return ss.resp.Body.Close()
	}
	return nil
}

func parseSnapshotEvent(lines []string) *SnapshotEvent {
	kv := make(map[string]string)
	var rawLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Content-") {
			continue
		}
		rawLines = append(rawLines, line)
		idx := strings.Index(line, "=")
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		kv[key] = val
	}

	code := kv["Code"]
	if code == "" && len(rawLines) == 0 {
		return nil
	}

	return &SnapshotEvent{
		Code:   code,
		Action: kv["action"],
		Raw:    strings.Join(rawLines, "\n"),
	}
}

// Subscribe opens a long-lived multipart stream for snapshot events on the
// given channel. The heartbeat parameter specifies the keep-alive interval in
// seconds. Events is a list of event codes to subscribe to.
// CGI: snapManager.cgi?action=attachFileProc&channel=N&heartbeat=H&Ede[0].Code=X&...
func (s *SnapshotService) Subscribe(ctx context.Context, channel int, heartbeat int, events []string) (<-chan SnapshotEvent, *SnapshotStream, error) {
	params := url.Values{
		"action":    {"attachFileProc"},
		"channel":   {fmt.Sprintf("%d", channel)},
		"heartbeat": {fmt.Sprintf("%d", heartbeat)},
	}
	for i, code := range events {
		params.Set(fmt.Sprintf("Ede[%d].Code", i), code)
	}

	resp, err := s.client.get(ctx, "/cgi-bin/snapManager.cgi", params)
	if err != nil {
		return nil, nil, fmt.Errorf("amcrest: snapshot subscribe: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    "failed to subscribe to snapshot events",
		}
	}

	stream := &SnapshotStream{
		resp:    resp,
		scanner: bufio.NewScanner(resp.Body),
	}

	ch := make(chan SnapshotEvent)
	go func() {
		defer close(ch)
		for {
			evt, ok := stream.Next()
			if !ok {
				return
			}
			select {
			case ch <- *evt:
			case <-ctx.Done():
				return
			}
		}
	}()

	return ch, stream, nil
}
