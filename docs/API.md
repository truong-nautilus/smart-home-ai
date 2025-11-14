# API Documentation

## Claude Realtime API Integration

### WebSocket Connection

```go
claudeClient := claude.NewRealtimeClient(config)
err := claudeClient.Connect()
```

### Sending Audio

```go
// Send audio data (PCM 16-bit, 16kHz, Mono)
err := claudeClient.SendAudio(audioData)

// Commit buffer to trigger processing
err := claudeClient.CommitAudioBuffer()
```

### Receiving Commands

```go
commandChan := claudeClient.GetCommandChannel()
for cmd := range commandChan {
    // Process command
}
```

## Device Controllers

### Tapo Devices

```go
// Initialize
device := devices.NewTapoDevice(ip, model, config)

// Turn on/off
device.TurnOn()
device.TurnOff()

// Set brightness (1-100)
device.SetBrightness(80)

// Set color (hue: 0-360, saturation: 0-100)
device.SetColor(120, 100)

// Set color temperature (2500-6500K)
device.SetColorTemp(3000)

// Get device info
info, err := device.GetDeviceInfo()
```

### Broadlink Devices

```go
// Initialize
device := devices.NewBroadlinkDevice(ip, port)

// Discover device
err := device.Discover(5 * time.Second)

// Authenticate
err := device.Auth()

// Send IR command
err := device.SendIRCommand("260050000001...")

// Learn IR command
code, err := device.LearnIRCommand(30 * time.Second)
```

### MQTT Devices

```go
// Initialize client
client := devices.NewMQTTClient(config)
err := client.Connect()

// Control light
err := client.TurnOnLight("home/living/light")
err := client.TurnOffLight("home/living/light")
err := client.SetBrightness("home/living/light", 80)

// Control switch
err := client.TurnOnSwitch("shellies/shelly1/relay")
err := client.ToggleSwitch("shellies/shelly1/relay")

// Publish custom message
err := client.Publish("home/device/command", "ON")

// Subscribe to topic
err := client.Subscribe("home/device/state", handler)
```

### Xiaomi Devices

```go
// Vacuum Robot
vacuum, err := devices.NewVacuumRobot(ip, token)
err := vacuum.Start()
err := vacuum.Stop()
err := vacuum.Home()
status, err := vacuum.GetStatus()

// Smart Light
light, err := devices.NewXiaomiLight(ip, token)
err := light.TurnOn()
err := light.SetBrightness(80)
err := light.SetRGB(255, 0, 0)

// Air Purifier
purifier, err := devices.NewXiaomiAirPurifier(ip, token)
err := purifier.TurnOn()
err := purifier.SetMode("auto")
```

### HTTP Devices

```go
// Initialize
device := devices.NewHTTPDevice(baseURL, headers)

// Send requests
response, err := device.Get("/status")
response, err := device.Post("/control", data)
response, err := device.Put("/update", data)
```

## Command Router

### Initialize Router

```go
router := core.NewCommandRouter(config)
err := router.Initialize(tapoConfig, mqttConfig)
```

### Execute Command

```go
cmd := &core.Command{
    Action: "light.on",
    Device: "phong_khach",
    Value:  100,
}

err := router.ExecuteCommand(cmd)
```

## Security Manager

### Validate Command

```go
security := core.NewSecurityManager()
err := security.ValidateCommand(cmd)
```

### Command Logging

```go
security.LogCommand(cmd, success, err)
logs := security.GetCommandLog(100)
```

### Manage Allowed Commands

```go
security.AddAllowedCommand("custom.action")
security.RemoveAllowedCommand("dangerous.action")
isAllowed := security.IsCommandAllowed("light.on")
```

## Audio Recorder

### Initialize Recorder

```go
recorder, err := audio.NewRecorder(
    16000, // Sample rate
    1,     // Channels (mono)
    3200,  // Buffer size
)
```

### Start/Stop Recording

```go
err := recorder.Start()
defer recorder.Stop()
```

### Get Audio Stream

```go
audioChan := recorder.GetAudioChannel()
for audioData := range audioChan {
    // Process audio
}
```

## Configuration

### Load Configuration

```go
config, err := core.LoadConfig("config.json")
```

### Save Configuration

```go
err := core.SaveConfig("config.json", config)
```

## Error Handling

All functions return errors that should be checked:

```go
if err := device.TurnOn(); err != nil {
    log.Printf("Error: %v", err)
}
```

## Concurrency

The system uses goroutines and channels:

```go
// Audio streaming
go func() {
    for audioData := range audioChan {
        claudeClient.SendAudio(audioData)
    }
}()

// Command processing
go func() {
    for cmd := range commandChan {
        router.ExecuteCommand(cmd)
    }
}()
```

## Best Practices

1. Always check errors
2. Use defer for cleanup
3. Close connections when done
4. Handle panics in goroutines
5. Use contexts for cancellation
6. Validate input data
7. Log important events
8. Test with timeouts
