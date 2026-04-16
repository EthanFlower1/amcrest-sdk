package amcrest

import (
	"context"
	"testing"
)

func TestTraffic(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	requireCapability(t, hasTraffic, "Traffic Management")

	t.Run("FindBlackList", func(t *testing.T) {
		v, err := c.Traffic.FindRecord(ctx, "TrafficBlackList")
		if err != nil {
			t.Fatalf("FindRecord(TrafficBlackList): %v", err)
		}
		t.Logf("FindRecord TrafficBlackList:\n%s", v)
	})

	t.Run("FindRedList", func(t *testing.T) {
		v, err := c.Traffic.FindRecord(ctx, "TrafficRedList")
		if err != nil {
			t.Logf("FindRecord(TrafficRedList) not available: %v", err)
			return
		}
		t.Logf("FindRecord TrafficRedList:\n%s", v)
	})

	t.Run("InsertAndRemoveRecord", func(t *testing.T) {
		err := c.Traffic.InsertRecord(ctx, "TrafficBlackList", map[string]string{
			"PlateNumber": "TEST123",
			"MasterOfCar": "TestOwner",
		})
		if err != nil {
			t.Fatalf("InsertRecord: %v", err)
		}
		t.Log("InsertRecord succeeded")

		err = c.Traffic.RemoveRecord(ctx, "TrafficBlackList", 0)
		if err != nil {
			t.Fatalf("RemoveRecord: %v", err)
		}
		t.Log("RemoveRecord succeeded")
	})
}
