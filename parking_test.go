package amcrest

import (
	"context"
	"testing"
)

func TestParking(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	requireCapability(t, hasParking, "Parking Management")

	t.Run("GetSpaceStatus", func(t *testing.T) {
		v, err := c.Parking.GetSpaceStatus(ctx, 0, 0)
		if err != nil {
			t.Fatalf("GetSpaceStatus: %v", err)
		}
		t.Logf("SpaceStatus: %s", v)
	})

	t.Run("GetAllSpaceStatus", func(t *testing.T) {
		v, err := c.Parking.GetAllSpaceStatus(ctx, 0)
		if err != nil {
			t.Fatalf("GetAllSpaceStatus: %v", err)
		}
		t.Logf("AllSpaceStatus: %s", v)
	})
}
