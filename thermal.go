package amcrest

import (
	"context"
	"fmt"
	"io"
	"net/url"
)

// ThermalService handles thermal camera related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 481-504 (Section 11)
type ThermalService struct {
	client *Client
}

// GetCaps retrieves thermal radiometry capabilities via
// RadiometryManager.cgi?action=getCaps. Returns parsed key-value pairs.
func (s *ThermalService) GetCaps(ctx context.Context) (map[string]string, error) {
	body, err := s.client.cgiGet(ctx, "RadiometryManager.cgi", "getCaps", nil)
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}

// GetThermographyOptions retrieves the ThermographyOptions configuration
// table without stripping key prefixes.
func (s *ThermalService) GetThermographyOptions(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "ThermographyOptions")
}

// GetRadiometryCaps retrieves radiometry capability details via
// RadiometryManager.cgi?action=getRadiometryCaps. Returns the raw response body.
func (s *ThermalService) GetRadiometryCaps(ctx context.Context) (string, error) {
	return s.client.cgiGet(ctx, "RadiometryManager.cgi", "getRadiometryCaps", nil)
}

// GetFireWarningConfig retrieves the FireWarning configuration table without
// stripping key prefixes.
func (s *ThermalService) GetFireWarningConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "FireWarning")
}

// SetThermographyOptions updates ThermographyOptions configuration values.
// Keys should be prefixed with "ThermographyOptions." (e.g.,
// "ThermographyOptions.Emissivity").
func (s *ThermalService) SetThermographyOptions(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetExternSystemInfo retrieves external system information via
// ThermographyManager.cgi?action=getExternSystemInfo.
func (s *ThermalService) GetExternSystemInfo(ctx context.Context) (string, error) {
	return s.client.cgiGet(ctx, "ThermographyManager.cgi", "getExternSystemInfo", nil)
}

// GetPresetModeInfo retrieves preset mode parameters via
// ThermographyManager.cgi?action=getPresetParam.
func (s *ThermalService) GetPresetModeInfo(ctx context.Context) (string, error) {
	return s.client.cgiGet(ctx, "ThermographyManager.cgi", "getPresetParam", nil)
}

// EnableShutter enables the thermal camera shutter via
// ThermographyManager.cgi?action=enableShutter.
func (s *ThermalService) EnableShutter(ctx context.Context) error {
	return s.client.cgiAction(ctx, "ThermographyManager.cgi", "enableShutter", nil)
}

// DoFFC performs a flat-field correction (FFC) via
// ThermographyManager.cgi?action=doFFC.
func (s *ThermalService) DoFFC(ctx context.Context) error {
	return s.client.cgiAction(ctx, "ThermographyManager.cgi", "doFFC", nil)
}

// GetThermometryConfig retrieves the HeatImagingThermometry configuration table
// without stripping key prefixes.
func (s *ThermalService) GetThermometryConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "HeatImagingThermometry")
}

// SetThermometryConfig updates HeatImagingThermometry configuration values.
// Keys should be prefixed with "HeatImagingThermometry.".
func (s *ThermalService) SetThermometryConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetThermometryRule retrieves the ThermometryRule configuration table without
// stripping key prefixes.
func (s *ThermalService) GetThermometryRule(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "ThermometryRule")
}

// SetThermometryRule updates ThermometryRule configuration values.
// Keys should be prefixed with "ThermometryRule.".
func (s *ThermalService) SetThermometryRule(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetTemperEventConfig retrieves the HeatImagingTemper configuration table
// without stripping key prefixes.
func (s *ThermalService) GetTemperEventConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "HeatImagingTemper")
}

// SetTemperEventConfig updates HeatImagingTemper configuration values.
// Keys should be prefixed with "HeatImagingTemper.".
func (s *ThermalService) SetTemperEventConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetPointTemperature retrieves a random point temperature reading via
// RadiometryManager.cgi?action=getRandomPointTemper. Additional parameters
// (e.g., coordinate info) are merged into the query.
func (s *ThermalService) GetPointTemperature(ctx context.Context, channel int, params map[string]string) (string, error) {
	v := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	for k, val := range params {
		v.Set(k, val)
	}
	return s.client.cgiGet(ctx, "RadiometryManager.cgi", "getRandomPointTemper", v)
}

