package amcrest

import (
	"context"
	"errors"
	"testing"
)

func TestAudioGetInputChannels(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	n, err := c.Audio.GetInputChannels(ctx)
	if err != nil {
		t.Fatalf("GetInputChannels: %v", err)
	}
	if n < 0 {
		t.Fatalf("expected >= 0 input channels, got %d", n)
	}
	t.Logf("InputChannels: %d", n)
}

func TestAudioGetOutputChannels(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	n, err := c.Audio.GetOutputChannels(ctx)
	if err != nil {
		t.Fatalf("GetOutputChannels: %v", err)
	}
	if n < 0 {
		t.Fatalf("expected >= 0 output channels, got %d", n)
	}
	t.Logf("OutputChannels: %d", n)
}

func TestAudioGetVolume(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	vol, err := c.Audio.GetVolume(ctx)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 400 {
			t.Skip("AudioOutputVolume not supported on this device")
		}
		t.Fatalf("GetVolume: %v", err)
	}
	if vol == nil {
		t.Fatal("expected non-nil volume map")
	}
	for k, v := range vol {
		t.Logf("%s = %s", k, v)
	}
}
