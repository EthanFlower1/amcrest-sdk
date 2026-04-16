package amcrest

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Client is the root Amcrest API client. Domain-specific methods are accessed
// through the embedded service fields (e.g., client.System.Reboot(ctx)).
type Client struct {
	baseURL    string
	httpClient *http.Client

	// Domain services
	System        *SystemService
	User          *UserService
	Network       *NetworkService
	Video         *VideoService
	Snapshot      *SnapshotService
	Audio         *AudioService
	PTZ           *PTZService
	Event         *EventService
	Recording     *RecordingService
	Log           *LogService
	Storage       *StorageService
	Camera        *CameraService
	Display       *DisplayService
	Privacy       *PrivacyService
	Motion        *MotionService
	Upgrade       *UpgradeService
	Upload        *UploadService
	Analytics     *AnalyticsService
	Face          *FaceService
	People        *PeopleService
	WorkSuit      *WorkSuitService
	Traffic       *TrafficService
	Parking       *ParkingService
	Thermal       *ThermalService
	AccessControl *AccessControlService
	Building      *BuildingService
	DVR           *DVRService
	Peripheral    *PeripheralService
}

// Option configures the Client.
type Option func(*clientConfig)

type clientConfig struct {
	https      bool
	port       int
	httpClient *http.Client
	timeout    time.Duration
}

// WithHTTPS enables HTTPS for camera communication.
func WithHTTPS() Option {
	return func(c *clientConfig) { c.https = true }
}

// WithPort sets a custom port for camera communication.
func WithPort(port int) Option {
	return func(c *clientConfig) { c.port = port }
}

// WithHTTPClient sets a custom http.Client as the base transport.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *clientConfig) { c.httpClient = hc }
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *clientConfig) { c.timeout = d }
}

// NewClient creates a new Amcrest API client.
func NewClient(host, username, password string, opts ...Option) (*Client, error) {
	cfg := &clientConfig{
		timeout: 30 * time.Second,
	}
	for _, o := range opts {
		o(cfg)
	}

	scheme := "http"
	if cfg.https {
		scheme = "https"
	}

	port := cfg.port
	if port == 0 {
		if cfg.https {
			port = 443
		} else {
			port = 80
		}
	}

	baseURL := fmt.Sprintf("%s://%s:%d", scheme, host, port)

	var baseTransport http.RoundTripper
	if cfg.httpClient != nil {
		baseTransport = cfg.httpClient.Transport
	}
	if baseTransport == nil {
		baseTransport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	digestRT := newDigestTransport(username, password, baseTransport)

	httpClient := &http.Client{
		Transport: digestRT,
		Timeout:   cfg.timeout,
	}

	c := &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}

	// Initialize all services
	c.System = &SystemService{client: c}
	c.User = &UserService{client: c}
	c.Network = &NetworkService{client: c}
	c.Video = &VideoService{client: c}
	c.Snapshot = &SnapshotService{client: c}
	c.Audio = &AudioService{client: c}
	c.PTZ = &PTZService{client: c}
	c.Event = &EventService{client: c}
	c.Recording = &RecordingService{client: c}
	c.Log = &LogService{client: c}
	c.Storage = &StorageService{client: c}
	c.Camera = &CameraService{client: c}
	c.Display = &DisplayService{client: c}
	c.Privacy = &PrivacyService{client: c}
	c.Motion = &MotionService{client: c}
	c.Upgrade = &UpgradeService{client: c}
	c.Upload = &UploadService{client: c}
	c.Analytics = &AnalyticsService{client: c}
	c.Face = &FaceService{client: c}
	c.People = &PeopleService{client: c}
	c.WorkSuit = &WorkSuitService{client: c}
	c.Traffic = &TrafficService{client: c}
	c.Parking = &ParkingService{client: c}
	c.Thermal = &ThermalService{client: c}
	c.AccessControl = &AccessControlService{client: c}
	c.Building = &BuildingService{client: c}
	c.DVR = &DVRService{client: c}
	c.Peripheral = &PeripheralService{client: c}

	return c, nil
}

// get performs a GET request to the given path with optional query parameters.
func (c *Client) get(ctx context.Context, path string, params url.Values) (*http.Response, error) {
	u := c.baseURL + path
	if params != nil {
		u += "?" + params.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("amcrest: creating request: %w", err)
	}
	return c.httpClient.Do(req)
}

// postJSON performs a POST request with a JSON body and decodes the response into result.
func (c *Client) postJSON(ctx context.Context, path string, body interface{}, result interface{}) error {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("amcrest: marshaling JSON: %w", err)
	}

	u := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("amcrest: creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("amcrest: executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return &APIError{StatusCode: resp.StatusCode}
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("amcrest: decoding JSON response: %w", err)
		}
	}
	return nil
}

// postRaw performs a POST request with a JSON body and returns the raw response as a string.
func (c *Client) postRaw(ctx context.Context, path string, body interface{}) (string, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("amcrest: marshaling JSON: %w", err)
	}

	u := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("amcrest: creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("amcrest: executing request: %w", err)
	}
	return readBody(resp)
}

// cgiGet performs a GET to /cgi-bin/<cgi>?action=<action> with optional extra params,
// reads the response body, and returns it as a string.
func (c *Client) cgiGet(ctx context.Context, cgi, action string, params url.Values) (string, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("action", action)

	path := "/cgi-bin/" + cgi
	resp, err := c.get(ctx, path, params)
	if err != nil {
		return "", err
	}
	return readBody(resp)
}

// cgiAction performs a CGI GET and checks that the response is "OK".
func (c *Client) cgiAction(ctx context.Context, cgi, action string, params url.Values) error {
	if params == nil {
		params = url.Values{}
	}
	params.Set("action", action)

	path := "/cgi-bin/" + cgi
	resp, err := c.get(ctx, path, params)
	if err != nil {
		return err
	}
	return checkOK(resp)
}
