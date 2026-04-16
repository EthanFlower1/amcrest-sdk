//go:build dangerous

package amcrest

import (
	"context"
	"testing"
)

func TestStorageFormatSDCard(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	if err := c.Storage.FormatSDCard(ctx, "/mnt/sd"); err != nil {
		t.Fatalf("FormatSDCard: %v", err)
	}
	t.Log("FormatSDCard command sent successfully")
}
