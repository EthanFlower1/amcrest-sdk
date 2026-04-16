package amcrest

import (
	"context"
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
