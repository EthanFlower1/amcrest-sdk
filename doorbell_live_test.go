//go:build live

package amcrest

import (
	"fmt"
	"context"
	"testing"
	"time"
)

// Run with: go test -run TestLiveEventListener -v -tags live -timeout 120s
//
// Subscribes to ALL events and prints them in real time.
// Walk in front of the camera to trigger events. Ctrl+C to stop.
func TestLiveEventListener(t *testing.T) {
	c := testClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	fmt.Println("")
	fmt.Println("========================================")
	fmt.Println("  LIVE EVENT LISTENER")
	fmt.Println("  Listening for 90 seconds...")
	fmt.Println("  Walk in front of the camera now!")
	fmt.Println("========================================")
	fmt.Println("")

	events, stream, err := c.Event.Subscribe(ctx, []string{"All"}, 5)
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}
	defer stream.Close()

	count := 0
	for event := range events {
		if event.Code == "Heartbeat" {
			fmt.Printf("  [heartbeat]\n")
			continue
		}

		count++
		ts := time.Now().Format("15:04:05")
		fmt.Printf("\n>>> EVENT #%d at %s <<<\n", count, ts)
		fmt.Printf("  Code:   %s\n", event.Code)
		fmt.Printf("  Action: %s\n", event.Action)
		fmt.Printf("  Index:  %d\n", event.Index)
		if jsonData, ok := event.Data["data"]; ok {
			fmt.Printf("  Data:\n%s\n", jsonData)
		}
	}

	fmt.Printf("\n========================================\n")
	fmt.Printf("  Done. Received %d events.\n", count)
	fmt.Printf("========================================\n")
}
