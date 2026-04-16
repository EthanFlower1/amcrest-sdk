package amcrest

// PTZService handles PTZ (pan-tilt-zoom) related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 284-305 (Section 6.1)
type PTZService struct {
	client *Client
}
