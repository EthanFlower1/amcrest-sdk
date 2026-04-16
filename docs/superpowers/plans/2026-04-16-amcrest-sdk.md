# Amcrest Go SDK Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a complete Go SDK for the Amcrest HTTP API V3.26, covering all 15 chapters (~300+ endpoints) with integration tests against a real camera.

**Architecture:** Single `amcrest` package with a root `Client` holding embedded domain services (e.g., `client.System.Reboot(ctx)`). A shared HTTP layer handles digest authentication, key-value response parsing, JSON marshaling, and multipart stream reading. All services are flat in the root package.

**Tech Stack:** Go 1.21+, standard library only (net/http, crypto/md5, encoding/json, mime/multipart), no third-party dependencies. Integration tests use `testing` package with env-var-driven camera configuration.

**API Specification:** `docs/HTTP_API_V3.26.pdf` (650 pages). Each task references specific PDF page ranges.

**Test Camera:** `192.168.1.218`, username `admin`, password `Gsd4life.`

---

## Phase 1: Project Scaffolding & Core Infrastructure

These tasks MUST be completed sequentially before any service tasks. They establish the Go module, HTTP client, auth, parsing, and test helpers that every service depends on.

---

### Task 1: Go Module Init & Project Scaffolding

**Files:**
- Create: `go.mod`
- Create: `.env`
- Create: `.env.example`
- Create: `testhelper_test.go`

- [ ] **Step 1: Initialize Go module**

Run:
```bash
cd /Users/ethanflower/personal_projects/amcrest-sdk
go mod init github.com/EthanFlower/amcrest
```
Expected: `go.mod` created with `module github.com/EthanFlower/amcrest` and `go 1.21`

- [ ] **Step 2: Create .env file (gitignored)**

```
AMCREST_HOST=192.168.1.218
AMCREST_USERNAME=admin
AMCREST_PASSWORD=Gsd4life.
```

- [ ] **Step 3: Create .env.example file**

```
AMCREST_HOST=
AMCREST_USERNAME=
AMCREST_PASSWORD=
```

- [ ] **Step 4: Update .gitignore to include .env**

Verify `.env` is already in `.gitignore` (it is, per the existing file). No action needed.

- [ ] **Step 5: Create test helper file**

Create `testhelper_test.go`:

```go
package amcrest

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

func loadEnv() {
	f, err := os.Open(".env")
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			if os.Getenv(parts[0]) == "" {
				os.Setenv(parts[0], parts[1])
			}
		}
	}
}

func testClient(t *testing.T) *Client {
	t.Helper()
	loadEnv()
	host := os.Getenv("AMCREST_HOST")
	user := os.Getenv("AMCREST_USERNAME")
	pass := os.Getenv("AMCREST_PASSWORD")
	if host == "" || user == "" || pass == "" {
		t.Skip("AMCREST_HOST, AMCREST_USERNAME, AMCREST_PASSWORD must be set")
	}
	client, err := NewClient(host, user, pass)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return client
}
```

- [ ] **Step 6: Commit**

```bash
git add go.mod .env .env.example testhelper_test.go
git commit -m "feat: initialize Go module and test scaffolding"
```

---

### Task 2: Error Types

**Files:**
- Create: `errors.go`

**PDF Reference:** pp. 26-29 (Section 3.3 -- error response formats for both key-value and JSON APIs)

- [ ] **Step 1: Create errors.go**

```go
package amcrest

import "fmt"

// APIError represents an error response from the Amcrest API.
type APIError struct {
	StatusCode int
	Code       int
	Message    string
}

func (e *APIError) Error() string {
	if e.Code != 0 {
		return fmt.Sprintf("amcrest: HTTP %d, error code %d: %s", e.StatusCode, e.Code, e.Message)
	}
	if e.Message != "" {
		return fmt.Sprintf("amcrest: HTTP %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("amcrest: HTTP %d", e.StatusCode)
}
```

- [ ] **Step 2: Commit**

```bash
git add errors.go
git commit -m "feat: add API error type"
```

---

### Task 3: Digest Authentication Transport

**Files:**
- Create: `auth.go`

**PDF Reference:** pp. 30-31 (Section 3.4 Authentication -- digest auth per RFC 7616)

- [ ] **Step 1: Create auth.go**

```go
package amcrest

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

// digestTransport implements http.RoundTripper with HTTP Digest Authentication.
type digestTransport struct {
	username  string
	password  string
	transport http.RoundTripper

	mu    sync.Mutex
	realm string
	nonce string
	qop   string
	nc    int
}

func newDigestTransport(username, password string, transport http.RoundTripper) *digestTransport {
	if transport == nil {
		transport = http.DefaultTransport
	}
	return &digestTransport{
		username:  username,
		password:  password,
		transport: transport,
	}
}

func (t *digestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Try with cached credentials first
	t.mu.Lock()
	if t.nonce != "" {
		t.nc++
		authHeader := t.buildAuthHeader(req.Method, req.URL.RequestURI())
		req = req.Clone(req.Context())
		req.Header.Set("Authorization", authHeader)
		t.mu.Unlock()
		resp, err := t.transport.RoundTrip(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusUnauthorized {
			return resp, nil
		}
		// Nonce may have gone stale, fall through to re-auth
		resp.Body.Close()
	} else {
		t.mu.Unlock()
	}

	// Send initial request without auth to get WWW-Authenticate challenge
	initialReq := req.Clone(req.Context())
	initialReq.Header.Del("Authorization")

	// If the original request has a body, we need to handle it carefully
	resp, err := t.transport.RoundTrip(initialReq)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusUnauthorized {
		return resp, nil
	}

	challenge := resp.Header.Get("WWW-Authenticate")
	resp.Body.Close()

	if challenge == "" {
		return nil, fmt.Errorf("amcrest: 401 response without WWW-Authenticate header")
	}

	t.mu.Lock()
	t.parseChallenge(challenge)
	t.nc = 1
	authHeader := t.buildAuthHeader(req.Method, req.URL.RequestURI())
	t.mu.Unlock()

	authReq := req.Clone(req.Context())
	authReq.Header.Set("Authorization", authHeader)

	// Reset body if present
	if req.GetBody != nil {
		body, err := req.GetBody()
		if err != nil {
			return nil, fmt.Errorf("amcrest: failed to reset request body: %w", err)
		}
		authReq.Body = body
	}

	return t.transport.RoundTrip(authReq)
}

func (t *digestTransport) parseChallenge(challenge string) {
	challenge = strings.TrimPrefix(challenge, "Digest ")
	parts := splitChallenge(challenge)
	for k, v := range parts {
		switch strings.ToLower(k) {
		case "realm":
			t.realm = v
		case "nonce":
			t.nonce = v
		case "qop":
			t.qop = v
		}
	}
}

func splitChallenge(challenge string) map[string]string {
	result := make(map[string]string)
	parts := strings.Split(challenge, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		eqIdx := strings.Index(part, "=")
		if eqIdx < 0 {
			continue
		}
		key := strings.TrimSpace(part[:eqIdx])
		val := strings.TrimSpace(part[eqIdx+1:])
		val = strings.Trim(val, `"`)
		result[key] = val
	}
	return result
}

func (t *digestTransport) buildAuthHeader(method, uri string) string {
	ha1 := md5Hash(fmt.Sprintf("%s:%s:%s", t.username, t.realm, t.password))
	ha2 := md5Hash(fmt.Sprintf("%s:%s", method, uri))
	nc := fmt.Sprintf("%08x", t.nc)
	cnonce := md5Hash(fmt.Sprintf("%d", t.nc))

	var response string
	if t.qop == "auth" || t.qop == "auth-int" {
		response = md5Hash(fmt.Sprintf("%s:%s:%s:%s:%s:%s", ha1, t.nonce, nc, cnonce, t.qop, ha2))
	} else {
		response = md5Hash(fmt.Sprintf("%s:%s:%s", ha1, t.nonce, ha2))
	}

	parts := []string{
		fmt.Sprintf(`username="%s"`, t.username),
		fmt.Sprintf(`realm="%s"`, t.realm),
		fmt.Sprintf(`nonce="%s"`, t.nonce),
		fmt.Sprintf(`uri="%s"`, uri),
		fmt.Sprintf(`response="%s"`, response),
	}
	if t.qop != "" {
		parts = append(parts,
			fmt.Sprintf(`qop=%s`, t.qop),
			fmt.Sprintf(`nc=%s`, nc),
			fmt.Sprintf(`cnonce="%s"`, cnonce),
		)
	}
	return "Digest " + strings.Join(parts, ", ")
}

