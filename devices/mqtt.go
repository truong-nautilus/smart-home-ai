package devices

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQTTConfig holds MQTT configuration
type MQTTConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	ClientID string
}

// MQTTClient handles MQTT device communication
type MQTTClient struct {
	config MQTTConfig
	client mqtt.Client
}

// NewMQTTClient creates a new MQTT client
func NewMQTTClient(config MQTTConfig) *MQTTClient {
	return &MQTTClient{
		config: config,
	}
}

// Connect connects to MQTT broker
func (m *MQTTClient) Connect() error {
	opts := mqtt.NewClientOptions()

	broker := fmt.Sprintf("tcp://%s:%d", m.config.Host, m.config.Port)
	opts.AddBroker(broker)
	opts.SetClientID(m.config.ClientID)

	if m.config.Username != "" {
		opts.SetUsername(m.config.Username)
		opts.SetPassword(m.config.Password)
	}

	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)

	// TLS configuration
	opts.SetTLSConfig(&tls.Config{
		InsecureSkipVerify: true,
	})

	// Connection callbacks
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		log.Println("Connected to MQTT broker")
	})

	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Printf("MQTT connection lost: %v", err)
	})

	m.client = mqtt.NewClient(opts)

	token := m.client.Connect()
	token.Wait()

	if err := token.Error(); err != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %w", err)
	}

	log.Printf("Connected to MQTT broker: %s", broker)
	return nil
}

// Disconnect disconnects from MQTT broker
func (m *MQTTClient) Disconnect() {
	if m.client != nil && m.client.IsConnected() {
		m.client.Disconnect(250)
		log.Println("Disconnected from MQTT broker")
	}
}

// Publish publishes a message to a topic
func (m *MQTTClient) Publish(topic string, payload interface{}) error {
	if !m.client.IsConnected() {
		return fmt.Errorf("not connected to MQTT broker")
	}

	token := m.client.Publish(topic, 0, false, payload)
	token.Wait()

	if err := token.Error(); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published to %s: %v", topic, payload)
	return nil
}

// Subscribe subscribes to a topic
func (m *MQTTClient) Subscribe(topic string, callback mqtt.MessageHandler) error {
	if !m.client.IsConnected() {
		return fmt.Errorf("not connected to MQTT broker")
	}

	token := m.client.Subscribe(topic, 0, callback)
	token.Wait()

	if err := token.Error(); err != nil {
		return fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	log.Printf("Subscribed to topic: %s", topic)
	return nil
}

// Unsubscribe unsubscribes from a topic
func (m *MQTTClient) Unsubscribe(topic string) error {
	if !m.client.IsConnected() {
		return fmt.Errorf("not connected to MQTT broker")
	}

	token := m.client.Unsubscribe(topic)
	token.Wait()

	if err := token.Error(); err != nil {
		return fmt.Errorf("failed to unsubscribe from topic: %w", err)
	}

	log.Printf("Unsubscribed from topic: %s", topic)
	return nil
}

// IsConnected returns connection status
func (m *MQTTClient) IsConnected() bool {
	return m.client != nil && m.client.IsConnected()
}

// Common MQTT device control methods

// TurnOnLight turns on a light via MQTT
func (m *MQTTClient) TurnOnLight(topic string) error {
	return m.Publish(topic+"/set", "ON")
}

// TurnOffLight turns off a light via MQTT
func (m *MQTTClient) TurnOffLight(topic string) error {
	return m.Publish(topic+"/set", "OFF")
}

// SetBrightness sets light brightness (0-100)
func (m *MQTTClient) SetBrightness(topic string, brightness int) error {
	payload := fmt.Sprintf(`{"state":"ON","brightness":%d}`, brightness)
	return m.Publish(topic+"/set", payload)
}

// SetColor sets light color (RGB)
func (m *MQTTClient) SetColor(topic string, r, g, b int) error {
	payload := fmt.Sprintf(`{"state":"ON","color":{"r":%d,"g":%d,"b":%d}}`, r, g, b)
	return m.Publish(topic+"/set", payload)
}

// TurnOnSwitch turns on a switch via MQTT
func (m *MQTTClient) TurnOnSwitch(topic string) error {
	return m.Publish(topic+"/relay/0", "on")
}

// TurnOffSwitch turns off a switch via MQTT
func (m *MQTTClient) TurnOffSwitch(topic string) error {
	return m.Publish(topic+"/relay/0", "off")
}

// ToggleSwitch toggles a switch via MQTT
func (m *MQTTClient) ToggleSwitch(topic string) error {
	return m.Publish(topic+"/relay/0/command", "toggle")
}

// GetState gets device state
func (m *MQTTClient) GetState(topic string, callback mqtt.MessageHandler) error {
	// Subscribe to state topic
	if err := m.Subscribe(topic+"/state", callback); err != nil {
		return err
	}

	// Request state update
	return m.Publish(topic+"/get", "")
}

// ShellyDevice represents a Shelly device
type ShellyDevice struct {
	Topic  string
	client *MQTTClient
}

// NewShellyDevice creates a new Shelly device
func NewShellyDevice(topic string, client *MQTTClient) *ShellyDevice {
	return &ShellyDevice{
		Topic:  topic,
		client: client,
	}
}

// TurnOn turns on the Shelly device
func (s *ShellyDevice) TurnOn() error {
	return s.client.TurnOnSwitch(s.Topic)
}

// TurnOff turns off the Shelly device
func (s *ShellyDevice) TurnOff() error {
	return s.client.TurnOffSwitch(s.Topic)
}

// Toggle toggles the Shelly device
func (s *ShellyDevice) Toggle() error {
	return s.client.ToggleSwitch(s.Topic)
}

// SonoffDevice represents a Sonoff device
type SonoffDevice struct {
	Topic  string
	client *MQTTClient
}

// NewSonoffDevice creates a new Sonoff device
func NewSonoffDevice(topic string, client *MQTTClient) *SonoffDevice {
	return &SonoffDevice{
		Topic:  topic,
		client: client,
	}
}

// TurnOn turns on the Sonoff device
func (s *SonoffDevice) TurnOn() error {
	return s.client.Publish(s.Topic+"/cmnd/POWER", "ON")
}

// TurnOff turns off the Sonoff device
func (s *SonoffDevice) TurnOff() error {
	return s.client.Publish(s.Topic+"/cmnd/POWER", "OFF")
}

// Toggle toggles the Sonoff device
func (s *SonoffDevice) Toggle() error {
	return s.client.Publish(s.Topic+"/cmnd/POWER", "TOGGLE")
}

// ESP32Device represents a custom ESP32 device
type ESP32Device struct {
	Topic  string
	client *MQTTClient
}

// NewESP32Device creates a new ESP32 device
func NewESP32Device(topic string, client *MQTTClient) *ESP32Device {
	return &ESP32Device{
		Topic:  topic,
		client: client,
	}
}

// SendCommand sends a custom command to ESP32
func (e *ESP32Device) SendCommand(command string, value interface{}) error {
	topic := fmt.Sprintf("%s/%s", e.Topic, command)
	return e.client.Publish(topic, value)
}

// TurnOn turns on the ESP32 device
func (e *ESP32Device) TurnOn() error {
	return e.SendCommand("power", "on")
}

// TurnOff turns off the ESP32 device
func (e *ESP32Device) TurnOff() error {
	return e.SendCommand("power", "off")
}
