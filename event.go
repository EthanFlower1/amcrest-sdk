package amcrest

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// EventService handles event-related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 157-176 (Section 4.9)
type EventService struct {
	client *Client
}

// Subscribe opens a long-lived event stream for the given event codes.
// The heartbeat parameter controls how often (in seconds) the camera sends
// heartbeat events on the stream. Pass 0 to use the camera default.
// Returns a channel that delivers parsed events and the underlying EventStream
// which the caller must Close() when done. Cancel the context to stop streaming.
//
// CGI: eventManager.cgi?action=attach&codes=[Code1,Code2]&heartbeat=N
func (s *EventService) Subscribe(ctx context.Context, codes []string, heartbeat int) (<-chan Event, *EventStream, error) {
	params := url.Values{
		"action": {"attach"},
		"codes":  {fmt.Sprintf("[%s]", strings.Join(codes, ","))},
	}
	if heartbeat > 0 {
		params.Set("heartbeat", strconv.Itoa(heartbeat))
	}

	resp, err := s.client.get(ctx, "/cgi-bin/eventManager.cgi", params)
	if err != nil {
		return nil, nil, fmt.Errorf("amcrest: subscribing to events: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		resp.Body.Close()
		return nil, nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    "failed to subscribe to events",
		}
	}

	es := &EventStream{
		resp:    resp,
		scanner: bufio.NewScanner(resp.Body),
	}

	ch := make(chan Event, 16)
	go func() {
		defer close(ch)
		var block []string
		inBlock := false
		scanner := es.scanner

		for scanner.Scan() {
			line := scanner.Text()

			if strings.HasPrefix(line, "--") && !inBlock {
				inBlock = true
				block = nil
				continue
			}

			if inBlock {
				if strings.HasPrefix(line, "--") {
					evt := parseEventBlock(block)
					if evt != nil {
						select {
						case ch <- *evt:
						case <-ctx.Done():
							return
						}
					}
					block = nil
					continue
				}
				block = append(block, line)
			}
		}

		// Try to parse any remaining block
		if len(block) > 0 {
			evt := parseEventBlock(block)
			if evt != nil {
				select {
				case ch <- *evt:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return ch, es, nil
}

// parseEventBlock parses a multipart block into an Event, including heartbeats.
func parseEventBlock(lines []string) *Event {
	// First try the standard event parser
	evt := parseEvent(lines)
	if evt != nil {
		return evt
	}

	// Check for heartbeat: the block contains just "Heartbeat" (plus Content- headers)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Content-") {
			continue
		}
		if line == "Heartbeat" {
			return &Event{
				Code: "Heartbeat",
				Data: map[string]string{},
				Raw:  "Heartbeat",
			}
		}
	}

	return nil
}

// GetEventIndexes returns the event indexes for a given event code.
// CGI: eventManager.cgi?action=getEventIndexes&code=X
func (s *EventService) GetEventIndexes(ctx context.Context, code string) (string, error) {
	body, err := s.client.cgiGet(ctx, "eventManager.cgi", "getEventIndexes", url.Values{
		"code": {code},
	})
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(body), nil
}

// GetCaps returns the event manager capabilities as key-value pairs.
// CGI: eventManager.cgi?action=getCaps
func (s *EventService) GetCaps(ctx context.Context) (map[string]string, error) {
	body, err := s.client.cgiGet(ctx, "eventManager.cgi", "getCaps", nil)
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}

// GetSupportedEvents returns a list of event types the camera supports.
// CGI: eventManager.cgi?action=getExposureEvents
// Response format: events[0]=VideoMotion, events[1]=... etc.
func (s *EventService) GetSupportedEvents(ctx context.Context) ([]string, error) {
	body, err := s.client.cgiGet(ctx, "eventManager.cgi", "getExposureEvents", nil)
	if err != nil {
		return nil, err
	}

	kv := parseKV(body)
	var events []string
	for i := 0; ; i++ {
		key := fmt.Sprintf("events[%d]", i)
		val, ok := kv[key]
		if !ok {
			break
		}
		events = append(events, val)
	}
	return events, nil
}

// GetAlarmConfig returns the Alarm configuration table.
// CGI: configManager.cgi?action=getConfig&name=Alarm
func (s *EventService) GetAlarmConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getConfig(ctx, "Alarm")
}

// GetAlarmOutConfig returns the AlarmOut configuration table.
// CGI: configManager.cgi?action=getConfig&name=AlarmOut
func (s *EventService) GetAlarmOutConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getConfig(ctx, "AlarmOut")
}

// GetAlarmInputChannels returns the number of alarm input slots.
// CGI: alarm.cgi?action=getInSlots
func (s *EventService) GetAlarmInputChannels(ctx context.Context) (int, error) {
	body, err := s.client.cgiGet(ctx, "alarm.cgi", "getInSlots", nil)
	if err != nil {
		return 0, err
	}
	kv := parseKV(body)
	val, ok := kv["result"]
	if !ok {
		return 0, fmt.Errorf("amcrest: getInSlots: missing result in response: %s", body)
	}
	n, err := strconv.Atoi(strings.TrimSpace(val))
	if err != nil {
		return 0, fmt.Errorf("amcrest: getInSlots: parsing result %q: %w", val, err)
	}
	return n, nil
}

// GetAlarmOutputChannels returns the number of alarm output slots.
// CGI: alarm.cgi?action=getOutSlots
func (s *EventService) GetAlarmOutputChannels(ctx context.Context) (int, error) {
	body, err := s.client.cgiGet(ctx, "alarm.cgi", "getOutSlots", nil)
	if err != nil {
		return 0, err
	}
	kv := parseKV(body)
	val, ok := kv["result"]
	if !ok {
		return 0, fmt.Errorf("amcrest: getOutSlots: missing result in response: %s", body)
	}
	n, err := strconv.Atoi(strings.TrimSpace(val))
	if err != nil {
		return 0, fmt.Errorf("amcrest: getOutSlots: parsing result %q: %w", val, err)
	}
	return n, nil
}

