package amcrest

import (
	"context"
	"testing"
	"time"
)

func TestLogFind(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	end := time.Now().Format("2006-01-02 15:04:05")
	start := time.Now().Add(-24 * time.Hour).Format("2006-01-02 15:04:05")

	entries, err := c.Log.Find(ctx, start, end, "")
	if err != nil {
		t.Fatalf("Find: %v", err)
	}
	t.Logf("Found %d log entries in last 24h", len(entries))
	for i, e := range entries {
		if i >= 5 {
			t.Logf("  ... (%d more entries)", len(entries)-5)
			break
		}
		t.Logf("  [%d] Time=%s Type=%s User=%s Detail=%s", i, e.Time, e.Type, e.User, e.Detail)
	}
}

func TestLogBackup(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	end := time.Now().Format("2006-01-02 15:04:05")
	start := time.Now().Add(-24 * time.Hour).Format("2006-01-02 15:04:05")

	data, err := c.Log.Backup(ctx, start, end)
	if err != nil {
		t.Fatalf("Backup: %v", err)
	}
	t.Logf("Backup returned %d bytes", len(data))
}
