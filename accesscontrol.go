package amcrest

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
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

// ---------------------------------------------------------------------------
// Section 12.1 - Core operations
// ---------------------------------------------------------------------------

// CaptureFace triggers a face capture command.
// accessControl.cgi?action=captureCmd&...
func (s *AccessControlService) CaptureFace(ctx context.Context, params map[string]string) error {
	v := url.Values{}
	for k, val := range params {
		v.Set(k, val)
	}
	return s.client.cgiAction(ctx, "accessControl.cgi", "captureCmd", v)
}

// GetMeasureTemperatureConfig retrieves the MeasureTemperature configuration.
func (s *AccessControlService) GetMeasureTemperatureConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "MeasureTemperature")
}

// SetMeasureTemperatureConfig updates the MeasureTemperature configuration.
func (s *AccessControlService) SetMeasureTemperatureConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetCitizenCompareConfig retrieves the CitizenPictureCompare configuration.
func (s *AccessControlService) GetCitizenCompareConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "CitizenPictureCompare")
}

// SetCitizenCompareConfig updates the CitizenPictureCompare configuration.
func (s *AccessControlService) SetCitizenCompareConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// ---------------------------------------------------------------------------
// Section 12.2 - Access Control Manager (POST JSON APIs)
// ---------------------------------------------------------------------------

// AddSubController adds a sub-controller.
// POST /cgi-bin/api/AccessControl/addSubController
func (s *AccessControlService) AddSubController(ctx context.Context, body interface{}) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/AccessControl/addSubController", body, nil)
}

// ModifySubController modifies a sub-controller.
// POST /cgi-bin/api/AccessControl/modifySubController
func (s *AccessControlService) ModifySubController(ctx context.Context, body interface{}) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/AccessControl/modifySubController", body, nil)
}

// RemoveSubController removes a sub-controller.
// POST /cgi-bin/api/AccessControl/removeSubController
func (s *AccessControlService) RemoveSubController(ctx context.Context, body interface{}) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/AccessControl/removeSubController", body, nil)
}

// GetSubControllerInfo retrieves sub-controller information.
// POST /cgi-bin/api/AccessControl/getSubControllerInfo
func (s *AccessControlService) GetSubControllerInfo(ctx context.Context, body interface{}) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/AccessControl/getSubControllerInfo", body)
}

// GetSubControllerStates retrieves sub-controller states.
// POST /cgi-bin/api/AccessControl/getSubControllerStates
func (s *AccessControlService) GetSubControllerStates(ctx context.Context, body interface{}) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/AccessControl/getSubControllerStates", body)
}

// SetRepeatEnterRoute sets the repeat enter route configuration.
// POST /cgi-bin/api/AccessControl/setRepeatEnterRoute
func (s *AccessControlService) SetRepeatEnterRoute(ctx context.Context, body interface{}) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/AccessControl/setRepeatEnterRoute", body, nil)
}

// GetRepeatEnterRoute retrieves the repeat enter route configuration.
// POST /cgi-bin/api/AccessControl/getRepeatEnterRoute
func (s *AccessControlService) GetRepeatEnterRoute(ctx context.Context) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/AccessControl/getRepeatEnterRoute", map[string]interface{}{})
}

// SetABLockRoute sets the AB lock route configuration.
// POST /cgi-bin/api/AccessControl/setABLockRoute
func (s *AccessControlService) SetABLockRoute(ctx context.Context, body interface{}) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/AccessControl/setABLockRoute", body, nil)
}

// GetABLockRoute retrieves the AB lock route configuration.
// POST /cgi-bin/api/AccessControl/getABLockRoute
func (s *AccessControlService) GetABLockRoute(ctx context.Context) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/AccessControl/getABLockRoute", map[string]interface{}{})
}

// GetLogStatus retrieves the access control log status.
// POST /cgi-bin/api/AccessControl/getLogStatus
func (s *AccessControlService) GetLogStatus(ctx context.Context) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/AccessControl/getLogStatus", map[string]interface{}{})
}

// SyncOfflineLog triggers offline log synchronization.
// POST /cgi-bin/api/AccessControl/syncOfflineLog
func (s *AccessControlService) SyncOfflineLog(ctx context.Context) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/AccessControl/syncOfflineLog", map[string]interface{}{}, nil)
}

