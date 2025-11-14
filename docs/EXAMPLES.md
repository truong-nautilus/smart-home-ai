# Examples

## Basic Usage

### Starting the System

```bash
# Make sure .env and config.json are configured
make run
```

### Voice Commands Examples

**Vietnamese:**
- "Bật đèn phòng khách"
- "Tắt đèn phòng ngủ"
- "Đặt độ sáng đèn 80%"
- "Bật điều hòa 26 độ"
- "Tắt quạt"
- "Bắt đầu hút bụi"

**English:**
- "Turn on living room light"
- "Turn off bedroom light"
- "Set light brightness to 80%"
- "Set air conditioner to 26 degrees"
- "Turn off fan"
- "Start vacuum cleaning"

## Command JSON Format

Claude AI will respond with JSON commands:

### Light Commands

```json
// Turn on
{"action": "light.on", "device": "phong_khach"}

// Turn off
{"action": "light.off", "device": "phong_khach"}

// Set brightness
{"action": "light.brightness", "device": "phong_khach", "value": 80}

// Set color
{"action": "light.color", "device": "phong_khach", "value": {"hue": 120, "saturation": 100}}

// Set color temperature
{"action": "light.color_temp", "device": "phong_khach", "value": 3000}
```

### Switch Commands

```json
// Turn on
{"action": "switch.on", "device": "quat_phong_khach"}

// Turn off
{"action": "switch.off", "device": "quat_phong_khach"}

// Toggle
{"action": "switch.toggle", "device": "quat_phong_khach"}
```

### AC Commands

```json
// Turn on
{"action": "ac.on", "device": "dieu_hoa_phong_khach"}

// Turn off
{"action": "ac.off", "device": "dieu_hoa_phong_khach"}

// Set temperature
{"action": "ac.set_temp", "device": "dieu_hoa_phong_khach", "value": 26}
```

### Vacuum Commands

```json
// Start cleaning
{"action": "vacuum.start", "device": "robot_hut_bui"}

// Stop
{"action": "vacuum.stop", "device": "robot_hut_bui"}

// Pause
{"action": "vacuum.pause", "device": "robot_hut_bui"}

// Go home
{"action": "vacuum.home", "device": "robot_hut_bui"}

// Set fan speed
{"action": "vacuum.fan_speed", "device": "robot_hut_bui", "value": 60}
```

## Programmatic Usage

### Initialize System

```go
package main

import (
    "github.com/truong-nautilus/smart-home-ai/core"
    "github.com/truong-nautilus/smart-home-ai/devices"
)

func main() {
    // Load config
    config, _ := core.LoadConfig("config.json")
    
    // Create router
    router := core.NewCommandRouter(config)
    
    // Initialize devices
    tapoConfig := devices.TapoConfig{
        Email: "user@example.com",
        Password: "password",
    }
    
    mqttConfig := devices.MQTTConfig{
        Host: "192.168.1.100",
        Port: 1883,
    }
    
    router.Initialize(tapoConfig, mqttConfig)
    
    // Execute command
    cmd := &core.Command{
        Action: "light.on",
        Device: "phong_khach",
    }
    
    router.ExecuteCommand(cmd)
}
```

### Control Tapo Device

```go
package main

import (
    "github.com/truong-nautilus/smart-home-ai/devices"
)

func main() {
    config := devices.TapoConfig{
        Email: "user@example.com",
        Password: "password",
    }
    
    device := devices.NewTapoDevice("192.168.1.10", "L530", config)
    
    // Handshake
    device.Handshake()
    
    // Login
    device.Login()
    
    // Turn on
    device.TurnOn()
    
    // Set brightness
    device.SetBrightness(80)
    
    // Set color
    device.SetColor(120, 100)
}
```

### MQTT Control

```go
package main

import (
    "github.com/truong-nautilus/smart-home-ai/devices"
)

func main() {
    config := devices.MQTTConfig{
        Host: "192.168.1.100",
        Port: 1883,
    }
    
    client := devices.NewMQTTClient(config)
    client.Connect()
    defer client.Disconnect()
    
    // Turn on light
    client.TurnOnLight("home/living/light")
    
    // Set brightness
    client.SetBrightness("home/living/light", 80)
    
    // Subscribe to state
    client.Subscribe("home/living/light/state", func(client mqtt.Client, msg mqtt.Message) {
        println("State:", string(msg.Payload()))
    })
}
```

### Broadlink IR Control

