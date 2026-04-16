package amcrest

import (
	"context"
	"fmt"
)

// ParkingService handles parking-related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 460-465 (Section 10.5)
type ParkingService struct {
	client *Client
}

// GetSpaceStatus retrieves the status of a specific parking space via a JSON
// POST to /cgi-bin/api/TrafficParking/getSpaceStatus. The raw JSON response
// body is returned as a string.
func (s *ParkingService) GetSpaceStatus(ctx context.Context, channel, spaceNo int) (string, error) {
	body := map[string]interface{}{
		"channel": channel,
		"spaceNo": spaceNo,
	}
	var result interface{}
	err := s.client.postJSON(ctx, "/cgi-bin/api/TrafficParking/getSpaceStatus", body, &result)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", result), nil
}

// GetAllSpaceStatus retrieves the status of all parking spaces on a channel
// via a JSON POST to /cgi-bin/api/TrafficParking/getAllSpaceStatus. The raw
// JSON response body is returned as a string.
func (s *ParkingService) GetAllSpaceStatus(ctx context.Context, channel int) (string, error) {
	body := map[string]interface{}{
		"channel": channel,
	}
	var result interface{}
	err := s.client.postJSON(ctx, "/cgi-bin/api/TrafficParking/getAllSpaceStatus", body, &result)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", result), nil
}
