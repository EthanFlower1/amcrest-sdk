package amcrest

import (
	"context"
	"testing"
	"time"
)

func TestRecordingGetCaps(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
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
}

func TestRecordingGetRecordConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
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
}

func TestRecordingGetRecordMode(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
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
}

func TestRecordingGetMediaGlobal(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
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
}

func TestRecordingFindFiles(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24*time.Hour - time.Second)

	opts := FindFilesOpts{
		Channel:   1,
		StartTime: startOfDay.Format("2006-01-02 15:04:05"),
		EndTime:   endOfDay.Format("2006-01-02 15:04:05"),
	}

	files, err := c.Recording.FindFiles(ctx, opts)
	if err != nil {
		t.Fatalf("FindFiles: %v", err)
	}

	t.Logf("Found %d files for today", len(files))
	for i, f := range files {
		t.Logf("  [%d] Channel=%d Start=%s End=%s Type=%s Path=%s Length=%d Duration=%d",
			i, f.Channel, f.StartTime, f.EndTime, f.Type, f.FilePath, f.Length, f.Duration)
		if i >= 9 {
			t.Logf("  ... (showing first 10 of %d)", len(files))
			break
		}
	}
}