```go
package main

import (
    "github.com/truong-nautilus/smart-home-ai/devices"
    "time"
)

func main() {
    device := devices.NewBroadlinkDevice("192.168.1.30", 80)
    
    // Discover
    device.Discover(5 * time.Second)
    
    // Authenticate
    device.Auth()
    
    // Send IR command
    device.SendIRCommand("260050000001...")
    
    // Or learn a new command
    code, _ := device.LearnIRCommand(30 * time.Second)
    println("Learned code:", code)
}
```

### Xiaomi Vacuum Control

```go
package main

import (
    "github.com/truong-nautilus/smart-home-ai/devices"
)

func main() {
    vacuum, _ := devices.NewVacuumRobot(
        "192.168.1.40",
        "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
    )
    
    // Start cleaning
    vacuum.Start()
    
    // Set fan speed
    vacuum.SetFanSpeed(60)
    
    // Get status
    status, _ := vacuum.GetStatus()
    println("Battery:", status["battery"])
    
    // Go home
    vacuum.Home()
}
```

## Testing

### Test Individual Device

```go
package main

import (
    "testing"
    "github.com/truong-nautilus/smart-home-ai/devices"
)

func TestTapoDevice(t *testing.T) {
    config := devices.TapoConfig{
        Email: "user@example.com",
        Password: "password",
    }
    
    device := devices.NewTapoDevice("192.168.1.10", "L530", config)
    
    if err := device.Handshake(); err != nil {
        t.Fatalf("Handshake failed: %v", err)
    }
    
    if err := device.Login(); err != nil {
        t.Fatalf("Login failed: %v", err)
    }
    
    if err := device.TurnOn(); err != nil {
        t.Fatalf("TurnOn failed: %v", err)
    }
}
```

### Test Command Parsing

```go
package main

import (
    "testing"
    "github.com/truong-nautilus/smart-home-ai/core"
)

func TestParseCommand(t *testing.T) {
    text := `{"action":"light.on","device":"phong_khach","value":100}`
    
    cmd, err := core.ParseCommand(text)
    if err != nil {
        t.Fatalf("Parse failed: %v", err)
    }
    
    if cmd.Action != "light.on" {
        t.Errorf("Expected action light.on, got %s", cmd.Action)
    }
    
    if cmd.Device != "phong_khach" {
        t.Errorf("Expected device phong_khach, got %s", cmd.Device)
    }
}
```

## Debugging

### Enable Debug Logging

```go
// Set environment variable
DEBUG=true

// Or in code
log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
```

### Monitor MQTT Messages

```bash
# Subscribe to all topics
mosquitto_sub -h 192.168.1.100 -t "#" -v

# Subscribe to specific device
mosquitto_sub -h 192.168.1.100 -t "home/living/light/#" -v
```

### Test Claude API

```bash
# Test WebSocket connection
wscat -c "wss://api.anthropic.com/v1/messages" \
  -H "anthropic-version: 2023-06-01" \
  -H "x-api-key: $CLAUDE_API_KEY"
```

### Check Audio Input

```bash
# List audio devices (macOS)
system_profiler SPAudioDataType

# Test microphone
rec -r 16000 -c 1 test.wav
```

## Advanced Examples

### Custom Scene

```go
func executeScene(router *core.CommandRouter, sceneName string) error {
    scenes := map[string][]core.Command{
        "goodnight": {
            {Action: "light.off", Device: "phong_khach"},
            {Action: "light.off", Device: "phong_ngu"},
            {Action: "ac.off", Device: "dieu_hoa"},
        },
        "morning": {
            {Action: "light.on", Device: "phong_khach", Value: 100},
            {Action: "light.on", Device: "bep", Value: 80},
        },
    }
    
    commands, ok := scenes[sceneName]
    if !ok {
        return fmt.Errorf("scene not found: %s", sceneName)
    }
    
    for _, cmd := range commands {
        if err := router.ExecuteCommand(&cmd); err != nil {
            log.Printf("Error executing command: %v", err)
        }
    }
    
    return nil
}
```

### Voice Response

```go
// Play audio response
func playAudioResponse(audioData []byte) {
    // Use portaudio or similar to play audio
    // Implementation depends on audio library
}
```

## Troubleshooting Examples

### Check Device Connectivity

```go
func checkDeviceConnectivity(ip string) bool {
    conn, err := net.DialTimeout("tcp", ip+":80", 3*time.Second)
    if err != nil {
        return false
    }
    conn.Close()
    return true
}
```

### Retry Failed Commands

```go
func executeWithRetry(router *core.CommandRouter, cmd *core.Command, maxRetries int) error {
    var err error
    for i := 0; i < maxRetries; i++ {
        err = router.ExecuteCommand(cmd)
        if err == nil {
            return nil
        }
        time.Sleep(time.Second * time.Duration(i+1))
    }
    return err
}
```