func md5Hash(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}
```

- [ ] **Step 2: Commit**

```bash
git add auth.go
git commit -m "feat: add HTTP digest authentication transport"
```

---

### Task 4: Key-Value Response Parser

**Files:**
- Create: `parse.go`
- Create: `parse_test.go`

**PDF Reference:** pp. 26-27 (Section 3.3.1 key=value format -- response line format: `table.Name[index].Param=value`)

- [ ] **Step 1: Write the failing test**

Create `parse_test.go`:

```go
package amcrest

import "testing"

func TestParseKV(t *testing.T) {
	input := "table.General.MachineName=TestCam\ntable.General.LocalNo=1\n"
	result := parseKV(input)
	if result["table.General.MachineName"] != "TestCam" {
		t.Errorf("expected TestCam, got %s", result["table.General.MachineName"])
	}
	if result["table.General.LocalNo"] != "1" {
		t.Errorf("expected 1, got %s", result["table.General.LocalNo"])
	}
}

func TestParseKVSingleValue(t *testing.T) {
	input := "result=1\n"
	result := parseKV(input)
	if result["result"] != "1" {
		t.Errorf("expected 1, got %s", result["result"])
	}
}

func TestParseKVWithSpaces(t *testing.T) {
	input := "result = 2011-7-3 21:02:32\n"
	result := parseKV(input)
	if result["result"] != "2011-7-3 21:02:32" {
		t.Errorf("expected time string, got %s", result["result"])
	}
}

func TestStripTablePrefix(t *testing.T) {
	input := "table.General.MachineName=TestCam\ntable.General.LocalNo=1\n"
	result := parseKVWithPrefix(input, "table.General.")
	if result["MachineName"] != "TestCam" {
		t.Errorf("expected TestCam, got %s", result["MachineName"])
	}
	if result["LocalNo"] != "1" {
		t.Errorf("expected 1, got %s", result["LocalNo"])
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/ethanflower/personal_projects/amcrest-sdk && go test -run TestParseKV -v`
Expected: FAIL (functions not defined)

- [ ] **Step 3: Write minimal implementation**

Create `parse.go`:

```go
package amcrest

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// parseKV parses key=value response lines into a map.
func parseKV(body string) map[string]string {
	result := make(map[string]string)
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		idx := strings.Index(line, "=")
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		result[key] = val
	}
	return result
}

// parseKVWithPrefix parses key=value lines and strips the given prefix from keys.
func parseKVWithPrefix(body, prefix string) map[string]string {
	raw := parseKV(body)
	result := make(map[string]string)
	for k, v := range raw {
		if strings.HasPrefix(k, prefix) {
			result[strings.TrimPrefix(k, prefix)] = v
		} else {
			result[k] = v
		}
	}
	return result
}

// readBody reads and returns the entire response body as a string.
// Returns an error if the response indicates a failure.
func readBody(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("amcrest: reading response body: %w", err)
	}
	body := string(data)

	if resp.StatusCode >= 400 {
		return "", &APIError{
			StatusCode: resp.StatusCode,
			Message:    strings.TrimSpace(body),
		}
	}

	// Check for "Error" in the body text (some 200 responses contain errors)
	trimmed := strings.TrimSpace(body)
	if strings.HasPrefix(trimmed, "Error") {
		return "", &APIError{
			StatusCode: resp.StatusCode,
			Message:    trimmed,
		}
	}

	return body, nil
}

// checkOK reads the body and verifies it contains "OK".
func checkOK(resp *http.Response) error {
	body, err := readBody(resp)
	if err != nil {
		return err
	}
	if strings.TrimSpace(body) != "OK" {
		return fmt.Errorf("amcrest: expected OK, got: %s", body)
	}
	return nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd /Users/ethanflower/personal_projects/amcrest-sdk && go test -run TestParseKV -v && go test -run TestStripTablePrefix -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add parse.go parse_test.go
git commit -m "feat: add key-value response parser"
```

---

### Task 5: Client & Config Helpers

**Files:**
- Create: `amcrest.go`
- Create: `config.go`

**PDF Reference:** pp. 35-38 (Section 4.2 Configure Manager -- getConfig/setConfig pattern)

- [ ] **Step 1: Create amcrest.go**

```go
package amcrest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is the Amcrest API client with embedded domain services.
type Client struct {
	host     string
	baseURL  string
	httpClient *http.Client

	System       *SystemService
	User         *UserService
	Network      *NetworkService
	Video        *VideoService
	Snapshot     *SnapshotService
	Audio        *AudioService
	PTZ          *PTZService
	Event        *EventService
	Recording    *RecordingService
	Log          *LogService
	Storage      *StorageService
	Camera       *CameraService
	Display      *DisplayService
	Privacy      *PrivacyService
	Motion       *MotionService
	Upgrade      *UpgradeService
	Upload       *UploadService
	Analytics    *AnalyticsService
	Face         *FaceService
	People       *PeopleService
	WorkSuit     *WorkSuitService
	Traffic      *TrafficService
	Parking      *ParkingService
	Thermal      *ThermalService
	AccessControl *AccessControlService
	Building     *BuildingService
	DVR          *DVRService
	Peripheral   *PeripheralService
}

// Option configures the Client.
type Option func(*clientConfig)

type clientConfig struct {
	scheme     string
	port       int
	httpClient *http.Client
	timeout    time.Duration
}

// WithHTTPS configures the client to use HTTPS.
func WithHTTPS() Option {
	return func(c *clientConfig) {
		c.scheme = "https"
		if c.port == 80 {
			c.port = 443
		}
	}
}

// WithPort sets a custom port.
func WithPort(port int) Option {
	return func(c *clientConfig) { c.port = port }
}

// WithHTTPClient sets a custom base HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(c *clientConfig) { c.httpClient = client }
}

// WithTimeout sets the default request timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *clientConfig) { c.timeout = d }
}

// NewClient creates a new Amcrest API client.
func NewClient(host, username, password string, opts ...Option) (*Client, error) {
	cfg := &clientConfig{
		scheme:  "http",
		port:    80,
		timeout: 30 * time.Second,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	baseClient := cfg.httpClient
	if baseClient == nil {
		baseClient = &http.Client{}
	}

	transport := baseClient.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	httpClient := &http.Client{
		Transport: newDigestTransport(username, password, transport),
		Timeout:   cfg.timeout,
	}

	baseURL := fmt.Sprintf("%s://%s:%d", cfg.scheme, host, cfg.port)

	c := &Client{
		host:       host,
		baseURL:    baseURL,
		httpClient: httpClient,
	}

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

// get performs a GET request to the given CGI path with query parameters.
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

// postJSON performs a POST request with a JSON body to the given API path.
func (c *Client) postJSON(ctx context.Context, path string, body any, result any) error {
	u := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("amcrest: marshaling request body: %w", err)
		}
		reqBody = strings.NewReader(string(data))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, reqBody)
	if err != nil {
		return fmt.Errorf("amcrest: creating request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("amcrest: executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    strings.TrimSpace(string(respBody)),
		}
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("amcrest: decoding response: %w", err)
		}
	}

	return nil
}

// cgiGet performs a GET to /cgi-bin/<cgi>?action=<action>&<extra params> and returns the body.
func (c *Client) cgiGet(ctx context.Context, cgi, action string, params url.Values) (string, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("action", action)
	resp, err := c.get(ctx, "/cgi-bin/"+cgi, params)
	if err != nil {
		return "", err
	}
	return readBody(resp)
}

// cgiAction performs a GET to a CGI endpoint and checks for "OK" response.
func (c *Client) cgiAction(ctx context.Context, cgi, action string, params url.Values) error {
	if params == nil {
		params = url.Values{}
	}
	params.Set("action", action)
	resp, err := c.get(ctx, "/cgi-bin/"+cgi, params)
	if err != nil {
		return err
	}
	return checkOK(resp)
}
```

- [ ] **Step 2: Create config.go**

```go
package amcrest

import (
	"context"
	"net/url"
)

// getConfig retrieves a configuration by name and returns the parsed key-value map
// with the "table.<name>." prefix stripped.
func (c *Client) getConfig(ctx context.Context, name string) (map[string]string, error) {
	params := url.Values{}
	params.Set("name", name)
	body, err := c.cgiGet(ctx, "configManager.cgi", "getConfig", params)
	if err != nil {
		return nil, err
	}
	return parseKVWithPrefix(body, "table."+name+"."), nil
}

// getConfigIndexed retrieves an indexed configuration (e.g., Encode[0]) and returns
// the parsed key-value map with the "table.<name>[<index>]." prefix stripped.
func (c *Client) getConfigIndexed(ctx context.Context, name string, index int) (map[string]string, error) {
	params := url.Values{}
	params.Set("name", name)
	body, err := c.cgiGet(ctx, "configManager.cgi", "getConfig", params)
	if err != nil {
		return nil, err
	}
	prefix := "table." + name + "["
	return parseKVWithPrefix(body, prefix), nil
}

