package amcrest

import (
	"context"
	"errors"
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
		n, err := c.Video.GetVideoOutputChannels(ctx)
		if err != nil {
			var apiErr *APIError
			if errors.As(err, &apiErr) && (apiErr.StatusCode == 400 || apiErr.StatusCode == 501) {
				t.Skip("devVideoOutput not supported on this device")
			}
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

		testName := "SDK-Test-Title"
		if err := c.Video.SetChannelTitle(ctx, 0, testName); err != nil {
			var apiErr *APIError
			if errors.As(err, &apiErr) && (apiErr.StatusCode == 501 || apiErr.StatusCode == 400) {
				t.Skip("SetChannelTitle not supported on this device")
			}
			t.Fatalf("SetChannelTitle: %v", err)
		}
		t.Logf("Set ChannelTitle to: %s", testName)

		updatedCfg, err := c.Video.GetChannelTitle(ctx)
		if err != nil {
			t.Fatalf("GetChannelTitle (verify): %v", err)
		}
		if updatedCfg["table.ChannelTitle[0].Name"] != testName {
			t.Fatalf("expected %q, got %q", testName, updatedCfg["table.ChannelTitle[0].Name"])
		}

		if err := c.Video.SetChannelTitle(ctx, 0, origName); err != nil {
			t.Fatalf("SetChannelTitle (restore): %v", err)
		}
		t.Logf("Restored ChannelTitle to: %s", origName)
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
		cfg, err := c.Video.GetSmartEncodeConfig(ctx)
		if err != nil {
			var apiErr *APIError
			if errors.As(err, &apiErr) && (apiErr.StatusCode == 400 || apiErr.StatusCode == 501) {
				t.Skip("SmartEncode not supported on this device")
			}
			t.Fatalf("GetSmartEncodeConfig: %v", err)
		}
		for k, v := range cfg {
			t.Logf("SmartEncode.%s = %s", k, v)
		}
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
		caps, err := c.Video.GetEncodeConfigCaps(ctx, 0)
		if err != nil {
			var apiErr *APIError
			if errors.As(err, &apiErr) && (apiErr.StatusCode == 400 || apiErr.StatusCode == 501) {
				t.Skip("GetEncodeConfigCaps not supported on this device")
			}
			t.Fatalf("GetEncodeConfigCaps: %v", err)
		}
		for k, v := range caps {
			t.Logf("EncodeConfigCaps.%s = %s", k, v)
		}
	})

	t.Run("GetVideoEncodeROI", func(t *testing.T) {
		cfg, err := c.Video.GetVideoEncodeROI(ctx)
		if err != nil {
			var apiErr *APIError
			if errors.As(err, &apiErr) && (apiErr.StatusCode == 400 || apiErr.StatusCode == 501) {
				t.Skip("VideoEncodeROI not supported on this device")
			}
			t.Fatalf("GetVideoEncodeROI: %v", err)
		}
		for k, v := range cfg {
			t.Logf("ROI.%s = %s", k, v)
		}
	})
}
