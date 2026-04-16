package amcrest

import (
	"context"
	"testing"
)

func TestPeople(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	requireCapability(t, hasPeople, "People Counting")

	t.Run("GetSummary", func(t *testing.T) {
		body, err := c.People.GetSummary(ctx)
		if err != nil {
			t.Fatalf("GetSummary: %v", err)
		}
		t.Logf("GetSummary response:\n%s", body)
	})

	t.Run("GetCrowdCaps", func(t *testing.T) {
		body, err := c.People.GetCrowdCaps(ctx)
		if err != nil {
			t.Fatalf("GetCrowdCaps: %v", err)
		}
		t.Logf("GetCrowdCaps response:\n%s", body)
	})
}