// SyncSubControllerTime synchronizes sub-controller time.
// POST /cgi-bin/api/AccessControl/syncSubControllerTime
func (s *AccessControlService) SyncSubControllerTime(ctx context.Context) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/AccessControl/syncSubControllerTime", map[string]interface{}{}, nil)
}

// ---------------------------------------------------------------------------
// Section 12.3 - V1 Face management (CGI-based)
// ---------------------------------------------------------------------------

// AddUserFace adds a user face record.
// POST /cgi-bin/FaceInfoManager.cgi?action=add
func (s *AccessControlService) AddUserFace(ctx context.Context, body interface{}) error {
	return s.client.postJSON(ctx, "/cgi-bin/FaceInfoManager.cgi?action=add", body, nil)
}

// ModifyUserFace modifies a user face record.
// POST /cgi-bin/FaceInfoManager.cgi?action=modify
func (s *AccessControlService) ModifyUserFace(ctx context.Context, body interface{}) error {
	return s.client.postJSON(ctx, "/cgi-bin/FaceInfoManager.cgi?action=modify", body, nil)
}

// DeleteUserFace removes a user face record by user ID.
// FaceInfoManager.cgi?action=remove&UserID=X
func (s *AccessControlService) DeleteUserFace(ctx context.Context, userID string) error {
	params := url.Values{
		"UserID": {userID},
	}
	return s.client.cgiAction(ctx, "FaceInfoManager.cgi", "remove", params)
}

// StartFindUserFace starts a user face search and returns the search token.
// FaceInfoManager.cgi?action=startFind&...
func (s *AccessControlService) StartFindUserFace(ctx context.Context, params map[string]string) (string, error) {
	v := url.Values{}
	for k, val := range params {
		v.Set(k, val)
	}
	return s.client.cgiGet(ctx, "FaceInfoManager.cgi", "startFind", v)
}

// DoFindUserFace retrieves a page of user face search results.
// FaceInfoManager.cgi?action=doFind&token=T&beginNumber=O&count=C
func (s *AccessControlService) DoFindUserFace(ctx context.Context, token string, offset, count int) (string, error) {
	params := url.Values{
		"token":       {token},
		"beginNumber": {fmt.Sprintf("%d", offset)},
		"count":       {fmt.Sprintf("%d", count)},
	}
	return s.client.cgiGet(ctx, "FaceInfoManager.cgi", "doFind", params)
}

// StopFindUserFace stops a user face search.
// FaceInfoManager.cgi?action=stopFind&token=T
func (s *AccessControlService) StopFindUserFace(ctx context.Context, token string) error {
	params := url.Values{
		"token": {token},
	}
	return s.client.cgiAction(ctx, "FaceInfoManager.cgi", "stopFind", params)
}

// UpdateCard updates an existing access control card record.
// recordUpdater.cgi?action=update&name=AccessControlCard&recno=N&...
func (s *AccessControlService) UpdateCard(ctx context.Context, recno int, params map[string]string) error {
	v := url.Values{
		"name":  {"AccessControlCard"},
		"recno": {fmt.Sprintf("%d", recno)},
	}
	for k, val := range params {
		v.Set(k, val)
	}
	return s.client.cgiAction(ctx, "recordUpdater.cgi", "update", v)
}

// GetCardCount returns the total number of access control card records.
// recordFinder.cgi?action=getQuerySize&name=AccessControlCard
func (s *AccessControlService) GetCardCount(ctx context.Context) (int, error) {
	params := url.Values{
		"name": {"AccessControlCard"},
	}
	body, err := s.client.cgiGet(ctx, "recordFinder.cgi", "getQuerySize", params)
	if err != nil {
		return 0, err
	}
	kv := parseKV(body)
	count, err := strconv.Atoi(kv["count"])
	if err != nil {
		return 0, fmt.Errorf("amcrest: parsing card count: %w", err)
	}
	return count, nil
}

// ---------------------------------------------------------------------------
// Section 12.4 - V2 User management (POST JSON APIs)
// ---------------------------------------------------------------------------

// AddUserV2 adds a user via the V2 API.
// POST /cgi-bin/api/AccessControl/addUser
func (s *AccessControlService) AddUserV2(ctx context.Context, body interface{}) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/AccessControl/addUser", body, nil)
}

// ModifyUserV2 modifies a user via the V2 API.
// POST /cgi-bin/api/AccessControl/modifyUser
func (s *AccessControlService) ModifyUserV2(ctx context.Context, body interface{}) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/AccessControl/modifyUser", body, nil)
}

