package amcrest

import (
	"context"
	"testing"
	"time"
)

func TestRecording(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	t.Run("GetCaps", func(t *testing.T) {
		if !hasRecordCaps {
			t.Skip("camera does not support recordManager.cgi getCaps")
		}
		caps, err := c.Recording.GetCaps(ctx)
		if err != nil {
			t.Fatalf("GetCaps: %v", err)
		}
		if len(caps) == 0 {
			t.Fatal("expected non-empty caps")
		}
		for k, v := range caps {
			t.Logf("Caps.%s = %s", k, v)
		}
	})

	t.Run("GetRecordConfig", func(t *testing.T) {
		if !hasMediaFileFind {
			t.Skip("camera does not support Record config")
		}
		cfg, err := c.Recording.GetRecordConfig(ctx)
		if err != nil {
			t.Fatalf("GetRecordConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty Record config")
		}
		for k, v := range cfg {
			t.Logf("Record.%s = %s", k, v)
		}
	})

	t.Run("SetRecordConfig", func(t *testing.T) {
		if !hasMediaFileFind {
			t.Skip("camera does not support Record config")
		}
		original, err := c.Recording.GetRecordConfig(ctx)
		if err != nil {
			t.Fatalf("GetRecordConfig (save): %v", err)
		}

		// Find PreRecord key.
		preRecordKey := ""
		origVal := ""
		for k, v := range original {
			if contains(k, "PreRecord") {
				preRecordKey = k
				origVal = v
				break
			}
		}
		if preRecordKey == "" {
			t.Skip("no PreRecord key found in RecordConfig")
		}
		t.Logf("Original %s = %s", preRecordKey, origVal)

		setKey := preRecordKey
		if len(setKey) > 6 && setKey[:6] == "table." {
			setKey = setKey[6:]
		}

		defer func() {
			_ = c.Recording.SetRecordConfig(ctx, map[string]string{
				setKey: origVal,
			})
		}()

		newVal := "5"
		if origVal == "5" {
			newVal = "4"
		}
		err = c.Recording.SetRecordConfig(ctx, map[string]string{
			setKey: newVal,
		})
		skipOnSetError(t, err, "SetRecordConfig")

		updated, err := c.Recording.GetRecordConfig(ctx)
		if err != nil {
			t.Fatalf("GetRecordConfig (verify): %v", err)
		}
		if updated[preRecordKey] != newVal {
			t.Fatalf("expected %s=%q, got %q", preRecordKey, newVal, updated[preRecordKey])
		}
		t.Logf("Verified PreRecord changed to %q", newVal)
	})

	t.Run("GetRecordMode", func(t *testing.T) {
		if !hasMediaFileFind {
			t.Skip("camera does not support RecordMode config")
		}
		cfg, err := c.Recording.GetRecordMode(ctx)
		if err != nil {
			t.Fatalf("GetRecordMode: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty RecordMode config")
		}
		for k, v := range cfg {
			t.Logf("RecordMode.%s = %s", k, v)
		}
	})

	t.Run("SetRecordMode", func(t *testing.T) {
		if !hasMediaFileFind {
			t.Skip("camera does not support RecordMode config")
		}
		original, err := c.Recording.GetRecordMode(ctx)
		if err != nil {
			t.Fatalf("GetRecordMode (save): %v", err)
		}

		// Find Mode key.
		modeKey := ""
		origVal := ""
		for k, v := range original {
			if contains(k, "Mode") {
				modeKey = k
				origVal = v
				break
			}
		}
		if modeKey == "" {
			t.Skip("no Mode key found in RecordMode config")
		}
		t.Logf("Original %s = %s", modeKey, origVal)

		setKey := modeKey
		if len(setKey) > 6 && setKey[:6] == "table." {
			setKey = setKey[6:]
		}

		defer func() {
			_ = c.Recording.SetRecordMode(ctx, map[string]string{
				setKey: origVal,
			})
		}()

		// Toggle between known modes (0=Auto, 1=Manual, 2=Off).
		newVal := "0"
		if origVal == "0" {
			newVal = "1"
		}
		err = c.Recording.SetRecordMode(ctx, map[string]string{
			setKey: newVal,
		})
		skipOnSetError(t, err, "SetRecordMode")

		updated, err := c.Recording.GetRecordMode(ctx)
		if err != nil {
			t.Fatalf("GetRecordMode (verify): %v", err)
		}
		if updated[modeKey] != newVal {
			t.Fatalf("expected %s=%q, got %q", modeKey, newVal, updated[modeKey])
		}
		t.Logf("Verified RecordMode changed to %q", newVal)
	})

	t.Run("GetMediaGlobal", func(t *testing.T) {
		if !hasMediaFileFind {
			t.Skip("camera does not support MediaGlobal config")
		}
		cfg, err := c.Recording.GetMediaGlobal(ctx)
		if err != nil {
			t.Fatalf("GetMediaGlobal: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty MediaGlobal config")
		}
		for k, v := range cfg {
			t.Logf("MediaGlobal.%s = %s", k, v)
		}
	})

	t.Run("FindFiles", func(t *testing.T) {
		if !hasMediaFileFind {
			t.Skip("camera does not support media file find")
		}
		now := time.Now()
		start := now.Add(-24 * time.Hour)

		opts := FindFilesOpts{
			Channel:   1,
			StartTime: start.Format("2006-01-02 15:04:05"),
			EndTime:   now.Format("2006-01-02 15:04:05"),
		}

		files, err := c.Recording.FindFiles(ctx, opts)
		if err != nil {
			t.Fatalf("FindFiles: %v", err)
		}

		t.Logf("Found %d files in last 24h", len(files))
		for i, f := range files {
			t.Logf("  [%d] Channel=%d Start=%s End=%s Type=%s Path=%s",
				i, f.Channel, f.StartTime, f.EndTime, f.Type, f.FilePath)
			if i >= 9 {
				t.Logf("  ... (showing first 10 of %d)", len(files))
				break
			}
		}
	})
}
