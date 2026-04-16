package amcrest

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
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

// GetAudioStream opens a long-lived audio stream from the camera.
// The caller is responsible for closing the response body.
// httpType should be "singlepart" or "multipart".
// CGI: audio.cgi?action=getAudio&httptype=X&channel=N
func (s *AudioService) GetAudioStream(ctx context.Context, channel int, httpType string) (*http.Response, error) {
	params := url.Values{
		"action":   {"getAudio"},
		"httptype": {httpType},
		"channel":  {strconv.Itoa(channel)},
	}
	resp, err := s.client.get(ctx, "/cgi-bin/audio.cgi", params)
	if err != nil {
		return nil, fmt.Errorf("amcrest: getting audio stream: %w", err)
	}
	if resp.StatusCode >= 400 {
		resp.Body.Close()
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    "failed to get audio stream",
		}
	}
	return resp, nil
}

// GetPostAudioURL returns the full URL for posting audio data to the camera.
// The caller should POST audio data to this URL.
// httpType should be "singlepart" or "multipart".
// CGI: audio.cgi?action=postAudio&httptype=X&channel=N
func (s *AudioService) GetPostAudioURL(ctx context.Context, channel int, httpType string) string {
	params := url.Values{
		"action":   {"postAudio"},
		"httptype": {httpType},
		"channel":  {strconv.Itoa(channel)},
	}
	return s.client.baseURL + "/cgi-bin/audio.cgi?" + params.Encode()
}

// GetAudioAnalyseConfig returns the audio analysis configuration for the given channel.
// POST: /cgi-bin/api/AudioAnalyseManager/getConfig
func (s *AudioService) GetAudioAnalyseConfig(ctx context.Context, channel int) (string, error) {
	reqBody := struct {
		AudioChannel int `json:"AudioChannel"`
	}{AudioChannel: channel}
	return s.client.postRaw(ctx, "/cgi-bin/api/AudioAnalyseManager/getConfig", reqBody)
}

// SetAudioAnalyseConfig sets the audio analysis configuration.
// POST: /cgi-bin/api/AudioAnalyseManager/setConfig
func (s *AudioService) SetAudioAnalyseConfig(ctx context.Context, body interface{}) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/AudioAnalyseManager/setConfig", body, nil)
}

// GetAudioAnalyseClassConfig returns the audio analysis class configuration.
// POST: /cgi-bin/api/AudioAnalyseManager/getClassConfig
func (s *AudioService) GetAudioAnalyseClassConfig(ctx context.Context, className string, channel int) (string, error) {
	reqBody := struct {
		ClassName    string `json:"ClassName"`
		AudioChannel int    `json:"AudioChannel"`
	}{ClassName: className, AudioChannel: channel}
	return s.client.postRaw(ctx, "/cgi-bin/api/AudioAnalyseManager/getClassConfig", reqBody)
}

// GetInputCaps returns the capabilities of the given audio input channel.
// CGI: devAudioInput.cgi?action=getCaps&channel=N
func (s *AudioService) GetInputCaps(ctx context.Context, channel int) (map[string]string, error) {
	body, err := s.client.cgiGet(ctx, "devAudioInput.cgi", "getCaps", url.Values{
		"channel": {strconv.Itoa(channel)},
	})
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}

// GetOutputCaps returns the capabilities of the given audio output channel.
// CGI: devAudioOutput.cgi?action=getCaps&channel=N
func (s *AudioService) GetOutputCaps(ctx context.Context, channel int) (map[string]string, error) {
	body, err := s.client.cgiGet(ctx, "devAudioOutput.cgi", "getCaps", url.Values{
		"channel": {strconv.Itoa(channel)},
	})
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}

// GetInputConfig returns the AudioInput configuration table.
// CGI: configManager.cgi?action=getConfig&name=AudioInput
func (s *AudioService) GetInputConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "AudioInput")
}

// SetInputConfig sets AudioInput configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *AudioService) SetInputConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetOutputConfig returns the AudioOutput configuration table.
// CGI: configManager.cgi?action=getConfig&name=AudioOutput
func (s *AudioService) GetOutputConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "AudioOutput")
}

// SetOutputConfig sets AudioOutput configuration values.
// CGI: configManager.cgi?action=setConfig
func (s *AudioService) SetOutputConfig(ctx context.Context, params map[string]string) error {
	return s.client.setConfig(ctx, params)
}

// GetEncodeCapability returns the AudioEncodeCapability configuration table.
// CGI: configManager.cgi?action=getConfig&name=AudioEncodeCapability
func (s *AudioService) GetEncodeCapability(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "AudioEncodeCapability")
}

// SetAudioAnalyseClassConfig sets the audio analysis class configuration.
// POST: /cgi-bin/api/AudioAnalyseManager/setClassConfig
func (s *AudioService) SetAudioAnalyseClassConfig(ctx context.Context, body interface{}) error {
	return s.client.postJSON(ctx, "/cgi-bin/api/AudioAnalyseManager/setClassConfig", body, nil)
}

// GetAudioAnalyseCaps returns the audio analysis capabilities.
// POST: /cgi-bin/api/AudioAnalyseManager/getCaps
func (s *AudioService) GetAudioAnalyseCaps(ctx context.Context) (string, error) {
	return s.client.postRaw(ctx, "/cgi-bin/api/AudioAnalyseManager/getCaps", struct{}{})
}
