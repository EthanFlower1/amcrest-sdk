package amcrest

import (
	"context"
	"fmt"
	"io"
	"net/url"
)

// SnapshotService handles snapshot-related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 66-70 (Section 4.4)
type SnapshotService struct {
	client *Client
}

// Get captures a snapshot from the given channel and returns the raw JPEG bytes.
// Channel is 1-based (e.g., 1 for the first video channel).
func (s *SnapshotService) Get(ctx context.Context, channel int) ([]byte, error) {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	resp, err := s.client.get(ctx, "/cgi-bin/snapshot.cgi", params)
	if err != nil {
		return nil, fmt.Errorf("amcrest: snapshot get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, &APIError{StatusCode: resp.StatusCode}
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("amcrest: reading snapshot body: %w", err)
	}

	return data, nil
}

// GetSnapConfig retrieves the snapshot configuration table.
func (s *SnapshotService) GetSnapConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getConfig(ctx, "Snap")
}
