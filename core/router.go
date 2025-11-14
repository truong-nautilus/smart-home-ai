package core

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/truong-nautilus/smart-home-ai/devices"
)

// Config represents the application configuration
type Config struct {
	Devices DevicesConfig `json:"devices"`
	Claude  ClaudeConfig  `json:"claude"`
	Audio   AudioConfig   `json:"audio"`
}

// DevicesConfig holds all device configurations
type DevicesConfig struct {
	Lights    map[string]DeviceInfo   `json:"lights"`
	Switches  map[string]DeviceInfo   `json:"switches"`
	IRDevices map[string]IRDeviceInfo `json:"ir_devices"`
	Vacuum    map[string]DeviceInfo   `json:"vacuum"`
}

// DeviceInfo holds basic device information
type DeviceInfo struct {
	Type  string `json:"type"`
	Model string `json:"model"`
	IP    string `json:"ip"`
	Topic string `json:"topic,omitempty"`
	Name  string `json:"name"`
}

// IRDeviceInfo holds IR device information
type IRDeviceInfo struct {
	Type     string            `json:"type"`
	DeviceIP string            `json:"device_ip"`
	Commands map[string]string `json:"commands"`
	Name     string            `json:"name"`
}

// ClaudeConfig holds Claude configuration
type ClaudeConfig struct {
	Model        string  `json:"model"`
	WebSocketURL string  `json:"websocket_url"`
	MaxTokens    int     `json:"max_tokens"`
	Temperature  float64 `json:"temperature"`
	SystemPrompt string  `json:"system_prompt"`
}

// AudioConfig holds audio configuration
type AudioConfig struct {
	SampleRate int `json:"sample_rate"`
	Channels   int `json:"channels"`
	BitDepth   int `json:"bit_depth"`
	BufferSize int `json:"buffer_size"`
}

// Command represents a parsed command
type Command struct {
	Action string      `json:"action"`
	Device string      `json:"device"`
	Value  interface{} `json:"value,omitempty"`
}

// ParseCommand parses a command from text or JSON
func ParseCommand(text string) (*Command, error) {
	text = strings.TrimSpace(text)

	// Try parsing as JSON
	var cmd Command
	if err := json.Unmarshal([]byte(text), &cmd); err == nil {
		if cmd.Action != "" {
			return &cmd, nil
		}
	}

	// Try extracting JSON from text
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start >= 0 && end > start {
		jsonStr := text[start : end+1]
		if err := json.Unmarshal([]byte(jsonStr), &cmd); err == nil {
			if cmd.Action != "" {
				return &cmd, nil
			}
		}
	}

	return nil, fmt.Errorf("not a valid command")
}

// CommandRouter routes commands to appropriate device controllers
type CommandRouter struct {
	config        *Config
	tapoDevices   map[string]*devices.TapoDevice
	broadlink     map[string]*devices.BroadlinkDevice
	mqttClient    *devices.MQTTClient
	xiaomiDevices map[string]interface{}
	httpDevices   map[string]*devices.HTTPDevice
}

// NewCommandRouter creates a new command router
func NewCommandRouter(config *Config) *CommandRouter {
	return &CommandRouter{
		config:        config,
		tapoDevices:   make(map[string]*devices.TapoDevice),
		broadlink:     make(map[string]*devices.BroadlinkDevice),
		xiaomiDevices: make(map[string]interface{}),
		httpDevices:   make(map[string]*devices.HTTPDevice),
	}
}

