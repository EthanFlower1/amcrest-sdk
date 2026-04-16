package amcrest

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// PeripheralService handles peripheral device related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 305-325, 618-650 (Sections 6.2, 14.1)
type PeripheralService struct {
	client *Client
}

// WiperStart starts continuous wiper movement on the given channel.
// rainBrush.cgi?action=moveContinuously&channel=N&interval=I
func (s *PeripheralService) WiperStart(ctx context.Context, channel, interval int) error {
	params := url.Values{
		"channel":  {fmt.Sprintf("%d", channel)},
		"interval": {fmt.Sprintf("%d", interval)},
	}
	return s.client.cgiAction(ctx, "rainBrush.cgi", "moveContinuously", params)
}

// WiperStop stops wiper movement on the given channel.
// rainBrush.cgi?action=stopMove&channel=N
func (s *PeripheralService) WiperStop(ctx context.Context, channel int) error {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	return s.client.cgiAction(ctx, "rainBrush.cgi", "stopMove", params)
}

// WiperOnce performs a single wiper sweep on the given channel.
// rainBrush.cgi?action=moveOnce&channel=N
func (s *PeripheralService) WiperOnce(ctx context.Context, channel int) error {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	return s.client.cgiAction(ctx, "rainBrush.cgi", "moveOnce", params)
}

// ControlCoaxialIO sends a coaxial control IO command.
// coaxialControlIO.cgi?action=control&channel=N&info[0].Type=T&info[0].IO=I&info[0].TriggerMode=M
// Uses a raw query string to preserve bracket characters.
func (s *PeripheralService) ControlCoaxialIO(ctx context.Context, channel int, ioType, io, triggerMode int) error {
	rawQuery := fmt.Sprintf(
		"action=control&channel=%d&info[0].Type=%d&info[0].IO=%d&info[0].TriggerMode=%d",
		channel, ioType, io, triggerMode,
	)
	u := s.client.baseURL + "/cgi-bin/coaxialControlIO.cgi?" + rawQuery

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return fmt.Errorf("amcrest: creating request: %w", err)
	}
	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("amcrest: executing request: %w", err)
	}
	return checkOK(resp)
}

// GetCoaxialIOStatus retrieves the coaxial IO status for the given channel.
// coaxialControlIO.cgi?action=getstatus&channel=N
func (s *PeripheralService) GetCoaxialIOStatus(ctx context.Context, channel int) (map[string]string, error) {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	body, err := s.client.cgiGet(ctx, "coaxialControlIO.cgi", "getStatus", params)
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}

// GetFlashlightConfig retrieves the FlashLight configuration.
func (s *PeripheralService) GetFlashlightConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "FlashLight")
}

// GetGPSConfig retrieves the GPS configuration.
func (s *PeripheralService) GetGPSConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "GPS")
}

// GetFishEyeConfig retrieves the FishEye configuration.
func (s *PeripheralService) GetFishEyeConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "FishEye")
}

// SetFlashlightConfig updates the FlashLight configuration.
func (s *PeripheralService) SetFlashlightConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetPIRConfig retrieves the PIR alarm configuration for the given channel.
// pirAlarm.cgi?action=getPirParam&channel=N
func (s *PeripheralService) GetPIRConfig(ctx context.Context, channel int) (string, error) {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	return s.client.cgiGet(ctx, "pirAlarm.cgi", "getPirParam", params)
}

// SetPIRConfig updates the PIR alarm configuration for the given channel.
// pirAlarm.cgi?action=setPirParam&channel=N&...
func (s *PeripheralService) SetPIRConfig(ctx context.Context, channel int, params map[string]string) error {
	vals := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	for k, v := range params {
		vals.Set(k, v)
	}
	return s.client.cgiAction(ctx, "pirAlarm.cgi", "setPirParam", vals)
}

// SCADAGetAttribute retrieves SCADA attributes.
func (s *PeripheralService) SCADAGetAttribute(ctx context.Context, body interface{}) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/SCADA/getAttribute", body)
}

// SCADASetAttribute updates SCADA attributes.
func (s *PeripheralService) SCADASetAttribute(ctx context.Context, body interface{}) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/SCADA/setAttribute", body, nil)
}

// SCADAGetData retrieves SCADA data.
func (s *PeripheralService) SCADAGetData(ctx context.Context, body interface{}) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/SCADA/get", body)
}

// SCADASetData updates SCADA data.
func (s *PeripheralService) SCADASetData(ctx context.Context, body interface{}) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/SCADA/set", body, nil)
}

// SCADAStartFind starts a SCADA find operation.
func (s *PeripheralService) SCADAStartFind(ctx context.Context, body interface{}) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/SCADA/startFind", body)
}

// SCADADoFind performs a SCADA find operation.
func (s *PeripheralService) SCADADoFind(ctx context.Context, body interface{}) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/SCADA/doFind", body)
}

