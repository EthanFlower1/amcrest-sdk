package amcrest

import (
	"context"
	"errors"
	"testing"
)

func TestTrafficFindRecord(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	v, err := c.Traffic.FindRecord(ctx, "TrafficBlackList")
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			t.Skipf("TrafficBlackList not supported on this device (HTTP %d)", apiErr.StatusCode)
		}
		t.Fatalf("FindRecord: %v", err)
	}
	t.Logf("FindRecord TrafficBlackList:\n%s", v)
}

func TestTrafficInsertAndRemoveRecord(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	// Attempt to insert a test record.
	err := c.Traffic.InsertRecord(ctx, "TrafficBlackList", map[string]string{
		"PlateNumber":  "TEST123",
		"MasterOfCar":  "TestOwner",
	})
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			t.Skipf("TrafficBlackList insert not supported on this device (HTTP %d)", apiErr.StatusCode)
		}
		t.Fatalf("InsertRecord: %v", err)
	}
	t.Log("InsertRecord succeeded")

	// Attempt to remove the record (recno 0 as a best-effort cleanup).
	err = c.Traffic.RemoveRecord(ctx, "TrafficBlackList", 0)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			t.Skipf("TrafficBlackList remove not supported on this device (HTTP %d)", apiErr.StatusCode)
		}
		t.Fatalf("RemoveRecord: %v", err)
	}
	t.Log("RemoveRecord succeeded")
}

func TestTrafficFindRedList(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	v, err := c.Traffic.FindRecord(ctx, "TrafficRedList")
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) {
			t.Skipf("TrafficRedList not supported on this device (HTTP %d)", apiErr.StatusCode)
		}
		t.Fatalf("FindRecord: %v", err)
	}
	t.Logf("FindRecord TrafficRedList:\n%s", v)
}