// getRawConfig retrieves a configuration by name and returns the raw key-value map
// without stripping any prefix.
func (c *Client) getRawConfig(ctx context.Context, name string) (map[string]string, error) {
	params := url.Values{}
	params.Set("name", name)
	body, err := c.cgiGet(ctx, "configManager.cgi", "getConfig", params)
	if err != nil {
		return nil, err
	}
	return parseKV(body), nil
}

// setConfig sets configuration values. The params map keys should be
// in the form "ConfigName.Param=value" or "ConfigName[index].Param=value".
func (c *Client) setConfig(ctx context.Context, params map[string]string) error {
	qp := url.Values{}
	for k, v := range params {
		qp.Set(k, v)
	}
	return c.cgiAction(ctx, "configManager.cgi", "setConfig", qp)
}
```

- [ ] **Step 3: Commit**

```bash
git add amcrest.go config.go
git commit -m "feat: add Client, HTTP helpers, and config manager"
```

---

### Task 6: Multipart Event Stream Reader

**Files:**
- Create: `stream.go`

**PDF Reference:** pp. 171-173 (Section 4.9.17 Subscribe to Event Message -- multipart/x-mixed-replace response format with heartbeat)

- [ ] **Step 1: Create stream.go**

```go
package amcrest

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Event represents an event received from the event subscription stream.
type Event struct {
	Code   string
	Action string // "Start", "Stop", or "Pulse"
	Index  int
	Data   json.RawMessage
}

// EventStream reads multipart event data from a long-lived HTTP response.
type EventStream struct {
	resp   *http.Response
	reader *bufio.Reader
	cancel context.CancelFunc
}

// Subscribe opens a long-lived connection to the event stream and returns
// a channel of Events. Cancel the context to close the stream.
func (c *Client) subscribe(ctx context.Context, path string, params map[string]string) (<-chan Event, *EventStream, error) {
	ctx, cancel := context.WithCancel(ctx)

	qp := make(map[string]string)
	for k, v := range params {
		qp[k] = v
	}

	urlStr := c.baseURL + path + "?"
	pairs := []string{}
	for k, v := range qp {
		pairs = append(pairs, k+"="+v)
	}
	urlStr += strings.Join(pairs, "&")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		cancel()
		return nil, nil, fmt.Errorf("amcrest: creating stream request: %w", err)
	}

	// Use a client without timeout for long-lived connections
	streamClient := &http.Client{
		Transport: c.httpClient.Transport,
	}

	resp, err := streamClient.Do(req)
	if err != nil {
		cancel()
		return nil, nil, fmt.Errorf("amcrest: opening stream: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		cancel()
		return nil, nil, &APIError{StatusCode: resp.StatusCode}
	}

	es := &EventStream{
		resp:   resp,
		reader: bufio.NewReader(resp.Body),
		cancel: cancel,
	}

	ch := make(chan Event, 16)
	go es.readLoop(ch)
	return ch, es, nil
}

func (es *EventStream) readLoop(ch chan<- Event) {
	defer close(ch)
	defer es.resp.Body.Close()

	for {
		line, err := es.reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				// Context canceled or connection closed
			}
			return
		}
		line = strings.TrimSpace(line)

		if line == "" || line == "Heartbeat" || strings.HasPrefix(line, "--") {
			continue
		}

		if strings.HasPrefix(line, "Code=") {
			event := parseEventLine(line)
			// Try to read data lines
			for {
				next, err := es.reader.ReadString('\n')
				if err != nil {
					ch <- event
					return
				}
				next = strings.TrimSpace(next)
				if next == "" || strings.HasPrefix(next, "--") {
					break
				}
				if strings.HasPrefix(next, "data=") {
					event.Data = json.RawMessage(strings.TrimPrefix(next, "data="))
				}
			}
			ch <- event
		}
	}
}

func parseEventLine(line string) Event {
	var event Event
	parts := strings.Split(line, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		idx := strings.Index(part, "=")
		if idx < 0 {
			continue
		}
		key := part[:idx]
		val := part[idx+1:]
		switch key {
		case "Code":
			event.Code = val
		case "action":
			event.Action = val
		case "index":
			fmt.Sscanf(val, "%d", &event.Index)
		}
	}
	return event
}

// Close closes the event stream.
func (es *EventStream) Close() {
	es.cancel()
}
```

- [ ] **Step 2: Commit**

```bash
git add stream.go
git commit -m "feat: add multipart event stream reader"
```

---

### Task 7: Service Stubs (All 25 Services)

**Files:**
- Create: `system.go`, `users.go`, `network.go`, `video.go`, `snapshot.go`, `audio.go`, `ptz.go`, `events.go`, `recording.go`, `logs.go`, `storage.go`, `camera.go`, `display.go`, `privacy.go`, `motion.go`, `upgrade.go`, `upload.go`, `analytics.go`, `face.go`, `people.go`, `worksuit.go`, `traffic.go`, `parking.go`, `thermal.go`, `accesscontrol.go`, `building.go`, `dvr.go`, `peripheral.go`

- [ ] **Step 1: Create all service stub files**

Each file follows the same pattern. Example for `system.go`:

```go
package amcrest

// SystemService handles system-related API calls.
// PDF Reference: docs/HTTP_API_V3.26.pdf pp. 108-128 (Section 4.6)
type SystemService struct {
	client *Client
}
```

Create the same struct pattern for each service file:

| File | Struct | PDF Pages |
|------|--------|-----------|
| `system.go` | `SystemService` | pp. 108-128 |
| `users.go` | `UserService` | pp. 129-137 |
| `network.go` | `NetworkService` | pp. 138-156 |
| `video.go` | `VideoService` | pp. 71-93 |
| `snapshot.go` | `SnapshotService` | pp. 66-70 |
| `audio.go` | `AudioService` | pp. 39-65 |
| `ptz.go` | `PTZService` | pp. 284-305 |
| `events.go` | `EventService` | pp. 157-176 |
| `recording.go` | `RecordingService` | pp. 177-208 |
| `logs.go` | `LogService` | pp. 208-212 |
| `storage.go` | `StorageService` | pp. 259-278 |
| `camera.go` | `CameraService` | pp. 235-258 |
| `display.go` | `DisplayService` | pp. 279-283 |
| `privacy.go` | `PrivacyService` | pp. 94-107 |
| `motion.go` | `MotionService` | pp. 99-107, 416-419 |
| `upgrade.go` | `UpgradeService` | pp. 213-218 |
| `upload.go` | `UploadService` | pp. 219-234 |
| `analytics.go` | `AnalyticsService` | pp. 381-401 |
| `face.go` | `FaceService` | pp. 334-361 |
| `people.go` | `PeopleService` | pp. 362-381 |
| `worksuit.go` | `WorkSuitService` | pp. 402-416 |
| `traffic.go` | `TrafficService` | pp. 440-466 |
| `parking.go` | `ParkingService` | pp. 460-465 |
| `thermal.go` | `ThermalService` | pp. 481-504 |
| `accesscontrol.go` | `AccessControlService` | pp. 505-584 |
| `building.go` | `BuildingService` | pp. 585-608 |
| `dvr.go` | `DVRService` | pp. 609-617 |
| `peripheral.go` | `PeripheralService` | pp. 305-325, 618-650 |

- [ ] **Step 2: Verify it compiles**

Run: `cd /Users/ethanflower/personal_projects/amcrest-sdk && go build ./...`
Expected: Success, no errors

- [ ] **Step 3: Run integration test to verify auth works**

Create `amcrest_test.go`:

```go
package amcrest

import (
	"context"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := testClient(t)
	// Basic connectivity test -- hit magicBox to get device type
	ctx := context.Background()
	body, err := c.cgiGet(ctx, "magicBox.cgi", "getDeviceType", nil)
	if err != nil {
		t.Fatalf("getDeviceType failed: %v", err)
	}
	result := parseKV(body)
	if result["type"] == "" {
		t.Fatalf("expected device type, got empty. Full body: %s", body)
	}
	t.Logf("Device type: %s", result["type"])
}
```

Run: `cd /Users/ethanflower/personal_projects/amcrest-sdk && go test -run TestNewClient -v`
Expected: PASS with device type logged

- [ ] **Step 4: Commit**

```bash
git add system.go users.go network.go video.go snapshot.go audio.go ptz.go events.go recording.go logs.go storage.go camera.go display.go privacy.go motion.go upgrade.go upload.go analytics.go face.go people.go worksuit.go traffic.go parking.go thermal.go accesscontrol.go building.go dvr.go peripheral.go amcrest_test.go
git commit -m "feat: add all service stubs and connectivity test"
```

---

## Phase 2: Core Services

After Phase 1, these services can be implemented **in parallel** since they are independent. Each task is a complete service implementation with integration tests.

Each agent implementing a service task MUST:
1. Read the PDF pages listed for their service
2. Implement all endpoints documented in those pages
3. Write integration tests for read-only endpoints (at minimum)
4. Write integration tests for write endpoints where safe (with cleanup/restore)

---

### Task 8: SystemService

**Files:**
- Modify: `system.go`
- Create: `system_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 108-128 (Section 4.6)

