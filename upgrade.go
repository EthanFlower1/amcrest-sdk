package amcrest

// UpgradeService handles firmware upgrade related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 213-218 (Section 5.3)
type UpgradeService struct {
	client *Client
}
