package amcrest

import (
	"context"
	"fmt"
)

// WorkSuitService handles work suit detection related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 402-416 (Section 9.7)
type WorkSuitService struct {
	client *Client
}

// FindGroup returns the raw response listing all work suit comparison groups.
// POST /cgi-bin/api/WorkSuitCompareServer/findGroup with an empty JSON body.
func (s *WorkSuitService) FindGroup(ctx context.Context) (string, error) {
	body, err := s.client.postRaw(ctx, "/cgi-bin/api/WorkSuitCompareServer/findGroup", struct{}{})
	if err != nil {
		return "", fmt.Errorf("WorkSuitService.FindGroup: %w", err)
	}
	return body, nil
}

// GetGroup returns the work suit group assigned to the given video channel.
// POST /cgi-bin/api/WorkSuitCompareServer/getGroup with {"channel":N}.
func (s *WorkSuitService) GetGroup(ctx context.Context, channel int) (string, error) {
	reqBody := map[string]int{"channel": channel}
	body, err := s.client.postRaw(ctx, "/cgi-bin/api/WorkSuitCompareServer/getGroup", reqBody)
	if err != nil {
		return "", fmt.Errorf("WorkSuitService.GetGroup: %w", err)
	}
	return body, nil
}
