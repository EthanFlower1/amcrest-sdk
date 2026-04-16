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

// GetEncodeConfigCaps returns the encoding configuration capabilities for the given channel.
// CGI: encode.cgi?action=getConfigCaps&channel=N
func (s *VideoService) GetEncodeConfigCaps(ctx context.Context, channel int) (map[string]string, error) {
	body, err := s.client.cgiGet(ctx, "encode.cgi", "getConfigCaps", url.Values{
		"channel": {strconv.Itoa(channel)},
	})
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}

// SetEncodeConfig sets Encode configuration values. Keys should be prefixed
// with "Encode." (e.g., "Encode[0].MainFormat[0].Video.Compression").
func (s *VideoService) SetEncodeConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetVideoEncodeROI returns the raw VideoEncodeROI configuration table.
func (s *VideoService) GetVideoEncodeROI(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "VideoEncodeROI")
}

// SetVideoEncodeROI sets VideoEncodeROI configuration values.
func (s *VideoService) SetVideoEncodeROI(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetMaxRemoteInputChannels returns the maximum number of remote input channels.
// CGI: magicBox.cgi?action=getProductDefinition&name=MaxRemoteInputChannels
func (s *VideoService) GetMaxRemoteInputChannels(ctx context.Context) (int, error) {
	body, err := s.client.cgiGet(ctx, "magicBox.cgi", "getProductDefinition", url.Values{
		"name": {"MaxRemoteInputChannels"},
	})
	if err != nil {
		return 0, err
	}
	kv := parseKV(body)
	val, ok := kv["table.MaxRemoteInputChannels"]
	if !ok {
		return 0, fmt.Errorf("amcrest: MaxRemoteInputChannels not found in response")
	}
	n, err := strconv.Atoi(strings.TrimSpace(val))
	if err != nil {
		return 0, fmt.Errorf("amcrest: parsing MaxRemoteInputChannels: %w", err)
	}
	return n, nil
}

// SetVideoStandard sets the video standard (e.g., "PAL" or "NTSC").
// CGI: configManager.cgi?action=setConfig&VideoStandard=X
func (s *VideoService) SetVideoStandard(ctx context.Context, standard string) error {
	return s.client.setConfig(ctx, map[string]string{
		"VideoStandard": standard,
	})
}

// SetVideoWidget sets VideoWidget configuration values.
func (s *VideoService) SetVideoWidget(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetCurrentWindow returns the current video input window for the given channel.
// CGI: devVideoInput.cgi?action=getCurrentWindow&channel=N
func (s *VideoService) GetCurrentWindow(ctx context.Context, channel int) (map[string]string, error) {
	body, err := s.client.cgiGet(ctx, "devVideoInput.cgi", "getCurrentWindow", url.Values{
		"channel": {strconv.Itoa(channel)},
	})
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}

// SetCurrentWindow sets the current video input window for the given channel.
// The rect parameter specifies [left, top, right, bottom] coordinates.
// CGI: devVideoInput.cgi?action=setCurrentWindow&channel=N&region[0]=L&region[1]=T&region[2]=R&region[3]=B
func (s *VideoService) SetCurrentWindow(ctx context.Context, channel int, rect [4]int) error {
	params := url.Values{
		"channel":   {strconv.Itoa(channel)},
		"region[0]": {strconv.Itoa(rect[0])},
		"region[1]": {strconv.Itoa(rect[1])},
		"region[2]": {strconv.Itoa(rect[2])},
		"region[3]": {strconv.Itoa(rect[3])},
	}
	return s.client.cgiAction(ctx, "devVideoInput.cgi", "setCurrentWindow", params)
}

// GetVideoOutConfig returns the raw VideoOut configuration table.
func (s *VideoService) GetVideoOutConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "VideoOut")
}

// SetVideoOutConfig sets VideoOut configuration values.
func (s *VideoService) SetVideoOutConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// SetSmartEncodeConfig sets SmartEncode configuration values.
func (s *VideoService) SetSmartEncodeConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetDecoderCaps returns the video decoder capabilities.
// CGI: DevVideoDec.cgi?action=getCaps
func (s *VideoService) GetDecoderCaps(ctx context.Context) (map[string]string, error) {
	body, err := s.client.cgiGet(ctx, "DevVideoDec.cgi", "getCaps", nil)
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}

// GetLAEConfig returns the raw LAEConfig configuration table.
func (s *VideoService) GetLAEConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "LAEConfig")
}

// SetLAEConfig sets LAEConfig configuration values.
func (s *VideoService) SetLAEConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}
