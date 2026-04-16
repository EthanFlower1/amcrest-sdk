package amcrest

import (
	"context"
	"testing"
)

func TestNetworkGetInterfaces(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	v, err := c.Network.GetInterfaces(ctx)
	if err != nil {
		t.Skipf("GetInterfaces not supported: %v", err)
	}
	if v == "" {
		t.Fatal("expected non-empty interfaces response")
	}
	t.Logf("Interfaces:\n%s", v)
}

func TestNetworkGetAccessFilter(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Network.GetAccessFilter(ctx)
	if err != nil {
		t.Skipf("GetAccessFilter not supported: %v", err)
	}
	if len(cfg) == 0 {
		t.Fatal("expected non-empty AccessFilter config")
	}
	for k, v := range cfg {
		t.Logf("AccessFilter.%s = %s", k, v)
	}
}

func TestNetworkGetNetworkConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Network.GetNetworkConfig(ctx)
	if err != nil {
		t.Skipf("GetNetworkConfig not supported: %v", err)
	}
	if len(cfg) == 0 {
		t.Fatal("expected non-empty Network config")
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}

func TestNetworkGetDDNSConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Network.GetDDNSConfig(ctx)
	if err != nil {
		t.Skipf("GetDDNSConfig not supported: %v", err)
	}
	if len(cfg) == 0 {
		t.Fatal("expected non-empty DDNS config")
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}

func TestNetworkGetEmailConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Network.GetEmailConfig(ctx)
	if err != nil {
		t.Skipf("GetEmailConfig not supported: %v", err)
	}
	if len(cfg) == 0 {
		t.Fatal("expected non-empty Email config")
	}
	for k, v := range cfg {
		t.Logf("Email.%s = %s", k, v)
	}
}

func TestNetworkGetWLanConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Network.GetWLanConfig(ctx)
	if err != nil {
		t.Skipf("GetWLanConfig not supported: %v", err)
	}
	if len(cfg) == 0 {
		t.Skip("WLan config is empty (WiFi may not be available)")
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}

func TestNetworkGetUPnPConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Network.GetUPnPConfig(ctx)
	if err != nil {
		t.Skipf("GetUPnPConfig not supported: %v", err)
	}
	if len(cfg) == 0 {
		t.Fatal("expected non-empty UPnP config")
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}

func TestNetworkGetUPnPStatus(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Network.GetUPnPStatus(ctx)
	if err != nil {
		t.Skipf("GetUPnPStatus not supported: %v", err)
	}
	if len(cfg) == 0 {
		t.Fatal("expected non-empty UPnP status")
	}
	for k, v := range cfg {
		t.Logf("%s = %s", k, v)
	}
}

func TestNetworkGetNTPConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Network.GetNTPConfig(ctx)
	if err != nil {
		t.Skipf("GetNTPConfig not supported: %v", err)
	}
	if len(cfg) == 0 {
		t.Fatal("expected non-empty NTP config")
	}
	for k, v := range cfg {
		t.Logf("NTP.%s = %s", k, v)
	}
}

func TestNetworkGetRTSPConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Network.GetRTSPConfig(ctx)
	if err != nil {
		t.Skipf("GetRTSPConfig not supported: %v", err)
	}
	if len(cfg) == 0 {
		t.Fatal("expected non-empty RTSP config")
	}
	for k, v := range cfg {
		t.Logf("RTSP.%s = %s", k, v)
	}
}

func TestNetworkGetAlarmServerConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Network.GetAlarmServerConfig(ctx)
	if err != nil {
		t.Skipf("GetAlarmServerConfig not supported: %v", err)
	}
	if len(cfg) == 0 {
		t.Fatal("expected non-empty AlarmServer config")
	}
	for k, v := range cfg {
		t.Logf("AlarmServer.%s = %s", k, v)
	}
}

func TestNetworkGetSSHDConfig(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	cfg, err := c.Network.GetSSHDConfig(ctx)
	if err != nil {
		t.Skipf("GetSSHDConfig not supported: %v", err)
	}
	if len(cfg) == 0 {
		t.Fatal("expected non-empty SSHD config")
	}
	for k, v := range cfg {
		t.Logf("SSHD.%s = %s", k, v)
	}
}

func TestNetworkScanWLanDevices(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()
	v, err := c.Network.ScanWLanDevices(ctx)
	if err != nil {
		t.Skipf("ScanWLanDevices not supported: %v", err)
	}
	t.Logf("WLan scan:\n%s", v)
}
