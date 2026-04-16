package amcrest

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"time"
)

// TestMultiCamera runs tests that require features only available on specific
// cameras. It automatically selects whichever camera supports the feature.
//
// Requires both AMCREST_HOST and AMCREST_HOST_2 configured in .env.
func TestMultiCamera(t *testing.T) {
	// Ensure primary client + caps are initialized
	c1 := testClient(t)
	initCaps(t, c1)

	ctx := context.Background()

	// ---------------------------------------------------------------
	// Two-way Audio (AD410 has audio output; IP5M does not)
	// ---------------------------------------------------------------
	t.Run("TwoWayAudio", func(t *testing.T) {
		c := testClientWithFeature(t, "audio output", func(client *Client) bool {
			n, err := client.Audio.GetOutputChannels(context.Background())
			return err == nil && n > 0
		})
		ctx := context.Background()

		t.Run("GetVolume", func(t *testing.T) {
			vol, err := c.Audio.GetVolume(ctx)
			if err != nil {
				t.Fatalf("GetVolume: %v", err)
			}
			for k, v := range vol {
				t.Logf("Volume.%s = %s", k, v)
			}
		})

		t.Run("SetVolume_SaveRestore", func(t *testing.T) {
			original, err := c.Audio.GetVolume(ctx)
			if err != nil {
				t.Fatalf("GetVolume (save): %v", err)
			}

			origVal := ""
			for k, v := range original {
				if strings.Contains(k, "AudioOutputVolume") {
					origVal = v
					break
				}
			}
			origInt, _ := strconv.Atoi(origVal)
			defer c.Audio.SetVolume(ctx, 0, origInt)

			newVol := 50
			if origInt == 50 {
				newVol = 45
			}
			err = c.Audio.SetVolume(ctx, 0, newVol)
			skipOnSetError(t, err, "SetVolume")

			updated, err := c.Audio.GetVolume(ctx)
			if err != nil {
				t.Fatalf("GetVolume (verify): %v", err)
			}
			for _, v := range updated {
				if v == strconv.Itoa(newVol) {
					t.Logf("Verified volume changed to %d", newVol)
					return
				}
			}
			t.Fatalf("volume not updated to %d", newVol)
		})
	})

	// ---------------------------------------------------------------
	// WiFi (AD410 has WiFi; IP5M does not)
	// ---------------------------------------------------------------
	t.Run("WiFi", func(t *testing.T) {
		c := testClientWithFeature(t, "WiFi", func(client *Client) bool {
			return probeConfig(context.Background(), client, "WLan")
		})
		ctx := context.Background()

		t.Run("GetWLanConfig", func(t *testing.T) {
			cfg, err := c.Network.GetWLanConfig(ctx)
			if err != nil {
				t.Fatalf("GetWLanConfig: %v", err)
			}
			for k, v := range cfg {
				t.Logf("WLan.%s = %s", k, v)
			}
		})

		t.Run("ScanWLanDevices", func(t *testing.T) {
			body, err := c.Network.ScanWLanDevices(ctx)
			if err != nil {
				t.Fatalf("ScanWLanDevices: %v", err)
			}
			t.Logf("Scan result:\n%s", body)
		})
	})

	// ---------------------------------------------------------------
	// Smart Motion Detection (IP5M has it; AD410 does not)
	// ---------------------------------------------------------------
	t.Run("SmartMotionDetect", func(t *testing.T) {
		c := testClientWithFeature(t, "SmartMotionDetect", func(client *Client) bool {
			return probeConfig(context.Background(), client, "SmartMotionDetect")
		})
		ctx := context.Background()

		t.Run("GetSmartMotionConfig", func(t *testing.T) {
			cfg, err := c.Motion.GetSmartMotionConfig(ctx)
			if err != nil {
				t.Fatalf("GetSmartMotionConfig: %v", err)
			}
			for k, v := range cfg {
				t.Logf("SMD.%s = %s", k, v)
			}
		})

		t.Run("SetSmartMotionConfig_SaveRestore", func(t *testing.T) {
			original, err := c.Motion.GetSmartMotionConfig(ctx)
			if err != nil {
				t.Fatalf("GetSmartMotionConfig (save): %v", err)
			}

			origSens := original["table.SmartMotionDetect[0].Sensitivity"]
			defer func() {
				c.Motion.SetSmartMotionConfig(ctx, map[string]string{
					"SmartMotionDetect[0].Sensitivity": origSens,
				})
			}()

			newSens := "Low"
			if origSens == "Low" {
				newSens = "Middle"
			}
			err = c.Motion.SetSmartMotionConfig(ctx, map[string]string{
				"SmartMotionDetect[0].Sensitivity": newSens,
			})
			skipOnSetError(t, err, "SetSmartMotionConfig")

			updated, err := c.Motion.GetSmartMotionConfig(ctx)
			if err != nil {
				t.Fatalf("GetSmartMotionConfig (verify): %v", err)
			}
			got := updated["table.SmartMotionDetect[0].Sensitivity"]
			if got != newSens {
				t.Fatalf("expected Sensitivity=%s, got %s", newSens, got)
			}
			t.Logf("Verified Sensitivity changed to %s", newSens)
		})
	})

	// ---------------------------------------------------------------
	// User Management (IP5M has 'user' group; AD410 only has 'admin')
	// ---------------------------------------------------------------
	t.Run("UserManagement", func(t *testing.T) {
		c := testClientWithFeature(t, "user management", func(client *Client) bool {
			body, err := client.User.GetAllGroups(context.Background())
			if err != nil {
				return false
			}
			return strings.Contains(strings.ToLower(body), "user")
		})
		ctx := context.Background()

		t.Run("FullUserCRUD", func(t *testing.T) {
			const testUser = "sdkcrudtest"
			const testPass = "CrudPass123!"
			const testGroup = "user"

			// Create
			err := c.User.AddUser(ctx, testUser, testPass, testGroup)
			if err != nil {
				t.Fatalf("AddUser: %v", err)
			}
			defer c.User.DeleteUser(ctx, testUser)
			t.Logf("Created user %q", testUser)

			// Read
			info, err := c.User.GetUserInfo(ctx, testUser)
			if err != nil {
				t.Fatalf("GetUserInfo: %v", err)
			}
			if info["user.Name"] != testUser {
				t.Fatalf("expected Name=%s, got %s", testUser, info["user.Name"])
			}
			t.Logf("Verified user exists: %s", info["user.Name"])

			// Update
			err = c.User.ModifyUser(ctx, testUser, map[string]string{
				"user.Memo": "SDK CRUD test account",
			})
			if err != nil {
				t.Fatalf("ModifyUser: %v", err)
			}
			info, _ = c.User.GetUserInfo(ctx, testUser)
			t.Logf("Modified memo: %s", info["user.Memo"])

			// Modify password
			err = c.User.ModifyPassword(ctx, testUser, testPass, "NewPass456!")
			if err != nil {
				t.Fatalf("ModifyPassword: %v", err)
			}
			t.Logf("Password changed successfully")

			// Delete happens via defer
		})
	})

	// ---------------------------------------------------------------
	// Email Config (IP5M has it; AD410 does not)
	// ---------------------------------------------------------------
	t.Run("EmailConfig", func(t *testing.T) {
		c := testClientWithFeature(t, "Email config", func(client *Client) bool {
			return probeConfig(context.Background(), client, "Email")
		})
		ctx := context.Background()

		t.Run("GetSetEmailConfig", func(t *testing.T) {
			original, err := c.Network.GetEmailConfig(ctx)
			if err != nil {
				t.Fatalf("GetEmailConfig: %v", err)
			}

			origEnable := original["Enable"]
			defer func() {
				c.Network.SetEmailConfig(ctx, map[string]string{
					"Email.Enable": origEnable,
				})
			}()

			// Toggle enable
			newVal := "true"
			if origEnable == "true" {
				newVal = "false"
			}
			err = c.Network.SetEmailConfig(ctx, map[string]string{
				"Email.Enable": newVal,
			})
			skipOnSetError(t, err, "SetEmailConfig")

			updated, err := c.Network.GetEmailConfig(ctx)
			if err != nil {
				t.Fatalf("GetEmailConfig (verify): %v", err)
			}
			if updated["Enable"] != newVal {
				t.Fatalf("expected Enable=%s, got %s", newVal, updated["Enable"])
			}
			t.Logf("Verified Email.Enable changed to %s", newVal)
		})
	})

	// ---------------------------------------------------------------
	// CrossLineDetection event (IP5M supports it; AD410 does not)
	// ---------------------------------------------------------------
	t.Run("CrossLineDetection", func(t *testing.T) {
		c := testClientWithFeature(t, "CrossLineDetection event", func(client *Client) bool {
			evts, err := client.Event.GetSupportedEvents(context.Background())
			if err != nil {
				return false
			}
			for _, e := range evts {
				if e == "CrossLineDetection" {
					return true
				}
			}
			return false
		})
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		events, stream, err := c.Event.Subscribe(ctx, []string{"CrossLineDetection"}, 2)
		if err != nil {
			t.Fatalf("Subscribe: %v", err)
		}
		defer stream.Close()

		// Just verify subscription works -- wait for heartbeat
		for event := range events {
			t.Logf("Received: Code=%s Action=%s", event.Code, event.Action)
			break // got at least one message (heartbeat or event)
		}
	})

	// ---------------------------------------------------------------
	// Doorbell-specific events (AD410 has these; IP5M does not)
	// ---------------------------------------------------------------
	t.Run("DoorbellEvents", func(t *testing.T) {
		c := testClientWithFeature(t, "doorbell (CrossRegionDetection)", func(client *Client) bool {
			evts, err := client.Event.GetSupportedEvents(context.Background())
			if err != nil {
				return false
			}
			// AD410 has CrossRegionDetection but not CrossLineDetection
			hasCross := false
			hasLine := false
			for _, e := range evts {
				if e == "CrossRegionDetection" {
					hasCross = true
				}
				if e == "CrossLineDetection" {
					hasLine = true
				}
			}
			return hasCross && !hasLine // likely a doorbell
		})
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		events, stream, err := c.Event.Subscribe(ctx, []string{"CrossRegionDetection"}, 2)
		if err != nil {
			t.Fatalf("Subscribe: %v", err)
		}
		defer stream.Close()

		for event := range events {
			t.Logf("Received: Code=%s Action=%s", event.Code, event.Action)
			break
		}
	})

	// ---------------------------------------------------------------
	// Both cameras: Run snapshot on each
	// ---------------------------------------------------------------
	t.Run("SnapshotBothCameras", func(t *testing.T) {
		c2 := testClient2(t)
		if c2 == nil {
			t.Skip("second camera not configured")
		}

		t.Run("Camera1", func(t *testing.T) {
			snap, err := c1.Snapshot.Get(ctx, 1)
			if err != nil {
				t.Fatalf("Camera1 Snapshot: %v", err)
			}
			if len(snap) < 2 || snap[0] != 0xFF || snap[1] != 0xD8 {
				t.Fatal("Camera1: not a valid JPEG")
			}
			t.Logf("Camera1 snapshot: %d bytes", len(snap))
		})

		t.Run("Camera2", func(t *testing.T) {
			snap, err := c2.Snapshot.Get(ctx, 1)
			if err != nil {
				t.Fatalf("Camera2 Snapshot: %v", err)
			}
			if len(snap) < 2 || snap[0] != 0xFF || snap[1] != 0xD8 {
				t.Fatal("Camera2: not a valid JPEG")
			}
			t.Logf("Camera2 snapshot: %d bytes", len(snap))
		})
	})

	// ---------------------------------------------------------------
	// Both cameras: Compare capabilities
	// ---------------------------------------------------------------
	t.Run("CompareCapabilities", func(t *testing.T) {
		c2 := testClient2(t)
		if c2 == nil {
			t.Skip("second camera not configured")
		}

		type camInfo struct {
			name   string
			client *Client
		}

		cams := []camInfo{
			{"Camera1", c1},
			{"Camera2", c2},
		}

		for _, cam := range cams {
			t.Run(cam.name, func(t *testing.T) {
				ctx := context.Background()

				devType, _ := cam.client.System.GetDeviceType(ctx)
				serial, _ := cam.client.System.GetSerialNumber(ctx)
				firmware, _ := cam.client.System.GetSoftwareVersion(ctx)
				audioIn, _ := cam.client.Audio.GetInputChannels(ctx)
				audioOut, _ := cam.client.Audio.GetOutputChannels(ctx)
				videoIn, _ := cam.client.Video.GetVideoInputChannels(ctx)
				events, _ := cam.client.Event.GetSupportedEvents(ctx)

				t.Logf("Type:     %s", devType)
				t.Logf("Serial:   %s", serial)
				t.Logf("Firmware: %s", firmware)
				t.Logf("Video In: %d, Audio In: %d, Audio Out: %d", videoIn, audioIn, audioOut)
				t.Logf("Events:   %s", strings.Join(events, ", "))

				smd := probeConfig(ctx, cam.client, "SmartMotionDetect")
				wifi := probeConfig(ctx, cam.client, "WLan")
				t.Logf("SmartMotion: %v, WiFi: %v", smd, wifi)
			})
		}
	})
}
