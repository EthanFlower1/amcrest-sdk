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

// GetDiskRecordType returns the SupportDiskRecordType configuration as key-value pairs.
// CGI: configManager.cgi?action=getConfig&name=SupportDiskRecordType
func (s *StorageService) GetDiskRecordType(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "SupportDiskRecordType")
}

// SetDiskRecordType updates SupportDiskRecordType configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *StorageService) SetDiskRecordType(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetDetailedDiskInfo returns detailed device information for the given volume.
// API: POST /api/StorageDeviceManager/getDeviceInfos
func (s *StorageService) GetDetailedDiskInfo(ctx context.Context, volume string) (string, error) {
	body := map[string]interface{}{
		"volume": volume,
	}
	return s.client.postRaw(ctx, "/api/StorageDeviceManager/getDeviceInfos", body)
}

// SetNASConfig updates NAS configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *StorageService) SetNASConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetRecordStoragePoint returns the RecordStoragePoint configuration as key-value pairs.
// CGI: configManager.cgi?action=getConfig&name=RecordStoragePoint
func (s *StorageService) GetRecordStoragePoint(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "RecordStoragePoint")
}

// SetRecordStoragePoint updates RecordStoragePoint configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *StorageService) SetRecordStoragePoint(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// SetStorageGroupConfig updates StorageGroup configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *StorageService) SetStorageGroupConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// EncryptSDCard encrypts the SD card with the given password.
// CGI: SDEncrypt.cgi?action=encrypt&deviceName=X&password=Y
func (s *StorageService) EncryptSDCard(ctx context.Context, deviceName, password string) error {
	return s.client.cgiAction(ctx, "SDEncrypt.cgi", "encrypt", url.Values{
		"deviceName": {deviceName},
		"password":   {password},
	})
}

// DecryptSDCard decrypts the SD card with the given password.
// CGI: SDEncrypt.cgi?action=decrypt&deviceName=X&password=Y
func (s *StorageService) DecryptSDCard(ctx context.Context, deviceName, password string) error {
	return s.client.cgiAction(ctx, "SDEncrypt.cgi", "decrypt", url.Values{
		"deviceName": {deviceName},
		"password":   {password},
	})
}

// ClearSDCardPassword clears the encryption password for the SD card.
// CGI: SDEncrypt.cgi?action=clearPassword&deviceName=X&password=Y
func (s *StorageService) ClearSDCardPassword(ctx context.Context, deviceName, password string) error {
	return s.client.cgiAction(ctx, "SDEncrypt.cgi", "clearPassword", url.Values{
		"deviceName": {deviceName},
		"password":   {password},
	})
}

// ModifySDCardPassword changes the encryption password for the SD card.
// CGI: SDEncrypt.cgi?action=modifyPassword&deviceName=X&oldPassword=Y&newPassword=Z
func (s *StorageService) ModifySDCardPassword(ctx context.Context, deviceName, oldPassword, newPassword string) error {
	return s.client.cgiAction(ctx, "SDEncrypt.cgi", "modifyPassword", url.Values{
		"deviceName":  {deviceName},
		"oldPassword": {oldPassword},
		"newPassword": {newPassword},
	})
}

// GetSDCardErrorPolicy returns the operate error policy for the SD card.
// CGI: SDEncrypt.cgi?action=getOperateErrorPolicy&deviceName=X&operate=Y
func (s *StorageService) GetSDCardErrorPolicy(ctx context.Context, deviceName, operate string) (map[string]string, error) {
	body, err := s.client.cgiGet(ctx, "SDEncrypt.cgi", "getOperateErrorPolicy", url.Values{
		"deviceName": {deviceName},
		"operate":    {operate},
	})
	if err != nil {
		return nil, fmt.Errorf("storage GetSDCardErrorPolicy: %w", err)
	}
	return parseKV(body), nil
}

// SetStorageHealthAlarm updates StorageHealthAlarm configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *StorageService) SetStorageHealthAlarm(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}
