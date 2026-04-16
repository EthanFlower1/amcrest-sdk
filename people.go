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

// StartFindCount starts a people counting search. The body should contain the
// search parameters (e.g., channel, startTime, endTime). Returns a token.
// PDF 9.3.2: POST /cgi-bin/api/VideoStatServer/startFind
func (s *PeopleService) StartFindCount(ctx context.Context, body interface{}) (string, error) {
	resp, err := s.client.postRaw(ctx, "/cgi-bin/api/VideoStatServer/startFind", body)
	if err != nil {
		return "", fmt.Errorf("PeopleService.StartFindCount: %w", err)
	}
	return resp, nil
}

// DoFindCount retrieves a page of people counting search results.
// PDF 9.3.2: POST /cgi-bin/api/VideoStatServer/doFind
func (s *PeopleService) DoFindCount(ctx context.Context, body interface{}) (string, error) {
	resp, err := s.client.postRaw(ctx, "/cgi-bin/api/VideoStatServer/doFind", body)
	if err != nil {
		return "", fmt.Errorf("PeopleService.DoFindCount: %w", err)
	}
	return resp, nil
}

// StopFindCount stops a people counting search and releases the token.
// PDF 9.3.2: POST /cgi-bin/api/VideoStatServer/stopFind
func (s *PeopleService) StopFindCount(ctx context.Context, token int) error {
	reqBody := map[string]int{"token": token}
	return s.client.postJSON(ctx, "/cgi-bin/api/VideoStatServer/stopFind", reqBody, nil)
}

// ClearCount clears the people counting statistics for the given channel.
// PDF 9.3.3: POST /cgi-bin/api/VideoStatServer/clearCount
func (s *PeopleService) ClearCount(ctx context.Context, channel int) error {
	reqBody := map[string]int{"channel": channel}
	return s.client.postJSON(ctx, "/cgi-bin/api/VideoStatServer/clearCount", reqBody, nil)
}

// SubscribeCount subscribes to real-time people counting statistics.
// PDF 9.3.4: POST /cgi-bin/api/VideoStatServer/subscribeStat
func (s *PeopleService) SubscribeCount(ctx context.Context) (string, error) {
	body, err := s.client.postRaw(ctx, "/cgi-bin/api/VideoStatServer/subscribeStat", struct{}{})
	if err != nil {
		return "", fmt.Errorf("PeopleService.SubscribeCount: %w", err)
	}
	return body, nil
}

// GetHeatMap returns heat map statistics.
// PDF 9.4.1: POST /cgi-bin/api/HeatMapStat/getHeatMapInfo
func (s *PeopleService) GetHeatMap(ctx context.Context, body interface{}) (string, error) {
	resp, err := s.client.postRaw(ctx, "/cgi-bin/api/HeatMapStat/getHeatMapInfo", body)
	if err != nil {
		return "", fmt.Errorf("PeopleService.GetHeatMap: %w", err)
	}
	return resp, nil
}

// GetPeopleHeatMap returns people counting heat map data.
// PDF 9.4.2: POST /cgi-bin/api/VideoStatServer/getHeatMap
func (s *PeopleService) GetPeopleHeatMap(ctx context.Context, body interface{}) (string, error) {
	resp, err := s.client.postRaw(ctx, "/cgi-bin/api/VideoStatServer/getHeatMap", body)
	if err != nil {
		return "", fmt.Errorf("PeopleService.GetPeopleHeatMap: %w", err)
	}
	return resp, nil
}

// GetCurrentCrowdStat returns the current crowd distribution statistics.
// PDF 9.5.3: POST /cgi-bin/api/CrowdDistriMap/getCurrentStat
func (s *PeopleService) GetCurrentCrowdStat(ctx context.Context, body interface{}) (string, error) {
	resp, err := s.client.postRaw(ctx, "/cgi-bin/api/CrowdDistriMap/getCurrentStat", body)
	if err != nil {
		return "", fmt.Errorf("PeopleService.GetCurrentCrowdStat: %w", err)
	}
	return resp, nil
}