// GetTemperature retrieves temperature data via
// RadiometryManager.cgi?action=getTemper. Additional parameters are merged
// into the query.
func (s *ThermalService) GetTemperature(ctx context.Context, channel int, params map[string]string) (string, error) {
	v := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	for k, val := range params {
		v.Set(k, val)
	}
	return s.client.cgiGet(ctx, "RadiometryManager.cgi", "getTemper", v)
}

// StartFindTemperature begins a temperature search session via
// RadiometryManager.cgi?action=startFind. Returns the raw response containing
// the search token.
func (s *ThermalService) StartFindTemperature(ctx context.Context, params map[string]string) (string, error) {
	v := url.Values{}
	for k, val := range params {
		v.Set(k, val)
	}
	return s.client.cgiGet(ctx, "RadiometryManager.cgi", "startFind", v)
}

// DoFindTemperature retrieves a batch of temperature search results via
// RadiometryManager.cgi?action=doFind using the given token and count.
func (s *ThermalService) DoFindTemperature(ctx context.Context, token string, count int) (string, error) {
	return s.client.cgiGet(ctx, "RadiometryManager.cgi", "doFind", url.Values{
		"token": {token},
		"count": {fmt.Sprintf("%d", count)},
	})
}

// StopFindTemperature ends a temperature search session via
// RadiometryManager.cgi?action=stopFind.
func (s *ThermalService) StopFindTemperature(ctx context.Context, token string) error {
	return s.client.cgiAction(ctx, "RadiometryManager.cgi", "stopFind", url.Values{
		"token": {token},
	})
}

// SetFireWarningConfig updates FireWarning configuration values.
// Keys should be prefixed with "FireWarning.".
func (s *ThermalService) SetFireWarningConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetFireWarningModeConfig retrieves the FireWarningMode configuration table
// without stripping key prefixes.
func (s *ThermalService) GetFireWarningModeConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "FireWarningMode")
}

// SetFireWarningModeConfig updates FireWarningMode configuration values.
// Keys should be prefixed with "FireWarningMode.".
func (s *ThermalService) SetFireWarningModeConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetCurrentHotColdSpot retrieves the current hot/cold spot data via
// TemperCorrection.cgi?action=getCurrentHotColdSpot. Returns parsed key-value
// pairs.
func (s *ThermalService) GetCurrentHotColdSpot(ctx context.Context, channel int) (map[string]string, error) {
	body, err := s.client.cgiGet(ctx, "TemperCorrection.cgi", "getCurrentHotColdSpot", url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	})
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}

// GetPreAlarmEventConfig retrieves the PreAlarmEvent configuration table
// without stripping key prefixes.
func (s *ThermalService) GetPreAlarmEventConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "PreAlarmEvent")
}

// SetPreAlarmEventConfig updates PreAlarmEvent configuration values.
// Keys should be prefixed with "PreAlarmEvent.".
func (s *ThermalService) SetPreAlarmEventConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetHeatMapInfo downloads raw heat-map binary data via
// RadiometryManager.cgi?action=getHeatMapsDirectly for the given channel.
func (s *ThermalService) GetHeatMapInfo(ctx context.Context, channel int) ([]byte, error) {
	params := url.Values{
		"action":  {"getHeatMapsDirectly"},
		"channel": {fmt.Sprintf("%d", channel)},
	}

	resp, err := s.client.get(ctx, "/cgi-bin/RadiometryManager.cgi", params)
	if err != nil {
		return nil, fmt.Errorf("thermal getHeatMapsDirectly: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, &APIError{StatusCode: resp.StatusCode}
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("thermal getHeatMapsDirectly: reading body: %w", err)
	}
	return data, nil
}

// SetEnvironmentTemp sets the environment temperature used for thermal
// correction via TemperCustom.cgi?action=setEnvTemp.
func (s *ThermalService) SetEnvironmentTemp(ctx context.Context, temp int) error {
	return s.client.cgiAction(ctx, "TemperCustom.cgi", "setEnvTemp", url.Values{
		"EnvironmentTemp": {fmt.Sprintf("%d", temp)},
	})
}
