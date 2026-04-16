package amcrest

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// BuildingService handles building intercom related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 585-608 (Section 13)
type BuildingService struct {
	client *Client
}

// GetSIPConfig retrieves the SIP configuration.
func (s *BuildingService) GetSIPConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "SIPConfig")
}

// GetRoomNumberCount returns the number of VideoTalkContact records.
// recordFinder.cgi?action=getQuerySize&name=VideoTalkContact
func (s *BuildingService) GetRoomNumberCount(ctx context.Context) (int, error) {
	params := url.Values{
		"name": {"VideoTalkContact"},
	}
	body, err := s.client.cgiGet(ctx, "recordFinder.cgi", "getQuerySize", params)
	if err != nil {
		return 0, err
	}
	// Response format: count=N
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "count=") {
			val := strings.TrimPrefix(line, "count=")
			n, err := strconv.Atoi(strings.TrimSpace(val))
			if err != nil {
				return 0, fmt.Errorf("amcrest: parsing count: %w", err)
			}
			return n, nil
		}
	}
	return 0, fmt.Errorf("amcrest: count not found in response: %s", strings.TrimSpace(body))
}

// SetSIPConfig updates the SIP configuration.
func (s *BuildingService) SetSIPConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetRegistarConfig retrieves the Registar configuration.
func (s *BuildingService) GetRegistarConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "Registar")
}

// SetRegistarConfig updates the Registar configuration.
func (s *BuildingService) SetRegistarConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// AddRoomNumber inserts a new VideoTalkContact record.
// recordUpdater.cgi?action=insert&name=VideoTalkContact
func (s *BuildingService) AddRoomNumber(ctx context.Context, params map[string]string) error {
	vals := url.Values{
		"name": {"VideoTalkContact"},
	}
	for k, v := range params {
		vals.Set(k, v)
	}
	return s.client.cgiAction(ctx, "recordUpdater.cgi", "insert", vals)
}

// FindRoomNumbers searches for VideoTalkContact records matching the given parameters.
// recordFinder.cgi?action=find&name=VideoTalkContact
func (s *BuildingService) FindRoomNumbers(ctx context.Context, params map[string]string) (string, error) {
	vals := url.Values{
		"name": {"VideoTalkContact"},
	}
	for k, v := range params {
		vals.Set(k, v)
	}
	return s.client.cgiGet(ctx, "recordFinder.cgi", "find", vals)
}

// GetRoomNumber retrieves a specific VideoTalkContact record by record number.
// recordUpdater.cgi?action=get&name=VideoTalkContact&recno=N
func (s *BuildingService) GetRoomNumber(ctx context.Context, recno int) (string, error) {
	params := url.Values{
		"name":  {"VideoTalkContact"},
		"recno": {strconv.Itoa(recno)},
	}
	return s.client.cgiGet(ctx, "recordUpdater.cgi", "get", params)
}

// UpdateRoomNumber updates a specific VideoTalkContact record by record number.
// recordUpdater.cgi?action=update&name=VideoTalkContact&recno=N
func (s *BuildingService) UpdateRoomNumber(ctx context.Context, recno int, params map[string]string) error {
	vals := url.Values{
		"name":  {"VideoTalkContact"},
		"recno": {strconv.Itoa(recno)},
	}
	for k, v := range params {
		vals.Set(k, v)
	}
	return s.client.cgiAction(ctx, "recordUpdater.cgi", "update", vals)
}

// DeleteRoomNumber removes a specific VideoTalkContact record by record number.
// recordUpdater.cgi?action=remove&name=VideoTalkContact&recno=N
func (s *BuildingService) DeleteRoomNumber(ctx context.Context, recno int) error {
	params := url.Values{
		"name":  {"VideoTalkContact"},
		"recno": {strconv.Itoa(recno)},
	}
	return s.client.cgiAction(ctx, "recordUpdater.cgi", "remove", params)
}

// ClearRoomNumbers removes all VideoTalkContact records.
// recordUpdater.cgi?action=clear&name=VideoTalkContact
func (s *BuildingService) ClearRoomNumbers(ctx context.Context) error {
	params := url.Values{
		"name": {"VideoTalkContact"},
	}
	return s.client.cgiAction(ctx, "recordUpdater.cgi", "clear", params)
}

// InviteVideoTalk initiates a video talk invitation.
// VideoTalkPeer.cgi?action=invite
func (s *BuildingService) InviteVideoTalk(ctx context.Context, params map[string]string) error {
	vals := url.Values{}
	for k, v := range params {
		vals.Set(k, v)
	}
	return s.client.cgiAction(ctx, "VideoTalkPeer.cgi", "invite", vals)
}

