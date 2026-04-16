package amcrest

import (
	"context"
	"strconv"
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

	t.Run("SetVolume", func(t *testing.T) {
		if audioOutChans <= 0 {
			t.Skip("camera has no audio output channels, skipping volume test")
		}
		original, err := c.Audio.GetVolume(ctx)
		if err != nil {
			t.Fatalf("GetVolume (save): %v", err)
		}

		// Find the volume key for channel 0.
		volKey := ""
		origVal := ""
		for k, v := range original {
			if contains(k, "AudioOutputVolume[0]") || contains(k, "AudioOutputVolume") {
				volKey = k
				origVal = v
				break
			}
		}
		if volKey == "" {
			t.Skip("no AudioOutputVolume key found")
		}
		t.Logf("Original %s = %s", volKey, origVal)

		origInt, _ := strconv.Atoi(origVal)

		defer func() {
			_ = c.Audio.SetVolume(ctx, 0, origInt)
		}()

		// Set to 50 (or 45 if already 50).
		newVol := 50
		if origInt == 50 {
			newVol = 45
		}
		err = c.Audio.SetVolume(ctx, 0, newVol)
		skipOnSetError(t, err, "SetVolume")

		updated, err := c.Audio.GetVolume(ctx)
		if err != nil {
			t.Fatalf("GetVolume (verify): %v", err)
		}
		updatedVal := ""
		for k, v := range updated {
			if contains(k, "AudioOutputVolume[0]") || contains(k, "AudioOutputVolume") {
				updatedVal = v
				break
			}
		}
		expectedStr := strconv.Itoa(newVol)
		if updatedVal != expectedStr {
			t.Fatalf("expected volume=%q, got %q", expectedStr, updatedVal)
		}
		t.Logf("Verified volume changed to %d", newVol)
	})
}
