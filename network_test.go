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
