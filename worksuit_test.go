package amcrest

import (
	"context"
	"testing"
)

func TestWorkSuitFindGroup(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	body, err := c.WorkSuit.FindGroup(ctx)
	if err != nil {
		skipIfUnsupported(t, err)
		t.Fatalf("FindGroup: %v", err)
	}
	t.Logf("FindGroup response:\n%s", body)
}

func TestWorkSuitGetGroup(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	body, err := c.WorkSuit.GetGroup(ctx, 0)
	if err != nil {
		skipIfUnsupported(t, err)
		t.Fatalf("GetGroup: %v", err)
	}
	t.Logf("GetGroup(0) response:\n%s", body)
}
