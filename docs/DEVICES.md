# Device Configuration Guide

## Tapo Devices (TP-Link)

### Supported Models
- **P100/P105**: Smart Plug
- **L530**: Smart Bulb (Color)
- **L510**: Smart Bulb (White)

### Configuration

```json
{
  "lights": {
    "phong_khach": {
      "type": "tapo",
      "model": "L530",
      "ip": "192.168.1.10",
      "name": "Đèn Phòng Khách"
    }
  },
  "switches": {
    "quat": {
      "type": "tapo",
      "model": "P100",
      "ip": "192.168.1.20",
      "name": "Quạt"
    }
  }
}
```

### Environment Variables

```bash
TAPO_USER=your_tapo_email@example.com
TAPO_PASS=your_tapo_password
```

### Setup Steps
1. Install Tapo app
2. Add devices to app
3. Get device local IP address
4. Use same email/password in .env

---

## Broadlink Devices

### Supported Models
- **RM4 Pro**: IR/RF Controller
- **RM Mini**: IR Controller

### Configuration

```json
{
  "ir_devices": {
    "dieu_hoa": {
      "type": "broadlink",
      "device_ip": "192.168.1.30",
      "commands": {
        "on": "26005000000126...",
        "off": "26005000000126...",
        "temp_18": "26005000000126...",
        "temp_26": "26005000000126..."
      },
      "name": "Điều Hòa"
    }
  }
}
```

### Learning IR Codes

```bash
# Use the Broadlink app or python-broadlink to learn codes
pip install broadlink
python -c "import broadlink; ..."
```

Or create a learning tool:

```go
device := devices.NewBroadlinkDevice("192.168.1.30", 80)
device.Auth()
code, err := device.LearnIRCommand(30 * time.Second)
fmt.Println("IR Code:", code)
```

---

## MQTT Devices

### Supported Platforms
- **Shelly**: Smart switches and relays
- **Sonoff**: Smart switches (Tasmota firmware)
- **ESP32/ESP8266**: Custom devices

### Configuration

```json
{
  "lights": {
    "bep": {
      "type": "mqtt",
      "topic": "home/kitchen/light",
      "name": "Đèn Bếp"
    }
  }
}
```

### Environment Variables

```bash
MQTT_HOST=192.168.1.100
MQTT_PORT=1883
MQTT_USER=mqtt_username
MQTT_PASS=mqtt_password
```

### MQTT Broker Setup

**Using Mosquitto:**

```bash
# Install
brew install mosquitto  # macOS
sudo apt install mosquitto  # Ubuntu

# Start
brew services start mosquitto
# or
sudo systemctl start mosquitto

# Test
mosquitto_pub -t "test/topic" -m "hello"
mosquitto_sub -t "test/topic"
```

### Device-Specific Topics

**Shelly:**
```
shellies/shelly1/relay/0/command  -> "on" / "off"
shellies/shelly1/relay/0          <- status
```

**Sonoff (Tasmota):**
```
cmnd/sonoff1/POWER  -> "ON" / "OFF"
stat/sonoff1/POWER  <- status
```

**Custom ESP32:**
```
home/device/power   -> "on" / "off"
home/device/state   <- status
```

---

## Xiaomi Devices

### Supported Models
- **Vacuum Robots**: Mi Robot Vacuum, Roborock
- **Smart Lights**: Yeelight
- **Air Purifiers**: Mi Air Purifier

### Getting Token

**Method 1: Mi Home App (iOS)**
1. Install older version of Mi Home (3.5.8)
2. Login and add device
3. Get token from app logs

**Method 2: python-miio**
```bash
pip install python-miio
mirobo --ip <device_ip> --token 00000000000000000000000000000000 info
```

**Method 3: Android**
1. Use MiHome Mod app
2. Add device
3. View token in device settings

### Configuration

```json
{
  "vacuum": {
    "robot_hut_bui": {
      "type": "xiaomi",
      "ip": "192.168.1.40",
      "name": "Robot Hút Bụi"
    }
  }
}
```