// GetBlindDetectConfig returns the BlindDetect (video tampering) configuration.
// CGI: configManager.cgi?action=getConfig&name=BlindDetect
func (s *EventService) GetBlindDetectConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getConfig(ctx, "BlindDetect")
}

// GetLossDetectConfig returns the LossDetect (video loss) configuration.
// CGI: configManager.cgi?action=getConfig&name=LossDetect
func (s *EventService) GetLossDetectConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getConfig(ctx, "LossDetect")
}

// GetEventHandlerConfig returns the event handler configuration for the given name.
// The name parameter is e.g. "Alarm[0].EventHandler".
// CGI: configManager.cgi?action=getConfig&name=<name>
func (s *EventService) GetEventHandlerConfig(ctx context.Context, name string) (map[string]string, error) {
	return s.client.getRawConfig(ctx, name)
}

// SetEventHandlerConfig sets event handler configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *EventService) SetEventHandlerConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// SetAlarmConfig sets Alarm configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *EventService) SetAlarmConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// SetAlarmOutConfig sets AlarmOut configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *EventService) SetAlarmOutConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetAlarmInputStates returns the current alarm input states as a bitmask.
// CGI: alarm.cgi?action=getInState
func (s *EventService) GetAlarmInputStates(ctx context.Context) (int, error) {
	body, err := s.client.cgiGet(ctx, "alarm.cgi", "getInState", nil)
	if err != nil {
		return 0, err
	}
	kv := parseKV(body)
	val, ok := kv["result"]
	if !ok {
		return 0, fmt.Errorf("amcrest: getInState: missing result in response: %s", body)
	}
	n, err := strconv.Atoi(strings.TrimSpace(val))
	if err != nil {
		return 0, fmt.Errorf("amcrest: getInState: parsing result %q: %w", val, err)
	}
	return n, nil
}

// GetAlarmOutputStates returns the current alarm output states as a bitmask.
// CGI: alarm.cgi?action=getOutState
func (s *EventService) GetAlarmOutputStates(ctx context.Context) (int, error) {
	body, err := s.client.cgiGet(ctx, "alarm.cgi", "getOutState", nil)
	if err != nil {
		return 0, err
	}
	kv := parseKV(body)
	val, ok := kv["result"]
	if !ok {
		return 0, fmt.Errorf("amcrest: getOutState: missing result in response: %s", body)
	}
	n, err := strconv.Atoi(strings.TrimSpace(val))
	if err != nil {
		return 0, fmt.Errorf("amcrest: getOutState: parsing result %q: %w", val, err)
	}
	return n, nil
}

// SetBlindDetectConfig sets BlindDetect (video tampering) configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *EventService) SetBlindDetectConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// SetLossDetectConfig sets LossDetect (video loss) configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *EventService) SetLossDetectConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetLoginFailureAlarmConfig returns the LoginFailureAlarm configuration.
// CGI: configManager.cgi?action=getConfig&name=LoginFailureAlarm
func (s *EventService) GetLoginFailureAlarmConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "LoginFailureAlarm")
}

// SetLoginFailureAlarmConfig sets LoginFailureAlarm configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *EventService) SetLoginFailureAlarmConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetStorageNotExistConfig returns the StorageNotExist configuration.
// CGI: configManager.cgi?action=getConfig&name=StorageNotExist
func (s *EventService) GetStorageNotExistConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "StorageNotExist")
}

// SetStorageNotExistConfig sets StorageNotExist configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *EventService) SetStorageNotExistConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetStorageFailureConfig returns the StorageFailure configuration.
// CGI: configManager.cgi?action=getConfig&name=StorageFailure
func (s *EventService) GetStorageFailureConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "StorageFailure")
}

// SetStorageFailureConfig sets StorageFailure configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *EventService) SetStorageFailureConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetStorageLowSpaceConfig returns the StorageLowSpace configuration.
// CGI: configManager.cgi?action=getConfig&name=StorageLowSpace
func (s *EventService) GetStorageLowSpaceConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "StorageLowSpace")
}

// SetStorageLowSpaceConfig sets StorageLowSpace configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *EventService) SetStorageLowSpaceConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetNetAbortConfig returns the NetAbort configuration.
// CGI: configManager.cgi?action=getConfig&name=NetAbort
func (s *EventService) GetNetAbortConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "NetAbort")
}

// SetNetAbortConfig sets NetAbort configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *EventService) SetNetAbortConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetIPConflictConfig returns the IPConflict configuration.
// CGI: configManager.cgi?action=getConfig&name=IPConflict
func (s *EventService) GetIPConflictConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "IPConflict")
}

// SetIPConflictConfig sets IPConflict configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *EventService) SetIPConflictConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetNetAlarmConfig returns the NetAlarm configuration.
// CGI: configManager.cgi?action=getConfig&name=NetAlarm
func (s *EventService) GetNetAlarmConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "NetAlarm")
}

// SetNetAlarmConfig sets NetAlarm configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *EventService) SetNetAlarmConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// SetNetAlarmState sets the net alarm state for a given channel.
// CGI: netAlarm.cgi?action=setState&channel=N&alarm=true/false
func (s *EventService) SetNetAlarmState(ctx context.Context, channel int, alarm bool) error {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
		"alarm":   {fmt.Sprintf("%t", alarm)},
	}
	return s.client.cgiAction(ctx, "netAlarm.cgi", "setState", params)
}
