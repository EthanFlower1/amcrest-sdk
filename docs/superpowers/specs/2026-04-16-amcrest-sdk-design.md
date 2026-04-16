# Amcrest Go SDK Design Spec

**Date:** 2026-04-16
**Module:** `github.com/EthanFlower/amcrest`
**Go Version:** 1.21+
**API Spec:** Amcrest HTTP API V3.26 (650 pages, 15 chapters, ~300+ endpoints)

## Overview

A complete Go SDK for the Amcrest HTTP API V3.26, covering all 15 chapters of the specification. The SDK provides a single `Client` with embedded domain services for every API category -- from basic camera operations to traffic analytics, thermal imaging, and access control.

## Architecture

### Client Pattern: Embedded Services

A root `Client` holds a configured HTTP client (with digest auth) and exposes domain services as fields:

```go
client, err := amcrest.NewClient("192.168.1.218", "admin", "password")

info, err := client.System.GetSoftwareVersion(ctx)
snap, err := client.Snapshot.Get(ctx, 1)
err = client.PTZ.GotoPreset(ctx, 1, 1)
client.Face.CreateGroup(ctx, "visitors", "Visitor faces")
client.Thermal.GetTemperature(ctx, 1, point)
client.AccessControl.OpenDoor(ctx, 1)
```

### HTTP Layer

The Amcrest API uses two styles:

1. **CGI key-value style** (majority): `GET /cgi-bin/<module>.cgi?action=<action>&param=value` with `table.Foo.Bar=value` responses
2. **JSON API style** (newer): `POST /cgi-bin/api/<Module>/<action>` with JSON request/response bodies

**Digest Authentication:** Implements `http.RoundTripper` per RFC 7616. On 401 response, parses `WWW-Authenticate`, computes digest, retries. Caches nonce for subsequent requests.

**Key-Value Parser:** Parses `table.Foo[0].Bar=value` response lines into Go structs using reflection and struct tags:

```go
type EncodeConfig struct {
    BitRate     int     `kv:"BitRate"`
    Compression string  `kv:"Compression"`
    FPS         float64 `kv:"FPS"`
}
```

**Generic Config Helpers:** Since `configManager.cgi` covers ~50% of endpoints:

```go
func (c *Client) getConfig(ctx context.Context, name string, out any) error
func (c *Client) setConfig(ctx context.Context, params map[string]string) error
```

**Multipart Stream Reader:** For long-lived connections (event subscription, audio, snapshot subscription). Returns channels that yield parsed events:

```go
events, err := client.Events.Subscribe(ctx, []string{"VideoMotion", "FaceDetection"})
for event := range events {
    fmt.Println(event.Code, event.Action, event.Data)
}
```

**Error Types:**

```go
type APIError struct {
    StatusCode int
    Code       int    // from JSON ErrorCode
    Message    string // from JSON ErrorMsg or body text
}
```

## Package Structure

All files in root package `amcrest`. No subpackages.

### Core Infrastructure (6 files)

| File | Purpose |
|------|---------|
| `amcrest.go` | Client struct, NewClient(), Option funcs, service init |
| `auth.go` | Digest auth round-tripper (RFC 7616) |
| `parse.go` | Key-value response parser with struct tags |
| `config.go` | Generic configManager.cgi get/set helpers |
| `stream.go` | Multipart stream reader for events/audio/snapshots |
| `errors.go` | APIError type and error handling |

### Service Files (25 files)

| File | Service | API Chapters | Description |
|------|---------|-------------|-------------|
| `system.go` | SystemService | 4.6 | Device info, time, reboot, version, language, auto-maintain |
| `users.go` | UserService | 4.7 | User/group CRUD, passwords, auth policy, export |
| `network.go` | NetworkService | 4.8 | Interfaces, DDNS, email, WiFi, NTP, UPnP, RTSP, PPPoE, SSHD, cellular |
| `video.go` | VideoService | 4.5 | Encode config/caps, video input/output, channel titles, widgets, smart encode, video standard |
| `snapshot.go` | SnapshotService | 4.4 | Snap config, capture, subscribe to snapshot events |
| `audio.go` | AudioService | 4.3 | Input/output channels, streaming, volume, audio analysis |
| `ptz.go` | PTZService | 8.1 | Movement, presets, tours, scans, patterns, electronic PTZ, view range |
| `events.go` | EventService | 4.9 | Subscribe, alarm in/out config, blind/loss detect, event caps, net alarm |
| `recording.go` | RecordingService | 4.10 | Record config/mode, media file search (basic + face/traffic/IVS/etc), download |
| `logs.go` | LogService | 4.11 | Find, clear, backup, seek, export encrypted |
| `storage.go` | StorageService | 6 | Disk info, format, NAS, storage groups/points, SD encrypt, health alarm |
| `camera.go` | CameraService | 5 | Image settings, exposure, backlight, white balance, day/night, zoom/focus, lighting, video-in options |
| `display.go` | DisplayService | 7 | GUI settings, split screen, monitor tour |
| `privacy.go` | PrivacyService | 4.5.18-28 | Privacy masking CRUD, enable/disable, goto |
| `motion.go` | MotionService | 4.5.25, 9.8 | Motion detection config, smart motion detection, SMD data search |
| `upgrade.go` | UpgradeService | 4.12 | Firmware upload, upgrade by URL, cloud upgrade, state, cancel |
| `upload.go` | UploadService | 4.13 | HTTP uploading config (picture, event, report data) |
| `analytics.go` | AnalyticsService | 9.6 | Video analyse config, rules, capabilities, scene management, intelligent tour |
| `face.go` | FaceService | 9.2 | Face groups, persons, recognition config, search by picture, database export/import |
| `people.go` | PeopleService | 9.3, 9.4, 9.5 | People counting, heatmaps, crowd distribution, traces |
| `worksuit.go` | WorkSuitService | 9.7 | Compliance library, workwear detection |
| `traffic.go` | TrafficService | 10 | Traffic events, flow stats, record management, snap operations, vehicle distribution |
| `parking.go` | ParkingService | 10.5 | Parking space status, light control, access filter, overline |
| `thermal.go` | ThermalService | 11 | Thermography, radiometry, temperature measurement, fire warning, heat maps |
| `accesscontrol.go` | AccessControlService | 12 | Door control, status, config, events, user accounts (V1 & V2), cards, fingerprints, faces, admin passwords |
| `building.go` | BuildingService | 13 | Video talk, SIP, room numbers, elevator |
| `dvr.go` | DVRService | 14 | File finder, record protection, bandwidth, file transfer |
| `peripheral.go` | PeripheralService | 8.2-8.8, 15 | Wiper, illuminator, flashlight, coaxial IO, PIR, SCADA, gyro, GPS, lens, fisheye, radar, water quality, advertisement |

