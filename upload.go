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

// GetReportUploadConfig returns the report HTTP upload configuration.
// Uses configManager getRawConfig with name ReportHttpUpload.
// PDF 4.13.5
func (s *UploadService) GetReportUploadConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "ReportHttpUpload")
}

// SetReportUploadConfig updates the report HTTP upload configuration.
// PDF 4.13.5
func (s *UploadService) SetReportUploadConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// Device-push upload format specifications (PDF 4.13.2, 4.13.4, 4.13.6-4.13.12):
//
// These sections define the format of HTTP requests that the camera POSTs to an
// external server. They are not endpoints the SDK calls on the camera.
//
//   4.13.2  - Picture upload format: the camera POSTs picture data as multipart/form-data
//             to the URL configured in PictureHttpUpload.
//   4.13.4  - Event upload format: the camera POSTs event JSON/text to the URL configured
//             in EventHttpUpload.
//   4.13.6  - Face snapshot push format: face detection image data pushed to server.
//   4.13.7  - Face recognition push format: face recognition result data pushed to server.
//   4.13.8  - People counting push format: people counting statistics pushed to server.
//   4.13.9  - Heat map push format: heat map data pushed to server.
//   4.13.10 - Crowd distribution push format: crowd distribution statistics pushed to server.
//   4.13.11 - Vehicle detection push format: vehicle detection data pushed to server.
//   4.13.12 - Work suit detection push format: work suit detection data pushed to server.