// SCADAStopFind stops a SCADA find operation by token.
func (s *PeripheralService) SCADAStopFind(ctx context.Context, token int) error {
	body := map[string]interface{}{"token": token}
	return s.client.postJSON(ctx, "/cgi-bin/api/SCADA/stopFind", body, nil)
}

// SCADAGetDeviceList retrieves the list of SCADA devices.
func (s *PeripheralService) SCADAGetDeviceList(ctx context.Context) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/SCADA/getDeviceList", map[string]interface{}{})
}

// GetGyroData retrieves gyroscope data.
func (s *PeripheralService) GetGyroData(ctx context.Context) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/Gyro/getData", map[string]interface{}{})
}

// SetGPSConfig updates the GPS configuration.
func (s *PeripheralService) SetGPSConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetGPSStatus retrieves the current GPS status.
// GpsControl.cgi?action=getStatus
func (s *PeripheralService) GetGPSStatus(ctx context.Context) (string, error) {
	return s.client.cgiGet(ctx, "GpsControl.cgi", "getStatus", nil)
}

// GetGPSCaps retrieves GPS capabilities.
// GpsControl.cgi?action=getCaps
func (s *PeripheralService) GetGPSCaps(ctx context.Context) (string, error) {
	return s.client.cgiGet(ctx, "GpsControl.cgi", "getCaps", nil)
}

// SetFishEyeConfig updates the FishEye configuration.
func (s *PeripheralService) SetFishEyeConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetFishEyeCaps retrieves FishEye capabilities for the given channel.
// devVideoInput.cgi?action=getCapsEx&name=VideoInFishEye&channel=N
func (s *PeripheralService) GetFishEyeCaps(ctx context.Context, channel int) (string, error) {
	params := url.Values{
		"name":    {"VideoInFishEye"},
		"channel": {fmt.Sprintf("%d", channel)},
	}
	return s.client.cgiGet(ctx, "devVideoInput.cgi", "getCapsEx", params)
}

// SetIlluminatorConfig updates the SignLight (illuminator) configuration.
func (s *PeripheralService) SetIlluminatorConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetLensCaps retrieves lens function capabilities.
// LensFunc.cgi?action=getCaps
func (s *PeripheralService) GetLensCaps(ctx context.Context) (string, error) {
	return s.client.cgiGet(ctx, "LensFunc.cgi", "getCaps", nil)
}

// AdjustAngleContinuously starts continuous angle adjustment on the given channel.
// LensFunc.cgi?action=adjustAngleContinuously&channel=N&...
func (s *PeripheralService) AdjustAngleContinuously(ctx context.Context, channel int, params map[string]string) error {
	vals := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	for k, v := range params {
		vals.Set(k, v)
	}
	return s.client.cgiAction(ctx, "LensFunc.cgi", "adjustAngleContinuously", vals)
}

// StopAdjustAngle stops angle adjustment on the given channel.
// LensFunc.cgi?action=stopAdjustAngle&channel=N
func (s *PeripheralService) StopAdjustAngle(ctx context.Context, channel int) error {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	return s.client.cgiAction(ctx, "LensFunc.cgi", "stopAdjustAngle", params)
}

// AdjustDepthField adjusts the depth of field on the given channel.
// LensFunc.cgi?action=adjustDepthField&channel=N&focus=F&zoom=Z
func (s *PeripheralService) AdjustDepthField(ctx context.Context, channel, focus, zoom int) error {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
		"focus":   {fmt.Sprintf("%d", focus)},
		"zoom":    {fmt.Sprintf("%d", zoom)},
	}
	return s.client.cgiAction(ctx, "LensFunc.cgi", "adjustDepthField", params)
}

// AdjustDepthFieldContinuously starts continuous depth of field adjustment on the given channel.
// LensFunc.cgi?action=adjustDepthFieldContinuously&channel=N&focus=F&zoom=Z
func (s *PeripheralService) AdjustDepthFieldContinuously(ctx context.Context, channel, focus, zoom int) error {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
		"focus":   {fmt.Sprintf("%d", focus)},
		"zoom":    {fmt.Sprintf("%d", zoom)},
	}
	return s.client.cgiAction(ctx, "LensFunc.cgi", "adjustDepthFieldContinuously", params)
}

// GetDepthFieldStatus retrieves the current depth of field status.
// LensFunc.cgi?action=getDepthFieldStatus
func (s *PeripheralService) GetDepthFieldStatus(ctx context.Context) (map[string]string, error) {
	body, err := s.client.cgiGet(ctx, "LensFunc.cgi", "getDepthFieldStatus", nil)
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}