## Client API

### Constructor

```go
func NewClient(host, username, password string, opts ...Option) (*Client, error)
```

### Options

```go
func WithHTTPS() Option               // Use HTTPS (default: HTTP)
func WithPort(port int) Option         // Custom port (default: 80/443)
func WithHTTPClient(c *http.Client) Option  // Custom HTTP client
func WithTimeout(d time.Duration) Option    // Request timeout
```

### Conventions

- `context.Context` as first parameter on every method
- Channel numbers are 1-based in the SDK (matching request convention); 0-based response mapping handled internally
- Streaming methods return Go channels
- All methods return `error` as last return value
- Config get methods return typed structs
- Config set methods accept typed structs or option params

## API Patterns

### Pattern 1: ConfigManager Get/Set

~50% of endpoints use `configManager.cgi`. Each service wraps this:

```go
func (s *NetworkService) GetNTPConfig(ctx context.Context) (*NTPConfig, error)
func (s *NetworkService) SetNTPConfig(ctx context.Context, cfg *NTPConfig) error
```

### Pattern 2: Simple Action

```go
func (s *SystemService) Reboot(ctx context.Context) error
func (s *SystemService) GetSerialNumber(ctx context.Context) (string, error)
```

### Pattern 3: Stateful Finder (Media Files)

The SDK abstracts the create/findFile/findNextFile/close/destroy lifecycle:

```go
func (s *RecordingService) FindFiles(ctx context.Context, opts FindFilesOpts) ([]MediaFile, error)
```

Internally manages the object lifecycle. For large result sets, an iterator variant:

```go
func (s *RecordingService) FindFilesIter(ctx context.Context, opts FindFilesOpts) *MediaFileIterator
```

### Pattern 4: Token-Based Search (Logs, Records)

Similar abstraction over startFind/doFind/stopFind:

```go
func (s *LogService) Find(ctx context.Context, opts LogFindOpts) ([]LogEntry, error)
```

### Pattern 5: Event Subscription (Long-lived Stream)

```go
func (s *EventService) Subscribe(ctx context.Context, codes []string) (<-chan Event, error)
```

Returns a channel. Canceling the context closes the connection and channel. Heartbeats are handled internally.

## Testing Strategy

### Integration Tests Against Real Camera

Camera on network at `192.168.1.218`, credentials `admin` / `Gsd4life.`

### Configuration

- `.env` file (gitignored) for local dev
- `.env.example` committed with placeholders
- Environment variables: `AMCREST_HOST`, `AMCREST_USERNAME`, `AMCREST_PASSWORD`
- Env vars take precedence over `.env` file
- Tests skip via `t.Skip()` if not configured

### Test Categories

1. **Read-only tests** (default) -- `GetSoftwareVersion`, `GetSerialNo`, `GetSnapshot`, `GetCurrentTime`, etc. Safe to run anytime.

2. **Write tests with cleanup** -- Modify state, verify, restore original. E.g., `SetCurrentTime` saves current time, sets new, verifies, restores.

3. **Destructive tests** (`//go:build dangerous`) -- `Reboot`, `FactoryReset`, `FormatSDCard`. Never run automatically.

4. **Streaming tests with timeout** -- Connect, read a few messages or heartbeat, disconnect. Short context timeout.

### Running

```bash
go test ./...                    # Safe tests
go test -tags dangerous ./...    # Including destructive tests
```

### Test Helper

```go
func testClient(t *testing.T) *amcrest.Client {
    // Loads .env, checks env vars, skips if not set
    // Returns configured client
}
```

## Key Design Decisions

1. **Flat package** -- All services in root `amcrest` package. No subpackages for domains.
2. **Embedded services** -- `client.PTZ.Control(...)` style for discoverability and namespacing.
3. **1-based channels** -- SDK uses 1-based channel numbers everywhere, matching the request convention. Response 0-based indexing is mapped internally.
4. **No code generation** -- Hand-written for idiomatic Go.
5. **Integration tests only** -- Real camera on network, no mocks.
6. **Context everywhere** -- All methods take `context.Context` for cancellation and timeout.
