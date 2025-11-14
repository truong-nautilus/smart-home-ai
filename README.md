# Jarvis AI - Smart Home Voice Assistant

🏠 Hệ thống trợ lý giọng nói thông minh sử dụng Claude AI Realtime API để điều khiển thiết bị nhà thông minh theo thời gian thực.

## 🌟 Tính năng chính

- ✅ **Claude AI Realtime**: Sử dụng Claude 3.5 Sonnet qua WebSocket
- 🎤 **Voice Input**: Nhận lệnh giọng nói từ microphone (PCM 16-bit, 16kHz)
- 🔊 **Real-time Audio Streaming**: Xử lý audio theo thời gian thực
- 🏡 **Multi-Device Support**: Hỗ trợ nhiều loại thiết bị:
  - 💡 **Tapo** (P100 switches, L530 smart bulbs)
  - 📡 **Broadlink** (RM4 IR controller)
  - 🌐 **MQTT** (Shelly, Sonoff, ESP32)
  - 🤖 **Xiaomi Miio** (Robot vacuum, smart lights)
  - 🔌 **HTTP Generic** (Any REST API device)
- 🛡️ **Security**: Rate limiting, command validation, audit logging
- 🚀 **Concurrent Processing**: Go routines & channels for parallel execution

## 📁 Cấu trúc dự án

```
smart-home-ai/
├── audio/              # Audio recording & processing
│   ├── recorder.go     # Microphone input handler
│   └── recorder_test.go
├── claude/             # Claude Realtime API client
│   ├── client.go       # WebSocket client
│   └── client_test.go
├── devices/            # Device controllers
│   ├── tapo.go        # Tapo devices (P100, L530)
│   ├── broadlink.go   # Broadlink IR/RF
│   ├── mqtt.go        # MQTT devices (Shelly, Sonoff, ESP32)
│   └── xiaomi.go      # Xiaomi Miio devices
├── core/              # Core logic
│   ├── router.go      # Command router
│   ├── security.go    # Security manager
│   └── config.go      # Configuration loader
├── main.go            # Application entry point
├── config.json        # Device configuration
├── .env.example       # Environment variables template
├── Makefile          # Build & run commands
└── README.md         # This file
```

## 🚀 Quick Start

### 1. Prerequisites

- Go 1.22 or higher
- macOS (for audio capture)
- Claude API key
- Smart home devices (Tapo, Broadlink, MQTT broker, etc.)

### 2. Installation

```bash
# Clone repository
git clone https://github.com/truong-nautilus/smart-home-ai.git
cd smart-home-ai

# Install dependencies
make deps

# Setup environment
make setup
```

### 3. Configuration

**Edit `.env` file:**

```bash
# Claude AI Configuration
CLAUDE_API_KEY=your_claude_api_key_here

# MQTT Configuration
MQTT_HOST=192.168.1.100
MQTT_PORT=1883
MQTT_USER=
MQTT_PASS=

# Tapo Configuration
TAPO_USER=your_tapo_email@example.com
TAPO_PASS=your_tapo_password

# Xiaomi Configuration
XIAOMI_TOKEN=your_xiaomi_token_here
XIAOMI_IP=192.168.1.102
```

**Edit `config.json`:** Cấu hình thiết bị của bạn

```json
{
  "devices": {
    "lights": {
      "phong_khach": {
        "type": "tapo",
        "model": "L530",
        "ip": "192.168.1.10",
        "name": "Đèn Phòng Khách"
      }
    }
  }
}
```

### 4. Run

```bash
# Run directly
make run

# Or build and run
make build
./bin/jarvis
```

## 🎯 Usage Examples

Sau khi chạy, bạn có thể nói các lệnh sau:

### 💡 Điều khiển đèn
- "Bật đèn phòng khách"
- "Tắt đèn phòng ngủ"
- "Đặt độ sáng đèn phòng khách 80%"
- "Đổi màu đèn sang đỏ"

### ❄️ Điều khiển điều hòa
- "Bật điều hòa"
- "Đặt nhiệt độ điều hòa 26 độ"
- "Tắt điều hòa phòng khách"

### 🤖 Điều khiển robot hút bụi
- "Bắt đầu hút bụi"
- "Dừng robot hút bụi"
- "Về sạc"

### 📺 Điều khiển TV (qua IR)
- "Bật TV"
- "Tăng âm lượng"
- "Giảm âm lượng"

