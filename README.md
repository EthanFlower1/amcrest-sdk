# amcrest

Go SDK for the Amcrest HTTP API V3.26. Covers all 15 chapters of the API specification with 28 domain services.

## Install

```bash
go get github.com/EthanFlower/amcrest
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/EthanFlower/amcrest"
)

func main() {
    client, err := amcrest.NewClient("192.168.1.100", "admin", "password")
    if err != nil {
        log.Fatal(err)
    }
    ctx := context.Background()

    // Get device info
    version, _ := client.System.GetSoftwareVersion(ctx)
    serial, _ := client.System.GetSerialNumber(ctx)
    fmt.Printf("Camera: %s (SN: %s)\n", version, serial)

    // Take a snapshot
    jpeg, _ := client.Snapshot.Get(ctx, 1) // channel 1
    os.WriteFile("snapshot.jpg", jpeg, 0644)

    // Subscribe to motion events
    events, stream, _ := client.Event.Subscribe(ctx, []string{"VideoMotion"}, 5)
    defer stream.Close()
    for event := range events {
        fmt.Printf("Event: %s action=%s\n", event.Code, event.Action)
    }
}
```

## Services

All services are accessed through the `Client` struct:

| Service | Access | Description |
|---------|--------|-------------|
| **System** | `client.System` | Device info, time, reboot, version, language |
| **User** | `client.User` | User/group CRUD, passwords |
| **Network** | `client.Network` | Interfaces, DDNS, email, WiFi, NTP, UPnP, RTSP |
| **Video** | `client.Video` | Encode config, video input/output, channel titles |
| **Snapshot** | `client.Snapshot` | Capture snapshots, snap config |
| **Audio** | `client.Audio` | Audio channels, volume, streaming |
| **PTZ** | `client.PTZ` | Pan/tilt/zoom, presets, tours, patterns |
| **Event** | `client.Event` | Subscribe to events, alarm config |
| **Recording** | `client.Recording` | Record config, media file search, download |
| **Log** | `client.Log` | Find, clear, backup logs |
| **Storage** | `client.Storage` | Disk info, format, NAS, SD encryption |
| **Camera** | `client.Camera` | Image, exposure, white balance, day/night, focus |
| **Display** | `client.Display` | GUI, split screen, monitor tour |
| **Privacy** | `client.Privacy` | Privacy masking CRUD |
| **Motion** | `client.Motion` | Motion detection, smart motion detection |
| **Upgrade** | `client.Upgrade` | Firmware upgrade, cloud update |
| **Upload** | `client.Upload` | HTTP uploading config |
| **Analytics** | `client.Analytics` | Video analysis rules, scene management |
| **Face** | `client.Face` | Face recognition groups |
| **People** | `client.People` | People counting, crowd stats |
| **WorkSuit** | `client.WorkSuit` | Workwear compliance detection |
| **Traffic** | `client.Traffic` | Traffic records, flow stats |
| **Parking** | `client.Parking` | Parking space status |
| **Thermal** | `client.Thermal` | Thermography, radiometry |
| **AccessControl** | `client.AccessControl` | Door control, access config |
| **Building** | `client.Building` | SIP, room numbers |
| **DVR** | `client.DVR` | File finder, record protection |
| **Peripheral** | `client.Peripheral` | Wiper, coaxial IO, GPS, fisheye, flashlight |

## Client Options

```go
// HTTPS
client, _ := amcrest.NewClient("192.168.1.100", "admin", "pass", amcrest.WithHTTPS())

// Custom port
client, _ := amcrest.NewClient("192.168.1.100", "admin", "pass", amcrest.WithPort(8080))

// Custom timeout
client, _ := amcrest.NewClient("192.168.1.100", "admin", "pass", amcrest.WithTimeout(10*time.Second))
```

## Examples

### Search and download recordings

```go
files, _ := client.Recording.FindFiles(ctx, amcrest.FindFilesOpts{
    Channel:   1,
    StartTime: "2024-01-15 00:00:00",
    EndTime:   "2024-01-15 23:59:59",
})

for _, f := range files {
    fmt.Printf("%s (%s) %d bytes\n", f.FilePath, f.Type, f.Length)
}

// Download a file
data, _ := client.Recording.DownloadFile(ctx, files[0].FilePath)
```

### PTZ control

```go
// Move camera right
client.PTZ.Control(ctx, 1, amcrest.PTZRight, 0, 1, 0)
time.Sleep(2 * time.Second)
client.PTZ.Stop(ctx, 1, amcrest.PTZRight)

// Go to preset
client.PTZ.GotoPreset(ctx, 1, 1)
```

### Search logs

```go
entries, _ := client.Log.Find(ctx, "2024-01-15 00:00:00", "2024-01-15 23:59:59", "")
for _, e := range entries {
    fmt.Printf("[%s] %s: %s\n", e.Time, e.Type, e.Detail)
}
```

## Testing

Tests run against a real Amcrest camera. Configure via environment variables or `.env` file:

```bash
AMCREST_HOST=192.168.1.100
AMCREST_USERNAME=admin
AMCREST_PASSWORD=your_password
```

```bash
go test ./...                    # Safe read-only tests
go test -tags dangerous ./...    # Including destructive tests (reboot, format, etc.)
```

## License

Apache 2.0
