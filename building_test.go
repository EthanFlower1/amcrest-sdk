package amcrest

import (
	"context"
	"testing"
)

func TestBuilding(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	requireCapability(t, hasBuilding, "Building Intercom")

	t.Run("GetSIPConfig", func(t *testing.T) {
		cfg, err := c.Building.GetSIPConfig(ctx)
		if err != nil {
			t.Fatalf("GetSIPConfig: %v", err)
		}
		for k, v := range cfg {
			t.Logf("SIP.%s = %s", k, v)
		}
	})

	t.Run("GetRoomNumberCount", func(t *testing.T) {
		count, err := c.Building.GetRoomNumberCount(ctx)
		if err != nil {
			t.Fatalf("GetRoomNumberCount: %v", err)
		}
		t.Logf("Room number count: %d", count)
	})
}
