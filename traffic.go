package amcrest

import (
	"context"
	"fmt"
	"net/url"
)

// TrafficService handles traffic-related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 440-466 (Section 10)
type TrafficService struct {
	client *Client
}

// InsertRecord inserts a record into a traffic list (e.g., TrafficBlackList or
// TrafficRedList) via recordUpdater.cgi. The params map may include keys such
// as PlateNumber, MasterOfCar, etc.
func (s *TrafficService) InsertRecord(ctx context.Context, name string, params map[string]string) error {
	qv := url.Values{}
	qv.Set("name", name)
	for k, v := range params {
		qv.Set(k, v)
	}
	return s.client.cgiAction(ctx, "recordUpdater.cgi", "insert", qv)
}

// RemoveRecord removes a record from a traffic list by record number.
func (s *TrafficService) RemoveRecord(ctx context.Context, name string, recno int) error {
	return s.client.cgiAction(ctx, "recordUpdater.cgi", "remove", url.Values{
		"name":  {name},
		"recno": {fmt.Sprintf("%d", recno)},
	})
}

// FindRecord searches for records in the given traffic list and returns the
// raw response body.
func (s *TrafficService) FindRecord(ctx context.Context, name string) (string, error) {
	return s.client.cgiGet(ctx, "recordFinder.cgi", "find", url.Values{
		"name": {name},
	})
}
