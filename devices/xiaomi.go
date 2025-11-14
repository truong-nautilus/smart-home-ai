package devices

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// XiaomiConfig holds Xiaomi device configuration
type XiaomiConfig struct {
	Token string
	IP    string
}

// XiaomiDevice represents a Xiaomi Miio device
type XiaomiDevice struct {
	IP       string
	Token    []byte
	DeviceID uint32
	Stamp    time.Time
	client   *http.Client
}

// XiaomiRequest represents a request to Xiaomi device
type XiaomiRequest struct {
	ID     int           `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params,omitempty"`
}

// XiaomiResponse represents a response from Xiaomi device
type XiaomiResponse struct {
	ID     int           `json:"id"`
	Result []interface{} `json:"result,omitempty"`
	Error  interface{}   `json:"error,omitempty"`
}

// NewXiaomiDevice creates a new Xiaomi device
func NewXiaomiDevice(ip, token string) (*XiaomiDevice, error) {
	tokenBytes, err := hex.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if len(tokenBytes) != 16 {
		return nil, fmt.Errorf("token must be 16 bytes")
	}

	return &XiaomiDevice{
		IP:    ip,
		Token: tokenBytes,
		Stamp: time.Now(),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

// Discover discovers Xiaomi devices on the network
func (x *XiaomiDevice) Discover() error {
	// Send handshake packet
	packet := make([]byte, 32)
	packet[0] = 0x21
	packet[1] = 0x31

	// Add timestamp
	now := time.Now().Unix()
	packet[12] = byte(now >> 24)
	packet[13] = byte(now >> 16)
	packet[14] = byte(now >> 8)
	packet[15] = byte(now)

	// Calculate checksum
	checksum := md5.Sum(packet)
	copy(packet[16:], checksum[:])

	// Send packet via UDP
	conn, err := net.DialTimeout("udp", fmt.Sprintf("%s:54321", x.IP), 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	_, err = conn.Write(packet)
	if err != nil {
		return fmt.Errorf("failed to send discovery packet: %w", err)
	}

	// Read response
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return fmt.Errorf("no response from device: %w", err)
	}

	if n < 32 {
		return fmt.Errorf("invalid response length")
	}

	// Extract device ID
	response := buffer[:n]
	x.DeviceID = uint32(response[8])<<24 | uint32(response[9])<<16 | uint32(response[10])<<8 | uint32(response[11])

	return nil
}

// SendCommand sends a command to Xiaomi device
func (x *XiaomiDevice) SendCommand(method string, params []interface{}) ([]interface{}, error) {
	req := XiaomiRequest{
		ID:     1,
		Method: method,
		Params: params,
	}

	reqJSON, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Build packet
	packet := x.buildPacket(reqJSON)

	// Send via UDP
	conn, err := net.DialTimeout("udp", fmt.Sprintf("%s:54321", x.IP), 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	_, err = conn.Write(packet)
	if err != nil {
		return nil, fmt.Errorf("failed to send command: %w", err)
	}

	// Read response
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("no response from device: %w", err)
	}

	// Parse response
	respJSON := buffer[32:n]
	var resp XiaomiResponse
	if err := json.Unmarshal(respJSON, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("device error: %v", resp.Error)
	}

	return resp.Result, nil
}

// buildPacket builds a Miio protocol packet
func (x *XiaomiDevice) buildPacket(data []byte) []byte {
	// Header
	packet := make([]byte, 32+len(data))
	packet[0] = 0x21
	packet[1] = 0x31

	// Length
	length := uint16(len(packet))
	packet[2] = byte(length >> 8)
	packet[3] = byte(length)

	// Unknown
	packet[4] = 0x00
	packet[5] = 0x00
	packet[6] = 0x00
	packet[7] = 0x00

	// Device ID
	packet[8] = byte(x.DeviceID >> 24)
	packet[9] = byte(x.DeviceID >> 16)
	packet[10] = byte(x.DeviceID >> 8)
	packet[11] = byte(x.DeviceID)

	// Timestamp
	now := time.Now().Unix()
	packet[12] = byte(now >> 24)
	packet[13] = byte(now >> 16)
	packet[14] = byte(now >> 8)
	packet[15] = byte(now)

	// Checksum placeholder
	copy(packet[16:32], x.Token)

	// Data
	copy(packet[32:], data)

	// Calculate checksum
	checksum := md5.Sum(packet)
	copy(packet[16:32], checksum[:])

	return packet
}

// VacuumRobot represents a Xiaomi vacuum robot
type VacuumRobot struct {
	device *XiaomiDevice
}

// NewVacuumRobot creates a new vacuum robot controller
func NewVacuumRobot(ip, token string) (*VacuumRobot, error) {
	device, err := NewXiaomiDevice(ip, token)
	if err != nil {
		return nil, err
	}

	return &VacuumRobot{
		device: device,
	}, nil
}

// Start starts cleaning
func (v *VacuumRobot) Start() error {
	_, err := v.device.SendCommand("app_start", nil)
	return err
}

// Stop stops cleaning
func (v *VacuumRobot) Stop() error {
	_, err := v.device.SendCommand("app_stop", nil)
	return err
}

// Pause pauses cleaning
func (v *VacuumRobot) Pause() error {
	_, err := v.device.SendCommand("app_pause", nil)
	return err
}

// Home sends robot to charging dock
func (v *VacuumRobot) Home() error {
	_, err := v.device.SendCommand("app_charge", nil)
	return err
}

// Spot starts spot cleaning
func (v *VacuumRobot) Spot() error {
	_, err := v.device.SendCommand("app_spot", nil)
	return err
}

// SetFanSpeed sets fan speed (silent=38, standard=60, medium=77, turbo=90)
func (v *VacuumRobot) SetFanSpeed(speed int) error {
	_, err := v.device.SendCommand("set_custom_mode", []interface{}{speed})
	return err
}

// GetStatus gets vacuum status
func (v *VacuumRobot) GetStatus() (map[string]interface{}, error) {
	result, err := v.device.SendCommand("get_status", nil)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("empty status response")
	}

	status, ok := result[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid status format")
	}

	return status, nil
}

// FindMe makes the robot emit a sound
func (v *VacuumRobot) FindMe() error {
	_, err := v.device.SendCommand("find_me", nil)
	return err
}

// XiaomiLight represents a Xiaomi smart light
type XiaomiLight struct {
	device *XiaomiDevice
}

// NewXiaomiLight creates a new Xiaomi light controller
func NewXiaomiLight(ip, token string) (*XiaomiLight, error) {
	device, err := NewXiaomiDevice(ip, token)
	if err != nil {
		return nil, err
	}

	return &XiaomiLight{
		device: device,
	}, nil
}

// TurnOn turns on the light
func (l *XiaomiLight) TurnOn() error {
	_, err := l.device.SendCommand("set_power", []interface{}{"on"})
	return err
}

// TurnOff turns off the light
func (l *XiaomiLight) TurnOff() error {
	_, err := l.device.SendCommand("set_power", []interface{}{"off"})
	return err
}

// SetBrightness sets brightness (1-100)
func (l *XiaomiLight) SetBrightness(brightness int) error {
	if brightness < 1 || brightness > 100 {
		return fmt.Errorf("brightness must be between 1 and 100")
	}
	_, err := l.device.SendCommand("set_bright", []interface{}{brightness})
	return err
}

// SetColorTemp sets color temperature (1700-6500K)
func (l *XiaomiLight) SetColorTemp(temp int) error {
	if temp < 1700 || temp > 6500 {
		return fmt.Errorf("color temperature must be between 1700 and 6500")
	}
	_, err := l.device.SendCommand("set_ct_abx", []interface{}{temp, "smooth", 500})
	return err
}

// SetRGB sets RGB color
func (l *XiaomiLight) SetRGB(r, g, b int) error {
	// Convert RGB to decimal
	rgb := (r << 16) | (g << 8) | b
	_, err := l.device.SendCommand("set_rgb", []interface{}{rgb})
	return err
}

// XiaomiAirPurifier represents a Xiaomi air purifier
type XiaomiAirPurifier struct {
	device *XiaomiDevice
}

// NewXiaomiAirPurifier creates a new air purifier controller
func NewXiaomiAirPurifier(ip, token string) (*XiaomiAirPurifier, error) {
	device, err := NewXiaomiDevice(ip, token)
	if err != nil {
		return nil, err
	}

	return &XiaomiAirPurifier{
		device: device,
	}, nil
}

// TurnOn turns on the air purifier
func (a *XiaomiAirPurifier) TurnOn() error {
	_, err := a.device.SendCommand("set_power", []interface{}{"on"})
	return err
}

// TurnOff turns off the air purifier
func (a *XiaomiAirPurifier) TurnOff() error {
	_, err := a.device.SendCommand("set_power", []interface{}{"off"})
	return err
}

// SetMode sets operation mode (auto, silent, favorite)
func (a *XiaomiAirPurifier) SetMode(mode string) error {
	_, err := a.device.SendCommand("set_mode", []interface{}{mode})
	return err
}

// SetFavoriteLevel sets favorite level (0-14)
func (a *XiaomiAirPurifier) SetFavoriteLevel(level int) error {
	if level < 0 || level > 14 {
		return fmt.Errorf("level must be between 0 and 14")
	}
	_, err := a.device.SendCommand("set_level_favorite", []interface{}{level})
	return err
}

// HTTPDevice represents a generic HTTP-controlled device
type HTTPDevice struct {
	BaseURL string
	Headers map[string]string
	client  *http.Client
}

// NewHTTPDevice creates a new HTTP device
func NewHTTPDevice(baseURL string, headers map[string]string) *HTTPDevice {
	return &HTTPDevice{
		BaseURL: baseURL,
		Headers: headers,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendRequest sends an HTTP request
func (h *HTTPDevice) SendRequest(method, path string, body []byte) ([]byte, error) {
	url := h.BaseURL + path

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for key, value := range h.Headers {
		req.Header.Set(key, value)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// Get sends a GET request
func (h *HTTPDevice) Get(path string) ([]byte, error) {
	return h.SendRequest("GET", path, nil)
}

// Post sends a POST request
func (h *HTTPDevice) Post(path string, body []byte) ([]byte, error) {
	return h.SendRequest("POST", path, body)
}

// Put sends a PUT request
func (h *HTTPDevice) Put(path string, body []byte) ([]byte, error) {
	return h.SendRequest("PUT", path, body)
}
