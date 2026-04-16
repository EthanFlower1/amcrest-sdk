package amcrest

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// AccessControlService handles access control related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 505-584 (Section 12)
type AccessControlService struct {
	client *Client
}

// OpenDoor opens the door on the given channel.
// accessControl.cgi?action=openDoor&channel=N&UserID=0&Type=Remote
func (s *AccessControlService) OpenDoor(ctx context.Context, channel int) error {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
		"UserID":  {"0"},
		"Type":    {"Remote"},
	}
	return s.client.cgiAction(ctx, "accessControl.cgi", "openDoor", params)
}

// CloseDoor closes the door on the given channel.
// accessControl.cgi?action=closeDoor&channel=N
func (s *AccessControlService) CloseDoor(ctx context.Context, channel int) error {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	return s.client.cgiAction(ctx, "accessControl.cgi", "closeDoor", params)
}

// GetDoorStatus returns the door status for the given channel.
// accessControl.cgi?action=getDoorStatus&channel=N
func (s *AccessControlService) GetDoorStatus(ctx context.Context, channel int) (string, error) {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	body, err := s.client.cgiGet(ctx, "accessControl.cgi", "getDoorStatus", params)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(body), nil
}

// GetGeneralConfig retrieves the AccessControlGeneral configuration.
func (s *AccessControlService) GetGeneralConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "AccessControlGeneral")
}

// GetControlConfig retrieves the AccessControl configuration.
func (s *AccessControlService) GetControlConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "AccessControl")
}