### Environment Variables

```bash
XIAOMI_TOKEN=a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6
XIAOMI_IP=192.168.1.40
```

---

## HTTP Generic Devices

### Configuration

```json
{
  "custom": {
    "smart_lock": {
      "type": "http",
      "base_url": "http://192.168.1.50",
      "headers": {
        "Authorization": "Bearer token123",
        "Content-Type": "application/json"
      },
      "name": "Smart Lock"
    }
  }
}
```

### Usage

```go
device := devices.NewHTTPDevice(
    "http://192.168.1.50",
    map[string]string{
        "Authorization": "Bearer token123",
    },
)

// Send command
resp, err := device.Post("/lock", []byte(`{"action":"unlock"}`))
```

---

## Network Setup

### Static IP Assignment

**Recommended:** Assign static IPs to all devices via router DHCP reservation

**Why:**
- Devices won't change IP after reboot
- Configuration remains valid
- Easier troubleshooting

### Network Requirements

- All devices must be on same network/VLAN
- Firewall must allow:
  - Port 80/443: HTTP/HTTPS
  - Port 1883/8883: MQTT
  - Port 54321: Xiaomi Miio
  - UDP broadcasts: Device discovery

### Testing Connectivity

```bash
# Ping device
ping 192.168.1.10

# Check port
nc -zv 192.168.1.10 80

# MQTT test
mosquitto_sub -h 192.168.1.100 -t "#" -v
```

---

## Troubleshooting

### Tapo Issues

**Problem:** Authentication failed
- Verify email/password
- Try logging out and in on Tapo app
- Check if account has 2FA enabled

**Problem:** Device not responding
- Check IP address is correct
- Ensure device is on same network
- Try power cycling device

### Broadlink Issues

**Problem:** Discovery fails
- Device must be in AP mode for first setup
- Check firewall allows UDP broadcasts
- Ensure device is on same network

**Problem:** IR command not working
- Re-learn the command
- Test command with official app first
- Check device battery (for remote)

### MQTT Issues

**Problem:** Connection refused
- Check MQTT broker is running
- Verify host/port/credentials
- Check firewall rules

**Problem:** Messages not received
- Verify topic names match
- Check QoS settings
- Use MQTT Explorer to debug

### Xiaomi Issues

**Problem:** Invalid token
- Token changes after factory reset
- Re-extract token from app
- Ensure token is 32 hex characters

**Problem:** Device not responding
- Check IP address
- Verify token is correct
- Try pinging device

---

## Security Best Practices

1. **Use strong passwords** for all devices
2. **Isolate smart home network** using VLAN
3. **Enable encryption** for MQTT (TLS/SSL)
4. **Regular firmware updates** for devices
5. **Disable UPnP** on router
6. **Use firewall rules** to restrict access
7. **Monitor logs** for suspicious activity

---

## Advanced Configuration

### Custom Device Types

Create new device controller:

```go
package devices

type CustomDevice struct {
    IP string
}

func (d *CustomDevice) TurnOn() error {
    // Implementation
}
```

Add to router:

```go
case "custom":
    return r.executeCustom(cmd.Device, action, cmd.Value)
```

### Scene Automation

Create scenes in config.json:

```json
{
  "scenes": {
    "goodnight": {
      "commands": [
        {"action": "light.off", "device": "phong_khach"},
        {"action": "light.off", "device": "phong_ngu"},
        {"action": "ac.off", "device": "dieu_hoa"}
      ]
    }
  }
}
```

### Scheduling

Use cron for scheduled tasks:

```bash
# Crontab entry
0 22 * * * curl http://localhost:8080/scene/goodnight
```

---

## References

- [Tapo Protocol](https://github.com/fishbigger/TapoP100)
- [Broadlink Protocol](https://github.com/mjg59/python-broadlink)
- [MQTT Specification](https://mqtt.org/)
- [Xiaomi Miio](https://github.com/rytilahti/python-miio)
