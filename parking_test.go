package amcrest

import (
	"context"
	"errors"
	"testing"
)

func TestParkingGetSpaceStatus(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	v, err := c.Parking.GetSpaceStatus(ctx, 0, 0)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			t.Skipf("TrafficParking/getSpaceStatus not supported on this device (HTTP %d)", apiErr.StatusCode)
		}
		t.Fatalf("GetSpaceStatus: %v", err)
	}
	t.Logf("SpaceStatus: %s", v)
}

func TestParkingGetAllSpaceStatus(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	v, err := c.Parking.GetAllSpaceStatus(ctx, 0)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			t.Skipf("TrafficParking/getAllSpaceStatus not supported on this device (HTTP %d)", apiErr.StatusCode)
		}
		t.Fatalf("GetAllSpaceStatus: %v", err)
	}
	t.Logf("AllSpaceStatus: %s", v)
}
