package amcrest

import (
	"context"
	"testing"
)

func TestWorkSuit(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	requireCapability(t, hasWorkSuit, "WorkSuit Detection")

	t.Run("FindGroup", func(t *testing.T) {
		body, err := c.WorkSuit.FindGroup(ctx)
		if err != nil {
			t.Fatalf("FindGroup: %v", err)
		}
		t.Logf("FindGroup response:\n%s", body)
	})

	t.Run("GetGroup", func(t *testing.T) {
		body, err := c.WorkSuit.GetGroup(ctx, 0)
		if err != nil {
			t.Fatalf("GetGroup: %v", err)
		}
		t.Logf("GetGroup(0) response:\n%s", body)
	})
}
