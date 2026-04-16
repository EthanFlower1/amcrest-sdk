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

// ---------------------------------------------------------------------------
// Section 15.6.11-15.6.19 - Radar configuration
// ---------------------------------------------------------------------------

// GetMapParaConfig retrieves the MapPara configuration.
func (s *PeripheralService) GetMapParaConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "MapPara")
}

// SetMapParaConfig updates the MapPara configuration.
func (s *PeripheralService) SetMapParaConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetRadarAnalyseRuleConfig retrieves the RadarAnalyseRule configuration.
func (s *PeripheralService) GetRadarAnalyseRuleConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "RadarAnalyseRule")
}

// SetRadarAnalyseRuleConfig updates the RadarAnalyseRule configuration.
func (s *PeripheralService) SetRadarAnalyseRuleConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetRadarCalibrationConfig retrieves the RadarCalibration configuration.
func (s *PeripheralService) GetRadarCalibrationConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "RadarCalibration")
}

// SetRadarCalibrationConfig updates the RadarCalibration configuration.
func (s *PeripheralService) SetRadarCalibrationConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetRadarGuardLineConfig retrieves the RadarGuardLine configuration.
func (s *PeripheralService) GetRadarGuardLineConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "RadarGuardLine")
}

// SetRadarGuardLineConfig updates the RadarGuardLine configuration.
func (s *PeripheralService) SetRadarGuardLineConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetRadarLinkConfig retrieves the RadarLink configuration.
func (s *PeripheralService) GetRadarLinkConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "RadarLink")
}

// SetRadarLinkConfig updates the RadarLink configuration.
func (s *PeripheralService) SetRadarLinkConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetRadarLinkDeviceConfig retrieves the RadarLinkDevice configuration.
func (s *PeripheralService) GetRadarLinkDeviceConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "RadarLinkDevice")
}

// SetRadarLinkDeviceConfig updates the RadarLinkDevice configuration.
func (s *PeripheralService) SetRadarLinkDeviceConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetRadarParaConfig retrieves the RadarPara configuration.
func (s *PeripheralService) GetRadarParaConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "RadarPara")
}

// SetRadarParaConfig updates the RadarPara configuration.
func (s *PeripheralService) SetRadarParaConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetRadarTrackGlobalConfig retrieves the RadarTrackGlobal configuration.
func (s *PeripheralService) GetRadarTrackGlobalConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "RadarTrackGlobal")
}

// SetRadarTrackGlobalConfig updates the RadarTrackGlobal configuration.
func (s *PeripheralService) SetRadarTrackGlobalConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetRemoteSDLinkConfig retrieves the RemoteSDLink configuration.
func (s *PeripheralService) GetRemoteSDLinkConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "RemoteSDLink")
}

// SetRemoteSDLinkConfig updates the RemoteSDLink configuration.
func (s *PeripheralService) SetRemoteSDLinkConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// ---------------------------------------------------------------------------
// Section 15.6 - Radar operations
// ---------------------------------------------------------------------------

// SubscribeRadarAlarm subscribes to radar alarm point information on the given channel.
// radarAdaptor.cgi?action=attachAlarmPointInfo&channel=N
func (s *PeripheralService) SubscribeRadarAlarm(ctx context.Context, channel int) (string, error) {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	return s.client.cgiGet(ctx, "radarAdaptor.cgi", "attachAlarmPointInfo", params)
}

// ManualLocate triggers a manual radar locate operation.
// radarAdaptor.cgi?action=manualLocate&...
func (s *PeripheralService) ManualLocate(ctx context.Context, params map[string]string) error {
	v := url.Values{}
	for k, val := range params {
		v.Set(k, val)
	}
	return s.client.cgiAction(ctx, "radarAdaptor.cgi", "manualLocate", v)
}

// AddRadarLinkSD adds a radar-linked SD device.
// radarAdaptor.cgi?action=addRadarLinkSD&...
func (s *PeripheralService) AddRadarLinkSD(ctx context.Context, params map[string]string) error {
	v := url.Values{}
	for k, val := range params {
		v.Set(k, val)
	}
	return s.client.cgiAction(ctx, "radarAdaptor.cgi", "addRadarLinkSD", v)
}

// DelRadarLinkSD removes a radar-linked SD device.
// radarAdaptor.cgi?action=delRadarLinkSD&...
func (s *PeripheralService) DelRadarLinkSD(ctx context.Context, params map[string]string) error {
	v := url.Values{}
	for k, val := range params {
		v.Set(k, val)
	}
	return s.client.cgiAction(ctx, "radarAdaptor.cgi", "delRadarLinkSD", v)
}

// GetRadarLinkSDState retrieves the radar-linked SD device state.
// radarAdaptor.cgi?action=getLinkSDState
func (s *PeripheralService) GetRadarLinkSDState(ctx context.Context) (string, error) {
	return s.client.cgiGet(ctx, "radarAdaptor.cgi", "getLinkSDState", nil)
}

// ---------------------------------------------------------------------------
// Section 15.2 - Open Platform
// ---------------------------------------------------------------------------

// StartApp starts an installed application.
// installManager.cgi?action=start&name=AppName
func (s *PeripheralService) StartApp(ctx context.Context, appName string) error {
	params := url.Values{
		"name": {appName},
	}
	return s.client.cgiAction(ctx, "installManager.cgi", "start", params)
}

// StopApp stops a running application.
// installManager.cgi?action=stop&name=AppName
func (s *PeripheralService) StopApp(ctx context.Context, appName string) error {
	params := url.Values{
		"name": {appName},
	}
	return s.client.cgiAction(ctx, "installManager.cgi", "stop", params)
}

// UninstallApp uninstalls an application.
// dhop.cgi?action=uninstall&name=AppName
func (s *PeripheralService) UninstallApp(ctx context.Context, appName string) error {
	params := url.Values{
		"name": {appName},
	}
	return s.client.cgiAction(ctx, "dhop.cgi", "uninstall", params)
}

// ---------------------------------------------------------------------------
// Section 15.4.8 - Scene Correction
// ---------------------------------------------------------------------------

// CorrectScene performs scene correction via the lens function API.
// POST /cgi-bin/api/LensFunc/correctScene
func (s *PeripheralService) CorrectScene(ctx context.Context, body interface{}) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/LensFunc/correctScene", body, nil)
}
