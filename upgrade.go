package amcrest

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
)

// UpgradeService handles firmware upgrade related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 213-218 (Section 4.12)
type UpgradeService struct {
	client *Client
}

// GetState returns the current firmware upgrade state.
// CGI: upgrader.cgi?action=getState
func (s *UpgradeService) GetState(ctx context.Context) (map[string]string, error) {
	body, err := s.client.cgiGet(ctx, "upgrader.cgi", "getState", nil)
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}

// Cancel cancels an in-progress firmware upgrade.
// CGI: upgrader.cgi?action=cancel
func (s *UpgradeService) Cancel(ctx context.Context) error {
	return s.client.cgiAction(ctx, "upgrader.cgi", "cancel", nil)
}

// CheckCloudUpdate checks for available cloud firmware updates.
// POST /cgi-bin/api/CloudUpgrader/check with JSON body {"way":0}.
func (s *UpgradeService) CheckCloudUpdate(ctx context.Context) (map[string]string, error) {
	reqBody := map[string]int{"way": 0}
	var result map[string]string
	err := s.client.postJSON(ctx, "/cgi-bin/api/CloudUpgrader/check", reqBody, &result)
	if err != nil {
		return nil, fmt.Errorf("CheckCloudUpdate: %w", err)
	}
	return result, nil
}

// GetAutoUpgradeConfig returns the auto-upgrade configuration.
func (s *UpgradeService) GetAutoUpgradeConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "AutoUpgrade")
}

// SetAutoUpgradeConfig updates the auto-upgrade configuration.
func (s *UpgradeService) SetAutoUpgradeConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetCloudUpgradeMode returns the cloud upgrade mode configuration.
func (s *UpgradeService) GetCloudUpgradeMode(ctx context.Context) (map[string]string, error) {
	return s.client.getConfig(ctx, "CloudUpgrade")
}

// SetCloudUpgradeMode updates the cloud upgrade mode.
func (s *UpgradeService) SetCloudUpgradeMode(ctx context.Context, params url.Values) error {
	return s.client.cgiAction(ctx, "configManager.cgi", "setConfig", params)
}

// UploadFirmware uploads a firmware binary to the camera as multipart/form-data.
// PDF 4.12.1: POST /cgi-bin/upgrader.cgi?action=uploadFirmware
func (s *UpgradeService) UploadFirmware(ctx context.Context, firmware []byte, filename string) error {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, err := w.CreateFormFile("firmware", filename)
	if err != nil {
		return fmt.Errorf("UpgradeService.UploadFirmware: creating form file: %w", err)
	}
	if _, err := part.Write(firmware); err != nil {
		return fmt.Errorf("UpgradeService.UploadFirmware: writing firmware: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("UpgradeService.UploadFirmware: closing multipart writer: %w", err)
	}

	u := s.client.baseURL + "/cgi-bin/upgrader.cgi?action=uploadFirmware"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, &buf)
	if err != nil {
		return fmt.Errorf("UpgradeService.UploadFirmware: creating request: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("UpgradeService.UploadFirmware: executing request: %w", err)
	}
	return checkOK(resp)
}

// UpdateByURL instructs the camera to download and apply firmware from a URL.
// PDF 4.12.3: upgrader.cgi?action=updateFirmwareByUrl&Url=X&checkType=0&checkSum=Y
func (s *UpgradeService) UpdateByURL(ctx context.Context, firmwareURL, checkSum string) error {
	params := url.Values{
		"Url":       {firmwareURL},
		"checkType": {"0"},
		"checkSum":  {checkSum},
	}
	return s.client.cgiAction(ctx, "upgrader.cgi", "updateFirmwareByUrl", params)
}

// ExecuteCloudUpdate executes a cloud firmware update with the given way parameter.
// PDF 4.12.6: POST /cgi-bin/api/CloudUpgrader/execute with JSON {"way": N}.
func (s *UpgradeService) ExecuteCloudUpdate(ctx context.Context, way int) error {
	reqBody := map[string]int{"way": way}
	return s.client.postJSON(ctx, "/cgi-bin/api/CloudUpgrader/execute", reqBody, nil)
}

// CancelCloudUpdate cancels a cloud firmware update in progress.
// PDF 4.12.7: POST /cgi-bin/api/CloudUpgrader/cancel.
func (s *UpgradeService) CancelCloudUpdate(ctx context.Context) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/CloudUpgrader/cancel", struct{}{}, nil)
}
