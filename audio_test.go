package amcrest

import (
	"context"
	"testing"
)

func TestAudio(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	t.Run("GetInputChannels", func(t *testing.T) {
		n, err := c.Audio.GetInputChannels(ctx)
		if err != nil {
			t.Fatalf("GetInputChannels: %v", err)
		}
		if n < 0 {
			t.Fatalf("expected >= 0 input channels, got %d", n)
		}
		t.Logf("InputChannels: %d", n)
	})

	t.Run("GetOutputChannels", func(t *testing.T) {
		n, err := c.Audio.GetOutputChannels(ctx)
		if err != nil {
			t.Fatalf("GetOutputChannels: %v", err)
		}
		if n < 0 {
			t.Fatalf("expected >= 0 output channels, got %d", n)
		}
		t.Logf("OutputChannels: %d", n)
	})

	t.Run("GetVolume", func(t *testing.T) {
		if audioOutChans <= 0 {
			t.Skip("camera has no audio output channels, skipping volume test")
		}
		vol, err := c.Audio.GetVolume(ctx)
		if err != nil {
			t.Fatalf("GetVolume: %v", err)
		}
		if vol == nil {
			t.Fatal("expected non-nil volume map")
		}
		for k, v := range vol {
			t.Logf("Volume.%s = %s", k, v)
		}
	})
}
