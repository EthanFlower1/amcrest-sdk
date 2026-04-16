package amcrest

import (
	"context"
	"errors"
	"testing"
)

func skipIfNoAccessControl(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		return
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		t.Skipf("camera does not support access control (HTTP %d), skipping", apiErr.StatusCode)
	}
}

func TestAccessControlGetDoorStatus(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	status, err := c.AccessControl.GetDoorStatus(ctx, 0)
	skipIfNoAccessControl(t, err)
	if err != nil {
		t.Fatalf("GetDoorStatus: %v", err)
	}
	t.Logf("Door status: %s", status)
}

func TestAccessControlGetGeneralConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	cfg, err := c.AccessControl.GetGeneralConfig(ctx)
	skipIfNoAccessControl(t, err)
	if err != nil {
		t.Fatalf("GetGeneralConfig: %v", err)
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}

func TestAccessControlGetControlConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	cfg, err := c.AccessControl.GetControlConfig(ctx)
	skipIfNoAccessControl(t, err)
	if err != nil {
		t.Fatalf("GetControlConfig: %v", err)
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}
