package amcrest

import "context"

// MotionService handles motion detection related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 99-107, 416-419 (Sections 3.5, 8.3)
type MotionService struct {
	client *Client
}

// GetConfig returns the MotionDetect configuration as a raw key-value map.
// CGI: configManager.cgi?action=getConfig&name=MotionDetect
func (s *MotionService) GetConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "MotionDetect")
}

// SetConfig updates MotionDetect configuration values.
// Keys should include the full "MotionDetect" prefix, e.g. "MotionDetect[0].Enable".
// CGI: configManager.cgi?action=setConfig
func (s *MotionService) SetConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetSmartMotionConfig returns the SmartMotionDetect configuration as a raw key-value map.
// CGI: configManager.cgi?action=getConfig&name=SmartMotionDetect
func (s *MotionService) GetSmartMotionConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "SmartMotionDetect")
}

// SetSmartMotionConfig updates SmartMotionDetect configuration values.
// Keys should include the full "SmartMotionDetect" prefix, e.g. "SmartMotionDetect[0].Enable".
// CGI: configManager.cgi?action=setConfig
func (s *MotionService) SetSmartMotionConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}