Implement these endpoints by reading the PDF pages:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `GetDeviceType` | `magicBox.cgi?action=getDeviceType` | 4.6.8 (p.113) |
| `GetHardwareVersion` | `magicBox.cgi?action=getHardwareVersion` | 4.6.9 (p.113) |
| `GetSerialNumber` | `magicBox.cgi?action=getSerialNo` | 4.6.10 (p.114) |
| `GetMachineName` | `magicBox.cgi?action=getMachineName` | 4.6.11 (p.114) |
| `GetSystemInfo` | `magicBox.cgi?action=getSystemInfoNew` | 4.6.13 (p.115) |
| `GetVendor` | `magicBox.cgi?action=getVendor` | 4.6.14 (p.116) |
| `GetSoftwareVersion` | `magicBox.cgi?action=getSoftwareVersion` | 4.6.15 (p.116) |
| `GetDeviceClass` | `magicBox.cgi?action=getDeviceClass` | 4.6.18 (p.117) |
| `GetCurrentTime` | `global.cgi?action=getCurrentTime` | 4.6.2 (p.109) |
| `SetCurrentTime` | `global.cgi?action=setCurrentTime` | 4.6.3 (p.109) |
| `GetGeneralConfig` | `configManager.cgi?action=getConfig&name=General` | 4.6.1 (p.108) |
| `SetGeneralConfig` | `configManager.cgi?action=setConfig` | 4.6.1 (p.108) |
| `GetLocalesConfig` | `configManager.cgi?action=getConfig&name=Locales` | 4.6.4 (p.109) |
| `SetLocalesConfig` | `configManager.cgi?action=setConfig` | 4.6.4 (p.109) |
| `GetHolidayConfig` | `configManager.cgi?action=getConfig&name=Holiday` | 4.6.5 (p.111) |
| `GetLanguageCaps` | `magicBox.cgi?action=getLanguageCaps` | 4.6.6 (p.112) |
| `GetLanguage` | `configManager.cgi?action=getConfig&name=Language` | 4.6.7 (p.112) |
| `GetAutoMaintainConfig` | `configManager.cgi?action=getConfig&name=AutoMaintain` | 4.6.19 (p.117) |
| `SetAutoMaintainConfig` | `configManager.cgi?action=setConfig` | 4.6.19 (p.117) |
| `Reboot` | `magicBox.cgi?action=reboot` | 4.6.20 (p.119) |
| `Shutdown` | `magicBox.cgi?action=shutdown` | 4.6.21 (p.119) |
| `FactoryReset` | `magicBox.cgi?action=resetSystemEx` | 4.6.22 (p.119) |
| `GetOnvifVersion` | `IntervideoManager.cgi?action=getVersion&Name=Onvif` | 4.6.16 (p.117) |
| `GetHTTPAPIVersion` | `IntervideoManager.cgi?action=getVersion&Name=CGI` | 4.6.17 (p.117) |
| `GetCompleteMachineVersion` | `api/MagicBox/getCompleteMachineVersion` | 4.6.27 (p.128) |

- [ ] **Step 1: Read PDF pp. 108-128 and implement all endpoints above as methods on SystemService**

Each method should use the `client.cgiGet()`, `client.cgiAction()`, `client.getConfig()`, `client.setConfig()`, or `client.postJSON()` helpers as appropriate. Define typed response structs for complex responses (e.g., `SystemInfo`, `GeneralConfig`, `LocalesConfig`, `AutoMaintainConfig`).

- [ ] **Step 2: Write integration tests in system_test.go**

Test all read-only methods. For `SetCurrentTime`, save current time, set a new time, verify it changed, then restore original. Mark `Reboot`, `Shutdown`, and `FactoryReset` with `//go:build dangerous`.

- [ ] **Step 3: Run tests**

Run: `cd /Users/ethanflower/personal_projects/amcrest-sdk && go test -run TestSystem -v`
Expected: All tests PASS

- [ ] **Step 4: Commit**

```bash
git add system.go system_test.go
git commit -m "feat: implement SystemService with integration tests"
```

---

### Task 9: SnapshotService

**Files:**
- Modify: `snapshot.go`
- Create: `snapshot_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 66-70 (Section 4.4)

Implement these endpoints:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `GetSnapConfig` | `configManager.cgi?action=getConfig&name=Snap` | 4.4.1 (p.66) |
| `SetSnapConfig` | `configManager.cgi?action=setConfig` | 4.4.1 (p.66) |
| `Get` | `snapshot.cgi?channel=N&type=T` | 4.4.2 (p.68) |
| `Subscribe` | `snapManager.cgi?action=attachFileProc` | 4.4.3 (p.68) |

- [ ] **Step 1: Read PDF pp. 66-70 and implement all endpoints**

`Get` should return `[]byte` (JPEG data). `Subscribe` should return a channel for streaming snapshot events.

- [ ] **Step 2: Write integration tests**

Test `Get` by capturing a snapshot and verifying it starts with JPEG magic bytes (`0xFF 0xD8`). Test `Subscribe` with a short timeout context.

- [ ] **Step 3: Run tests and commit**

---

### Task 10: UserService

**Files:**
- Modify: `users.go`
- Create: `users_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 129-137 (Section 4.7)

Implement these endpoints:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `GetUserInfo` | `userManager.cgi?action=getUserInfo` | 4.7.1 (p.133) |
| `GetAllUsers` | `userManager.cgi?action=getUserInfoAll` | 4.7.2 (p.133) |
| `GetActiveUsers` | `userManager.cgi?action=getActiveUserInfoAll` | 4.7.3 (p.134) |
| `GetGroupInfo` | `userManager.cgi?action=getGroupInfo` | 4.7.4 (p.134) |
| `GetAllGroups` | `userManager.cgi?action=getGroupInfoAll` | 4.7.5 (p.134) |
| `AddUser` | `userManager.cgi?action=addUser` | 4.7.6 (p.135) |
| `DeleteUser` | `userManager.cgi?action=deleteUser` | 4.7.7 (p.135) |
| `ModifyUser` | `userManager.cgi?action=modifyUser` | 4.7.8 (p.136) |
| `ModifyPassword` | `userManager.cgi?action=modifyPassword` | 4.7.9 (p.136) |
| `ModifyPasswordByManager` | `userManager.cgi?action=modifyPasswordByManager` | 4.7.10 (p.136) |

- [ ] **Step 1: Read PDF pp. 129-137 and implement all endpoints**
- [ ] **Step 2: Write integration tests** (read-only tests for Get methods; AddUser/DeleteUser as a pair with cleanup)
- [ ] **Step 3: Run tests and commit**

---

### Task 11: EventService

