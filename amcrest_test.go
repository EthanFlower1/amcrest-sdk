package amcrest

import (
	"context"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	body, err := c.cgiGet(ctx, "magicBox.cgi", "getDeviceType", nil)
	if err != nil {
		t.Fatalf("getDeviceType failed: %v", err)
	}
	result := parseKV(body)
	if result["type"] == "" {
		t.Fatalf("expected device type, got empty. Full body: %s", body)
	}
	t.Logf("Device type: %s", result["type"])
}
