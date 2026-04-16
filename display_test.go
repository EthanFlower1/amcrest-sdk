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
		if !hasGUISet {
			t.Skip("camera does not support GUISet config")
		}
		cfg, err := c.Display.GetGUIConfig(ctx)
		if err != nil {
			t.Fatalf("GetGUIConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty config")
		}
		for k, v := range cfg {
			t.Logf("GUI.%s = %s", k, v)
		}
	})

	t.Run("GetSplitMode", func(t *testing.T) {
		if !hasSplitScreen {
			t.Skip("camera does not support split.cgi")
		}
		mode, err := c.Display.GetSplitMode(ctx, 0)
		if err != nil {
			t.Fatalf("GetSplitMode: %v", err)
		}
		for k, v := range mode {
			t.Logf("SplitMode.%s = %s", k, v)
		}
	})

	t.Run("GetMonitorTour", func(t *testing.T) {
		if !hasMonitorTour {
			t.Skip("camera does not support MonitorTour config")
		}
		cfg, err := c.Display.GetMonitorTour(ctx)
		if err != nil {
			t.Fatalf("GetMonitorTour: %v", err)
		}
		for k, v := range cfg {
			t.Logf("MonitorTour.%s = %s", k, v)
		}
	})

	t.Run("SetMonitorTour", func(t *testing.T) {
		if !hasMonitorTour {
			t.Skip("camera does not support MonitorTour config")
		}
		original, err := c.Display.GetMonitorTour(ctx)
		if err != nil {
			t.Fatalf("GetMonitorTour (save): %v", err)
		}

		// Find an Enable key to toggle.
		enableKey := ""
		origVal := ""
		for k, v := range original {
			if contains(k, "Enable") {
				enableKey = k
				origVal = v
				break
			}
		}
		if enableKey == "" {
			t.Skip("no Enable key found in MonitorTour config")
		}
		t.Logf("Original %s = %s", enableKey, origVal)

		setKey := enableKey
		if len(setKey) > 6 && setKey[:6] == "table." {
			setKey = setKey[6:]
		}

		defer func() {
			_ = c.Display.SetMonitorTour(ctx, map[string]string{
				setKey: origVal,
			})
		}()

		newVal := "true"
		if origVal == "true" {
			newVal = "false"
		}
		err = c.Display.SetMonitorTour(ctx, map[string]string{
			setKey: newVal,
		})
		skipOnSetError(t, err, "SetMonitorTour")

		updated, err := c.Display.GetMonitorTour(ctx)
		if err != nil {
			t.Fatalf("GetMonitorTour (verify): %v", err)
		}
		if updated[enableKey] != newVal {
			t.Fatalf("expected %s=%q, got %q", enableKey, newVal, updated[enableKey])
		}
		t.Logf("Verified MonitorTour Enable changed to %q", newVal)
	})

	t.Run("GetMonitorCollection", func(t *testing.T) {
		if !hasMonitorCollect {
			t.Skip("camera does not support MonitorCollection config")
		}
		cfg, err := c.Display.GetMonitorCollection(ctx)
		if err != nil {
			t.Fatalf("GetMonitorCollection: %v", err)
		}
		for k, v := range cfg {
			t.Logf("MonitorCollection.%s = %s", k, v)
		}
	})
}
