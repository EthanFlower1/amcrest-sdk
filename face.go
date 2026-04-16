package amcrest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Face event codes for use with EventService.Subscribe.
// Example: client.Event.Subscribe(ctx, []string{EventFaceDetection}, 5)
const (
	EventFaceDetection       = "FaceDetection"
	EventFaceRecognition     = "FaceRecognition"
	EventFaceFeatureAbstract = "FaceFeatureAbstract"
)

// FaceService handles face detection and recognition related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 334-361 (Section 9.2)
type FaceService struct {
	client *Client
}

// CreateGroup creates a new face recognition group with the given name and detail.
// Returns the groupID assigned by the camera.
// CGI: faceRecognitionServer.cgi?action=createGroup&groupName=X&groupDetail=Y
func (s *FaceService) CreateGroup(ctx context.Context, name, detail string) (string, error) {
	params := url.Values{
		"groupName":   {name},
		"groupDetail": {detail},
	}
	body, err := s.client.cgiGet(ctx, "faceRecognitionServer.cgi", "createGroup", params)
	if err != nil {
		return "", fmt.Errorf("FaceService.CreateGroup: %w", err)
	}
	kv := parseKV(body)
	id, ok := kv["groupID"]
	if !ok {
		return "", fmt.Errorf("FaceService.CreateGroup: groupID not found in response: %s", body)
	}
	return id, nil
}

// DeleteGroup deletes the face recognition group with the given groupID.
// CGI: faceRecognitionServer.cgi?action=deleteGroup&groupID=X
func (s *FaceService) DeleteGroup(ctx context.Context, groupID string) error {
	params := url.Values{
		"groupID": {groupID},
	}
	return s.client.cgiAction(ctx, "faceRecognitionServer.cgi", "deleteGroup", params)
}

// FindGroup returns the raw response listing all face recognition groups.
// CGI: faceRecognitionServer.cgi?action=findGroup
func (s *FaceService) FindGroup(ctx context.Context) (string, error) {
	body, err := s.client.cgiGet(ctx, "faceRecognitionServer.cgi", "findGroup", nil)
	if err != nil {
		return "", fmt.Errorf("FaceService.FindGroup: %w", err)
	}
	return body, nil
}

// GetGroupForChannel returns the raw response for the face group assigned to
// the given video channel.
// CGI: faceRecognitionServer.cgi?action=getGroup&channel=N
func (s *FaceService) GetGroupForChannel(ctx context.Context, channel int) (string, error) {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	body, err := s.client.cgiGet(ctx, "faceRecognitionServer.cgi", "getGroup", params)
	if err != nil {
		return "", fmt.Errorf("FaceService.GetGroupForChannel: %w", err)
	}
	return body, nil
}

// ModifyGroup modifies the name and detail of an existing face recognition group.
// PDF 9.2.2: faceRecognitionServer.cgi?action=modifyGroup
func (s *FaceService) ModifyGroup(ctx context.Context, groupID, name, detail string) error {
	params := url.Values{
		"groupID":     {groupID},
		"groupName":   {name},
		"groupDetail": {detail},
	}
	return s.client.cgiAction(ctx, "faceRecognitionServer.cgi", "modifyGroup", params)
}

// DeployGroup deploys a face recognition group to the specified channels with
// the given similarity threshold.
// PDF 9.2.4: faceRecognitionServer.cgi?action=putDisposition
func (s *FaceService) DeployGroup(ctx context.Context, groupID string, channels []int, similarity int) error {
	params := url.Values{
		"groupID":    {groupID},
		"similarity": {fmt.Sprintf("%d", similarity)},
	}
	chStrs := make([]string, len(channels))
	for i, ch := range channels {
		chStrs[i] = fmt.Sprintf("%d", ch)
	}
	params.Set("channel", strings.Join(chStrs, ","))
	return s.client.cgiAction(ctx, "faceRecognitionServer.cgi", "putDisposition", params)
}

// UndeployGroup removes a face recognition group from the specified channels.
// PDF 9.2.4: faceRecognitionServer.cgi?action=deleteDisposition
func (s *FaceService) UndeployGroup(ctx context.Context, groupID string, channels []int) error {
	params := url.Values{
		"groupID": {groupID},
	}
	chStrs := make([]string, len(channels))
	for i, ch := range channels {
		chStrs[i] = fmt.Sprintf("%d", ch)
	}
	params.Set("channel", strings.Join(chStrs, ","))
	return s.client.cgiAction(ctx, "faceRecognitionServer.cgi", "deleteDisposition", params)
}