**Files:**
- Modify: `events.go`
- Create: `events_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 157-176 (Section 4.9)

Implement these endpoints:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `Subscribe` | `eventManager.cgi?action=attach` | 4.9.17 (p.171) |
| `GetEventIndexes` | `eventManager.cgi?action=getEventIndexes` | 4.9.16 (p.170) |
| `GetCaps` | `eventManager.cgi?action=getCaps` | 4.9.18 (p.173) |
| `GetSupportedEvents` | `eventManager.cgi?action=getExposureEvents` | 4.9.21 (p.176) |
| `GetAlarmConfig` | `configManager.cgi?action=getConfig&name=Alarm` | 4.9.2 (p.161) |
| `SetAlarmConfig` | `configManager.cgi?action=setConfig` | 4.9.2 (p.161) |
| `GetAlarmOutConfig` | `configManager.cgi?action=getConfig&name=AlarmOut` | 4.9.3 (p.162) |
| `SetAlarmOutConfig` | `configManager.cgi?action=setConfig` | 4.9.3 (p.162) |
| `GetAlarmInputChannels` | `alarm.cgi?action=getInSlots` | 4.9.4 (p.163) |
| `GetAlarmOutputChannels` | `alarm.cgi?action=getOutSlots` | 4.9.5 (p.163) |
| `GetAlarmInputStates` | `alarm.cgi?action=getInState` | 4.9.6 (p.163) |
| `GetAlarmOutputStates` | `alarm.cgi?action=getOutState` | 4.9.7 (p.163) |
| `GetBlindDetectConfig` | `configManager.cgi?action=getConfig&name=BlindDetect` | 4.9.8 (p.164) |
| `GetLossDetectConfig` | `configManager.cgi?action=getConfig&name=LossDetect` | 4.9.9 (p.165) |
| `GetLoginFailureAlarmConfig` | `configManager.cgi?action=getConfig&name=LoginFailureAlarm` | 4.9.10 (p.165) |
| `GetStorageNotExistConfig` | `configManager.cgi?action=getConfig&name=StorageNotExist` | 4.9.11 (p.166) |
| `GetStorageFailureConfig` | `configManager.cgi?action=getConfig&name=StorageFailure` | 4.9.12 (p.167) |
| `GetStorageLowSpaceConfig` | `configManager.cgi?action=getConfig&name=StorageLowSpace` | 4.9.13 (p.168) |
| `GetNetAbortConfig` | `configManager.cgi?action=getConfig&name=NetAbort` | 4.9.14 (p.169) |
| `GetIPConflictConfig` | `configManager.cgi?action=getConfig&name=IPConflict` | 4.9.15 (p.169) |
| `SetNetAlarmState` | `netAlarm.cgi?action=setState` | 4.9.20 (p.176) |

- [ ] **Step 1: Read PDF pp. 157-176 and implement all endpoints**
- [ ] **Step 2: Write integration tests** (Subscribe with short timeout to get heartbeat; read-only for all Get methods)
- [ ] **Step 3: Run tests and commit**

---

### Task 12: PTZService

**Files:**
- Modify: `ptz.go`
- Create: `ptz_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 284-305 (Section 8.1)

Implement these endpoints:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `Control` | `ptz.cgi?action=start&code=<action>` | 8.1.5 (p.287) |
| `Stop` | `ptz.cgi?action=stop&code=<action>` | 8.1.5 (p.287) |
| `GetPresets` | `ptz.cgi?action=getPresets` | 8.1.6 (p.292) |
| `SetPreset` | `ptz.cgi?action=start&code=SetPreset` | 8.1.6 (p.292) |
| `GotoPreset` | `ptz.cgi?action=start&code=GotoPreset` | 8.1.6 (p.292) |
| `ClearPreset` | `ptz.cgi?action=start&code=ClearPreset` | 8.1.6 (p.292) |
| `GetTours` | `ptz.cgi?action=start&code=QueryTour` | 8.1.7 (p.294) |
| `StartTour` | `ptz.cgi?action=start&code=StartTour` | 8.1.7 (p.294) |
| `StopTour` | `ptz.cgi?action=start&code=StopTour` | 8.1.7 (p.294) |
| `GetStatus` | `ptz.cgi?action=getStatus` | 8.1.4 (p.287) |
| `GetConfig` | `configManager.cgi?action=getConfig&name=Ptz` | 8.1.1 (p.284) |
| `SetConfig` | `configManager.cgi?action=setConfig` | 8.1.1 (p.284) |
| `GetProtocolList` | `ptz.cgi?action=getProtocolList` | 8.1.2 (p.285) |
| `GetCaps` | `ptz.cgi?action=getCurrentProtocolCaps` | 8.1.3 (p.285) |
| `Restart` | `ptz.cgi?action=start&code=Restart` | 8.1.12 (p.301) |
| `Reset` | `ptz.cgi?action=start&code=Reset` | 8.1.13 (p.301) |

- [ ] **Step 1: Read PDF pp. 284-305 and implement all endpoints**

Define PTZ action constants: `Up`, `Down`, `Left`, `Right`, `ZoomTele`, `ZoomWide`, `FocusNear`, `FocusFar`, `LeftUp`, `RightUp`, `LeftDown`, `RightDown`.

- [ ] **Step 2: Write integration tests** (GetPresets, GetStatus, GetConfig are safe read-only tests)
- [ ] **Step 3: Run tests and commit**

---

### Task 13: NetworkService

**Files:**
- Modify: `network.go`
- Create: `network_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 138-156 (Section 4.8)

Implement these endpoints:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `GetInterfaces` | `netApp.cgi?action=getInterfaces` | 4.8.1 (p.138) |
| `GetAccessFilter` | `configManager.cgi?action=getConfig&name=AccessFilter` | 4.8.2 (p.138) |
| `GetNetworkConfig` | `configManager.cgi?action=getConfig&name=Network` | 4.8.3 (p.139) |
| `SetNetworkConfig` | `configManager.cgi?action=setConfig` | 4.8.3 (p.139) |
| `GetDDNSConfig` | `configManager.cgi?action=getConfig&name=DDNS` | 4.8.5 (p.141) |
| `SetDDNSConfig` | `configManager.cgi?action=setConfig` | 4.8.5 (p.141) |
| `GetEmailConfig` | `configManager.cgi?action=getConfig&name=Email` | 4.8.6 (p.144) |
| `SetEmailConfig` | `configManager.cgi?action=setConfig` | 4.8.6 (p.144) |
| `GetWLanConfig` | `configManager.cgi?action=getConfig&name=WLan` | 4.8.7 (p.145) |
| `ScanWLanDevices` | `wlan.cgi?action=scanWlanDevices` | 4.8.8 (p.146) |
| `GetUPnPConfig` | `configManager.cgi?action=getConfig&name=UPnP` | 4.8.9 (p.147) |
| `GetUPnPStatus` | `netApp.cgi?action=getUPnPStatus` | 4.8.10 (p.148) |
| `GetNTPConfig` | `configManager.cgi?action=getConfig&name=NTP` | 4.8.11 (p.148) |
| `SetNTPConfig` | `configManager.cgi?action=setConfig` | 4.8.11 (p.148) |
| `GetRTSPConfig` | `configManager.cgi?action=getConfig&name=RTSP` | 4.8.12 (p.149) |
| `SetRTSPConfig` | `configManager.cgi?action=setConfig` | 4.8.12 (p.149) |
| `GetAlarmServerConfig` | `configManager.cgi?action=getConfig&name=AlarmServer` | 4.8.13 (p.150) |
| `GetSSHDConfig` | `configManager.cgi?action=getConfig&name=SSHD` | 4.8.15 (p.151) |

- [ ] **Step 1: Read PDF pp. 138-156 and implement all endpoints**
- [ ] **Step 2: Write integration tests** (read-only tests for all Get methods)
- [ ] **Step 3: Run tests and commit**

---

### Task 14: VideoService

**Files:**
- Modify: `video.go`
- Create: `video_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 71-93 (Section 4.5)

Implement these endpoints:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `GetMaxExtraStreams` | `magicBox.cgi?action=getProductDefinition&name=MaxExtraStream` | 4.5.1 (p.71) |
| `GetEncodeCaps` | `encode.cgi?action=getCaps` | 4.5.2 (p.71) |
| `GetEncodeConfigCaps` | `encode.cgi?action=getConfigCaps` | 4.5.3 (p.72) |
| `GetEncodeConfig` | `configManager.cgi?action=getConfig&name=Encode` | 4.5.4 (p.76) |
| `SetEncodeConfig` | `configManager.cgi?action=setConfig` | 4.5.4 (p.76) |
| `GetVideoInputChannels` | `devVideoInput.cgi?action=getCollect` | 4.5.7 (p.83) |
| `GetVideoOutputChannels` | `devVideoOutput.cgi?action=getCollect` | 4.5.8 (p.83) |
| `GetVideoStandard` | `configManager.cgi?action=getConfig&name=VideoStandard` | 4.5.10 (p.83) |
| `GetVideoInputCaps` | `devVideoInput.cgi?action=getCaps` | 4.5.12 (p.87) |
| `GetChannelTitle` | `configManager.cgi?action=getConfig&name=ChannelTitle` | 4.5.6 (p.82) |
| `SetChannelTitle` | `configManager.cgi?action=setConfig` | 4.5.6 (p.82) |
| `GetVideoWidget` | `configManager.cgi?action=getConfig&name=VideoWidget` | 4.5.11 (p.84) |
| `GetSmartEncode` | `configManager.cgi?action=getConfig&name=SmartEncode` | 4.5.16 (p.92) |

- [ ] **Step 1: Read PDF pp. 71-93 and implement all endpoints**
- [ ] **Step 2: Write integration tests**
- [ ] **Step 3: Run tests and commit**

---

### Task 15: AudioService

