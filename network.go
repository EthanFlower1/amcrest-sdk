package amcrest

import (
	"context"
	"fmt"
	"net/url"
)

// NetworkService handles network-related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 138-156 (Section 4.8)
type NetworkService struct {
	client *Client
}

// GetInterfaces returns the raw body listing network interfaces.
// CGI: netApp.cgi?action=getInterfaces
func (n *NetworkService) GetInterfaces(ctx context.Context) (string, error) {
	return n.client.cgiGet(ctx, "netApp.cgi", "getInterfaces", nil)
}

// GetAccessFilter returns the AccessFilter configuration table with the
// "table.AccessFilter." prefix stripped from keys.
func (n *NetworkService) GetAccessFilter(ctx context.Context) (map[string]string, error) {
	return n.client.getConfig(ctx, "AccessFilter")
}

// GetNetworkConfig returns the Network configuration table without stripping
// prefixes, since keys contain interface names (e.g., "table.Network.eth0.IPAddress").
func (n *NetworkService) GetNetworkConfig(ctx context.Context) (map[string]string, error) {
	return n.client.getRawConfig(ctx, "Network")
}

// GetDDNSConfig returns the DDNS configuration table without stripping
// prefixes, since keys may contain provider-specific sub-tables.
func (n *NetworkService) GetDDNSConfig(ctx context.Context) (map[string]string, error) {
	return n.client.getRawConfig(ctx, "DDNS")
}

// GetEmailConfig returns the Email configuration table with the
// "table.Email." prefix stripped from keys.
func (n *NetworkService) GetEmailConfig(ctx context.Context) (map[string]string, error) {
	return n.client.getConfig(ctx, "Email")
}

// SetEmailConfig sets Email configuration values. Keys should be prefixed
// with "Email." (e.g., "Email.Enable", "Email.SendAddress").
func (n *NetworkService) SetEmailConfig(ctx context.Context, params map[string]string) error {
	return n.client.setConfig(ctx, params)
}

// GetWLanConfig returns the WLan configuration table without stripping
// prefixes, since keys may contain interface-specific sub-tables.
func (n *NetworkService) GetWLanConfig(ctx context.Context) (map[string]string, error) {
	return n.client.getRawConfig(ctx, "WLan")
}

// GetUPnPConfig returns the UPnP configuration table without stripping
// prefixes, since keys may contain interface-specific sub-tables.
func (n *NetworkService) GetUPnPConfig(ctx context.Context) (map[string]string, error) {
	return n.client.getRawConfig(ctx, "UPnP")
}

// GetUPnPStatus returns the current UPnP status as key-value pairs.
// CGI: netApp.cgi?action=getUPnPStatus
func (n *NetworkService) GetUPnPStatus(ctx context.Context) (map[string]string, error) {
	body, err := n.client.cgiGet(ctx, "netApp.cgi", "getUPnPStatus", nil)
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}

// GetNTPConfig returns the NTP configuration table with the
// "table.NTP." prefix stripped from keys.
func (n *NetworkService) GetNTPConfig(ctx context.Context) (map[string]string, error) {
	return n.client.getConfig(ctx, "NTP")
}

// SetNTPConfig sets NTP configuration values. Keys should be prefixed
// with "NTP." (e.g., "NTP.Enable", "NTP.Address").
func (n *NetworkService) SetNTPConfig(ctx context.Context, params map[string]string) error {
	return n.client.setConfig(ctx, params)
}

// GetRTSPConfig returns the RTSP configuration table with the
// "table.RTSP." prefix stripped from keys.
func (n *NetworkService) GetRTSPConfig(ctx context.Context) (map[string]string, error) {
	return n.client.getConfig(ctx, "RTSP")
}

// GetAlarmServerConfig returns the AlarmServer configuration table with the
// "table.AlarmServer." prefix stripped from keys.
func (n *NetworkService) GetAlarmServerConfig(ctx context.Context) (map[string]string, error) {
	return n.client.getConfig(ctx, "AlarmServer")
}

// GetSSHDConfig returns the SSHD configuration table with the
// "table.SSHD." prefix stripped from keys.
func (n *NetworkService) GetSSHDConfig(ctx context.Context) (map[string]string, error) {
	return n.client.getConfig(ctx, "SSHD")
}

// ScanWLanDevices triggers a WiFi scan and returns the raw body with
// discovered wireless networks.
// CGI: wlan.cgi?action=scanWlanDevices
func (n *NetworkService) ScanWLanDevices(ctx context.Context) (string, error) {
	return n.client.cgiGet(ctx, "wlan.cgi", "scanWlanDevices", nil)
}

// --- Convenience helpers (not part of the required API) ---

// GetNetworkConfigForInterface returns the Network configuration for a specific
// interface (e.g., "eth0") with the "table.Network.<iface>." prefix stripped.
func (n *NetworkService) GetNetworkConfigForInterface(ctx context.Context, iface string) (map[string]string, error) {
	body, err := n.client.cgiGet(ctx, "configManager.cgi", "getConfig", url.Values{
		"name": {"Network"},
	})
	if err != nil {
		return nil, err
	}
	prefix := fmt.Sprintf("table.Network.%s.", iface)
	return parseKVWithPrefix(body, prefix), nil
}
