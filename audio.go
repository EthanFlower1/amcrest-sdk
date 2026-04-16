package amcrest

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// AudioService handles audio-related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 39-65 (Section 4.3)
type AudioService struct {
	client *Client
}

// GetInputChannels returns the number of audio input channels.
// CGI: devAudioInput.cgi?action=getCollect
func (s *AudioService) GetInputChannels(ctx context.Context) (int, error) {
	body, err := s.client.cgiGet(ctx, "devAudioInput.cgi", "getCollect", nil)
	if err != nil {
		return 0, err
	}
	result := parseKV(body)["result"]
	n, err := strconv.Atoi(strings.TrimSpace(result))
	if err != nil {
		return 0, fmt.Errorf("amcrest: parsing audio input channels %q: %w", result, err)
	}
	return n, nil
}

// GetOutputChannels returns the number of audio output channels.
// CGI: devAudioOutput.cgi?action=getCollect
func (s *AudioService) GetOutputChannels(ctx context.Context) (int, error) {
	body, err := s.client.cgiGet(ctx, "devAudioOutput.cgi", "getCollect", nil)
	if err != nil {
		return 0, err
	}
	result := parseKV(body)["result"]
	n, err := strconv.Atoi(strings.TrimSpace(result))
	if err != nil {
		return 0, fmt.Errorf("amcrest: parsing audio output channels %q: %w", result, err)
	}
	return n, nil
}

// GetVolume returns the AudioOutputVolume configuration as a raw key-value map.
// CGI: configManager.cgi?action=getConfig&name=AudioOutputVolume
func (s *AudioService) GetVolume(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "AudioOutputVolume")
}

// SetVolume sets the audio output volume for a specific channel.
// CGI: configManager.cgi?action=setConfig&AudioOutputVolume[channel]=volume
func (s *AudioService) SetVolume(ctx context.Context, channel, volume int) error {
	key := fmt.Sprintf("AudioOutputVolume[%d]", channel)
	return s.client.setConfig(ctx, map[string]string{
		key: strconv.Itoa(volume),
	})
}