## 🔧 Device Integration

### Tapo Devices

```go
// Tự động được xử lý qua config.json
// Hỗ trợ: P100 (switch), L530 (smart bulb)
```

### Broadlink IR

```go
// Thêm IR codes vào config.json
"ir_devices": {
  "dieu_hoa_phong_khach": {
    "type": "broadlink",
    "device_ip": "192.168.1.30",
    "commands": {
      "on": "260050000001...",
      "temp_26": "260050000001..."
    }
  }
}
```

### MQTT Devices

```go
// Cấu hình topic trong config.json
"lights": {
  "bep": {
    "type": "mqtt",
    "topic": "home/kitchen/light",
    "name": "Đèn Bếp"
  }
}
```

### Xiaomi Devices

```go
// Cần token từ Mi Home app
// Xem hướng dẫn: https://github.com/rytilahti/python-miio
```

## 📊 Command Format

Claude AI sẽ trả về JSON command với format:

```json
{
  "action": "light.on",
  "device": "phong_khach",
  "value": 100
}
```

### Supported Actions

**Lights:**
- `light.on` - Bật đèn
- `light.off` - Tắt đèn
- `light.brightness` - Đặt độ sáng (1-100)
- `light.color` - Đổi màu (hue, saturation)
- `light.color_temp` - Đặt nhiệt độ màu (2500-6500K)

**Switches:**
- `switch.on` - Bật công tắc
- `switch.off` - Tắt công tắc
- `switch.toggle` - Đảo trạng thái

**AC (Air Conditioner):**
- `ac.on` - Bật điều hòa
- `ac.off` - Tắt điều hòa
- `ac.set_temp` - Đặt nhiệt độ (16-30)

**Vacuum:**
- `vacuum.start` - Bắt đầu hút
- `vacuum.stop` - Dừng
- `vacuum.pause` - Tạm dừng
- `vacuum.home` - Về sạc
- `vacuum.fan_speed` - Đặt tốc độ quạt

## 🔒 Security Features

- **Rate Limiting**: Giới hạn 10 lệnh/phút
- **Command Validation**: Kiểm tra lệnh hợp lệ
- **Audit Logging**: Ghi log tất cả lệnh
- **Allowed Commands**: Whitelist các lệnh được phép

## 🛠️ Development

### Run Tests

```bash
make test
```

### Format Code

```bash
make fmt
```

### Run Linter

```bash
make lint
```

### Live Reload (with air)

```bash
make tools  # Install air
make dev    # Run with live reload
```

## 🐛 Troubleshooting

### Audio Issues

```bash
# Check microphone permissions
# System Preferences > Security & Privacy > Microphone

# Test audio device
# Make sure your microphone is working
```

### Claude API Connection

```bash
# Verify API key
echo $CLAUDE_API_KEY

# Check internet connection
# Claude API requires stable internet
```

### Device Connection Issues

```bash
# Tapo: Verify email/password
# MQTT: Check broker is running
# Broadlink: Device must be on same network
# Xiaomi: Verify token is correct
```

## 📚 API References

- [Claude API Documentation](https://docs.anthropic.com/claude/reference/streaming)
- [Tapo Protocol](https://github.com/fishbigger/TapoP100)
- [Broadlink Protocol](https://github.com/mjg59/python-broadlink)
- [MQTT Protocol](https://mqtt.org/)
- [Xiaomi Miio Protocol](https://github.com/rytilahti/python-miio)

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📝 License

MIT License - see LICENSE file for details

## 👨‍💻 Author

Created by truong-nautilus

## 🙏 Acknowledgments

- Claude AI by Anthropic
- Go community
- Smart home device manufacturers

---

**Note**: This project requires Claude API access and compatible smart home devices. Make sure to configure all devices properly before running.

## 🔮 Future Improvements

- [ ] Add voice response playback
- [ ] Support more device types
- [ ] Web dashboard for monitoring
- [ ] Mobile app integration
- [ ] Scene automation
- [ ] Multi-language support
- [ ] Cloud sync for configurations

## 💬 Support

For issues and questions:
- GitHub Issues: [Create an issue](https://github.com/truong-nautilus/smart-home-ai/issues)
- Email: phamthetruong@example.com

---

**Happy Smart Home Automation! 🏠✨**
