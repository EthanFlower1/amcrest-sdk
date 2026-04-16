package amcrest

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// VideoService handles video-related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 71-93 (Section 4.5)
type VideoService struct {
	client *Client
}

// GetMaxExtraStreams returns the maximum number of extra streams supported.
// CGI: magicBox.cgi?action=getProductDefinition&name=MaxExtraStream
func (s *VideoService) GetMaxExtraStreams(ctx context.Context) (int, error) {
	body, err := s.client.cgiGet(ctx, "magicBox.cgi", "getProductDefinition", url.Values{
		"name": {"MaxExtraStream"},
	})
	if err != nil {
		return 0, err
	}
	kv := parseKV(body)
	val, ok := kv["table.MaxExtraStream"]
	if !ok {
		return 0, fmt.Errorf("amcrest: MaxExtraStream not found in response")
	}
	n, err := strconv.Atoi(strings.TrimSpace(val))
	if err != nil {
		return 0, fmt.Errorf("amcrest: parsing MaxExtraStream: %w", err)
	}
	return n, nil
}

// GetEncodeCaps returns the encoding capabilities of the device.
// CGI: encode.cgi?action=getCaps
func (s *VideoService) GetEncodeCaps(ctx context.Context) (map[string]string, error) {
	body, err := s.client.cgiGet(ctx, "encode.cgi", "getCaps", nil)
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}

// GetEncodeConfig returns the raw Encode configuration table.
// Uses configManager getRawConfig to preserve the complex indexed structure.
func (s *VideoService) GetEncodeConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "Encode")
}

// GetVideoInputChannels returns the number of video input channels.
// CGI: devVideoInput.cgi?action=getCollect
func (s *VideoService) GetVideoInputChannels(ctx context.Context) (int, error) {
	body, err := s.client.cgiGet(ctx, "devVideoInput.cgi", "getCollect", nil)
	if err != nil {
		return 0, err
	}
	kv := parseKV(body)
	val, ok := kv["result"]
	if !ok {
		return 0, fmt.Errorf("amcrest: result not found in response")
	}
	n, err := strconv.Atoi(strings.TrimSpace(val))
	if err != nil {
		return 0, fmt.Errorf("amcrest: parsing video input channels: %w", err)
	}
	return n, nil
}

// GetVideoOutputChannels returns the number of video output channels.
// CGI: devVideoOutput.cgi?action=getCollect
func (s *VideoService) GetVideoOutputChannels(ctx context.Context) (int, error) {
	body, err := s.client.cgiGet(ctx, "devVideoOutput.cgi", "getCollect", nil)
	if err != nil {
		return 0, err
	}
	kv := parseKV(body)
	val, ok := kv["result"]
	if !ok {
		return 0, fmt.Errorf("amcrest: result not found in response")
	}
	n, err := strconv.Atoi(strings.TrimSpace(val))
	if err != nil {
		return 0, fmt.Errorf("amcrest: parsing video output channels: %w", err)
	}
	return n, nil
}

// GetVideoStandard returns the video standard (e.g., "PAL" or "NTSC").
// Uses configManager getConfig VideoStandard.
func (s *VideoService) GetVideoStandard(ctx context.Context) (string, error) {
	cfg, err := s.client.getConfig(ctx, "VideoStandard")
	if err != nil {
		return "", err
	}
	if v, ok := cfg["Type"]; ok {
		return v, nil
	}
	// Fallback: return the first value found.
	for _, v := range cfg {
		return v, nil
	}
	return "", fmt.Errorf("amcrest: VideoStandard not found in response")
}

// GetChannelTitle returns the raw ChannelTitle configuration table.
func (s *VideoService) GetChannelTitle(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "ChannelTitle")
}

// SetChannelTitle sets the channel title name for the given channel index.
// CGI: configManager.cgi?action=setConfig&ChannelTitle[channel].Name=name
func (s *VideoService) SetChannelTitle(ctx context.Context, channel int, name string) error {
	key := fmt.Sprintf("ChannelTitle[%d].Name", channel)
	return s.client.setConfig(ctx, map[string]string{key: name})
}

// GetVideoWidget returns the raw VideoWidget configuration table.
func (s *VideoService) GetVideoWidget(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "VideoWidget")
}

// GetSmartEncodeConfig returns the raw SmartEncode configuration table.
func (s *VideoService) GetSmartEncodeConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "SmartEncode")
}

// GetVideoInputCaps returns the video input capabilities for the given channel.
// CGI: devVideoInput.cgi?action=getCaps&channel=N
func (s *VideoService) GetVideoInputCaps(ctx context.Context, channel int) (map[string]string, error) {
	body, err := s.client.cgiGet(ctx, "devVideoInput.cgi", "getCaps", url.Values{
		"channel": {strconv.Itoa(channel)},
	})
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}