// SetGroupForChannel assigns face recognition groups with similarity thresholds
// to a specific channel.
// PDF 9.2.4: faceRecognitionServer.cgi?action=setGroup
func (s *FaceService) SetGroupForChannel(ctx context.Context, channel int, groups []string, similarities []int) error {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	params.Set("groupID", strings.Join(groups, ","))
	simStrs := make([]string, len(similarities))
	for i, sim := range similarities {
		simStrs[i] = fmt.Sprintf("%d", sim)
	}
	params.Set("similarity", strings.Join(simStrs, ","))
	return s.client.cgiAction(ctx, "faceRecognitionServer.cgi", "setGroup", params)
}

// ReAbstractByGroup re-abstracts face features for all persons in the given group.
// PDF 9.2.6: faceRecognitionServer.cgi?action=reAbstract&groupID=X
func (s *FaceService) ReAbstractByGroup(ctx context.Context, groupID string) error {
	params := url.Values{
		"groupID": {groupID},
	}
	return s.client.cgiAction(ctx, "faceRecognitionServer.cgi", "reAbstract", params)
}

// AddPerson adds a person to a face recognition group. The params map should
// include keys like groupID, name, sex, birthday, etc. Returns the UID assigned.
// PDF 9.2.7: faceRecognitionServer.cgi?action=addPerson
func (s *FaceService) AddPerson(ctx context.Context, params map[string]string) (string, error) {
	qv := url.Values{}
	for k, v := range params {
		qv.Set(k, v)
	}
	body, err := s.client.cgiGet(ctx, "faceRecognitionServer.cgi", "addPerson", qv)
	if err != nil {
		return "", fmt.Errorf("FaceService.AddPerson: %w", err)
	}
	kv := parseKV(body)
	uid, ok := kv["uid"]
	if !ok {
		return "", fmt.Errorf("FaceService.AddPerson: uid not found in response: %s", body)
	}
	return uid, nil
}

// ModifyPerson modifies an existing person in a face recognition group.
// PDF 9.2.8: faceRecognitionServer.cgi?action=modifyPerson
func (s *FaceService) ModifyPerson(ctx context.Context, params map[string]string) error {
	qv := url.Values{}
	for k, v := range params {
		qv.Set(k, v)
	}
	return s.client.cgiAction(ctx, "faceRecognitionServer.cgi", "modifyPerson", qv)
}

// DeletePerson deletes a person from a face recognition group.
// PDF 9.2.9: faceRecognitionServer.cgi?action=deletePerson&groupID=X&uid=Y
func (s *FaceService) DeletePerson(ctx context.Context, groupID, uid string) error {
	params := url.Values{
		"groupID": {groupID},
		"uid":     {uid},
	}
	return s.client.cgiAction(ctx, "faceRecognitionServer.cgi", "deletePerson", params)
}

// StartFindPerson begins a person search and returns a search token.
// PDF 9.2.10: faceRecognitionServer.cgi?action=startFindPerson
func (s *FaceService) StartFindPerson(ctx context.Context, params map[string]string) (string, error) {
	qv := url.Values{}
	for k, v := range params {
		qv.Set(k, v)
	}
	body, err := s.client.cgiGet(ctx, "faceRecognitionServer.cgi", "startFindPerson", qv)
	if err != nil {
		return "", fmt.Errorf("FaceService.StartFindPerson: %w", err)
	}
	kv := parseKV(body)
	token, ok := kv["token"]
	if !ok {
		return "", fmt.Errorf("FaceService.StartFindPerson: token not found in response: %s", body)
	}
	return token, nil
}

// DoFindPerson retrieves a page of person search results using the given token.
// PDF 9.2.10: faceRecognitionServer.cgi?action=doFindPerson&token=X&offset=N&count=M
func (s *FaceService) DoFindPerson(ctx context.Context, token string, offset, count int) (string, error) {
	params := url.Values{
		"token":  {token},
		"offset": {fmt.Sprintf("%d", offset)},
		"count":  {fmt.Sprintf("%d", count)},
	}
	body, err := s.client.cgiGet(ctx, "faceRecognitionServer.cgi", "doFindPerson", params)
	if err != nil {
		return "", fmt.Errorf("FaceService.DoFindPerson: %w", err)
	}
	return body, nil
}