// DeleteAllUsersV2 deletes all users via the V2 API.
// POST /cgi-bin/api/AccessControl/deleteAllUser
func (s *AccessControlService) DeleteAllUsersV2(ctx context.Context) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/AccessControl/deleteAllUser", map[string]interface{}{}, nil)
}

// DeleteUsersV2 deletes specified users via the V2 API.
// POST /cgi-bin/api/AccessControl/deleteUser
func (s *AccessControlService) DeleteUsersV2(ctx context.Context, body interface{}) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/AccessControl/deleteUser", body, nil)
}

// FindUsersV2 finds users via the V2 API.
// POST /cgi-bin/api/AccessControl/findUser
func (s *AccessControlService) FindUsersV2(ctx context.Context, body interface{}) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/AccessControl/findUser", body)
}

// StartFindUsersV2 starts a user search via the V2 API.
// POST /cgi-bin/api/AccessControl/startFindUser
func (s *AccessControlService) StartFindUsersV2(ctx context.Context, body interface{}) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/AccessControl/startFindUser", body)
}

// DoFindUsersV2 performs a paginated user search via the V2 API.
// POST /cgi-bin/api/AccessControl/doFindUser
func (s *AccessControlService) DoFindUsersV2(ctx context.Context, body interface{}) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/AccessControl/doFindUser", body)
}

// StopFindUsersV2 stops a user search via the V2 API.
// POST /cgi-bin/api/AccessControl/stopFindUser
func (s *AccessControlService) StopFindUsersV2(ctx context.Context, token int) error {
	body := map[string]interface{}{"token": token}
	return s.client.postJSON(ctx, "/cgi-bin/api/AccessControl/stopFindUser", body, nil)
}

// AddCardsV2 adds cards via the V2 API.
// POST /cgi-bin/api/AccessControl/addCard
func (s *AccessControlService) AddCardsV2(ctx context.Context, body interface{}) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/AccessControl/addCard", body, nil)
}

// ModifyCardsV2 modifies cards via the V2 API.
// POST /cgi-bin/api/AccessControl/modifyCard
func (s *AccessControlService) ModifyCardsV2(ctx context.Context, body interface{}) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/AccessControl/modifyCard", body, nil)
}

// DeleteAllCardsV2 deletes all cards via the V2 API.
// POST /cgi-bin/api/AccessControl/deleteAllCard
func (s *AccessControlService) DeleteAllCardsV2(ctx context.Context) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/AccessControl/deleteAllCard", map[string]interface{}{}, nil)
}

// DeleteCardsV2 deletes specified cards via the V2 API.
// POST /cgi-bin/api/AccessControl/deleteCard
func (s *AccessControlService) DeleteCardsV2(ctx context.Context, body interface{}) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/AccessControl/deleteCard", body, nil)
}

// FindCardsV2 finds cards via the V2 API.
// POST /cgi-bin/api/AccessControl/findCard
func (s *AccessControlService) FindCardsV2(ctx context.Context, body interface{}) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/AccessControl/findCard", body)
}

// AddFingerprintsV2 adds fingerprints via the V2 API.
// POST /cgi-bin/api/AccessControl/addFingerprint
func (s *AccessControlService) AddFingerprintsV2(ctx context.Context, body interface{}) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/AccessControl/addFingerprint", body, nil)
}

// ModifyFingerprintV2 modifies a fingerprint via the V2 API.
// POST /cgi-bin/api/AccessControl/modifyFingerprint
func (s *AccessControlService) ModifyFingerprintV2(ctx context.Context, body interface{}) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/AccessControl/modifyFingerprint", body, nil)
}

// DeleteAllFingerprintsV2 deletes all fingerprints via the V2 API.
// POST /cgi-bin/api/AccessControl/deleteAllFingerprint
func (s *AccessControlService) DeleteAllFingerprintsV2(ctx context.Context) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/AccessControl/deleteAllFingerprint", map[string]interface{}{}, nil)
}

// DeleteFingerprintsV2 deletes specified fingerprints via the V2 API.
// POST /cgi-bin/api/AccessControl/deleteFingerprint
func (s *AccessControlService) DeleteFingerprintsV2(ctx context.Context, body interface{}) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/AccessControl/deleteFingerprint", body, nil)
}

// FindFingerprintsV2 finds fingerprints via the V2 API.
// POST /cgi-bin/api/AccessControl/findFingerprint
func (s *AccessControlService) FindFingerprintsV2(ctx context.Context, body interface{}) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/AccessControl/findFingerprint", body)
}

