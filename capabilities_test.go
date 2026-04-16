package amcrest

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"
)

// CameraProfile captures the full capability report for a connected camera.
type CameraProfile struct {
	// Device identity
	DeviceType    string `json:"device_type"`
	DeviceClass   string `json:"device_class"`
	SerialNumber  string `json:"serial_number"`
	MachineID     string `json:"machine_name"`
	Vendor        string `json:"vendor"`
	Firmware      string `json:"firmware_version"`
	Hardware      string `json:"hardware_version"`
	LanguageCaps  string `json:"language_caps"`
	VideoStandard string `json:"video_standard"`

	// Channels
	VideoInputChannels  int `json:"video_input_channels"`
	VideoOutputChannels int `json:"video_output_channels"`
	AudioInputChannels  int `json:"audio_input_channels"`
	AudioOutputChannels int `json:"audio_output_channels"`
	MaxExtraStreams      int `json:"max_extra_streams"`

	// Events
	SupportedEvents []string          `json:"supported_events"`
	EventCaps       map[string]string `json:"event_caps"`

	// Video input capabilities
	VideoInputCaps map[string]string `json:"video_input_caps"`

	// Encoding capabilities
	EncodeCaps map[string]string `json:"encode_caps"`

	// Recording capabilities
	RecordCaps map[string]string `json:"record_caps"`

	// Storage capabilities
	StorageCaps map[string]string `json:"storage_caps"`

	// Feature support
	Features map[string]bool `json:"features"`
}

