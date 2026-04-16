package amcrest

import (
	"context"
	"testing"
)

func TestDisplayGetGUIConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Display.GetGUIConfig(ctx)
	if err != nil {
		t.Skip("GetGUIConfig not supported:", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config map")
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}

func TestDisplayGetSplitMode(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	mode, err := c.Display.GetSplitMode(ctx, 0)
	if err != nil {
		t.Skip("GetSplitMode not supported:", err)
	}
	if mode == nil {
		t.Fatal("expected non-nil mode map")
	}
	for k, v := range mode {
		t.Logf("%s = %s", k, v)
	}
}

func TestDisplayGetMonitorTour(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Display.GetMonitorTour(ctx)
	if err != nil {
		t.Skip("GetMonitorTour not supported:", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config map")
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}
