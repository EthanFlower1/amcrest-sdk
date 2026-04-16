package amcrest

import (
	"context"
	"fmt"
)

// PeopleService handles people counting related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 362-381 (Section 9.3-9.5)
type PeopleService struct {
	client *Client
}

// GetSummary returns the people counting summary via the VideoStatServer API.
// POST /cgi-bin/api/VideoStatServer/getSummary with an empty JSON body.
func (s *PeopleService) GetSummary(ctx context.Context) (string, error) {
	body, err := s.client.postRaw(ctx, "/cgi-bin/api/VideoStatServer/getSummary", struct{}{})
	if err != nil {
		return "", fmt.Errorf("PeopleService.GetSummary: %w", err)
	}
	return body, nil
}

// GetCrowdCaps returns the crowd distribution map channel capabilities.
// POST /cgi-bin/api/CrowdDistriMap/getChannelCaps with an empty JSON body.
func (s *PeopleService) GetCrowdCaps(ctx context.Context) (string, error) {
	body, err := s.client.postRaw(ctx, "/cgi-bin/api/CrowdDistriMap/getChannelCaps", struct{}{})
	if err != nil {
		return "", fmt.Errorf("PeopleService.GetCrowdCaps: %w", err)
	}
	return body, nil
}