**Files:**
- Modify: `audio.go`
- Create: `audio_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 39-65 (Section 4.3)

Implement these endpoints:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `GetInputChannels` | `devAudioInput.cgi?action=getCollect` | 4.3.1 (p.39) |
| `GetOutputChannels` | `devAudioOutput.cgi?action=getCollect` | 4.3.2 (p.39) |
| `GetVolume` | `configManager.cgi?action=getConfig&name=AudioOutputVolume` | 4.3.5 (p.42) |
| `SetVolume` | `configManager.cgi?action=setConfig` | 4.3.5 (p.42) |

- [ ] **Step 1: Read PDF pp. 39-65 and implement endpoints**
- [ ] **Step 2: Write integration tests**
- [ ] **Step 3: Run tests and commit**

---

### Task 16: RecordingService

**Files:**
- Modify: `recording.go`
- Create: `recording_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 177-208 (Section 4.10)

Implement these endpoints:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `GetCaps` | `recordManager.cgi?action=getCaps` | 4.10.1 (p.177) |
| `GetRecordConfig` | `configManager.cgi?action=getConfig&name=Record` | 4.10.2 (p.177) |
| `SetRecordConfig` | `configManager.cgi?action=setConfig` | 4.10.2 (p.177) |
| `GetRecordMode` | `configManager.cgi?action=getConfig&name=RecordMode` | 4.10.3 (p.178) |
| `SetRecordMode` | `configManager.cgi?action=setConfig` | 4.10.3 (p.178) |
| `GetMediaGlobal` | `configManager.cgi?action=getConfig&name=MediaGlobal` | 4.10.4 (p.179) |
| `FindFiles` | `mediaFileFind.cgi` (full lifecycle) | 4.10.5 (p.180) |
| `DownloadFile` | `RPC_Loadfile/<filename>` | 4.10.12 (p.201) |
| `DownloadByTime` | `loadfile.cgi?action=startLoad` | 4.10.13 (p.207) |

- [ ] **Step 1: Read PDF pp. 177-208 and implement all endpoints**

`FindFiles` should abstract the stateful finder pattern (factory.create → findFile → findNextFile → close → destroy) into a single method that returns `[]MediaFile`. Define `MediaFile` struct and `FindFilesOpts` for search parameters.

- [ ] **Step 2: Write integration tests** (FindFiles with a recent time range, verify results; GetRecordConfig etc.)
- [ ] **Step 3: Run tests and commit**

---

### Task 17: LogService

**Files:**
- Modify: `logs.go`
- Create: `logs_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 208-212 (Section 4.11)

Implement these endpoints:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `Find` | `log.cgi` (startFind → doFind → stopFind) | 4.11.1 (p.209) |
| `Clear` | `log.cgi?action=clear` | 4.11.2 (p.210) |
| `Backup` | `Log.backup?action=All` | 4.11.3 (p.211) |

- [ ] **Step 1: Read PDF pp. 208-212 and implement all endpoints**

`Find` should abstract the token-based search pattern into a single method returning `[]LogEntry`.

- [ ] **Step 2: Write integration tests** (Find with recent time range; Clear with `dangerous` tag)
- [ ] **Step 3: Run tests and commit**

---

### Task 18: StorageService

**Files:**
- Modify: `storage.go`
- Create: `storage_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 259-278 (Section 6)

Implement these endpoints:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `GetDiskInfo` | `storageDevice.cgi?action=factory.getPortInfo` | 6.1.1 (p.259) |
| `GetDeviceNames` | `storageDevice.cgi?action=factory.getCollect` | 6.1.2 (p.260) |
| `GetAllDeviceInfo` | `storageDevice.cgi?action=getDeviceAllInfo` | 6.1.3 (p.260) |
| `GetCaps` | `storage.cgi?action=getCaps` | 6.1.4 (p.261) |
| `FormatSDCard` | `storageDevice.cgi?action=setDevice&type=FormatPatition` | 6.1.5 (p.261) |
| `GetNASConfig` | `configManager.cgi?action=getConfig&name=NAS` | 6.2.1 (p.271) |
| `SetNASConfig` | `configManager.cgi?action=setConfig` | 6.2.1 (p.271) |
| `GetStorageGroupConfig` | `configManager.cgi?action=getConfig&name=StorageGroup` | 6.3.2 (p.274) |
| `GetStorageHealthAlarm` | `configManager.cgi?action=getConfig&name=StorageHealthAlarm` | 6.4.6 (p.277) |

- [ ] **Step 1: Read PDF pp. 259-278 and implement all endpoints** (FormatSDCard behind `dangerous` tag)
- [ ] **Step 2: Write integration tests**
- [ ] **Step 3: Run tests and commit**

---

### Task 19: CameraService

**Files:**
- Modify: `camera.go`
- Create: `camera_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 235-258 (Section 5)

Implement these endpoints:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `GetImageConfig` | `configManager.cgi?action=getConfig&name=VideoColor` | 5.1.1 (p.235) |
| `SetImageConfig` | `configManager.cgi?action=setConfig` | 5.1.1 (p.235) |
| `GetExposureConfig` | `configManager.cgi?action=getConfig&name=VideoInExposure` | 5.2.1 (p.238) |
| `SetExposureConfig` | `configManager.cgi?action=setConfig` | 5.2.1 (p.238) |
| `GetBacklightConfig` | `configManager.cgi?action=getConfig&name=VideoInBacklight` | 5.3.1 (p.240) |
| `GetWhiteBalanceConfig` | `configManager.cgi?action=getConfig&name=VideoInWhiteBalance` | 5.4.1 (p.241) |
| `GetDayNightConfig` | `configManager.cgi?action=getConfig&name=VideoInDayNight` | 5.5.1 (p.242) |
| `AutoFocus` | `devVideoInput.cgi?action=autoFocus` | 5.6.3 (p.244) |
| `GetFocusStatus` | `devVideoInput.cgi?action=getFocusStatus` | 5.6.4 (p.245) |
| `GetLightingConfig` | `configManager.cgi?action=getConfig&name=Lighting` | 5.7.1 (p.247) |
| `GetVideoInOptions` | `configManager.cgi?action=getConfig&name=VideoInOptions` | 5.8.2 (p.250) |

- [ ] **Step 1: Read PDF pp. 235-258 and implement all endpoints**
- [ ] **Step 2: Write integration tests**
- [ ] **Step 3: Run tests and commit**

---

### Task 20: PrivacyService

**Files:**
- Modify: `privacy.go`
- Create: `privacy_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 94-107 (Sections 4.5.18-4.5.28)

Implement these endpoints:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `GetConfig` | `configManager.cgi?action=getConfig&name=PrivacyMasking` | 4.5.18 (p.94) |
| `SetConfig` | `configManager.cgi?action=setConfig` | 4.5.18 (p.94) |
| `GetMasking` | `PrivacyMasking.cgi?action=getPrivacyMasking` | 4.5.19 (p.96) |
| `SetMasking` | `PrivacyMasking.cgi?action=setPrivacyMasking` | 4.5.20 (p.97) |
| `GotoMasking` | `PrivacyMasking.cgi?action=gotoPrivacyMasking` | 4.5.21 (p.98) |
| `DeleteMasking` | `PrivacyMasking.cgi?action=deletePrivacyMasking` | 4.5.22 (p.98) |
| `ClearMasking` | `PrivacyMasking.cgi?action=clearPrivacyMasking` | 4.5.23 (p.99) |
| `GetEnable` | `PrivacyMasking.cgi?action=getPrivacyMaskingEnable` | 4.5.28 (p.107) |
| `SetEnable` | `PrivacyMasking.cgi?action=setPrivacyMaskingEnable` | 4.5.27 (p.107) |

- [ ] **Step 1: Read PDF pp. 94-107 and implement all endpoints**
- [ ] **Step 2: Write integration tests**
- [ ] **Step 3: Run tests and commit**

---

### Task 21: MotionService

**Files:**
- Modify: `motion.go`
- Create: `motion_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 99-107 (Section 4.5.25-26), pp. 416-419 (Section 9.8)

Implement these endpoints:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `GetConfig` | `configManager.cgi?action=getConfig&name=MotionDetect` | 4.5.25 (p.99) |
| `SetConfig` | `configManager.cgi?action=setConfig` | 4.5.25 (p.99) |
| `GetSmartMotionConfig` | `configManager.cgi?action=getConfig&name=SmartMotionDetect` | 9.8.1 (p.416) |
| `SetSmartMotionConfig` | `configManager.cgi?action=setConfig` | 9.8.1 (p.416) |

- [ ] **Step 1: Read PDF pages and implement all endpoints**
- [ ] **Step 2: Write integration tests**
- [ ] **Step 3: Run tests and commit**

---

### Task 22: DisplayService

**Files:**
- Modify: `display.go`
- Create: `display_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 279-283 (Section 7)

