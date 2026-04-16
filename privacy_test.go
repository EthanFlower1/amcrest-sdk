package amcrest

import (
	"context"
	"testing"
)

func TestPrivacyGetConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Privacy.GetConfig(ctx)
	if err != nil {
		t.Skip("Privacy.GetConfig not supported:", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config map")
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}

func TestPrivacyGetMasking(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	body, err := c.Privacy.GetMasking(ctx, 0, 0, 10)
	if err != nil {
		t.Skip("Privacy.GetMasking not supported:", err)
	}
	t.Logf("GetMasking response:\n%s", body)
}

func TestPrivacyGetEnable(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	enabled, err := c.Privacy.GetEnable(ctx, 0)
	if err != nil {
		t.Skip("Privacy.GetEnable not supported:", err)
	}
	t.Logf("PrivacyMasking enabled: %v", enabled)
}
