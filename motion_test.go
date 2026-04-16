package amcrest

import (
	"context"
	"testing"
)

func TestMotionGetConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Motion.GetConfig(ctx)
	if err != nil {
		t.Skip("Motion.GetConfig not supported:", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config map")
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}

func TestMotionGetSmartMotionConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Motion.GetSmartMotionConfig(ctx)
	if err != nil {
		t.Skip("Motion.GetSmartMotionConfig not supported:", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config map")
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}
