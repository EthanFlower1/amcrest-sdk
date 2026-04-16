package amcrest

import (
	"context"
	"strings"
	"testing"
)

func TestUsers(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	t.Run("GetAllUsers", func(t *testing.T) {
		body, err := c.User.GetAllUsers(ctx)
		if err != nil {
			t.Fatalf("GetAllUsers: %v", err)
		}
		if body == "" {
			t.Fatal("expected non-empty response")
		}
		t.Logf("GetAllUsers response:\n%s", body)
	})

	t.Run("GetActiveUsers", func(t *testing.T) {
		body, err := c.User.GetActiveUsers(ctx)
		if err != nil {
			t.Fatalf("GetActiveUsers: %v", err)
		}
		if body == "" {
			t.Fatal("expected non-empty response")
		}
		t.Logf("GetActiveUsers response:\n%s", body)
	})

	t.Run("GetAllGroups", func(t *testing.T) {
		body, err := c.User.GetAllGroups(ctx)
		if err != nil {
			t.Fatalf("GetAllGroups: %v", err)
		}
		if body == "" {
			t.Fatal("expected non-empty response")
		}
		t.Logf("GetAllGroups response:\n%s", body)
	})

	t.Run("GetUserInfoAdmin", func(t *testing.T) {
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
	})

	t.Run("AddDeleteUser", func(t *testing.T) {
		const testUser = "sdktest"
		const testPass = "TestPass123!"
		const testGroup = "user"

		if err := c.User.AddUser(ctx, testUser, testPass, testGroup); err != nil {
			t.Fatalf("AddUser: %v", err)
		}
		t.Logf("Created user %q", testUser)

		body, err := c.User.GetAllUsers(ctx)
		if err != nil {
			t.Fatalf("GetAllUsers after add: %v", err)
		}
		if !strings.Contains(body, testUser) {
			t.Fatalf("expected %q in user list, got:\n%s", testUser, body)
		}
		t.Logf("Verified user %q exists", testUser)

		if err := c.User.DeleteUser(ctx, testUser); err != nil {
			t.Fatalf("DeleteUser: %v", err)
		}
		t.Logf("Deleted user %q", testUser)

		body, err = c.User.GetAllUsers(ctx)
		if err != nil {
			t.Fatalf("GetAllUsers after delete: %v", err)
		}
		if strings.Contains(body, testUser) {
			t.Fatalf("user %q still present after deletion:\n%s", testUser, body)
		}
		t.Logf("Verified user %q removed", testUser)
	})
}
