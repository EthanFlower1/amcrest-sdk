package amcrest

// MotionService handles motion detection related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 99-107, 416-419 (Sections 3.5, 8.3)
type MotionService struct {
	client *Client
}