// Initialize initializes all device connections
func (r *CommandRouter) Initialize(tapoConfig devices.TapoConfig, mqttConfig devices.MQTTConfig) error {
	log.Println("Initializing device connections...")

	// Initialize Tapo devices
	for id, info := range r.config.Devices.Lights {
		if info.Type == "tapo" {
			device := devices.NewTapoDevice(info.IP, info.Model, tapoConfig)
			r.tapoDevices[id] = device
			log.Printf("Initialized Tapo device: %s (%s)", info.Name, id)
		}
	}

	for id, info := range r.config.Devices.Switches {
		if info.Type == "tapo" {
			device := devices.NewTapoDevice(info.IP, info.Model, tapoConfig)
			r.tapoDevices[id] = device
			log.Printf("Initialized Tapo device: %s (%s)", info.Name, id)
		}
	}

	// Initialize Broadlink devices
	for id, info := range r.config.Devices.IRDevices {
		if info.Type == "broadlink" {
			device := devices.NewBroadlinkDevice(info.DeviceIP, 80)
			r.broadlink[id] = device
			log.Printf("Initialized Broadlink device: %s (%s)", info.Name, id)
		}
	}

	// Initialize MQTT client
	if mqttConfig.Host != "" {
		r.mqttClient = devices.NewMQTTClient(mqttConfig)
		if err := r.mqttClient.Connect(); err != nil {
			log.Printf("Warning: Failed to connect to MQTT broker: %v", err)
		} else {
			log.Println("Connected to MQTT broker")
		}
	}

	log.Println("Device initialization complete")
	return nil
}

// ExecuteCommand executes a command
func (r *CommandRouter) ExecuteCommand(cmd *Command) error {
	log.Printf("Executing command: action=%s, device=%s, value=%v", cmd.Action, cmd.Device, cmd.Value)

	parts := strings.Split(cmd.Action, ".")
	if len(parts) != 2 {
		return fmt.Errorf("invalid action format: %s", cmd.Action)
	}

	deviceType := parts[0]
	action := parts[1]

	switch deviceType {
	case "light":
		return r.executeLight(cmd.Device, action, cmd.Value)
	case "switch":
		return r.executeSwitch(cmd.Device, action, cmd.Value)
	case "ac":
		return r.executeAC(cmd.Device, action, cmd.Value)
	case "vacuum":
		return r.executeVacuum(cmd.Device, action, cmd.Value)
	case "tv":
		return r.executeTV(cmd.Device, action, cmd.Value)
	default:
		return fmt.Errorf("unknown device type: %s", deviceType)
	}
}

// executeLight executes light commands
func (r *CommandRouter) executeLight(deviceID, action string, value interface{}) error {
	// Check if it's a Tapo device
	if device, ok := r.tapoDevices[deviceID]; ok {
		switch action {
		case "on":
			return device.TurnOn()
		case "off":
			return device.TurnOff()
		case "brightness":
			if brightness, ok := value.(float64); ok {
				return device.SetBrightness(int(brightness))
			}
			return fmt.Errorf("invalid brightness value")
		case "color":
			if colorMap, ok := value.(map[string]interface{}); ok {
				hue := int(colorMap["hue"].(float64))
				sat := int(colorMap["saturation"].(float64))
				return device.SetColor(hue, sat)
			}
			return fmt.Errorf("invalid color value")
		case "color_temp":
			if temp, ok := value.(float64); ok {
				return device.SetColorTemp(int(temp))
			}
			return fmt.Errorf("invalid color temperature value")
		default:
			return fmt.Errorf("unknown light action: %s", action)
		}
	}

	// Check if it's an MQTT device
	if info, ok := r.config.Devices.Lights[deviceID]; ok && info.Type == "mqtt" {
		if r.mqttClient == nil {
			return fmt.Errorf("MQTT client not initialized")
		}

		switch action {
		case "on":
			return r.mqttClient.TurnOnLight(info.Topic)
		case "off":
			return r.mqttClient.TurnOffLight(info.Topic)
		case "brightness":
			if brightness, ok := value.(float64); ok {
				return r.mqttClient.SetBrightness(info.Topic, int(brightness))
			}
			return fmt.Errorf("invalid brightness value")
		default:
			return fmt.Errorf("unknown light action: %s", action)
		}
	}

	return fmt.Errorf("device not found: %s", deviceID)
}

