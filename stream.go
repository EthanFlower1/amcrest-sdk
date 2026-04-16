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
	var rawLines []string
	var jsonLines []string
	inJSON := false
	braceDepth := 0

	// Separate the header line (Code=...;action=...;index=...;data={)
	// from any subsequent JSON data lines.
	var headerLine string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "Content-") {
			continue
		}
		rawLines = append(rawLines, line)

		if headerLine == "" && strings.HasPrefix(trimmed, "Code=") {
			headerLine = trimmed
			// Check if this line opens a JSON block
			if strings.HasSuffix(trimmed, "{") || strings.Contains(trimmed, "data={") {
				inJSON = true
				braceDepth = 1
			}
			continue
		}

		if inJSON {
			jsonLines = append(jsonLines, line)
			braceDepth += strings.Count(trimmed, "{") - strings.Count(trimmed, "}")
			if braceDepth <= 0 {
				inJSON = false
			}
		}
	}

	if headerLine == "" {
		// Check for Heartbeat
		for _, line := range lines {
			if strings.TrimSpace(line) == "Heartbeat" {
				return &Event{Code: "Heartbeat", Raw: "Heartbeat"}
			}
		}
		return nil
	}

	// Parse the header: Code=XXX;action=YYY;index=ZZZ;data={
	evt := &Event{
		Data: make(map[string]string),
		Raw:  strings.Join(rawLines, "\n"),
	}

	// Strip "data={" or "data=" suffix from header before parsing fields
	header := headerLine
	dataIdx := strings.Index(header, ";data=")
	if dataIdx >= 0 {
		header = header[:dataIdx]
	}

	// Parse semicolon-separated key=value pairs
	parts := strings.Split(header, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		eqIdx := strings.Index(part, "=")
		if eqIdx < 0 {
			continue
		}
		key := strings.TrimSpace(part[:eqIdx])
		val := strings.TrimSpace(part[eqIdx+1:])
		switch key {
		case "Code":
			evt.Code = val
		case "action":
			evt.Action = val
		case "index":
			fmt.Sscanf(val, "%d", &evt.Index)
		}
		evt.Data[key] = val
	}

	// Attach the JSON data if present
	if len(jsonLines) > 0 {
		evt.Data["data"] = strings.Join(jsonLines, "\n")
	}

	if evt.Code == "" {
		return nil
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