Implement these endpoints:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `GetGUIConfig` | `configManager.cgi?action=getConfig&name=GUISet` | 7.1.1 (p.279) |
| `SetGUIConfig` | `configManager.cgi?action=setConfig` | 7.1.1 (p.279) |
| `GetSplitMode` | `split.cgi?action=getMode` | 7.2.1 (p.280) |
| `SetSplitMode` | `split.cgi?action=setMode` | 7.2.1 (p.280) |
| `GetMonitorTour` | `configManager.cgi?action=getConfig&name=MonitorTour` | 7.3.1 (p.281) |
| `SetMonitorTour` | `configManager.cgi?action=setConfig` | 7.3.1 (p.281) |

- [ ] **Step 1: Read PDF pp. 279-283 and implement all endpoints**
- [ ] **Step 2: Write integration tests**
- [ ] **Step 3: Run tests and commit**

---

### Task 23: UpgradeService

**Files:**
- Modify: `upgrade.go`
- Create: `upgrade_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 213-218 (Section 4.12)

Implement these endpoints:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `GetState` | `upgrader.cgi?action=getState` | 4.12.2 (p.214) |
| `Cancel` | `upgrader.cgi?action=cancel` | 4.12.4 (p.215) |
| `CheckCloudUpdate` | `api/CloudUpgrader/check` | 4.12.5 (p.215) |

- [ ] **Step 1: Read PDF pp. 213-218 and implement endpoints** (UploadFirmware and ExecuteOnlineUpdate behind `dangerous` tag)
- [ ] **Step 2: Write integration tests** (GetState is safe; all write operations are dangerous)
- [ ] **Step 3: Run tests and commit**

---

### Task 24: UploadService

**Files:**
- Modify: `upload.go`
- Create: `upload_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 219-234 (Section 4.13)

Implement these endpoints:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `GetPictureUploadConfig` | `configManager.cgi?action=getConfig&name=PictureHttpUpload` | 4.13.1 (p.219) |
| `SetPictureUploadConfig` | `configManager.cgi?action=setConfig` | 4.13.1 (p.219) |
| `GetEventUploadConfig` | `configManager.cgi?action=getConfig&name=EventHttpUpload` | 4.13.3 (p.221) |
| `SetEventUploadConfig` | `configManager.cgi?action=setConfig` | 4.13.3 (p.221) |

- [ ] **Step 1: Read PDF pp. 219-234 and implement all endpoints**
- [ ] **Step 2: Write integration tests**
- [ ] **Step 3: Run tests and commit**

---

### Task 25: AnalyticsService

**Files:**
- Modify: `analytics.go`
- Create: `analytics_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 381-401 (Section 9.6)

Implement these endpoints:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `GetCaps` | `devVideoAnalyse.cgi?action=getCaps` | 9.6.1 (p.381) |
| `GetGlobalConfig` | `configManager.cgi?action=getConfig&name=VideoAnalyseGlobal` | 9.6.2 (p.382) |
| `SetGlobalConfig` | `configManager.cgi?action=setConfig` | 9.6.2 (p.382) |
| `GetRuleConfig` | `configManager.cgi?action=getConfig&name=VideoAnalyseRule` | 9.6.3 (p.384) |
| `SetRuleConfig` | `configManager.cgi?action=setConfig` | 9.6.3 (p.384) |
| `GetSceneList` | `devVideoAnalyse.cgi?action=getSceneList` | 9.6.15 (p.400) |
| `EnableScene` | `devVideoAnalyse.cgi?action=enableScene` | 9.6.16 (p.400) |
| `DisableScene` | `devVideoAnalyse.cgi?action=disableScene` | 9.6.16 (p.400) |

- [ ] **Step 1: Read PDF pp. 381-401 and implement all endpoints**
- [ ] **Step 2: Write integration tests**
- [ ] **Step 3: Run tests and commit**

---

### Task 26: FaceService

**Files:**
- Modify: `face.go`
- Create: `face_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 334-361 (Section 9.2)

Implement these endpoints:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `CreateGroup` | `faceRecognitionServer.cgi?action=createGroup` | 9.2.1 (p.334) |
| `ModifyGroup` | `faceRecognitionServer.cgi?action=modifyGroup` | 9.2.2 (p.334) |
| `DeleteGroup` | `faceRecognitionServer.cgi?action=deleteGroup` | 9.2.3 (p.334) |
| `FindGroup` | `faceRecognitionServer.cgi?action=findGroup` | 9.2.5 (p.337) |
| `DeployGroup` | `faceRecognitionServer.cgi?action=putDisposition` | 9.2.4 (p.335) |
| `GetGroupForChannel` | `faceRecognitionServer.cgi?action=getGroup` | 9.2.4 (p.335) |

- [ ] **Step 1: Read PDF pp. 334-361 and implement all endpoints**
- [ ] **Step 2: Write integration tests** (CreateGroup/DeleteGroup as a pair with cleanup)
- [ ] **Step 3: Run tests and commit**

---

### Task 27: PeopleService

**Files:**
- Modify: `people.go`
- Create: `people_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 362-381 (Sections 9.3, 9.4, 9.5)

Implement these endpoints:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `GetSummary` | `api/VideoStatServer/getSummary` | 9.3.1 (p.362) |
| `QueryCount` | `api/VideoStatServer/queryCount` | 9.3.2 (p.364) |
| `ClearCount` | `api/VideoStatServer/clearCount` | 9.3.3 (p.367) |
| `SubscribeCount` | `api/VideoStatServer/subscribeStat` | 9.3.4 (p.367) |
| `GetHeatMap` | `api/HeatMapStat/getHeatMapInfo` | 9.4.1 (p.373) |
| `GetCrowdCaps` | `api/CrowdDistriMap/getChannelCaps` | 9.5.1 (p.379) |
| `GetCrowdStat` | `api/CrowdDistriMap/getCurrentStat` | 9.5.3 (p.381) |

- [ ] **Step 1: Read PDF pp. 362-381 and implement all endpoints**
- [ ] **Step 2: Write integration tests**
- [ ] **Step 3: Run tests and commit**

---

### Task 28: WorkSuitService

**Files:**
- Modify: `worksuit.go`
- Create: `worksuit_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 402-416 (Section 9.7)

Implement these endpoints:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `CreateGroup` | `api/WorkSuitCompareServer/createGroup` | 9.7.1 (p.405) |
| `DeleteGroup` | `api/WorkSuitCompareServer/deleteGroup` | 9.7.2 (p.406) |
| `FindGroup` | `api/WorkSuitCompareServer/findGroup` | 9.7.3 (p.407) |
| `GetGroup` | `api/WorkSuitCompareServer/getGroup` | 9.7.4 (p.408) |
| `ModifyGroup` | `api/WorkSuitCompareServer/modifyGroup` | 9.7.5 (p.409) |
| `DeployGroup` | `api/WorkSuitCompareServer/setGroup` | 9.7.6 (p.410) |

- [ ] **Step 1: Read PDF pp. 402-416 and implement all endpoints**
- [ ] **Step 2: Write integration tests**
- [ ] **Step 3: Run tests and commit**

---

### Task 29: TrafficService

**Files:**
- Modify: `traffic.go`
- Create: `traffic_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 440-466 (Section 10)

Implement these endpoints:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `FindFlowHistory` | `recordFinder.cgi?action=find&name=TrafficFlow` | 10.2.2 (p.446) |
| `StartFlowSearch` | `api/trafficFlowStat/startFind` | 10.2.3 (p.447) |
| `DoFlowSearch` | `api/trafficFlowStat/doFind` | 10.2.4 (p.448) |
| `StopFlowSearch` | `api/trafficFlowStat/stopFind` | 10.2.5 (p.449) |
| `InsertRecord` | `recordUpdater.cgi?action=insert` | 10.3.1 (p.450) |
| `UpdateRecord` | `recordUpdater.cgi?action=update` | 10.3.2 (p.451) |
| `RemoveRecord` | `recordUpdater.cgi?action=remove` | 10.3.3 (p.451) |
| `FindRecord` | `recordFinder.cgi?action=find` | 10.3.4 (p.452) |
| `ManualSnap` | `api/trafficSnap/manualSnap` | 10.4.3 (p.459) |

- [ ] **Step 1: Read PDF pp. 440-466 and implement all endpoints**
- [ ] **Step 2: Write integration tests**
- [ ] **Step 3: Run tests and commit**

---

### Task 30: ParkingService

**Files:**
- Modify: `parking.go`
- Create: `parking_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 460-465 (Section 10.5)

Implement these endpoints:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `GetSpaceStatus` | `api/TrafficParking/getSpaceStatus` | 10.5.1 (p.460) |
| `GetAllSpaceStatus` | `api/TrafficParking/getAllSpaceStatus` | 10.5.2 (p.461) |
| `GetLightConfig` | `configManager.cgi?action=getConfig&name=TrafficParkingSpaceLightState` | 10.5.3 (p.461) |

