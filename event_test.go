package amcrest

import (
	"context"
	"testing"
	"time"
)

func TestGetSupportedEvents(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	events, err := c.Event.GetSupportedEvents(ctx)
	if err != nil {
		t.Fatalf("GetSupportedEvents: %v", err)
	}
	if len(events) == 0 {
		t.Fatal("expected non-empty event list")
	}
	t.Logf("Supported events (%d): %v", len(events), events)
}

func TestGetCaps(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	caps, err := c.Event.GetCaps(ctx)
	if err != nil {
		t.Fatalf("GetCaps: %v", err)
	}
	if len(caps) == 0 {
		t.Fatal("expected non-empty capabilities map")
	}
	t.Logf("Event caps (%d entries):", len(caps))
	for k, v := range caps {
		t.Logf("  %s = %s", k, v)
	}
}

func TestSubscribe(t *testing.T) {
	c := testClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ch, es, err := c.Event.Subscribe(ctx, []string{"All"}, 1)
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}
	defer es.Close()

	select {
	case evt, ok := <-ch:
		if !ok {
			t.Fatal("event channel closed without receiving an event")
		}
		t.Logf("Received event: Code=%s Action=%s Raw=%s", evt.Code, evt.Action, evt.Raw)
	case <-ctx.Done():
		t.Fatal("timed out waiting for event or heartbeat")
	}
}

func TestGetAlarmInputChannels(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	n, err := c.Event.GetAlarmInputChannels(ctx)
	if err != nil {
		t.Fatalf("GetAlarmInputChannels: %v", err)
	}
	t.Logf("Alarm input channels: %d", n)
}

func TestGetAlarmOutputChannels(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	n, err := c.Event.GetAlarmOutputChannels(ctx)
	if err != nil {
		t.Fatalf("GetAlarmOutputChannels: %v", err)
	}
	t.Logf("Alarm output channels: %d", n)
}
