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

// CreateGroup creates a new work suit comparison group.
// PDF 9.7.1: POST /cgi-bin/api/WorkSuitCompareServer/createGroup
func (s *WorkSuitService) CreateGroup(ctx context.Context, body interface{}) (string, error) {
	resp, err := s.client.postRaw(ctx, "/cgi-bin/api/WorkSuitCompareServer/createGroup", body)
	if err != nil {
		return "", fmt.Errorf("WorkSuitService.CreateGroup: %w", err)
	}
	return resp, nil
}

// DeleteGroup deletes a work suit comparison group by ID.
// PDF 9.7.2: POST /cgi-bin/api/WorkSuitCompareServer/deleteGroup
func (s *WorkSuitService) DeleteGroup(ctx context.Context, groupID string) error {
	reqBody := map[string]string{"groupID": groupID}
	return s.client.postJSON(ctx, "/cgi-bin/api/WorkSuitCompareServer/deleteGroup", reqBody, nil)
}

// ModifyGroup modifies an existing work suit comparison group.
// PDF 9.7.5: POST /cgi-bin/api/WorkSuitCompareServer/modifyGroup
func (s *WorkSuitService) ModifyGroup(ctx context.Context, body interface{}) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/WorkSuitCompareServer/modifyGroup", body, nil)
}

// DeployGroup deploys (assigns) a work suit comparison group to channels.
// PDF 9.7.6: POST /cgi-bin/api/WorkSuitCompareServer/setGroup
func (s *WorkSuitService) DeployGroup(ctx context.Context, body interface{}) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/WorkSuitCompareServer/setGroup", body, nil)
}

// StartFind begins a work suit comparison search. Returns the raw response
// containing a search token.
// PDF 9.7.7: POST /cgi-bin/api/WorkSuitCompareServer/startFind
func (s *WorkSuitService) StartFind(ctx context.Context, body interface{}) (string, error) {
	resp, err := s.client.postRaw(ctx, "/cgi-bin/api/WorkSuitCompareServer/startFind", body)
	if err != nil {
		return "", fmt.Errorf("WorkSuitService.StartFind: %w", err)
	}
	return resp, nil
}

// DoFind retrieves a page of work suit comparison search results.
// PDF 9.7.8: POST /cgi-bin/api/WorkSuitCompareServer/doFind
func (s *WorkSuitService) DoFind(ctx context.Context, body interface{}) (string, error) {
	resp, err := s.client.postRaw(ctx, "/cgi-bin/api/WorkSuitCompareServer/doFind", body)
	if err != nil {
		return "", fmt.Errorf("WorkSuitService.DoFind: %w", err)
	}
	return resp, nil
}

// StopFind stops a work suit comparison search and releases the token.
// PDF 9.7.9: POST /cgi-bin/api/WorkSuitCompareServer/stopFind
func (s *WorkSuitService) StopFind(ctx context.Context, token int) error {
	reqBody := map[string]int{"token": token}
	return s.client.postJSON(ctx, "/cgi-bin/api/WorkSuitCompareServer/stopFind", reqBody, nil)
}

// DeleteByUID deletes work suit entries by their UIDs.
// PDF 9.7.10: POST /cgi-bin/api/WorkSuitCompareServer/deleteByUID
func (s *WorkSuitService) DeleteByUID(ctx context.Context, uids []string) error {
	reqBody := map[string][]string{"uids": uids}
	return s.client.postJSON(ctx, "/cgi-bin/api/WorkSuitCompareServer/deleteByUID", reqBody, nil)
}

// ReAbstract re-abstracts work suit features for the given group.
// PDF 9.7.11: POST /cgi-bin/api/WorkSuitCompareServer/reAbstract
func (s *WorkSuitService) ReAbstract(ctx context.Context, groupID string) error {
	reqBody := map[string]string{"groupID": groupID}
	return s.client.postJSON(ctx, "/cgi-bin/api/WorkSuitCompareServer/reAbstract", reqBody, nil)
}

// StopReAbstract stops an in-progress re-abstraction.
// PDF 9.7.12: POST /cgi-bin/api/WorkSuitCompareServer/stopReAbstract
func (s *WorkSuitService) StopReAbstract(ctx context.Context) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/WorkSuitCompareServer/stopReAbstract", struct{}{}, nil)
}
