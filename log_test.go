package amcrest

import (
	"context"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	t.Run("Find", func(t *testing.T) {
		if !hasLogFind {
			t.Skip("camera does not support log.cgi find")
		}
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
	})

	t.Run("Backup", func(t *testing.T) {
		if !hasLogFind {
			t.Skip("camera does not support log.cgi")
		}
		end := time.Now().Format("2006-01-02 15:04:05")
		start := time.Now().Add(-24 * time.Hour).Format("2006-01-02 15:04:05")

		data, err := c.Log.Backup(ctx, start, end)
		if err != nil {
			t.Fatalf("Backup: %v", err)
		}
		t.Logf("Backup returned %d bytes", len(data))
	})
}
