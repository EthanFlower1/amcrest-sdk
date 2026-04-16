package amcrest

import (
	"context"
	"errors"
	"testing"
)

func skipIfNoBuilding(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		return
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		t.Skipf("camera does not support building intercom (HTTP %d), skipping", apiErr.StatusCode)
	}
}

func TestBuildingGetSIPConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	cfg, err := c.Building.GetSIPConfig(ctx)
	skipIfNoBuilding(t, err)
	if err != nil {
		t.Fatalf("GetSIPConfig: %v", err)
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}

func TestBuildingGetRoomNumberCount(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	count, err := c.Building.GetRoomNumberCount(ctx)
	skipIfNoBuilding(t, err)
	if err != nil {
		t.Fatalf("GetRoomNumberCount: %v", err)
	}
	t.Logf("Room number count: %d", count)
}
