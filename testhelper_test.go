package amcrest

import (
	"bufio"
	"context"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func loadEnv() {
	f, err := os.Open(".env")
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			if os.Getenv(parts[0]) == "" {
				os.Setenv(parts[0], parts[1])
			}
		}
	}
}

func testClient(t *testing.T) *Client {
	t.Helper()
	loadEnv()
	host := os.Getenv("AMCREST_HOST")
	user := os.Getenv("AMCREST_USERNAME")
	pass := os.Getenv("AMCREST_PASSWORD")
	if host == "" || user == "" || pass == "" {
		t.Skip("AMCREST_HOST, AMCREST_USERNAME, AMCREST_PASSWORD must be set")
	}
	client, err := NewClient(host, user, pass)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return client
}

// Capability cache - queried once per test run.
var (
	capsOnce        sync.Once
	supportedEvents []string
	eventCaps       map[string]string
	encodeCaps      map[string]string
	recordCaps      map[string]string
	storageCaps     map[string]string
	videoInputCaps  map[string]string
	ptzCapsRaw      string
	hasPTZ          bool
	hasAnalytics    bool
	hasThermal      bool
	hasAccessCtrl   bool
	hasFaceRec      bool
	hasGPS          bool
	hasBuilding     bool
	hasTraffic      bool
	hasParking      bool
	hasPeople       bool
	hasWorkSuit     bool
	audioOutChans   int
)

func initCaps(t *testing.T, c *Client) {
	t.Helper()
	capsOnce.Do(func() {
		ctx := context.Background()

		// Supported events
		if evts, err := c.Event.GetSupportedEvents(ctx); err == nil {
			supportedEvents = evts
		}

		// Event caps
		if caps, err := c.Event.GetCaps(ctx); err == nil {
			eventCaps = caps
		}

		// Encode caps
		if caps, err := c.Video.GetEncodeCaps(ctx); err == nil {
			encodeCaps = caps
		}

		// Record caps
		if caps, err := c.Recording.GetCaps(ctx); err == nil {
			recordCaps = caps
		}

		// Storage caps
		if caps, err := c.Storage.GetCaps(ctx); err == nil {
			storageCaps = caps
		}

		// Video input caps (channel 0)
		if caps, err := c.Video.GetVideoInputCaps(ctx, 0); err == nil {
			videoInputCaps = caps
		}

		// PTZ - probe via GetStatus; 400 means no PTZ support
		if _, err := c.PTZ.GetStatus(ctx, 0); err == nil {
			hasPTZ = true
		}
		if raw, err := c.PTZ.GetCaps(ctx, 0); err == nil {
			ptzCapsRaw = raw
		}

		// Analytics
		if _, err := c.Analytics.GetCaps(ctx, 0); err == nil {
			hasAnalytics = true
		}

		// Thermal
		if _, err := c.Thermal.GetCaps(ctx); err == nil {
			hasThermal = true
		}

		// Access Control
		if _, err := c.AccessControl.GetDoorStatus(ctx, 0); err == nil {
			hasAccessCtrl = true
		}

		// Face Recognition
		if _, err := c.Face.FindGroup(ctx); err == nil {
			hasFaceRec = true
		}

		// GPS
		if _, err := c.Peripheral.GetGPSCaps(ctx); err == nil {
			hasGPS = true
		}

		// Building
		if _, err := c.Building.GetSIPConfig(ctx); err == nil {
			hasBuilding = true
		}

		// Traffic
		if _, err := c.Traffic.FindRecord(ctx, "TrafficBlackList"); err == nil {
			hasTraffic = true
		}

		// Parking
		if _, err := c.Parking.GetSpaceStatus(ctx, 0, 0); err == nil {
			hasParking = true
		}

		// People counting
		if _, err := c.People.GetSummary(ctx); err == nil {
			hasPeople = true
		}

		// WorkSuit
		if _, err := c.WorkSuit.FindGroup(ctx); err == nil {
			hasWorkSuit = true
		}

		// Audio output channels
		if n, err := c.Audio.GetOutputChannels(ctx); err == nil {
			audioOutChans = n
		}
	})
}

func supportsEvent(event string) bool {
	for _, e := range supportedEvents {
		if e == event {
			return true
		}
	}
	return false
}

func requireCapability(t *testing.T, supported bool, name string) {
	t.Helper()
	if !supported {
		t.Skipf("camera does not support %s", name)
	}
}

func capsInt(caps map[string]string, key string) int {
	if caps == nil {
		return 0
	}
	v, ok := caps[key]
	if !ok {
		return 0
	}
	n, _ := strconv.Atoi(strings.TrimSpace(v))
	return n
}
