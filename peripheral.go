package amcrest

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// PeripheralService handles peripheral device related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 305-325, 618-650 (Sections 6.2, 14.1)
type PeripheralService struct {
	client *Client
}

// WiperStart starts continuous wiper movement on the given channel.
// rainBrush.cgi?action=moveContinuously&channel=N&interval=I
func (s *PeripheralService) WiperStart(ctx context.Context, channel, interval int) error {
	params := url.Values{
		"channel":  {fmt.Sprintf("%d", channel)},
		"interval": {fmt.Sprintf("%d", interval)},
	}
	return s.client.cgiAction(ctx, "rainBrush.cgi", "moveContinuously", params)
}

// WiperStop stops wiper movement on the given channel.
// rainBrush.cgi?action=stopMove&channel=N
func (s *PeripheralService) WiperStop(ctx context.Context, channel int) error {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	return s.client.cgiAction(ctx, "rainBrush.cgi", "stopMove", params)
}

// WiperOnce performs a single wiper sweep on the given channel.
// rainBrush.cgi?action=moveOnce&channel=N
func (s *PeripheralService) WiperOnce(ctx context.Context, channel int) error {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	return s.client.cgiAction(ctx, "rainBrush.cgi", "moveOnce", params)
}

// ControlCoaxialIO sends a coaxial control IO command.
// coaxialControlIO.cgi?action=control&channel=N&info[0].Type=T&info[0].IO=I&info[0].TriggerMode=M
// Uses a raw query string to preserve bracket characters.
func (s *PeripheralService) ControlCoaxialIO(ctx context.Context, channel int, ioType, io, triggerMode int) error {
	rawQuery := fmt.Sprintf(
		"action=control&channel=%d&info[0].Type=%d&info[0].IO=%d&info[0].TriggerMode=%d",
		channel, ioType, io, triggerMode,
	)
	u := s.client.baseURL + "/cgi-bin/coaxialControlIO.cgi?" + rawQuery

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return fmt.Errorf("amcrest: creating request: %w", err)
	}
	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("amcrest: executing request: %w", err)
	}
	return checkOK(resp)
}

// GetCoaxialIOStatus retrieves the coaxial IO status for the given channel.
// coaxialControlIO.cgi?action=getstatus&channel=N
func (s *PeripheralService) GetCoaxialIOStatus(ctx context.Context, channel int) (map[string]string, error) {
	params := url.Values{
		"channel": {fmt.Sprintf("%d", channel)},
	}
	body, err := s.client.cgiGet(ctx, "coaxialControlIO.cgi", "getStatus", params)
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}

// GetFlashlightConfig retrieves the FlashLight configuration.
func (s *PeripheralService) GetFlashlightConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "FlashLight")
}

// GetGPSConfig retrieves the GPS configuration.
func (s *PeripheralService) GetGPSConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "GPS")
}

// GetFishEyeConfig retrieves the FishEye configuration.
func (s *PeripheralService) GetFishEyeConfig(ctx context.Context) (map[string]string, error) {
	return s.client.getRawConfig(ctx, "FishEye")
}