// CancelVideoTalk cancels a pending video talk invitation.
// VideoTalkPeer.cgi?action=cancel
func (s *BuildingService) CancelVideoTalk(ctx context.Context, params map[string]string) error {
	vals := url.Values{}
	for k, v := range params {
		vals.Set(k, v)
	}
	return s.client.cgiAction(ctx, "VideoTalkPeer.cgi", "cancel", vals)
}

// AnswerVideoTalk accepts an incoming video talk.
// VideoTalkPeer.cgi?action=answer
func (s *BuildingService) AnswerVideoTalk(ctx context.Context, params map[string]string) error {
	vals := url.Values{}
	for k, v := range params {
		vals.Set(k, v)
	}
	return s.client.cgiAction(ctx, "VideoTalkPeer.cgi", "answer", vals)
}

// RefuseVideoTalk declines an incoming video talk.
// VideoTalkPeer.cgi?action=refuse
func (s *BuildingService) RefuseVideoTalk(ctx context.Context, params map[string]string) error {
	vals := url.Values{}
	for k, v := range params {
		vals.Set(k, v)
	}
	return s.client.cgiAction(ctx, "VideoTalkPeer.cgi", "refuse", vals)
}

// HangUpVideoTalk terminates an active video talk session.
// VideoTalkPeer.cgi?action=hangUp
func (s *BuildingService) HangUpVideoTalk(ctx context.Context, params map[string]string) error {
	vals := url.Values{}
	for k, v := range params {
		vals.Set(k, v)
	}
	return s.client.cgiAction(ctx, "VideoTalkPeer.cgi", "hangUp", vals)
}

// SubscribeVideoTalkState subscribes to video talk state changes.
// This is a long-lived stream. The raw response body is returned.
// CGI: VideoTalkPeer.cgi?action=attachState&heartbeat=N
func (s *BuildingService) SubscribeVideoTalkState(ctx context.Context, heartbeat int) (string, error) {
	return s.client.cgiGet(ctx, "VideoTalkPeer.cgi", "attachState", url.Values{
		"heartbeat": {strconv.Itoa(heartbeat)},
	})
}

// UnsubscribeVideoTalkState unsubscribes from video talk state changes.
// CGI: VideoTalkPeer.cgi?action=detachState
func (s *BuildingService) UnsubscribeVideoTalkState(ctx context.Context) error {
	return s.client.cgiAction(ctx, "VideoTalkPeer.cgi", "detachState", nil)
}

// QueryVideoTalkLog queries video talk log records with the given condition
// parameters.
// CGI: recordFinder.cgi?action=find&name=VideoTalkLog&...
func (s *BuildingService) QueryVideoTalkLog(ctx context.Context, params map[string]string) (string, error) {
	vals := url.Values{
		"name": {"VideoTalkLog"},
	}
	for k, v := range params {
		vals.Set(k, v)
	}
	return s.client.cgiGet(ctx, "recordFinder.cgi", "find", vals)
}

// InsertAnnouncement inserts a new announcement record.
// CGI: recordUpdater.cgi?action=insert&name=Announcement&...
func (s *BuildingService) InsertAnnouncement(ctx context.Context, params map[string]string) error {
	vals := url.Values{
		"name": {"Announcement"},
	}
	for k, v := range params {
		vals.Set(k, v)
	}
	return s.client.cgiAction(ctx, "recordUpdater.cgi", "insert", vals)
}

// QueryAlarmRecords queries alarm records with the given condition parameters.
// CGI: recordFinder.cgi?action=find&name=AlarmRecord&...
func (s *BuildingService) QueryAlarmRecords(ctx context.Context, params map[string]string) (string, error) {
	vals := url.Values{
		"name": {"AlarmRecord"},
	}
	for k, v := range params {
		vals.Set(k, v)
	}
	return s.client.cgiGet(ctx, "recordFinder.cgi", "find", vals)
}

// SetElevatorFloorInfo sets the elevator floor information.
// CGI: ElevatorFloorCounter.cgi?action=setElevatorFloorInfo&...
func (s *BuildingService) SetElevatorFloorInfo(ctx context.Context, params map[string]string) error {
	vals := url.Values{}
	for k, v := range params {
		vals.Set(k, v)
	}
	return s.client.cgiAction(ctx, "ElevatorFloorCounter.cgi", "setElevatorFloorInfo", vals)
}

// GetElevatorWorkInfo returns elevator work information.
// CGI: ElevatorFloorCounter.cgi?action=getElevatorWorkInfo
func (s *BuildingService) GetElevatorWorkInfo(ctx context.Context) (string, error) {
	return s.client.cgiGet(ctx, "ElevatorFloorCounter.cgi", "getElevatorWorkInfo", nil)
}

// GetElevatorCaps returns the elevator floor counter capabilities.
// CGI: ElevatorFloorCounter.cgi?action=getCaps
func (s *BuildingService) GetElevatorCaps(ctx context.Context) (string, error) {
	return s.client.cgiGet(ctx, "ElevatorFloorCounter.cgi", "getCaps", nil)
}
