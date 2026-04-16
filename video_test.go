package amcrest

import (
	"context"
	"errors"
	"testing"
)

func TestVideoGetMaxExtraStreams(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	n, err := c.Video.GetMaxExtraStreams(ctx)
	if err != nil {
		t.Fatalf("GetMaxExtraStreams: %v", err)
	}
	if n < 0 {
		t.Fatalf("expected non-negative MaxExtraStreams, got %d", n)
	}
	t.Logf("MaxExtraStreams: %d", n)
}

func TestVideoGetEncodeCaps(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
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
}

func TestVideoGetEncodeConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Video.GetEncodeConfig(ctx)
	if err != nil {
		t.Fatalf("GetEncodeConfig: %v", err)
	}
	if len(cfg) == 0 {
		t.Fatal("expected non-empty Encode config")
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}

func TestVideoGetVideoInputChannels(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	n, err := c.Video.GetVideoInputChannels(ctx)
	if err != nil {
		t.Fatalf("GetVideoInputChannels: %v", err)
	}
	if n < 1 {
		t.Fatalf("expected at least 1 video input channel, got %d", n)
	}
	t.Logf("VideoInputChannels: %d", n)
}

func TestVideoGetVideoOutputChannels(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	n, err := c.Video.GetVideoOutputChannels(ctx)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) && (apiErr.StatusCode == 400 || apiErr.StatusCode == 501) {
			t.Skip("devVideoOutput not supported on this device")
		}
		t.Fatalf("GetVideoOutputChannels: %v", err)
	}
	t.Logf("VideoOutputChannels: %d", n)
}

func TestVideoGetVideoStandard(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	v, err := c.Video.GetVideoStandard(ctx)
	if err != nil {
		t.Fatalf("GetVideoStandard: %v", err)
	}
	if v == "" {
		t.Fatal("expected non-empty video standard")
	}
	t.Logf("VideoStandard: %s", v)
}

func TestVideoGetChannelTitle(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Video.GetChannelTitle(ctx)
	if err != nil {
		t.Fatalf("GetChannelTitle: %v", err)
	}
	if len(cfg) == 0 {
		t.Fatal("expected non-empty ChannelTitle config")
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}

func TestVideoSetChannelTitle(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	// Save the original channel title.
	origCfg, err := c.Video.GetChannelTitle(ctx)
	if err != nil {
		t.Fatalf("GetChannelTitle (save): %v", err)
	}
	origName, ok := origCfg["table.ChannelTitle[0].Name"]
	if !ok {
		t.Fatal("could not find table.ChannelTitle[0].Name in config")
	}
	t.Logf("Original ChannelTitle[0].Name: %s", origName)

	// Set a new title.
	testName := "SDK-Test-Title"
	if err := c.Video.SetChannelTitle(ctx, 0, testName); err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) && (apiErr.StatusCode == 501 || apiErr.StatusCode == 400) {
			t.Skip("SetChannelTitle not supported on this device")
		}
		t.Fatalf("SetChannelTitle: %v", err)
	}
	t.Logf("Set ChannelTitle[0].Name to: %s", testName)

	// Verify the title was updated.
	updatedCfg, err := c.Video.GetChannelTitle(ctx)
	if err != nil {
		t.Fatalf("GetChannelTitle (verify): %v", err)
	}
	updatedName := updatedCfg["table.ChannelTitle[0].Name"]
	if updatedName != testName {
		t.Fatalf("expected ChannelTitle[0].Name=%q, got %q", testName, updatedName)
	}
	t.Logf("Verified ChannelTitle[0].Name: %s", updatedName)

	// Restore the original title.
	if err := c.Video.SetChannelTitle(ctx, 0, origName); err != nil {
		t.Fatalf("SetChannelTitle (restore): %v", err)
	}

	// Verify restoration.
	restoredCfg, err := c.Video.GetChannelTitle(ctx)
	if err != nil {
		t.Fatalf("GetChannelTitle (restore verify): %v", err)
	}
	restoredName := restoredCfg["table.ChannelTitle[0].Name"]
	if restoredName != origName {
		t.Fatalf("expected restored ChannelTitle[0].Name=%q, got %q", origName, restoredName)
	}
	t.Logf("Restored ChannelTitle[0].Name: %s", restoredName)
}

func TestVideoGetVideoWidget(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Video.GetVideoWidget(ctx)
	if err != nil {
		t.Fatalf("GetVideoWidget: %v", err)
	}
	if len(cfg) == 0 {
		t.Fatal("expected non-empty VideoWidget config")
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}

func TestVideoGetSmartEncodeConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Video.GetSmartEncodeConfig(ctx)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) && (apiErr.StatusCode == 400 || apiErr.StatusCode == 501) {
			t.Skip("SmartEncode not supported on this device")
		}
		t.Fatalf("GetSmartEncodeConfig: %v", err)
	}
	if len(cfg) == 0 {
		t.Fatal("expected non-empty SmartEncode config")
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}

func TestVideoGetVideoInputCaps(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
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
}

func TestVideoAllReadOnly(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	tests := []struct {
		name string
		fn   func() error
	}{
		{"GetMaxExtraStreams", func() error { _, err := c.Video.GetMaxExtraStreams(ctx); return err }},
		{"GetEncodeCaps", func() error { _, err := c.Video.GetEncodeCaps(ctx); return err }},
		{"GetEncodeConfig", func() error { _, err := c.Video.GetEncodeConfig(ctx); return err }},
		{"GetVideoInputChannels", func() error { _, err := c.Video.GetVideoInputChannels(ctx); return err }},
		{"GetVideoStandard", func() error { _, err := c.Video.GetVideoStandard(ctx); return err }},
		{"GetChannelTitle", func() error { _, err := c.Video.GetChannelTitle(ctx); return err }},
		{"GetVideoWidget", func() error { _, err := c.Video.GetVideoWidget(ctx); return err }},
		{"GetVideoInputCaps(0)", func() error { _, err := c.Video.GetVideoInputCaps(ctx, 0); return err }},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.fn(); err != nil {
				var apiErr *APIError
				if errors.As(err, &apiErr) && (apiErr.StatusCode == 400 || apiErr.StatusCode == 501) {
					t.Skipf("%s not supported on this device", tc.name)
				}
				t.Errorf("%s failed: %v", tc.name, err)
			}
		})
	}
}
