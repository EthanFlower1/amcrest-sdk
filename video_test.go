package amcrest

import (
	"context"
	"testing"
)

func TestVideo(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	t.Run("GetMaxExtraStreams", func(t *testing.T) {
		n, err := c.Video.GetMaxExtraStreams(ctx)
		if err != nil {
			t.Fatalf("GetMaxExtraStreams: %v", err)
		}
		if n < 0 {
			t.Fatalf("expected non-negative MaxExtraStreams, got %d", n)
		}
		t.Logf("MaxExtraStreams: %d", n)
	})

	t.Run("GetEncodeCaps", func(t *testing.T) {
		caps, err := c.Video.GetEncodeCaps(ctx)
		if err != nil {
			t.Fatalf("GetEncodeCaps: %v", err)
		}
		if len(caps) == 0 {
			t.Fatal("expected non-empty encode caps")
		}
		for k, v := range caps {
			t.Logf("EncodeCaps.%s = %s", k, v)
		}
	})

	t.Run("GetEncodeConfig", func(t *testing.T) {
		cfg, err := c.Video.GetEncodeConfig(ctx)
		if err != nil {
			t.Fatalf("GetEncodeConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty Encode config")
		}
		for k, v := range cfg {
			t.Logf("Encode.%s = %s", k, v)
		}
	})

	t.Run("GetVideoInputChannels", func(t *testing.T) {
		n, err := c.Video.GetVideoInputChannels(ctx)
		if err != nil {
			t.Fatalf("GetVideoInputChannels: %v", err)
		}
		if n < 1 {
			t.Fatalf("expected at least 1 video input channel, got %d", n)
		}
		t.Logf("VideoInputChannels: %d", n)
	})

	t.Run("GetVideoOutputChannels", func(t *testing.T) {
		if !hasVideoOutput {
			t.Skip("camera does not support devVideoOutput channels")
		}
		n, err := c.Video.GetVideoOutputChannels(ctx)
		if err != nil {
			t.Fatalf("GetVideoOutputChannels: %v", err)
		}
		t.Logf("VideoOutputChannels: %d", n)
	})

	t.Run("GetVideoStandard", func(t *testing.T) {
		v, err := c.Video.GetVideoStandard(ctx)
		if err != nil {
			t.Fatalf("GetVideoStandard: %v", err)
		}
		if v == "" {
			t.Fatal("expected non-empty video standard")
		}
		t.Logf("VideoStandard: %s", v)
	})

	t.Run("GetChannelTitle", func(t *testing.T) {
		cfg, err := c.Video.GetChannelTitle(ctx)
		if err != nil {
			t.Fatalf("GetChannelTitle: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty ChannelTitle config")
		}
		for k, v := range cfg {
			t.Logf("ChannelTitle.%s = %s", k, v)
		}
	})

	t.Run("SetChannelTitle", func(t *testing.T) {
		origCfg, err := c.Video.GetChannelTitle(ctx)
		if err != nil {
			t.Fatalf("GetChannelTitle (save): %v", err)
		}
		origName, ok := origCfg["table.ChannelTitle[0].Name"]
		if !ok {
			t.Fatal("could not find table.ChannelTitle[0].Name")
		}
		t.Logf("Original ChannelTitle[0].Name: %s", origName)

		defer func() {
			_ = c.Video.SetChannelTitle(ctx, 0, origName)
		}()

		testName := "SDK-Test-Title"
		err = c.Video.SetChannelTitle(ctx, 0, testName)
		skipOnSetError(t, err, "SetChannelTitle")

		updatedCfg, err := c.Video.GetChannelTitle(ctx)
		if err != nil {
			t.Fatalf("GetChannelTitle (verify): %v", err)
		}
		if updatedCfg["table.ChannelTitle[0].Name"] != testName {
			t.Fatalf("expected %q, got %q", testName, updatedCfg["table.ChannelTitle[0].Name"])
		}
		t.Logf("Verified ChannelTitle changed to %q", testName)
	})

	t.Run("GetVideoWidget", func(t *testing.T) {
		cfg, err := c.Video.GetVideoWidget(ctx)
		if err != nil {
			t.Fatalf("GetVideoWidget: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty VideoWidget config")
		}
		for k, v := range cfg {
			t.Logf("VideoWidget.%s = %s", k, v)
		}
	})

	t.Run("GetSmartEncodeConfig", func(t *testing.T) {
		if !hasSmartEncode {
			t.Skip("camera does not support SmartEncode config")
		}
		cfg, err := c.Video.GetSmartEncodeConfig(ctx)
		if err != nil {
			t.Fatalf("GetSmartEncodeConfig: %v", err)
		}
		for k, v := range cfg {
			t.Logf("SmartEncode.%s = %s", k, v)
		}
	})

	t.Run("SetSmartEncodeConfig", func(t *testing.T) {
		if !hasSmartEncode {
			t.Skip("camera does not support SmartEncode config")
		}
		original, err := c.Video.GetSmartEncodeConfig(ctx)
		if err != nil {
			t.Fatalf("GetSmartEncodeConfig (save): %v", err)
		}

		// Find the Enable key.
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
			t.Skip("no Enable key found in SmartEncode config")
		}
		t.Logf("Original %s = %s", enableKey, origVal)

		// Strip the "table." prefix for the set key.
		setKey := enableKey
		if len(setKey) > 6 && setKey[:6] == "table." {
			setKey = setKey[6:]
		}

		defer func() {
			_ = c.Video.SetSmartEncodeConfig(ctx, map[string]string{
				setKey: origVal,
			})
		}()

		newVal := "true"
		if origVal == "true" {
			newVal = "false"
		}
		err = c.Video.SetSmartEncodeConfig(ctx, map[string]string{
			setKey: newVal,
		})
		skipOnSetError(t, err, "SetSmartEncodeConfig")

		updated, err := c.Video.GetSmartEncodeConfig(ctx)
		if err != nil {
			t.Fatalf("GetSmartEncodeConfig (verify): %v", err)
		}
		if updated[enableKey] != newVal {
			t.Fatalf("expected %s=%q, got %q", enableKey, newVal, updated[enableKey])
		}
		t.Logf("Verified %s changed to %q", enableKey, newVal)
	})

	t.Run("GetVideoInputCaps", func(t *testing.T) {
		caps, err := c.Video.GetVideoInputCaps(ctx, 0)
		if err != nil {
			t.Fatalf("GetVideoInputCaps: %v", err)
		}
		if len(caps) == 0 {
			t.Fatal("expected non-empty video input caps")
		}
		for k, v := range caps {
			t.Logf("VideoInputCaps.%s = %s", k, v)
		}
	})

	t.Run("GetEncodeConfigCaps", func(t *testing.T) {
		if !hasEncodeConfCaps {
			t.Skip("camera does not support EncodeConfigCaps")
		}
		caps, err := c.Video.GetEncodeConfigCaps(ctx, 0)
		if err != nil {
			t.Fatalf("GetEncodeConfigCaps: %v", err)
		}
		for k, v := range caps {
			t.Logf("EncodeConfigCaps.%s = %s", k, v)
		}
	})

	t.Run("GetVideoEncodeROI", func(t *testing.T) {
		if !hasEncodeROI {
			t.Skip("camera does not support VideoEncodeROI config")
		}
		cfg, err := c.Video.GetVideoEncodeROI(ctx)
		if err != nil {
			t.Fatalf("GetVideoEncodeROI: %v", err)
		}
		for k, v := range cfg {
			t.Logf("ROI.%s = %s", k, v)
		}
	})
}
