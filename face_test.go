package amcrest

import (
	"context"
	"testing"
)

func TestFace(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	requireCapability(t, hasFaceRec, "Face Recognition")

	t.Run("FindGroup", func(t *testing.T) {
		body, err := c.Face.FindGroup(ctx)
		if err != nil {
			t.Fatalf("FindGroup: %v", err)
		}
		t.Logf("FindGroup response:\n%s", body)
	})

	t.Run("GetGroupForChannel", func(t *testing.T) {
		body, err := c.Face.GetGroupForChannel(ctx, 0)
		if err != nil {
			t.Fatalf("GetGroupForChannel: %v", err)
		}
		t.Logf("GetGroupForChannel(0):\n%s", body)
	})

	t.Run("CreateDeleteGroup", func(t *testing.T) {
		groupID, err := c.Face.CreateGroup(ctx, "sdk-test-group", "integration test")
		if err != nil {
			t.Fatalf("CreateGroup: %v", err)
		}
		t.Logf("Created group with ID: %s", groupID)

		if err := c.Face.DeleteGroup(ctx, groupID); err != nil {
			t.Fatalf("DeleteGroup(%s): %v", groupID, err)
		}
		t.Logf("Deleted group %s", groupID)
	})
}
