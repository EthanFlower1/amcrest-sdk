package amcrest

import (
	"context"
	"strings"
	"testing"
)

func TestUserGetAllUsers(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	body, err := c.User.GetAllUsers(ctx)
	if err != nil {
		t.Fatalf("GetAllUsers: %v", err)
	}
	if body == "" {
		t.Fatal("expected non-empty response")
	}
	t.Logf("GetAllUsers response:\n%s", body)
}

func TestUserGetActiveUsers(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	body, err := c.User.GetActiveUsers(ctx)
	if err != nil {
		t.Fatalf("GetActiveUsers: %v", err)
	}
	if body == "" {
		t.Fatal("expected non-empty response")
	}
	t.Logf("GetActiveUsers response:\n%s", body)
}

func TestUserGetAllGroups(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	body, err := c.User.GetAllGroups(ctx)
	if err != nil {
		t.Fatalf("GetAllGroups: %v", err)
	}
	if body == "" {
		t.Fatal("expected non-empty response")
	}
	t.Logf("GetAllGroups response:\n%s", body)
}

func TestUserGetUserInfoAdmin(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	info, err := c.User.GetUserInfo(ctx, "admin")
	if err != nil {
		t.Fatalf("GetUserInfo(admin): %v", err)
	}
	if len(info) == 0 {
		t.Fatal("expected non-empty user info for admin")
	}
	for k, v := range info {
		t.Logf("admin.%s = %s", k, v)
	}
}

func TestUserAddDeleteUser(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	const testUser = "sdktest"
	const testPass = "TestPass123!"
	const testGroup = "user"

	// Create the test user.
	if err := c.User.AddUser(ctx, testUser, testPass, testGroup); err != nil {
		t.Fatalf("AddUser: %v", err)
	}
	t.Logf("Created user %q", testUser)

	// Verify user exists in the full user list.
	body, err := c.User.GetAllUsers(ctx)
	if err != nil {
		t.Fatalf("GetAllUsers after add: %v", err)
	}
	if !strings.Contains(body, testUser) {
		t.Fatalf("expected %q in user list, got:\n%s", testUser, body)
	}
	t.Logf("Verified user %q exists in user list", testUser)

	// Delete the test user.
	if err := c.User.DeleteUser(ctx, testUser); err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}
	t.Logf("Deleted user %q", testUser)

	// Verify user is gone.
	body, err = c.User.GetAllUsers(ctx)
	if err != nil {
		t.Fatalf("GetAllUsers after delete: %v", err)
	}
	if strings.Contains(body, testUser) {
		t.Fatalf("user %q still present after deletion:\n%s", testUser, body)
	}
	t.Logf("Verified user %q no longer in user list", testUser)
}
