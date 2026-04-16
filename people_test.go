package amcrest

import (
	"context"
	"testing"
)

func TestPeopleGetSummary(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	body, err := c.People.GetSummary(ctx)
	if err != nil {
		skipIfUnsupported(t, err)
		t.Fatalf("GetSummary: %v", err)
	}
	t.Logf("GetSummary response:\n%s", body)
}

func TestPeopleGetCrowdCaps(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	body, err := c.People.GetCrowdCaps(ctx)
	if err != nil {
		skipIfUnsupported(t, err)
		t.Fatalf("GetCrowdCaps: %v", err)
	}
	t.Logf("GetCrowdCaps response:\n%s", body)
}
