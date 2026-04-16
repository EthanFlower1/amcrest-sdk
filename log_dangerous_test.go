//go:build dangerous

package amcrest

import (
	"context"
	"testing"
)

func TestLogClear(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	if err := c.Log.Clear(ctx); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	t.Log("Log clear command sent successfully")
}
