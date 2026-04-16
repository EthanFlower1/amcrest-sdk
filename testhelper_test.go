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

// probeConfig checks whether a named config exists on this camera.
// Returns true if getRawConfig succeeds (HTTP 200), false on any error.
func probeConfig(ctx context.Context, c *Client, name string) bool {
	_, err := c.getRawConfig(ctx, name)
	return err == nil
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

	// Feature-level booleans (probed via endpoints).
	hasPTZ        bool
	hasAnalytics  bool
	hasThermal    bool
	hasAccessCtrl bool
	hasFaceRec    bool
	hasGPS        bool
	hasBuilding   bool
	hasTraffic    bool
	hasParking    bool
	hasPeople     bool
	hasWorkSuit   bool
	audioOutChans int

	// Upgrade-related probes.
	hasUpgrade      bool // upgrader.cgi?action=getState
	hasCloudUpgrade bool // api/CloudUpgrader/check
	hasAutoUpgrade  bool // config AutoUpgrade
	hasCloudMode    bool // config CloudUpgrade

	// Camera feature probes from videoInputCaps.
	hasLighting     bool // caps.InfraRed == "true" or Lighting config exists
	hasFocusControl bool // caps.ElectricFocus == "true"
	hasPrivacyMask  bool // caps.CoverCount > 0

	// Config-based probes.
	hasUploadPicture  bool // PictureHttpUpload
	hasUploadEvent    bool // EventHttpUpload
	hasUploadReport   bool // ReportHttpUpload
	hasSmartMotion    bool // SmartMotionDetect
	hasLAEConfig      bool // LAEConfig
	hasGUISet         bool // GUISet
	hasSplitScreen    bool // split.cgi
	hasMonitorTour    bool // MonitorTour
	hasMonitorCollect bool // MonitorCollection
	hasPrivacyConfig  bool // PrivacyMasking
	hasVideoInOptions bool // VideoInOptions
	hasWLan           bool // WLan
	hasSSHD           bool // SSHD
	hasLossDetect     bool // LossDetect
	hasNetAbort       bool // NetAbort
	hasIPConflict     bool // IPConflict
	hasAlarmConfig    bool // Alarm
	hasAlarmOut       bool // AlarmOut
	hasCoaxialIO      bool // CoaxialControlIO caps
	hasFlashlight     bool // FlashLight config
	hasFishEye        bool // FishEye config
	hasDVR            bool // DVR media find
	hasBandwidthLimit bool // BandwidthLimit config (DVR)
	hasLanguage       bool // language.cgi
	hasOnvifVersion   bool // onvif version endpoint
	hasHTTPAPIVersion bool // HTTP API version endpoint
	hasSnapWithType   bool // snapshot with type parameter
	hasSmartEncode    bool // SmartEncode config
	hasEncodeROI      bool // VideoEncodeROI config
	hasVideoOutput    bool // devVideoOutput channels
	hasEncodeConfCaps bool // EncodeConfigCaps
	hasViewRange      bool // PTZ ViewRange
	hasEPTZ           bool // EPTZ config
	hasAutoMovement   bool // PTZ AutoMovement config
	hasRedList        bool // TrafficRedList
	hasUserMgmt       bool // AddUser/DeleteUser
	hasEventHandler   bool // EventHandler config for VideoMotion
	hasAlarmInputCh   bool // alarm input channels endpoint
	hasAlarmOutputCh  bool // alarm output channels endpoint
	hasLogFind        bool // log.cgi startFind
	hasEmail          bool // Email config
	hasPrivacyMaskCGI bool // PrivacyMasking.cgi endpoints
	hasRecordCaps     bool // recordManager.cgi getCaps
	hasMediaFileFind  bool // mediaFileFind
	hasStorageDevInfo bool // storageDevice.cgi getDeviceAllInfo
	hasStorageCollect bool // storageDevice.cgi factory.getCollect
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

			// Derive camera feature booleans from videoInputCaps.
			if strings.EqualFold(caps["caps.InfraRed"], "true") {
				hasLighting = true
			}
			if strings.EqualFold(caps["caps.ElectricFocus"], "true") {
				hasFocusControl = true
			}
			if capsInt(caps, "caps.CoverCount") > 0 {
				hasPrivacyMask = true
			}
		}

		// PTZ - probe via GetStatus; 400 means no PTZ support.
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

		// --- Upgrade probes ---
		if _, err := c.Upgrade.GetState(ctx); err == nil {
			hasUpgrade = true
		}
		if _, err := c.Upgrade.CheckCloudUpdate(ctx); err == nil {
			hasCloudUpgrade = true
		}

		// --- Config-based probes ---
		hasAutoUpgrade = probeConfig(ctx, c, "AutoUpgrade")
		hasCloudMode = probeConfig(ctx, c, "CloudUpgrade")
		hasUploadPicture = probeConfig(ctx, c, "PictureHttpUpload")
		hasUploadEvent = probeConfig(ctx, c, "EventHttpUpload")
		hasUploadReport = probeConfig(ctx, c, "ReportHttpUpload")
		hasSmartMotion = probeConfig(ctx, c, "SmartMotionDetect")
		hasLAEConfig = probeConfig(ctx, c, "LAEConfig")
		hasGUISet = probeConfig(ctx, c, "GUISet")
		hasMonitorTour = probeConfig(ctx, c, "MonitorTour")
		hasMonitorCollect = probeConfig(ctx, c, "MonitorCollection")
		hasPrivacyConfig = probeConfig(ctx, c, "PrivacyMasking")
		hasVideoInOptions = probeConfig(ctx, c, "VideoInOptions")
		hasWLan = probeConfig(ctx, c, "WLan")
		hasSSHD = probeConfig(ctx, c, "SSHD")
		hasLossDetect = probeConfig(ctx, c, "LossDetect")
		hasNetAbort = probeConfig(ctx, c, "NetAbort")
		hasIPConflict = probeConfig(ctx, c, "IPConflict")
		hasAlarmConfig = probeConfig(ctx, c, "Alarm")
		hasAlarmOut = probeConfig(ctx, c, "AlarmOut")
		hasFlashlight = probeConfig(ctx, c, "FlashLight")
		hasFishEye = probeConfig(ctx, c, "FishEye")
		hasSmartEncode = probeConfig(ctx, c, "SmartEncode")
		hasEncodeROI = probeConfig(ctx, c, "VideoEncodeROI")
		hasBandwidthLimit = probeConfig(ctx, c, "BandwidthLimit")

		// If lighting wasn't detected from caps, try config probe.
		if !hasLighting {
			hasLighting = probeConfig(ctx, c, "Lighting")
		}

		// Non-config endpoint probes.
		if _, err := c.Peripheral.GetCoaxialIOStatus(ctx, 0); err == nil {
			hasCoaxialIO = true
		}
		if _, err := c.Display.GetSplitMode(ctx, 0); err == nil {
			hasSplitScreen = true
		}
		if _, err := c.Video.GetVideoOutputChannels(ctx); err == nil {
			hasVideoOutput = true
		}
		if _, err := c.Video.GetEncodeConfigCaps(ctx, 0); err == nil {
			hasEncodeConfCaps = true
		}
		if _, err := c.DVR.StartFind(ctx, 0, "2024-01-01 00:00:00", "2024-01-01 00:01:00"); err == nil {
			hasDVR = true
		}
		if _, err := c.System.GetLanguage(ctx); err == nil {
			hasLanguage = true
		}
		if _, err := c.System.GetOnvifVersion(ctx); err == nil {
			hasOnvifVersion = true
		}
		if _, err := c.System.GetHTTPAPIVersion(ctx); err == nil {
			hasHTTPAPIVersion = true
		}
		if _, err := c.Snapshot.GetWithType(ctx, 1, 0); err == nil {
			hasSnapWithType = true
		}

		// PTZ sub-features (only if hasPTZ).
		if hasPTZ {
			if _, err := c.PTZ.GetViewRangeStatus(ctx, 0); err == nil {
				hasViewRange = true
			}
			hasEPTZ = probeConfig(ctx, c, "Ptz")
			if _, err := c.PTZ.GetAutoMovementConfig(ctx); err == nil {
				hasAutoMovement = true
			}
		}

		// Traffic sub-feature.
		if hasTraffic {
			if _, err := c.Traffic.FindRecord(ctx, "TrafficRedList"); err == nil {
				hasRedList = true
			}
		}

		// User management: probe by checking if "user" group exists.
		if body, err := c.User.GetAllGroups(ctx); err == nil {
			hasUserMgmt = strings.Contains(strings.ToLower(body), "user")
		}

		// EventHandler config for VideoMotion.
		if supportsEvent("VideoMotion") {
			if _, err := c.Event.GetEventHandlerConfig(ctx, "VideoMotion"); err == nil {
				hasEventHandler = true
			}
		}

		// Alarm channels (not all devices have alarm I/O).
		if _, err := c.Event.GetAlarmInputChannels(ctx); err == nil {
			hasAlarmInputCh = true
		}
		if _, err := c.Event.GetAlarmOutputChannels(ctx); err == nil {
			hasAlarmOutputCh = true
		}

		// Log find (some doorbells/small devices don't support log.cgi).
		hasLogFind = true
		if _, err := c.Log.Find(ctx, "2024-01-01 00:00:00", "2024-01-01 00:01:00", ""); err != nil {
			hasLogFind = false
		}

		// Email config.
		hasEmail = probeConfig(ctx, c, "Email")

		// PrivacyMasking.cgi endpoints (separate from config).
		if _, err := c.Privacy.GetEnable(ctx, 0); err == nil {
			hasPrivacyMaskCGI = true
		}

		// Recording caps (recordManager.cgi).
		if _, err := c.Recording.GetCaps(ctx); err == nil {
			hasRecordCaps = true
		}

		// MediaFileFind - probe with a very short time range to see if
		// the mediaFileFind CGI exists.
		if files, err := c.Recording.FindFiles(ctx, FindFilesOpts{
			Channel:   1,
			StartTime: "2024-01-01 00:00:00",
			EndTime:   "2024-01-01 00:01:00",
		}); err == nil {
			hasMediaFileFind = true
			_ = files
		}

		// Storage device info endpoints.
		if _, err := c.Storage.GetAllDeviceInfo(ctx); err == nil {
			hasStorageDevInfo = true
		}
		if body, err := c.Storage.GetDeviceNames(ctx); err == nil && body != "" {
			hasStorageCollect = true
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
