package amcrest

import (
	"context"
	"fmt"
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