// StopFindPerson stops a person search and releases the token.
// PDF 9.2.10: faceRecognitionServer.cgi?action=stopFindPerson&token=X
func (s *FaceService) StopFindPerson(ctx context.Context, token string) error {
	params := url.Values{
		"token": {token},
	}
	return s.client.cgiAction(ctx, "faceRecognitionServer.cgi", "stopFindPerson", params)
}

// GetFaceRecAlarmConfig returns the face recognition alarm configuration.
// PDF 9.2.12: getRawConfig FaceRecognitionAlarm
func (s *FaceService) GetFaceRecAlarmConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "FaceRecognitionAlarm")
}

// SetFaceRecAlarmConfig updates the face recognition alarm configuration.
// PDF 9.2.12
func (s *FaceService) SetFaceRecAlarmConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// StartFindByPic starts a face search using a picture. The picture data is
// POSTed as image/jpeg to the CGI with search parameters in the query string.
// Returns the raw multipart/JSON response body which contains token, progress,
// and totalCount fields.
// PDF 9.2.13: faceRecognitionServer.cgi?action=startFindByPic
func (s *FaceService) StartFindByPic(ctx context.Context, picData []byte, groupIDs []string, similarity int, maxCandidate int) (string, error) {
	qv := url.Values{
		"action":     {"startFindByPic"},
		"Similarity": {fmt.Sprintf("%d", similarity)},
	}
	for i, gid := range groupIDs {
		qv.Set(fmt.Sprintf("GroupID[%d]", i), gid)
	}
	if maxCandidate > 0 {
		qv.Set("MaxCandidate", fmt.Sprintf("%d", maxCandidate))
	}

	u := s.client.baseURL + "/cgi-bin/faceRecognitionServer.cgi?" + qv.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(picData))
	if err != nil {
		return "", fmt.Errorf("FaceService.StartFindByPic: creating request: %w", err)
	}
	req.Header.Set("Content-Type", "image/jpeg")

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("FaceService.StartFindByPic: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", &APIError{StatusCode: resp.StatusCode, Message: "startFindByPic failed"}
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("FaceService.StartFindByPic: reading body: %w", err)
	}
	return string(data), nil
}

// DoFindByPic retrieves a page of face-by-picture search results.
// The response is multipart; the first part is JSON with Found and Candidates,
// followed by JPEG image parts. Returned as raw string.
// PDF 9.2.13: faceRecognitionServer.cgi?action=doFindByPic
func (s *FaceService) DoFindByPic(ctx context.Context, token string, index, count int) (string, error) {
	params := url.Values{
		"token": {token},
		"index": {fmt.Sprintf("%d", index)},
		"count": {fmt.Sprintf("%d", count)},
	}
	body, err := s.client.cgiGet(ctx, "faceRecognitionServer.cgi", "doFindByPic", params)
	if err != nil {
		return "", fmt.Errorf("FaceService.DoFindByPic: %w", err)
	}
	return body, nil
}

// StopFindByPic stops a face-by-picture search session.
// PDF 9.2.13: faceRecognitionServer.cgi?action=stopFindByPic
func (s *FaceService) StopFindByPic(ctx context.Context, token string) error {
	params := url.Values{
		"token": {token},
	}
	return s.client.cgiAction(ctx, "faceRecognitionServer.cgi", "stopFindByPic", params)
}

// StartFindHistoryByPic starts a history face search using a picture. The
// picture data is POSTed as image/jpeg with search parameters in the query
// string. Returns the raw response body containing token, progress, and
// totalCount.
// PDF 9.2.14: faceRecognitionServer.cgi?action=startFindHistoryByPic
func (s *FaceService) StartFindHistoryByPic(ctx context.Context, picData []byte, params map[string]string) (string, error) {
	qv := url.Values{
		"action": {"startFindHistoryByPic"},
	}
	for k, v := range params {
		qv.Set(k, v)
	}

	u := s.client.baseURL + "/cgi-bin/faceRecognitionServer.cgi?" + qv.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(picData))
	if err != nil {
		return "", fmt.Errorf("FaceService.StartFindHistoryByPic: creating request: %w", err)
	}
	req.Header.Set("Content-Type", "image/jpeg")

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("FaceService.StartFindHistoryByPic: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", &APIError{StatusCode: resp.StatusCode, Message: "startFindHistoryByPic failed"}
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("FaceService.StartFindHistoryByPic: reading body: %w", err)
	}
	return string(data), nil
}

