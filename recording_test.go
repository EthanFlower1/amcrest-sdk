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
