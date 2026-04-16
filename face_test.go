package amcrest

import (
	"context"
	"errors"
	"testing"
)

func skipIfUnsupported(t *testing.T, err error) {
	t.Helper()
	var apiErr *APIError
	if errors.As(err, &apiErr) && (apiErr.StatusCode == 400 || apiErr.StatusCode == 404 || apiErr.StatusCode == 501) {
		t.Skipf("not supported on this device (HTTP %d)", apiErr.StatusCode)
	}
}

func TestFaceFindGroup(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	body, err := c.Face.FindGroup(ctx)
	if err != nil {
		skipIfUnsupported(t, err)
		t.Fatalf("FindGroup: %v", err)
	}
	t.Logf("FindGroup response:\n%s", body)
}

func TestFaceGetGroupForChannel(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	body, err := c.Face.GetGroupForChannel(ctx, 0)
	if err != nil {
		skipIfUnsupported(t, err)
		t.Fatalf("GetGroupForChannel: %v", err)
	}
	t.Logf("GetGroupForChannel(0) response:\n%s", body)
}

func TestFaceCreateDeleteGroup(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	// Create a test group.
	groupID, err := c.Face.CreateGroup(ctx, "sdk-test-group", "integration test")
	if err != nil {
		skipIfUnsupported(t, err)
		t.Fatalf("CreateGroup: %v", err)
	}
	t.Logf("Created group with ID: %s", groupID)

	// Clean up: delete the group.
	if err := c.Face.DeleteGroup(ctx, groupID); err != nil {
		t.Fatalf("DeleteGroup(%s): %v", groupID, err)
	}
	t.Logf("Deleted group %s", groupID)
}