// DoFindHistoryByPic retrieves a page of history face-by-picture search results.
// PDF 9.2.14: faceRecognitionServer.cgi?action=doFindHistoryByPic
func (s *FaceService) DoFindHistoryByPic(ctx context.Context, token string, index, count int) (string, error) {
	params := url.Values{
		"token": {token},
		"index": {fmt.Sprintf("%d", index)},
		"count": {fmt.Sprintf("%d", count)},
	}
	body, err := s.client.cgiGet(ctx, "faceRecognitionServer.cgi", "doFindHistoryByPic", params)
	if err != nil {
		return "", fmt.Errorf("FaceService.DoFindHistoryByPic: %w", err)
	}
	return body, nil
}

// StopFindHistoryByPic stops a history face-by-picture search session.
// PDF 9.2.14: faceRecognitionServer.cgi?action=stopFindHistoryByPic
func (s *FaceService) StopFindHistoryByPic(ctx context.Context, token string) error {
	params := url.Values{
		"token": {token},
	}
	return s.client.cgiAction(ctx, "faceRecognitionServer.cgi", "stopFindHistoryByPic", params)
}

// StopReAbstractByGroup stops a group re-abstract operation by its token.
// PDF 9.2.6: faceRecognitionServer.cgi?action=stopGroupReAbstract
func (s *FaceService) StopReAbstractByGroup(ctx context.Context, token string) error {
	params := url.Values{
		"token": {token},
	}
	return s.client.cgiAction(ctx, "faceRecognitionServer.cgi", "stopGroupReAbstract", params)
}

// ReAbstractByPerson re-abstracts face features for specific persons by UID.
// PDF 9.2.11: faceRecognitionServer.cgi?action=reAbstract&uid[0]=X&uid[1]=Y
func (s *FaceService) ReAbstractByPerson(ctx context.Context, uids []string) error {
	params := url.Values{}
	for i, uid := range uids {
		params.Set(fmt.Sprintf("uid[%d]", i), uid)
	}
	return s.client.cgiAction(ctx, "faceRecognitionServer.cgi", "reAbstract", params)
}

// StopReAbstractByPerson stops a person-level re-abstract operation.
// PDF 9.2.11: faceRecognitionServer.cgi?action=stopReAbstract
func (s *FaceService) StopReAbstractByPerson(ctx context.Context) error {
	return s.client.cgiAction(ctx, "faceRecognitionServer.cgi", "stopReAbstract", nil)
}

// SetFaceIDThreshold sets the Face-ID recognition comparison threshold.
// PDF 9.2.19: configManager setConfig CitizenPictureCompareRule.Threshold=N
func (s *FaceService) SetFaceIDThreshold(ctx context.Context, threshold int) error {
	return s.client.setConfig(ctx, map[string]string{
		"CitizenPictureCompareRule.Threshold": fmt.Sprintf("%d", threshold),
	})
}

// GetFaceRecEventHandler returns the face recognition event handler configuration.
// PDF 9.2.18: getRawConfig FaceRecognitionEventHandler
func (s *FaceService) GetFaceRecEventHandler(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "FaceRecognitionEventHandler")
}

// SetFaceRecEventHandler updates the face recognition event handler configuration.
// PDF 9.2.18
func (s *FaceService) SetFaceRecEventHandler(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// ExportFaceDB exports the face database for the given group as binary data.
// PDF 9.2.20: POST /cgi-bin/api/FaceLibImExport/export
func (s *FaceService) ExportFaceDB(ctx context.Context, groupID string) ([]byte, error) {
	reqBody := map[string]string{"groupID": groupID}
	body, err := s.client.postRaw(ctx, "/cgi-bin/api/FaceLibImExport/export", reqBody)
	if err != nil {
		return nil, fmt.Errorf("FaceService.ExportFaceDB: %w", err)
	}
	return []byte(body), nil
}

// ImportFaceDB imports a face database binary.
// PDF 9.2.21: POST /cgi-bin/api/FaceLibImExport/import
func (s *FaceService) ImportFaceDB(ctx context.Context, data []byte) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/FaceLibImExport/import", data, nil)
}