// executeSwitch executes switch commands
func (r *CommandRouter) executeSwitch(deviceID, action string, value interface{}) error {
	// Check if it's a Tapo device
	if device, ok := r.tapoDevices[deviceID]; ok {
		switch action {
		case "on":
			return device.TurnOn()
		case "off":
			return device.TurnOff()
		default:
			return fmt.Errorf("unknown switch action: %s", action)
		}
	}

	// Check if it's an MQTT device
	if info, ok := r.config.Devices.Switches[deviceID]; ok && info.Type == "mqtt" {
		if r.mqttClient == nil {
			return fmt.Errorf("MQTT client not initialized")
		}

		switch action {
		case "on":
			return r.mqttClient.TurnOnSwitch(info.Topic)
		case "off":
			return r.mqttClient.TurnOffSwitch(info.Topic)
		case "toggle":
			return r.mqttClient.ToggleSwitch(info.Topic)
		default:
			return fmt.Errorf("unknown switch action: %s", action)
		}
	}

	return fmt.Errorf("device not found: %s", deviceID)
}

// executeAC executes air conditioner commands
func (r *CommandRouter) executeAC(deviceID, action string, value interface{}) error {
	info, ok := r.config.Devices.IRDevices[deviceID]
	if !ok {
		return fmt.Errorf("AC device not found: %s", deviceID)
	}

	device, ok := r.broadlink[deviceID]
	if !ok {
		return fmt.Errorf("Broadlink device not initialized: %s", deviceID)
	}

	var irCode string
	switch action {
	case "on":
		irCode = info.Commands["on"]
	case "off":
		irCode = info.Commands["off"]
	case "set_temp":
		if temp, ok := value.(float64); ok {
			tempKey := fmt.Sprintf("temp_%d", int(temp))
			irCode = info.Commands[tempKey]
		} else {
			return fmt.Errorf("invalid temperature value")
		}
	default:
		// Try to find command in device commands
		if cmd, ok := info.Commands[action]; ok {
			irCode = cmd
		} else {
			return fmt.Errorf("unknown AC action: %s", action)
		}
	}

	if irCode == "" {
		return fmt.Errorf("IR code not found for action: %s", action)
	}

	return device.SendIRCommand(irCode)
}

// executeVacuum executes vacuum commands
func (r *CommandRouter) executeVacuum(deviceID, action string, value interface{}) error {
	info, ok := r.config.Devices.Vacuum[deviceID]
	if !ok {
		return fmt.Errorf("vacuum device not found: %s", deviceID)
	}

	if info.Type != "xiaomi" {
		return fmt.Errorf("unsupported vacuum type: %s", info.Type)
	}

	// Get or create Xiaomi vacuum device
	var vacuum *devices.VacuumRobot
	if dev, ok := r.xiaomiDevices[deviceID]; ok {
		vacuum = dev.(*devices.VacuumRobot)
	} else {
		// Note: Token should be from environment or config
		return fmt.Errorf("xiaomi device not initialized: %s", deviceID)
	}

	switch action {
	case "start":
		return vacuum.Start()
	case "stop":
		return vacuum.Stop()
	case "pause":
		return vacuum.Pause()
	case "home":
		return vacuum.Home()
	case "spot":
		return vacuum.Spot()
	case "fan_speed":
		if speed, ok := value.(float64); ok {
			return vacuum.SetFanSpeed(int(speed))
		}
		return fmt.Errorf("invalid fan speed value")
	default:
		return fmt.Errorf("unknown vacuum action: %s", action)
	}
}

// executeTV executes TV commands via IR
func (r *CommandRouter) executeTV(deviceID, action string, value interface{}) error {
	info, ok := r.config.Devices.IRDevices[deviceID]
	if !ok {
		return fmt.Errorf("TV device not found: %s", deviceID)
	}

	device, ok := r.broadlink[deviceID]
	if !ok {
		return fmt.Errorf("Broadlink device not initialized: %s", deviceID)
	}

	irCode, ok := info.Commands[action]
	if !ok {
		return fmt.Errorf("IR code not found for action: %s", action)
	}

	return device.SendIRCommand(irCode)
}

// Close closes all device connections
func (r *CommandRouter) Close() {
	log.Println("Closing device connections...")

	if r.mqttClient != nil {
		r.mqttClient.Disconnect()
	}

	log.Println("All device connections closed")
}