// AutoAdjustDepthField triggers automatic depth of field adjustment.
// LensFunc.cgi?action=autoAdjustDepthField
func (s *PeripheralService) AutoAdjustDepthField(ctx context.Context) error {
	return s.client.cgiAction(ctx, "LensFunc.cgi", "autoAdjustDepthField", nil)
}

// ResetAngle resets the lens angle to the default position.
// LensManager.cgi?action=resetAngle
func (s *PeripheralService) ResetAngle(ctx context.Context) error {
	return s.client.cgiAction(ctx, "LensManager.cgi", "resetAngle", nil)
}

// GetRadarCaps retrieves radar capabilities.
// radarAdaptor.cgi?action=getCaps
func (s *PeripheralService) GetRadarCaps(ctx context.Context) (string, error) {
	return s.client.cgiGet(ctx, "radarAdaptor.cgi", "getCaps", nil)
}

// GetRadarCapsEx retrieves extended radar capabilities.
// radarAdaptor.cgi?action=getCapsEx
func (s *PeripheralService) GetRadarCapsEx(ctx context.Context) (string, error) {
	return s.client.cgiGet(ctx, "radarAdaptor.cgi", "getCapsEx", nil)
}

// GetRadarStatus retrieves the current radar status.
// radarAdaptor.cgi?action=getStatus
func (s *PeripheralService) GetRadarStatus(ctx context.Context) (string, error) {
	return s.client.cgiGet(ctx, "radarAdaptor.cgi", "getStatus", nil)
}

// CalculateRealSize calculates the real size using radar data.
// radarAdaptor.cgi?action=calculateRealSize&...
func (s *PeripheralService) CalculateRealSize(ctx context.Context, params map[string]string) (string, error) {
	vals := url.Values{}
	for k, v := range params {
		vals.Set(k, v)
	}
	return s.client.cgiGet(ctx, "radarAdaptor.cgi", "calculateRealSize", vals)
}

// StartRadarCalibration starts radar calibration.
// radarAdaptor.cgi?action=startCalibration&...
func (s *PeripheralService) StartRadarCalibration(ctx context.Context, params map[string]string) error {
	vals := url.Values{}
	for k, v := range params {
		vals.Set(k, v)
	}
	return s.client.cgiAction(ctx, "radarAdaptor.cgi", "startCalibration", vals)
}

// GetWaterRadarCaps retrieves water radar capabilities.
func (s *PeripheralService) GetWaterRadarCaps(ctx context.Context) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/WaterRadar/getCaps", map[string]interface{}{})
}

// GetWaterRadarData retrieves water radar object information.
func (s *PeripheralService) GetWaterRadarData(ctx context.Context) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/WaterRadar/getObjectInfo", map[string]interface{}{})
}

// GetWaterQualityCaps retrieves water quality monitoring capabilities.
func (s *PeripheralService) GetWaterQualityCaps(ctx context.Context) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/WaterDataStatServer/getCaps", map[string]interface{}{})
}

// GetWaterQualityData retrieves water quality data.
func (s *PeripheralService) GetWaterQualityData(ctx context.Context) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/WaterDataStatServer/getWaterData", map[string]interface{}{})
}

// StartWaterQualityFind starts a water quality data find operation.
func (s *PeripheralService) StartWaterQualityFind(ctx context.Context, body interface{}) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/WaterDataStatServer/startFind", body)
}

// DoWaterQualityFind performs a water quality data find operation.
func (s *PeripheralService) DoWaterQualityFind(ctx context.Context, body interface{}) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/WaterDataStatServer/doFind", body)
}

// StopWaterQualityFind stops a water quality data find operation by token.
func (s *PeripheralService) StopWaterQualityFind(ctx context.Context, token int) error {
	body := map[string]interface{}{"token": token}
	return s.client.postJSON(ctx, "/cgi-bin/api/WaterDataStatServer/stopFind", body, nil)
}

// ChangeAdStayTime changes the advertisement display sustain time.
// VideoOutput.cgi?action=changeSustain&Sustain=N
func (s *PeripheralService) ChangeAdStayTime(ctx context.Context, sustain int) error {
	params := url.Values{
		"Sustain": {fmt.Sprintf("%d", sustain)},
	}
	return s.client.cgiAction(ctx, "VideoOutput.cgi", "changeSustain", params)
}

// QueryAdFiles queries the list of delivered advertisement files.
// VideoOutput.cgi?action=queryDeliveredFile
func (s *PeripheralService) QueryAdFiles(ctx context.Context) (string, error) {
	return s.client.cgiGet(ctx, "VideoOutput.cgi", "queryDeliveredFile", nil)
}

// DiscoverDevices triggers device discovery.
// deviceDiscovery.cgi?action=attach
func (s *PeripheralService) DiscoverDevices(ctx context.Context) (string, error) {
	return s.client.cgiGet(ctx, "deviceDiscovery.cgi", "attach", nil)
}
