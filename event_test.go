package amcrest

import (
	"context"
	"testing"
	"time"
)

func TestEvent(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	t.Run("GetSupportedEvents", func(t *testing.T) {
		events, err := c.Event.GetSupportedEvents(ctx)
		if err != nil {
			t.Fatalf("GetSupportedEvents: %v", err)
		}
		if len(events) == 0 {
			t.Fatal("expected non-empty event list")
		}
		t.Logf("Supported events (%d): %v", len(events), events)
	})

	t.Run("GetCaps", func(t *testing.T) {
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
	})

	t.Run("Subscribe", func(t *testing.T) {
		subCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		ch, es, err := c.Event.Subscribe(subCtx, []string{"All"}, 1)
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
		case <-subCtx.Done():
			t.Fatal("timed out waiting for event or heartbeat")
		}
	})

	t.Run("GetAlarmInputChannels", func(t *testing.T) {
		if !hasAlarmInputCh {
			t.Skip("camera does not support alarm input channels")
		}
		n, err := c.Event.GetAlarmInputChannels(ctx)
		if err != nil {
			t.Fatalf("GetAlarmInputChannels: %v", err)
		}
		t.Logf("Alarm input channels: %d", n)
	})

	t.Run("GetAlarmOutputChannels", func(t *testing.T) {
		if !hasAlarmOutputCh {
			t.Skip("camera does not support alarm output channels")
		}
		n, err := c.Event.GetAlarmOutputChannels(ctx)
		if err != nil {
			t.Fatalf("GetAlarmOutputChannels: %v", err)
		}
		t.Logf("Alarm output channels: %d", n)
	})

	t.Run("GetBlindDetectConfig", func(t *testing.T) {
		if !supportsEvent("VideoBlind") {
			t.Skip("camera does not support VideoBlind event")
		}
		cfg, err := c.Event.GetBlindDetectConfig(ctx)
		if err != nil {
			t.Fatalf("GetBlindDetectConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty BlindDetect config")
		}
		for k, v := range cfg {
			t.Logf("BlindDetect.%s = %s", k, v)
		}
	})

	t.Run("GetLossDetectConfig", func(t *testing.T) {
		if !hasLossDetect {
			t.Skip("camera does not support LossDetect config")
		}
		cfg, err := c.Event.GetLossDetectConfig(ctx)
		if err != nil {
			t.Fatalf("GetLossDetectConfig: %v", err)
		}
		for k, v := range cfg {
			t.Logf("LossDetect.%s = %s", k, v)
		}
	})

	t.Run("GetLoginFailureAlarmConfig", func(t *testing.T) {
		if !supportsEvent("LoginFailure") {
			t.Skip("camera does not support LoginFailure event")
		}
		cfg, err := c.Event.GetLoginFailureAlarmConfig(ctx)
		if err != nil {
			t.Fatalf("GetLoginFailureAlarmConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty LoginFailureAlarm config")
		}
		for k, v := range cfg {
			t.Logf("LoginFailureAlarm.%s = %s", k, v)
		}
	})

	t.Run("GetStorageNotExistConfig", func(t *testing.T) {
		if !supportsEvent("StorageNotExist") {
			t.Skip("camera does not support StorageNotExist event")
		}
		cfg, err := c.Event.GetStorageNotExistConfig(ctx)
		if err != nil {
			t.Fatalf("GetStorageNotExistConfig: %v", err)
		}
		for k, v := range cfg {
			t.Logf("StorageNotExist.%s = %s", k, v)
		}
	})

	t.Run("GetStorageFailureConfig", func(t *testing.T) {
		if !supportsEvent("StorageFailure") {
			t.Skip("camera does not support StorageFailure event")
		}
		cfg, err := c.Event.GetStorageFailureConfig(ctx)
		if err != nil {
			t.Fatalf("GetStorageFailureConfig: %v", err)
		}
		for k, v := range cfg {
			t.Logf("StorageFailure.%s = %s", k, v)
		}
	})

	t.Run("GetStorageLowSpaceConfig", func(t *testing.T) {
		if !supportsEvent("StorageLowSpace") {
			t.Skip("camera does not support StorageLowSpace event")
		}
		cfg, err := c.Event.GetStorageLowSpaceConfig(ctx)
		if err != nil {
			t.Fatalf("GetStorageLowSpaceConfig: %v", err)
		}
		for k, v := range cfg {
			t.Logf("StorageLowSpace.%s = %s", k, v)
		}
	})

	t.Run("GetNetAbortConfig", func(t *testing.T) {
		if !hasNetAbort {
			t.Skip("camera does not support NetAbort config")
		}
		cfg, err := c.Event.GetNetAbortConfig(ctx)
		if err != nil {
			t.Fatalf("GetNetAbortConfig: %v", err)
		}
		for k, v := range cfg {
			t.Logf("NetAbort.%s = %s", k, v)
		}
	})

	t.Run("GetIPConflictConfig", func(t *testing.T) {
		if !hasIPConflict {
			t.Skip("camera does not support IPConflict config")
		}
		cfg, err := c.Event.GetIPConflictConfig(ctx)
		if err != nil {
			t.Fatalf("GetIPConflictConfig: %v", err)
		}
		for k, v := range cfg {
			t.Logf("IPConflict.%s = %s", k, v)
		}
	})

	t.Run("GetEventHandlerConfig_VideoMotion", func(t *testing.T) {
		if !hasEventHandler {
			t.Skip("camera does not support EventHandler config for VideoMotion")
		}
		cfg, err := c.Event.GetEventHandlerConfig(ctx, "VideoMotion")
		if err != nil {
			t.Fatalf("GetEventHandlerConfig(VideoMotion): %v", err)
		}
		for k, v := range cfg {
			t.Logf("EventHandler.VideoMotion.%s = %s", k, v)
		}
	})

	t.Run("GetAlarmConfig", func(t *testing.T) {
		if !hasAlarmConfig {
			t.Skip("camera does not support Alarm config")
		}
		cfg, err := c.Event.GetAlarmConfig(ctx)
		if err != nil {
			t.Fatalf("GetAlarmConfig: %v", err)
		}
		for k, v := range cfg {
			t.Logf("Alarm.%s = %s", k, v)
		}
	})

	t.Run("GetAlarmOutConfig", func(t *testing.T) {
		if !hasAlarmOut {
			t.Skip("camera does not support AlarmOut config")
		}
		cfg, err := c.Event.GetAlarmOutConfig(ctx)
		if err != nil {
			t.Fatalf("GetAlarmOutConfig: %v", err)
		}
		for k, v := range cfg {
			t.Logf("AlarmOut.%s = %s", k, v)
		}
	})
}