func TestCameraProfile(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	profile := CameraProfile{
		Features: make(map[string]bool),
	}

	// Device identity
	if v, err := c.System.GetDeviceType(ctx); err == nil {
		profile.DeviceType = v
	}
	if v, err := c.System.GetDeviceClass(ctx); err == nil {
		profile.DeviceClass = v
	}
	if v, err := c.System.GetSerialNumber(ctx); err == nil {
		profile.SerialNumber = v
	}
	if v, err := c.System.GetMachineName(ctx); err == nil {
		profile.MachineID = v
	}
	if v, err := c.System.GetVendor(ctx); err == nil {
		profile.Vendor = v
	}
	if v, err := c.System.GetSoftwareVersion(ctx); err == nil {
		profile.Firmware = v
	}
	if v, err := c.System.GetHardwareVersion(ctx); err == nil {
		profile.Hardware = v
	}
	if v, err := c.System.GetLanguageCaps(ctx); err == nil {
		profile.LanguageCaps = v
	}
	if v, err := c.Video.GetVideoStandard(ctx); err == nil {
		profile.VideoStandard = v
	}

	// Channel counts
	if v, err := c.Video.GetVideoInputChannels(ctx); err == nil {
		profile.VideoInputChannels = v
	}
	if v, err := c.Video.GetVideoOutputChannels(ctx); err == nil {
		profile.VideoOutputChannels = v
	}
	if v, err := c.Audio.GetInputChannels(ctx); err == nil {
		profile.AudioInputChannels = v
	}
	if v, err := c.Audio.GetOutputChannels(ctx); err == nil {
		profile.AudioOutputChannels = v
	}
	if v, err := c.Video.GetMaxExtraStreams(ctx); err == nil {
		profile.MaxExtraStreams = v
	}

	// Events
	profile.SupportedEvents = supportedEvents
	profile.EventCaps = eventCaps

	// Capabilities maps
	profile.VideoInputCaps = videoInputCaps
	profile.EncodeCaps = encodeCaps
	profile.RecordCaps = recordCaps
	profile.StorageCaps = storageCaps

	// Feature support (from initCaps probing)
	profile.Features["ptz"] = hasPTZ
	profile.Features["analytics"] = hasAnalytics
	profile.Features["thermal"] = hasThermal
	profile.Features["access_control"] = hasAccessCtrl
	profile.Features["face_recognition"] = hasFaceRec
	profile.Features["gps"] = hasGPS
	profile.Features["building_intercom"] = hasBuilding
	profile.Features["traffic"] = hasTraffic
	profile.Features["parking"] = hasParking
	profile.Features["people_counting"] = hasPeople
	profile.Features["worksuit_detection"] = hasWorkSuit
	profile.Features["audio_output"] = audioOutChans > 0
	profile.Features["user_management"] = hasUserMgmt

	// Probe additional features and add to map
	profile.Features["smart_motion_detect"] = probeConfig(ctx, c, "SmartMotionDetect")
	profile.Features["privacy_masking"] = videoInputCaps != nil && capsInt(videoInputCaps, "caps.CoverCount") > 0
	profile.Features["electric_focus"] = videoInputCaps != nil && videoInputCaps["caps.ElectricFocus"] == "true"
	profile.Features["infrared"] = videoInputCaps != nil && videoInputCaps["caps.InfraRed"] == "true"
	profile.Features["day_night_color"] = videoInputCaps != nil && videoInputCaps["caps.DayNightColor"] == "true"
	profile.Features["flip"] = videoInputCaps != nil && videoInputCaps["caps.Flip"] == "true"
	profile.Features["mirror"] = videoInputCaps != nil && videoInputCaps["caps.Mirror"] == "true"
	profile.Features["rotate90"] = videoInputCaps != nil && videoInputCaps["caps.Rotate90"] == "true"
	profile.Features["wide_dynamic_range"] = videoInputCaps != nil && capsInt(videoInputCaps, "caps.WideDynamicRange") > 0
	profile.Features["defog"] = videoInputCaps != nil && videoInputCaps["caps.Defog"] == "true"
	profile.Features["white_balance"] = videoInputCaps != nil && capsInt(videoInputCaps, "caps.WhiteBalance") > 0
	profile.Features["backlight"] = videoInputCaps != nil && capsInt(videoInputCaps, "caps.Backlight") > 0
	profile.Features["gain_control"] = videoInputCaps != nil && videoInputCaps["caps.Gain"] == "true"
	profile.Features["image_stabilization"] = videoInputCaps != nil && videoInputCaps["caps.ImageStabilization"] == "true"
	profile.Features["fisheye"] = videoInputCaps != nil && videoInputCaps["caps.FishEye"] == "true"

	profile.Features["mail_notification"] = eventCaps != nil && eventCaps["caps.MailEnable"] == "true"
	profile.Features["record_on_event"] = eventCaps != nil && eventCaps["caps.RecordEnable"] == "true"
	profile.Features["snapshot_on_event"] = eventCaps != nil && eventCaps["caps.SnapshotEnable"] == "true"
	profile.Features["cloud_record"] = eventCaps != nil && eventCaps["caps.RecordCloudEnable"] == "true"
	profile.Features["light_control"] = eventCaps != nil && eventCaps["caps.SupportLightControl"] == "true"
	profile.Features["voice_alert"] = eventCaps != nil && eventCaps["caps.VoiceEnable"] == "true"

	profile.Features["sd_card"] = storageCaps != nil && storageCaps["caps.IsLocalStore"] == "true"
	profile.Features["remote_storage"] = storageCaps != nil && storageCaps["caps.IsRemoteStore"] == "true"
	profile.Features["holiday_recording"] = recordCaps != nil && recordCaps["caps.SupportHoliday"] == "true"

	// Print summary to test log
	t.Logf("=== Camera Profile ===")
	t.Logf("Device:   %s (%s)", profile.DeviceType, profile.DeviceClass)
	t.Logf("Serial:   %s", profile.SerialNumber)
	t.Logf("Name:     %s", profile.MachineID)
	t.Logf("Vendor:   %s", profile.Vendor)
	t.Logf("Firmware: %s", profile.Firmware)
	t.Logf("Hardware: %s", profile.Hardware)
	t.Logf("")
	t.Logf("Channels: %d video in, %d video out, %d audio in, %d audio out, %d extra streams",
		profile.VideoInputChannels, profile.VideoOutputChannels,
		profile.AudioInputChannels, profile.AudioOutputChannels,
		profile.MaxExtraStreams)
	t.Logf("")
	t.Logf("Supported Events (%d):", len(profile.SupportedEvents))
	for _, e := range profile.SupportedEvents {
		t.Logf("  - %s", e)
	}
	t.Logf("")
	t.Logf("Features:")
	for name, supported := range profile.Features {
		status := "NO"
		if supported {
			status = "YES"
		}
		t.Logf("  %-30s %s", name, status)
	}

	// Save to file
	data, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		t.Fatalf("marshaling profile: %v", err)
	}

	outFile := "camera_profile.json"
	if err := os.WriteFile(outFile, data, 0644); err != nil {
		t.Fatalf("writing profile: %v", err)
	}
	t.Logf("")
	t.Logf("Full profile saved to %s", outFile)
}
