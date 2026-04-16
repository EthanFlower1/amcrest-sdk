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

// GetLockStatus returns the lock status for the given channel.
// accessControl.cgi?action=getLockStatus&channel=N
func (s *AccessControlService) GetLockStatus(ctx context.Context, channel int) (map[string]string, error) {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	body, err := s.client.cgiGet(ctx, "accessControl.cgi", "getLockStatus", params)
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}

// CaptureFingerprint triggers a fingerprint capture on the given reader.
// accessControl.cgi?action=captureFingerprint&ReaderID=X
func (s *AccessControlService) CaptureFingerprint(ctx context.Context, readerID string) error {
	params := url.Values{
		"ReaderID": {readerID},
	}
	return s.client.cgiAction(ctx, "accessControl.cgi", "captureFingerprint", params)
}

// QueryRecords queries access control card records.
// recordFinder.cgi?action=find&name=AccessControlCardRec plus extra params.
func (s *AccessControlService) QueryRecords(ctx context.Context, params map[string]string) (string, error) {
	v := url.Values{
		"name": {"AccessControlCardRec"},
	}
	for k, val := range params {
		v.Set(k, val)
	}
	return s.client.cgiGet(ctx, "recordFinder.cgi", "find", v)
}

// QueryAlarmRecords queries access control alarm records.
// recordFinder.cgi?action=find&name=AccessControlAlarmRecord plus extra params.
func (s *AccessControlService) QueryAlarmRecords(ctx context.Context, params map[string]string) (string, error) {
	v := url.Values{
		"name": {"AccessControlAlarmRecord"},
	}
	for k, val := range params {
		v.Set(k, val)
	}
	return s.client.cgiGet(ctx, "recordFinder.cgi", "find", v)
}

// SetGeneralConfig updates the AccessControlGeneral configuration.
func (s *AccessControlService) SetGeneralConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// SetControlConfig updates the AccessControl configuration.
func (s *AccessControlService) SetControlConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetWiegandConfig retrieves the Wiegand configuration.
func (s *AccessControlService) GetWiegandConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "Wiegand")
}

// SetWiegandConfig updates the Wiegand configuration.
func (s *AccessControlService) SetWiegandConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetTimeSchedule retrieves the AccessTimeSchedule configuration.
func (s *AccessControlService) GetTimeSchedule(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "AccessTimeSchedule")
}

// SetTimeSchedule updates the AccessTimeSchedule configuration.
func (s *AccessControlService) SetTimeSchedule(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetSpecialDayGroup retrieves the SpecialDayGroup configuration.
func (s *AccessControlService) GetSpecialDayGroup(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "SpecialDayGroup")
}

// SetSpecialDayGroup updates the SpecialDayGroup configuration.
func (s *AccessControlService) SetSpecialDayGroup(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetSpecialDaysSchedule retrieves the SpecialDaysSchedule configuration.
func (s *AccessControlService) GetSpecialDaysSchedule(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "SpecialDaysSchedule")
}

// SetSpecialDaysSchedule updates the SpecialDaysSchedule configuration.
func (s *AccessControlService) SetSpecialDaysSchedule(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetCaps retrieves the access control capabilities.
// POST /cgi-bin/api/AccessControl/getCaps
func (s *AccessControlService) GetCaps(ctx context.Context) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/AccessControl/getCaps", map[string]interface{}{})
}

// InsertCard inserts an access control card record.
// recordUpdater.cgi?action=insert&name=AccessControlCard plus params.
func (s *AccessControlService) InsertCard(ctx context.Context, params map[string]string) error {
	v := url.Values{
		"name": {"AccessControlCard"},
	}
	for k, val := range params {
		v.Set(k, val)
	}
	return s.client.cgiAction(ctx, "recordUpdater.cgi", "insert", v)
}

// RemoveCard removes an access control card record by record number.
// recordUpdater.cgi?action=remove&name=AccessControlCard&recno=N
func (s *AccessControlService) RemoveCard(ctx context.Context, recno int) error {
	params := url.Values{
		"name":  {"AccessControlCard"},
		"recno": {fmt.Sprintf("%d", recno)},
	}
	return s.client.cgiAction(ctx, "recordUpdater.cgi", "remove", params)
}

// ClearCards removes all access control card records.
// recordUpdater.cgi?action=clear&name=AccessControlCard
func (s *AccessControlService) ClearCards(ctx context.Context) error {
	params := url.Values{
		"name": {"AccessControlCard"},
	}
	return s.client.cgiAction(ctx, "recordUpdater.cgi", "clear", params)
}

// FindCards queries access control card records.
// recordFinder.cgi?action=find&name=AccessControlCard plus extra params.
func (s *AccessControlService) FindCards(ctx context.Context, params map[string]string) (string, error) {
	v := url.Values{
		"name": {"AccessControlCard"},
	}
	for k, val := range params {
		v.Set(k, val)
	}
	return s.client.cgiGet(ctx, "recordFinder.cgi", "find", v)
}
