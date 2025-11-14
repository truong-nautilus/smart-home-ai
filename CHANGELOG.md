# Changelog

All notable changes to Jarvis AI Smart Home will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-01-14

### Added
- Initial release of Jarvis AI Smart Home
- Claude AI Realtime API integration
- Real-time voice input processing
- Multi-device support:
  - Tapo smart devices (P100, L530)
  - Broadlink IR/RF controllers
  - MQTT devices (Shelly, Sonoff, ESP32)
  - Xiaomi Miio devices
  - HTTP generic devices
- Audio recording from microphone (PCM 16-bit, 16kHz)
- WebSocket communication with Claude API
- Command routing and execution
- Security features:
  - Rate limiting
  - Command validation
  - Audit logging
- Concurrent processing with goroutines
- Configuration via JSON and environment variables
- Comprehensive documentation
- Example code and usage patterns
- Build system with Makefile
- Unit tests for core components

### Features
- **Voice Control**: Natural language voice commands
- **Multi-Device**: Support for multiple smart home platforms
- **Real-Time**: Low-latency audio streaming and command execution
- **Secure**: Built-in security validation and rate limiting
- **Extensible**: Easy to add new device types
- **Well-Documented**: Complete API and device setup documentation

### Technical Details
- Written in Go 1.22+
- Uses gorilla/websocket for real-time communication
- Uses malgo for cross-platform audio capture
- Uses paho.mqtt.golang for MQTT support
- Clean architecture with separation of concerns
- Thread-safe concurrent operations
- Comprehensive error handling

### Documentation
- README with quick start guide
- API documentation
- Device setup guide
- Configuration examples
- Troubleshooting guide
- Contributing guidelines

### Supported Commands
- Light control (on/off, brightness, color, temperature)
- Switch control (on/off, toggle)
- Air conditioner control (temperature, on/off)
- Vacuum robot control (start, stop, home, fan speed)
- TV control via IR (power, volume)
- Custom device integration

### Requirements
- Go 1.22 or higher
- macOS (for audio capture)
- Claude API key
- Compatible smart home devices

---

## [Unreleased]

### Planned Features
- Voice response playback
- Web dashboard for monitoring
- Mobile app integration
- Scene automation
- Multi-language support
- Cloud sync for configurations
- More device type support
- Scheduled tasks
- Energy monitoring
- Usage statistics

---

**Note**: This is version 1.0.0 - the initial stable release.
