package amcrest

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// BuildingService handles building intercom related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 585-608 (Section 13)
type BuildingService struct {
	client *Client
}

// GetSIPConfig retrieves the SIP configuration.
func (s *BuildingService) GetSIPConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "SIPConfig")
}

// GetRoomNumberCount returns the number of VideoTalkContact records.
// recordFinder.cgi?action=getQuerySize&name=VideoTalkContact
func (s *BuildingService) GetRoomNumberCount(ctx context.Context) (int, error) {
	params := url.Values{
		"name": {"VideoTalkContact"},
	}
	body, err := s.client.cgiGet(ctx, "recordFinder.cgi", "getQuerySize", params)
	if err != nil {
		return 0, err
	}
	// Response format: count=N
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "count=") {
			val := strings.TrimPrefix(line, "count=")
			n, err := strconv.Atoi(strings.TrimSpace(val))
			if err != nil {
				return 0, fmt.Errorf("amcrest: parsing count: %w", err)
			}
			return n, nil
		}
	}
	return 0, fmt.Errorf("amcrest: count not found in response: %s", strings.TrimSpace(body))
}
