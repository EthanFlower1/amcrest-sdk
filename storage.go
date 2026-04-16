package amcrest

import (
	"context"
	"fmt"
	"net/url"
)

// StorageService handles storage-related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 259-278 (Section 6)
type StorageService struct {
	client *Client
}

// GetDiskInfo returns storage device port information as key-value pairs.
// CGI: storageDevice.cgi?action=factory.getPortInfo
func (s *StorageService) GetDiskInfo(ctx context.Context) (map[string]string, error) {
	body, err := s.client.cgiGet(ctx, "storageDevice.cgi", "factory.getPortInfo", nil)
	if err != nil {
		return nil, fmt.Errorf("storage GetDiskInfo: %w", err)
	}
	return parseKV(body), nil
}

// GetDeviceNames returns the raw body from the device collection query.
// CGI: storageDevice.cgi?action=factory.getCollect
func (s *StorageService) GetDeviceNames(ctx context.Context) (string, error) {
	body, err := s.client.cgiGet(ctx, "storageDevice.cgi", "factory.getCollect", nil)
	if err != nil {
		return "", fmt.Errorf("storage GetDeviceNames: %w", err)
	}
	return body, nil
}

// GetAllDeviceInfo returns the raw body with all storage device information.
// CGI: storageDevice.cgi?action=getDeviceAllInfo
func (s *StorageService) GetAllDeviceInfo(ctx context.Context) (string, error) {
	body, err := s.client.cgiGet(ctx, "storageDevice.cgi", "getDeviceAllInfo", nil)
	if err != nil {
		return "", fmt.Errorf("storage GetAllDeviceInfo: %w", err)
	}
	return body, nil
}

// GetCaps returns the storage capabilities as key-value pairs.
// CGI: storage.cgi?action=getCaps
func (s *StorageService) GetCaps(ctx context.Context) (map[string]string, error) {
	body, err := s.client.cgiGet(ctx, "storage.cgi", "getCaps", nil)
	if err != nil {
		return nil, fmt.Errorf("storage GetCaps: %w", err)
	}
	return parseKV(body), nil
}

// FormatSDCard formats the SD card at the given path (e.g., "/mnt/sd").
// CGI: storageDevice.cgi?action=setDevice&type=FormatPatition&path=PATH
// Note: "FormatPatition" is the camera's actual API spelling.
func (s *StorageService) FormatSDCard(ctx context.Context, path string) error {
	return s.client.cgiAction(ctx, "storageDevice.cgi", "setDevice", url.Values{
		"type": {"FormatPatition"},
		"path": {path},
	})
}

// GetNASConfig returns the NAS configuration as key-value pairs.
// CGI: configManager.cgi?action=getConfig&name=NAS
func (s *StorageService) GetNASConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "NAS")
}

// GetStorageGroupConfig returns the StorageGroup configuration as key-value pairs.
// CGI: configManager.cgi?action=getConfig&name=StorageGroup
func (s *StorageService) GetStorageGroupConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "StorageGroup")
}

// GetStorageHealthAlarm returns the StorageHealthAlarm configuration as key-value pairs.
// CGI: configManager.cgi?action=getConfig&name=StorageHealthAlarm
func (s *StorageService) GetStorageHealthAlarm(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "StorageHealthAlarm")
}
