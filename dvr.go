package amcrest

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// DVRService handles DVR-related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 609-617 (Section 14)
type DVRService struct {
	client *Client
}

// StartFind begins a file search and returns the find ID.
// FileFindHelper.cgi?action=startFind&channel=N&startTime=X&endTime=Y
func (s *DVRService) StartFind(ctx context.Context, channel int, startTime, endTime string) (string, error) {
	params := url.Values{
		"channel":   {fmt.Sprintf("%d", channel)},
		"startTime": {startTime},
		"endTime":   {endTime},
	}
	body, err := s.client.cgiGet(ctx, "FileFindHelper.cgi", "startFind", params)
	if err != nil {
		return "", err
	}
	// Response format: findId=N
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "findId=") {
			return strings.TrimPrefix(line, "findId="), nil
		}
	}
	return "", fmt.Errorf("amcrest: findId not found in response: %s", strings.TrimSpace(body))
}

// FindNext retrieves the next batch of results for a file search.
// FileFindHelper.cgi?action=findNext&findId=X&count=N
func (s *DVRService) FindNext(ctx context.Context, findId string, count int) (string, error) {
	params := url.Values{
		"findId": {findId},
		"count":  {fmt.Sprintf("%d", count)},
	}
	return s.client.cgiGet(ctx, "FileFindHelper.cgi", "findNext", params)
}

// StopFind stops a file search.
// FileFindHelper.cgi?action=stopFind&findId=X
func (s *DVRService) StopFind(ctx context.Context, findId string) error {
	params := url.Values{
		"findId": {findId},
	}
	return s.client.cgiAction(ctx, "FileFindHelper.cgi", "stopFind", params)
}

// GetBandwidthLimit retrieves the bandwidth limit state.
// BandLimit.cgi?action=getLimitState
func (s *DVRService) GetBandwidthLimit(ctx context.Context) (map[string]string, error) {
	body, err := s.client.cgiGet(ctx, "BandLimit.cgi", "getLimitState", nil)
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}
