package amcrest

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
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

// GetLightConfig retrieves the parking space light state configuration.
func (s *ParkingService) GetLightConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "ParkingSpaceLightState")
}

// SetLightConfig updates the parking space light state configuration.
func (s *ParkingService) SetLightConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// SetOrderState sets the order state for a specific parking space.
func (s *ParkingService) SetOrderState(ctx context.Context, channel, spaceNo, state int) error {
	return s.client.cgiAction(ctx, "trafficParking.cgi", "setOrderState", url.Values{
		"channel": {strconv.Itoa(channel)},
		"spaceNo": {strconv.Itoa(spaceNo)},
		"state":   {strconv.Itoa(state)},
	})
}

// SetLightState sets the light state for a specific parking space with
// additional parameters.
func (s *ParkingService) SetLightState(ctx context.Context, channel, spaceNo int, params map[string]string) error {
	v := url.Values{
		"channel": {strconv.Itoa(channel)},
		"spaceNo": {strconv.Itoa(spaceNo)},
	}
	for k, val := range params {
		v.Set(k, val)
	}
	return s.client.cgiAction(ctx, "trafficParking.cgi", "setLightState", v)
}

// GetAccessFilter retrieves the parking space access filter configuration.
func (s *ParkingService) GetAccessFilter(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "ParkingSpaceAccessFilter")
}

// SetAccessFilter updates the parking space access filter configuration.
func (s *ParkingService) SetAccessFilter(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// SetOverLineState enables or disables the over-line state for a specific
// parking space.
func (s *ParkingService) SetOverLineState(ctx context.Context, channel, spaceNo int, enable bool) error {
	return s.client.cgiAction(ctx, "trafficParking.cgi", "setOverLineState", url.Values{
		"channel": {strconv.Itoa(channel)},
		"spaceNo": {strconv.Itoa(spaceNo)},
		"enable":  {strconv.FormatBool(enable)},
	})
}
