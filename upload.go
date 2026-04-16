package amcrest

import "context"

// UploadService handles picture and event upload configuration.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 219-234 (Section 4.13)
type UploadService struct {
	client *Client
}

// GetPictureUploadConfig returns the picture HTTP upload configuration.
// Uses configManager getRawConfig with name PictureHttpUpload.
func (s *UploadService) GetPictureUploadConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "PictureHttpUpload")
}

// SetPictureUploadConfig updates the picture HTTP upload configuration.
func (s *UploadService) SetPictureUploadConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetEventUploadConfig returns the event HTTP upload configuration.
// Uses configManager getRawConfig with name EventHttpUpload.
func (s *UploadService) GetEventUploadConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "EventHttpUpload")
}

// SetEventUploadConfig updates the event HTTP upload configuration.
func (s *UploadService) SetEventUploadConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}
