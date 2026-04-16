package amcrest

import (
	"context"
	"testing"
)

func TestDisplay(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	t.Run("GetGUIConfig", func(t *testing.T) {
		cfg, err := c.Display.GetGUIConfig(ctx)
		if err != nil {
			t.Logf("GetGUIConfig not available: %v", err)
			return
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty config")
		}
		for k, v := range cfg {
			t.Logf("GUI.%s = %s", k, v)
		}
	})

	t.Run("GetSplitMode", func(t *testing.T) {
		mode, err := c.Display.GetSplitMode(ctx, 0)
		if err != nil {
			t.Logf("GetSplitMode not available: %v", err)
			return
		}
		for k, v := range mode {
			t.Logf("SplitMode.%s = %s", k, v)
		}
	})

	t.Run("GetMonitorTour", func(t *testing.T) {
		cfg, err := c.Display.GetMonitorTour(ctx)
		if err != nil {
			t.Logf("GetMonitorTour not available: %v", err)
			return
		}
		for k, v := range cfg {
			t.Logf("MonitorTour.%s = %s", k, v)
		}
	})

	t.Run("GetMonitorCollection", func(t *testing.T) {
		cfg, err := c.Display.GetMonitorCollection(ctx)
		if err != nil {
			t.Logf("GetMonitorCollection not available: %v", err)
			return
		}
		for k, v := range cfg {
			t.Logf("MonitorCollection.%s = %s", k, v)
		}
	})
}
