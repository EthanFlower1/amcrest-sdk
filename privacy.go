package amcrest

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// PrivacyService handles privacy masking related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 94-107 (Section 4.5.18-28)
type PrivacyService struct {
	client *Client
}

// GetConfig returns the PrivacyMasking configuration as a raw key-value map.
// CGI: configManager.cgi?action=getConfig&name=PrivacyMasking
func (s *PrivacyService) GetConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "PrivacyMasking")
}

// GetMasking returns the raw privacy masking regions for the given channel.
// CGI: PrivacyMasking.cgi?action=getPrivacyMasking&channel=N&offset=O&limit=L
func (s *PrivacyService) GetMasking(ctx context.Context, channel, offset, limit int) (string, error) {
	return s.client.cgiGet(ctx, "PrivacyMasking.cgi", "getPrivacyMasking", url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
		"offset":  {fmt.Sprintf("%d", offset)},
		"limit":   {fmt.Sprintf("%d", limit)},
	})
}

// SetEnable enables or disables privacy masking for the given channel.
// CGI: PrivacyMasking.cgi?action=setPrivacyMaskingEnable&channel=N&Enable=true/false
func (s *PrivacyService) SetEnable(ctx context.Context, channel int, enable bool) error {
	return s.client.cgiAction(ctx, "PrivacyMasking.cgi", "setPrivacyMaskingEnable", url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
		"Enable":  {fmt.Sprintf("%t", enable)},
	})
}

// GetEnable returns whether privacy masking is enabled for the given channel.
// CGI: PrivacyMasking.cgi?action=getPrivacyMaskingEnable&channel=N
func (s *PrivacyService) GetEnable(ctx context.Context, channel int) (bool, error) {
	body, err := s.client.cgiGet(ctx, "PrivacyMasking.cgi", "getPrivacyMaskingEnable", url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	})
	if err != nil {
		return false, err
	}
	kv := parseKV(body)
	val, ok := kv["Enable"]
	if !ok {
		return false, fmt.Errorf("amcrest: Enable key not found in response")
	}
	return strings.EqualFold(strings.TrimSpace(val), "true"), nil
}

// ClearMasking removes all privacy masking regions for the given channel.
// CGI: PrivacyMasking.cgi?action=clearPrivacyMasking&channel=N
func (s *PrivacyService) ClearMasking(ctx context.Context, channel int) error {
	return s.client.cgiAction(ctx, "PrivacyMasking.cgi", "clearPrivacyMasking", url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	})
}

// SetConfig updates PrivacyMasking configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *PrivacyService) SetConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// SetMasking sets a privacy masking region for the given channel.
// CGI: PrivacyMasking.cgi?action=setPrivacyMasking&channel=N&PrivacyMasking.Index=X&...
func (s *PrivacyService) SetMasking(ctx context.Context, channel int, params map[string]string) error {
	vals := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	for k, v := range params {
		vals.Set(k, v)
	}
	return s.client.cgiAction(ctx, "PrivacyMasking.cgi", "setPrivacyMasking", vals)
}

// GotoMasking navigates to a specific privacy masking region.
// CGI: PrivacyMasking.cgi?action=gotoPrivacyMasking&channel=N&index=I
func (s *PrivacyService) GotoMasking(ctx context.Context, channel, index int) error {
	return s.client.cgiAction(ctx, "PrivacyMasking.cgi", "gotoPrivacyMasking", url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
		"index":   {fmt.Sprintf("%d", index)},
	})
}

// DeleteMasking deletes a specific privacy masking region.
// CGI: PrivacyMasking.cgi?action=deletePrivacyMasking&channel=N&index=I
func (s *PrivacyService) DeleteMasking(ctx context.Context, channel, index int) error {
	return s.client.cgiAction(ctx, "PrivacyMasking.cgi", "deletePrivacyMasking", url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
		"index":   {fmt.Sprintf("%d", index)},
	})
}

// GetRealRect returns the real rectangle coordinates for a privacy masking region.
// CGI: PrivacyMasking.cgi?action=getRealRect&channel=N&index=I
func (s *PrivacyService) GetRealRect(ctx context.Context, channel, index int) (map[string]string, error) {
	body, err := s.client.cgiGet(ctx, "PrivacyMasking.cgi", "getRealRect", url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
		"index":   {fmt.Sprintf("%d", index)},
	})
	if err != nil {
		return nil, fmt.Errorf("privacy GetRealRect: %w", err)
	}
	return parseKV(body), nil
}