// AddFacesV2 adds faces via the V2 API.
// POST /cgi-bin/api/AccessControl/addFace
func (s *AccessControlService) AddFacesV2(ctx context.Context, body interface{}) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/AccessControl/addFace", body, nil)
}

// UpdateFacesV2 updates faces via the V2 API.
// POST /cgi-bin/api/AccessControl/updateFace
func (s *AccessControlService) UpdateFacesV2(ctx context.Context, body interface{}) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/AccessControl/updateFace", body, nil)
}

// DeleteAllFacesV2 deletes all faces via the V2 API.
// POST /cgi-bin/api/AccessControl/deleteAllFace
func (s *AccessControlService) DeleteAllFacesV2(ctx context.Context) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/AccessControl/deleteAllFace", map[string]interface{}{}, nil)
}

// DeleteFacesV2 deletes specified faces via the V2 API.
// POST /cgi-bin/api/AccessControl/deleteFace
func (s *AccessControlService) DeleteFacesV2(ctx context.Context, body interface{}) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/AccessControl/deleteFace", body, nil)
}

// FindFacesV2 finds faces via the V2 API.
// POST /cgi-bin/api/AccessControl/findFace
func (s *AccessControlService) FindFacesV2(ctx context.Context, body interface{}) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/AccessControl/findFace", body)
}

// GetProtocolCaps retrieves the access control protocol capabilities.
// POST /cgi-bin/api/AccessControl/getProtocolCaps
func (s *AccessControlService) GetProtocolCaps(ctx context.Context) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/AccessControl/getProtocolCaps", map[string]interface{}{})
}

// ---------------------------------------------------------------------------
// Section 12.5 - Admin Password (record-based CGI)
// ---------------------------------------------------------------------------

// AddAdminPassword inserts an admin password record.
// recordUpdater.cgi?action=insert&name=AccessControlPasswordRecord&...
func (s *AccessControlService) AddAdminPassword(ctx context.Context, params map[string]string) error {
	v := url.Values{
		"name": {"AccessControlPasswordRecord"},
	}
	for k, val := range params {
		v.Set(k, val)
	}
	return s.client.cgiAction(ctx, "recordUpdater.cgi", "insert", v)
}

// ModifyAdminPassword updates an admin password record.
// recordUpdater.cgi?action=update&name=AccessControlPasswordRecord&...
func (s *AccessControlService) ModifyAdminPassword(ctx context.Context, params map[string]string) error {
	v := url.Values{
		"name": {"AccessControlPasswordRecord"},
	}
	for k, val := range params {
		v.Set(k, val)
	}
	return s.client.cgiAction(ctx, "recordUpdater.cgi", "update", v)
}

// DeleteAdminPassword removes an admin password record by record number.
// recordUpdater.cgi?action=remove&name=AccessControlPasswordRecord&recno=N
func (s *AccessControlService) DeleteAdminPassword(ctx context.Context, recno int) error {
	params := url.Values{
		"name":  {"AccessControlPasswordRecord"},
		"recno": {fmt.Sprintf("%d", recno)},
	}
	return s.client.cgiAction(ctx, "recordUpdater.cgi", "remove", params)
}

// FindAdminPassword queries admin password records.
// recordFinder.cgi?action=find&name=AccessControlPasswordRecord&...
func (s *AccessControlService) FindAdminPassword(ctx context.Context, params map[string]string) (string, error) {
	v := url.Values{
		"name": {"AccessControlPasswordRecord"},
	}
	for k, val := range params {
		v.Set(k, val)
	}
	return s.client.cgiGet(ctx, "recordFinder.cgi", "find", v)
}

// GetAdminPasswordCount returns the total number of admin password records.
// recordFinder.cgi?action=getQuerySize&name=AccessControlPasswordRecord
func (s *AccessControlService) GetAdminPasswordCount(ctx context.Context) (int, error) {
	params := url.Values{
		"name": {"AccessControlPasswordRecord"},
	}
	body, err := s.client.cgiGet(ctx, "recordFinder.cgi", "getQuerySize", params)
	if err != nil {
		return 0, err
	}
	kv := parseKV(body)
	count, err := strconv.Atoi(kv["count"])
	if err != nil {
		return 0, fmt.Errorf("amcrest: parsing admin password count: %w", err)
	}
	return count, nil
}