- [ ] **Step 1: Read PDF pp. 460-465 and implement all endpoints**
- [ ] **Step 2: Write integration tests**
- [ ] **Step 3: Run tests and commit**

---

### Task 31: ThermalService

**Files:**
- Modify: `thermal.go`
- Create: `thermal_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 481-504 (Section 11)

Implement these endpoints:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `GetCaps` | `RadiometryManager.cgi?action=getCaps` | 11.1.1 (p.481) |
| `GetThermographyOptions` | `configManager.cgi?action=getConfig&name=ThermographyOptions` | 11.1.2 (p.482) |
| `GetRadiometryCaps` | `RadiometryManager.cgi?action=getRadiometryCaps` | 11.2.1 (p.486) |
| `GetThermometryConfig` | `configManager.cgi?action=getConfig&name=HeatImagingThermometry` | 11.2.2 (p.487) |
| `GetThermometryRule` | `configManager.cgi?action=getConfig&name=ThermometryRule` | 11.2.3 (p.489) |
| `GetTemperature` | `RadiometryManager.cgi?action=getTemperature` | 11.2.5 (p.492) |
| `GetFireWarningConfig` | `configManager.cgi?action=getConfig&name=FireWarning` | 11.2.11 (p.498) |
| `GetCurrentHotColdSpot` | `TemperCorrection.cgi?action=getCurrentHotColdSpot` | 11.2.15 (p.501) |

- [ ] **Step 1: Read PDF pp. 481-504 and implement all endpoints**
- [ ] **Step 2: Write integration tests**
- [ ] **Step 3: Run tests and commit**

---

### Task 32: AccessControlService

**Files:**
- Modify: `accesscontrol.go`
- Create: `accesscontrol_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 505-584 (Section 12)

Implement these endpoints:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `OpenDoor` | `accessControl.cgi?action=openDoor` | 12.1.1 (p.505) |
| `CloseDoor` | `accessControl.cgi?action=closeDoor` | 12.1.2 (p.505) |
| `GetDoorStatus` | `accessControl.cgi?action=getDoorStatus` | 12.1.3 (p.506) |
| `GetLockStatus` | `accessControl.cgi?action=getLockStatus` | 12.1.4 (p.506) |
| `QueryRecords` | `recordFinder.cgi?action=find&name=AccessControlCardRec` | 12.1.7 (p.509) |
| `QueryAlarmRecords` | `recordFinder.cgi?action=find&name=AccessControlAlarmRecord` | 12.1.8 (p.511) |
| `GetGeneralConfig` | `configManager.cgi?action=getConfig&name=AccessControlGeneral` | 12.1.12 (p.517) |
| `GetControlConfig` | `configManager.cgi?action=getConfig&name=AccessControl` | 12.1.13 (p.520) |
| `GetCaps` | `api/AccessControl/getCaps` | 12.2.1 (p.533) |

- [ ] **Step 1: Read PDF pp. 505-584 and implement key endpoints** (focus on the most commonly used access control operations)
- [ ] **Step 2: Write integration tests**
- [ ] **Step 3: Run tests and commit**

---

### Task 33: BuildingService

**Files:**
- Modify: `building.go`
- Create: `building_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 585-608 (Section 13)

Implement these endpoints:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `GetSIPConfig` | `configManager.cgi?action=getConfig&name=SIPConfig` | 13.3.1 (p.593) |
| `SetSIPConfig` | `configManager.cgi?action=setConfig` | 13.3.1 (p.593) |
| `GetRoomNumberCount` | `recordFinder.cgi?action=getQuerySize&name=VideoTalkContact` | 13.4.7 (p.606) |
| `FindRoomNumbers` | `recordFinder.cgi?action=find&name=VideoTalkContact` | 13.4.2 (p.600) |

- [ ] **Step 1: Read PDF pp. 585-608 and implement key endpoints**
- [ ] **Step 2: Write integration tests**
- [ ] **Step 3: Run tests and commit**

---

### Task 34: DVRService

**Files:**
- Modify: `dvr.go`
- Create: `dvr_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 609-617 (Section 14)

Implement these endpoints:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `StartFind` | `FileFindHelper.cgi?action=startFind` | 14.1.1 (p.609) |
| `FindNext` | `FileFindHelper.cgi?action=findNext` | 14.1.3 (p.611) |
| `StopFind` | `FileFindHelper.cgi?action=stopFind` | 14.1.4 (p.612) |
| `GetBandwidthLimit` | `BandLimit.cgi?action=getLimitState` | 14.2.1 (p.613) |
| `DownloadFile` | `FileManager.cgi?action=downloadFile` | 14.3.4 (p.614) |

- [ ] **Step 1: Read PDF pp. 609-617 and implement all endpoints**
- [ ] **Step 2: Write integration tests**
- [ ] **Step 3: Run tests and commit**

---

### Task 35: PeripheralService

**Files:**
- Modify: `peripheral.go`
- Create: `peripheral_test.go`

**PDF Reference:** `docs/HTTP_API_V3.26.pdf` pp. 305-325 (Sections 8.2-8.8), pp. 618-650 (Section 15)

Implement these endpoints:

| Method | API Endpoint | PDF Section |
|--------|-------------|-------------|
| `WiperStart` | `rainBrush.cgi?action=moveContinuously` | 8.2.1 (p.306) |
| `WiperStop` | `rainBrush.cgi?action=stopMove` | 8.2.2 (p.307) |
| `WiperOnce` | `rainBrush.cgi?action=moveOnce` | 8.2.3 (p.307) |
| `ControlCoaxialIO` | `coaxialControlIO.cgi?action=control` | 8.5.1 (p.309) |
| `GetCoaxialIOStatus` | `coaxialControlIO.cgi?action=getstatus` | 8.5.2 (p.310) |
| `GetFlashlightConfig` | `configManager.cgi?action=getConfig&name=FlashLight` | 8.4.1 (p.308) |
| `SetFlashlightConfig` | `configManager.cgi?action=setConfig` | 8.4.1 (p.308) |
| `GetGPSStatus` | `api/GPS/getStatus` | 15.3.3 (p.623) |
| `GetGPSConfig` | `configManager.cgi?action=getConfig&name=GPS` | 15.3.2 (p.623) |
| `GetFishEyeConfig` | `configManager.cgi?action=getConfig&name=FishEye` | 15.5.2 (p.629) |

- [ ] **Step 1: Read PDF pages and implement all endpoints**
- [ ] **Step 2: Write integration tests**
- [ ] **Step 3: Run tests and commit**

---

## Phase 3: Final Integration & Polish

### Task 36: Full Integration Test Suite

**Files:**
- Create: `integration_test.go`

- [ ] **Step 1: Create a comprehensive integration test that exercises the full client**

```go
//go:build integration

package amcrest

import (
	"context"
	"testing"
	"time"
)

func TestFullIntegration(t *testing.T) {
	c := testClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// System
	t.Run("System/DeviceType", func(t *testing.T) {
		dt, err := c.System.GetDeviceType(ctx)
		if err != nil {
			t.Fatalf("GetDeviceType: %v", err)
		}
		t.Logf("Device type: %s", dt)
	})

	// Snapshot
	t.Run("Snapshot/Get", func(t *testing.T) {
		data, err := c.Snapshot.Get(ctx, 1)
		if err != nil {
			t.Fatalf("Snapshot.Get: %v", err)
		}
		if len(data) < 2 || data[0] != 0xFF || data[1] != 0xD8 {
			t.Fatal("not a valid JPEG")
		}
		t.Logf("Snapshot size: %d bytes", len(data))
	})

	// Network
	t.Run("Network/Interfaces", func(t *testing.T) {
		ifaces, err := c.Network.GetInterfaces(ctx)
		if err != nil {
			t.Fatalf("GetInterfaces: %v", err)
		}
		t.Logf("Interfaces: %+v", ifaces)
	})
}
```

- [ ] **Step 2: Run the full integration test suite**

Run: `cd /Users/ethanflower/personal_projects/amcrest-sdk && go test -tags integration -v`
Expected: All tests PASS

- [ ] **Step 3: Commit**

```bash
git add integration_test.go
git commit -m "feat: add full integration test suite"
```

---

### Task 37: README

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update README with installation, usage examples, and service listing**

Include: module install command, NewClient example, examples for 5-6 popular operations (snapshot, system info, PTZ, events, recording search), and a table of all services.

- [ ] **Step 2: Commit**

```bash
git add README.md
git commit -m "docs: update README with SDK usage and examples"
```
