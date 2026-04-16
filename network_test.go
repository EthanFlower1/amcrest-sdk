package amcrest

import (
	"context"
	"testing"
)

func TestNetwork(t *testing.T) {
	c := testClient(t)
	initCaps(t, c)
	ctx := context.Background()

	t.Run("GetInterfaces", func(t *testing.T) {
		v, err := c.Network.GetInterfaces(ctx)
		if err != nil {
			t.Fatalf("GetInterfaces: %v", err)
		}
		if v == "" {
			t.Fatal("expected non-empty interfaces response")
		}
		t.Logf("Interfaces:\n%s", v)
	})

	t.Run("GetAccessFilter", func(t *testing.T) {
		cfg, err := c.Network.GetAccessFilter(ctx)
		if err != nil {
			t.Fatalf("GetAccessFilter: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty AccessFilter config")
		}
		for k, v := range cfg {
			t.Logf("AccessFilter.%s = %s", k, v)
		}
	})

	t.Run("GetNetworkConfig", func(t *testing.T) {
		cfg, err := c.Network.GetNetworkConfig(ctx)
		if err != nil {
			t.Fatalf("GetNetworkConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty Network config")
		}
		for k, v := range cfg {
			t.Logf("%s = %s", k, v)
		}
	})

	t.Run("GetDDNSConfig", func(t *testing.T) {
		cfg, err := c.Network.GetDDNSConfig(ctx)
		if err != nil {
			t.Fatalf("GetDDNSConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty DDNS config")
		}
		for k, v := range cfg {
			t.Logf("DDNS.%s = %s", k, v)
		}
	})

	t.Run("GetEmailConfig", func(t *testing.T) {
		if !hasEmail {
			t.Skip("camera does not support Email config")
		}
		cfg, err := c.Network.GetEmailConfig(ctx)
		if err != nil {
			t.Fatalf("GetEmailConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty Email config")
		}
		for k, v := range cfg {
			t.Logf("Email.%s = %s", k, v)
		}
	})

	t.Run("SetEmailConfig", func(t *testing.T) {
		if !hasEmail {
			t.Skip("camera does not support Email config")
		}
		original, err := c.Network.GetEmailConfig(ctx)
		if err != nil {
			t.Fatalf("GetEmailConfig (save): %v", err)
		}
		origEnable := original["Enable"]
		t.Logf("Original Email.Enable: %s", origEnable)

		defer func() {
			_ = c.Network.SetEmailConfig(ctx, map[string]string{
				"Email.Enable": origEnable,
			})
		}()

		newEnable := "true"
		if origEnable == "true" {
			newEnable = "false"
		}
		err = c.Network.SetEmailConfig(ctx, map[string]string{
			"Email.Enable": newEnable,
		})
		skipOnSetError(t, err, "SetEmailConfig")

		updated, err := c.Network.GetEmailConfig(ctx)
		if err != nil {
			t.Fatalf("GetEmailConfig (verify): %v", err)
		}
		if updated["Enable"] != newEnable {
			t.Fatalf("expected Email.Enable=%q, got %q", newEnable, updated["Enable"])
		}
		t.Logf("Verified Email.Enable changed to %q", newEnable)
	})

	t.Run("GetWLanConfig", func(t *testing.T) {
		if !hasWLan {
			t.Skip("camera does not support WLan config")
		}
		cfg, err := c.Network.GetWLanConfig(ctx)
		if err != nil {
			t.Fatalf("GetWLanConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Log("WLan config is empty (WiFi may not be configured)")
		}
		for k, v := range cfg {
			t.Logf("WLan.%s = %s", k, v)
		}
	})

	t.Run("GetUPnPConfig", func(t *testing.T) {
		cfg, err := c.Network.GetUPnPConfig(ctx)
		if err != nil {
			t.Fatalf("GetUPnPConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty UPnP config")
		}
		for k, v := range cfg {
			t.Logf("UPnP.%s = %s", k, v)
		}
	})

	t.Run("GetUPnPStatus", func(t *testing.T) {
		cfg, err := c.Network.GetUPnPStatus(ctx)
		if err != nil {
			t.Fatalf("GetUPnPStatus: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty UPnP status")
		}
		for k, v := range cfg {
			t.Logf("UPnP.%s = %s", k, v)
		}
	})

	t.Run("GetNTPConfig", func(t *testing.T) {
		cfg, err := c.Network.GetNTPConfig(ctx)
		if err != nil {
			t.Fatalf("GetNTPConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty NTP config")
		}
		for k, v := range cfg {
			t.Logf("NTP.%s = %s", k, v)
		}
	})

	t.Run("SetNTPConfig", func(t *testing.T) {
		original, err := c.Network.GetNTPConfig(ctx)
		if err != nil {
			t.Fatalf("GetNTPConfig (save): %v", err)
		}
		origEnable := original["Enable"]
		t.Logf("Original NTP.Enable: %s", origEnable)

		defer func() {
			_ = c.Network.SetNTPConfig(ctx, map[string]string{
				"NTP.Enable": origEnable,
			})
		}()

		newEnable := "true"
		if origEnable == "true" {
			newEnable = "false"
		}
		err = c.Network.SetNTPConfig(ctx, map[string]string{
			"NTP.Enable": newEnable,
		})
		skipOnSetError(t, err, "SetNTPConfig")

		updated, err := c.Network.GetNTPConfig(ctx)
		if err != nil {
			t.Fatalf("GetNTPConfig (verify): %v", err)
		}
		if updated["Enable"] != newEnable {
			t.Fatalf("expected NTP.Enable=%q, got %q", newEnable, updated["Enable"])
		}
		t.Logf("Verified NTP.Enable changed to %q", newEnable)
	})

	t.Run("GetRTSPConfig", func(t *testing.T) {
		cfg, err := c.Network.GetRTSPConfig(ctx)
		if err != nil {
			t.Fatalf("GetRTSPConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty RTSP config")
		}
		for k, v := range cfg {
			t.Logf("RTSP.%s = %s", k, v)
		}
	})

	t.Run("SetRTSPConfig", func(t *testing.T) {
		original, err := c.Network.GetRTSPConfig(ctx)
		if err != nil {
			t.Fatalf("GetRTSPConfig (save): %v", err)
		}
		origPort := original["Port"]
		t.Logf("Original RTSP.Port: %s", origPort)

		defer func() {
			_ = c.Network.SetRTSPConfig(ctx, map[string]string{
				"RTSP.Port": origPort,
			})
		}()

		// Change to a different valid port.
		newPort := "554"
		if origPort == "554" {
			newPort = "555"
		}
		err = c.Network.SetRTSPConfig(ctx, map[string]string{
			"RTSP.Port": newPort,
		})
		skipOnSetError(t, err, "SetRTSPConfig")

		updated, err := c.Network.GetRTSPConfig(ctx)
		if err != nil {
			t.Fatalf("GetRTSPConfig (verify): %v", err)
		}
		if updated["Port"] != newPort {
			t.Fatalf("expected RTSP.Port=%q, got %q", newPort, updated["Port"])
		}
		t.Logf("Verified RTSP.Port changed to %q", newPort)
	})

	t.Run("GetAlarmServerConfig", func(t *testing.T) {
		cfg, err := c.Network.GetAlarmServerConfig(ctx)
		if err != nil {
			t.Fatalf("GetAlarmServerConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty AlarmServer config")
		}
		for k, v := range cfg {
			t.Logf("AlarmServer.%s = %s", k, v)
		}
	})

	t.Run("SetAlarmServerConfig", func(t *testing.T) {
		original, err := c.Network.GetAlarmServerConfig(ctx)
		if err != nil {
			t.Fatalf("GetAlarmServerConfig (save): %v", err)
		}
		origEnable := original["Enable"]
		t.Logf("Original AlarmServer.Enable: %s", origEnable)

		defer func() {
			_ = c.Network.SetAlarmServerConfig(ctx, map[string]string{
				"AlarmServer.Enable": origEnable,
			})
		}()

		newEnable := "true"
		if origEnable == "true" {
			newEnable = "false"
		}
		err = c.Network.SetAlarmServerConfig(ctx, map[string]string{
			"AlarmServer.Enable": newEnable,
		})
		skipOnSetError(t, err, "SetAlarmServerConfig")

		updated, err := c.Network.GetAlarmServerConfig(ctx)
		if err != nil {
			t.Fatalf("GetAlarmServerConfig (verify): %v", err)
		}
		if updated["Enable"] != newEnable {
			t.Fatalf("expected AlarmServer.Enable=%q, got %q", newEnable, updated["Enable"])
		}
		t.Logf("Verified AlarmServer.Enable changed to %q", newEnable)
	})

	t.Run("GetSSHDConfig", func(t *testing.T) {
		if !hasSSHD {
			t.Skip("camera does not support SSHD config")
		}
		cfg, err := c.Network.GetSSHDConfig(ctx)
		if err != nil {
			t.Fatalf("GetSSHDConfig: %v", err)
		}
		if len(cfg) == 0 {
			t.Fatal("expected non-empty SSHD config")
		}
		for k, v := range cfg {
			t.Logf("SSHD.%s = %s", k, v)
		}
	})

	t.Run("SetSSHDConfig", func(t *testing.T) {
		if !hasSSHD {
			t.Skip("camera does not support SSHD config")
		}
		original, err := c.Network.GetSSHDConfig(ctx)
		if err != nil {
			t.Fatalf("GetSSHDConfig (save): %v", err)
		}
		origEnable := original["Enable"]
		t.Logf("Original SSHD.Enable: %s", origEnable)

		defer func() {
			_ = c.Network.SetSSHDConfig(ctx, map[string]string{
				"SSHD.Enable": origEnable,
			})
		}()

		newEnable := "true"
		if origEnable == "true" {
			newEnable = "false"
		}
		err = c.Network.SetSSHDConfig(ctx, map[string]string{
			"SSHD.Enable": newEnable,
		})
		skipOnSetError(t, err, "SetSSHDConfig")

		updated, err := c.Network.GetSSHDConfig(ctx)
		if err != nil {
			t.Fatalf("GetSSHDConfig (verify): %v", err)
		}
		if updated["Enable"] != newEnable {
			t.Fatalf("expected SSHD.Enable=%q, got %q", newEnable, updated["Enable"])
		}
		t.Logf("Verified SSHD.Enable changed to %q", newEnable)
	})

	t.Run("ScanWLanDevices", func(t *testing.T) {
		if !hasWLan {
			t.Skip("camera does not support WLan")
		}
		v, err := c.Network.ScanWLanDevices(ctx)
		if err != nil {
			t.Fatalf("ScanWLanDevices: %v", err)
		}
		t.Logf("WLan scan:\n%s", v)
	})
}
