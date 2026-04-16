package amcrest

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Event represents a single event from the Amcrest event stream.
type Event struct {
	Code   string
	Action string
	Index  int
	Data   map[string]string
	Raw    string
}

// EventStream reads events from a multipart event stream.
type EventStream struct {
	resp    *http.Response
	scanner *bufio.Scanner
	err     error
}

// Next returns the next event from the stream. It returns false when the
// stream is exhausted or an error occurs. Call Err() to check for errors.
func (es *EventStream) Next() (*Event, bool) {
	var block []string
	inBlock := false

	for es.scanner.Scan() {
		line := es.scanner.Text()

		// Detect boundary lines
		if strings.HasPrefix(line, "--") && !inBlock {
			inBlock = true
			block = nil
			continue
		}

		if inBlock {
			// End of block on empty line after content or next boundary
			if strings.HasPrefix(line, "--") {
				evt := parseEvent(block)
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

	if err := es.scanner.Err(); err != nil {
		es.err = err
	}

	// Try to parse any remaining block
	if len(block) > 0 {
		evt := parseEvent(block)
		if evt != nil {
			return evt, true
		}
	}

	return nil, false
}

// Err returns any error encountered during stream reading.
func (es *EventStream) Err() error {
	return es.err
}

// Close closes the underlying HTTP response body.
func (es *EventStream) Close() error {
	if es.resp != nil && es.resp.Body != nil {
		return es.resp.Body.Close()
	}
	return nil
}

func parseEvent(lines []string) *Event {
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
	if code == "" {
		return nil
	}

	evt := &Event{
		Code:   code,
		Action: kv["action"],
		Data:   kv,
		Raw:    strings.Join(rawLines, "\n"),
	}

	if idxStr, ok := kv["index"]; ok {
		fmt.Sscanf(idxStr, "%d", &evt.Index)
	}

	return evt
}

// subscribe opens a long-lived event stream connection for the given event codes.
func (c *Client) subscribe(ctx context.Context, codes []string) (*EventStream, error) {
	params := url.Values{
		"action": {"attach"},
	}
	for _, code := range codes {
		params.Add("codes", fmt.Sprintf("[%s]", code))
	}

	path := "/cgi-bin/eventManager.cgi"
	resp, err := c.get(ctx, path, params)
	if err != nil {
		return nil, fmt.Errorf("amcrest: subscribing to events: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    "failed to subscribe to events",
		}
	}

	return &EventStream{
		resp:    resp,
		scanner: bufio.NewScanner(resp.Body),
	}, nil
}
