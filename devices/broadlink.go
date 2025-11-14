package devices

import (
	"encoding/hex"
	"fmt"
	"net"
	"time"
)

// BroadlinkDevice represents a Broadlink IR/RF device
type BroadlinkDevice struct {
	IP      string
	Port    int
	MAC     net.HardwareAddr
	DevType int
	Count   int
	Key     []byte
	IV      []byte
	ID      []byte
}

// NewBroadlinkDevice creates a new Broadlink device
func NewBroadlinkDevice(ip string, port int) *BroadlinkDevice {
	return &BroadlinkDevice{
		IP:    ip,
		Port:  port,
		Count: 0,
		Key:   []byte{0x09, 0x76, 0x28, 0x34, 0x3f, 0xe9, 0x9e, 0x23, 0x76, 0x5c, 0x15, 0x13, 0xac, 0xcf, 0x8b, 0x02},
		IV:    []byte{0x56, 0x2e, 0x17, 0x99, 0x6d, 0x09, 0x3d, 0x28, 0xdd, 0xb3, 0xba, 0x69, 0x5a, 0x2e, 0x6f, 0x58},
		ID:    []byte{0, 0, 0, 0},
	}
}

// Discover discovers Broadlink devices on the network
func (b *BroadlinkDevice) Discover(timeout time.Duration) error {
	// Create UDP socket
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 0})
	if err != nil {
		return fmt.Errorf("failed to create UDP socket: %w", err)
	}
	defer conn.Close()

	// Set timeout
	conn.SetReadDeadline(time.Now().Add(timeout))

	// Get local IP
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	localIP := localAddr.IP.To4()

	// Build discovery packet
	packet := make([]byte, 0x30)

	// Header
	packet[0x00] = 0x5a
	packet[0x01] = 0xa5
	packet[0x02] = 0xaa
	packet[0x03] = 0x55
	packet[0x04] = 0x5a
	packet[0x05] = 0xa5
	packet[0x06] = 0xaa
	packet[0x07] = 0x55

	// Timestamp
	now := time.Now()
	packet[0x08] = byte(now.Year() & 0xff)
	packet[0x09] = byte(now.Year() >> 8)
	packet[0x0a] = byte(now.Minute())
	packet[0x0b] = byte(now.Hour())
	packet[0x0c] = byte(now.Year() % 100)
	packet[0x0d] = byte(now.Weekday())
	packet[0x0e] = byte(now.Day())
	packet[0x0f] = byte(now.Month())

	// Local IP
	copy(packet[0x18:0x1c], localIP)

	// Local port
	packet[0x1c] = byte(localAddr.Port & 0xff)
	packet[0x1d] = byte(localAddr.Port >> 8)

	// Checksum
	checksum := 0xbeaf
	for i := 0; i < len(packet); i++ {
		checksum += int(packet[i])
	}
	packet[0x20] = byte(checksum & 0xff)
	packet[0x21] = byte(checksum >> 8)

	// Broadcast discovery packet
	broadcastAddr := &net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: 80,
	}

	_, err = conn.WriteToUDP(packet, broadcastAddr)
	if err != nil {
		return fmt.Errorf("failed to send discovery packet: %w", err)
	}

	// Wait for response
	buffer := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return fmt.Errorf("no response from Broadlink device: %w", err)
	}

	if n < 0x30 {
		return fmt.Errorf("invalid response length")
	}

	// Parse response
	response := buffer[:n]
	b.IP = remoteAddr.IP.String()
	b.Port = int(response[0x1c]) | int(response[0x1d])<<8
	b.MAC = net.HardwareAddr(response[0x3a:0x40])
	b.DevType = int(response[0x34]) | int(response[0x35])<<8

	return nil
}

// Auth authenticates with the Broadlink device
func (b *BroadlinkDevice) Auth() error {
	payload := make([]byte, 0x50)
	payload[0x04] = 0x31
	payload[0x05] = 0x31
	payload[0x06] = 0x31
	payload[0x07] = 0x31
	payload[0x08] = 0x31
	payload[0x09] = 0x31
	payload[0x0a] = 0x31
	payload[0x0b] = 0x31
	payload[0x0c] = 0x31
	payload[0x0d] = 0x31
	payload[0x0e] = 0x31
	payload[0x0f] = 0x31
	payload[0x10] = 0x31
	payload[0x11] = 0x31
	payload[0x12] = 0x31
	payload[0x1e] = 0x01
	payload[0x2d] = 0x01
	payload[0x30] = byte('T')
	payload[0x31] = byte('e')
	payload[0x32] = byte('s')
	payload[0x33] = byte('t')

	response, err := b.sendPacket(0x65, payload)
	if err != nil {
		return err
	}

	if len(response) < 0x38 {
		return fmt.Errorf("invalid auth response")
	}

	// Extract device ID and key
	b.ID = response[0x00:0x04]
	b.Key = response[0x04:0x14]

	return nil
}

// SendIRCommand sends an IR command
func (b *BroadlinkDevice) SendIRCommand(data string) error {
	// Decode hex string to bytes
	irData, err := hex.DecodeString(data)
	if err != nil {
		return fmt.Errorf("invalid IR data: %w", err)
	}

	// Build command payload
	payload := make([]byte, 4+len(irData))
	payload[0] = 0x02 // IR command
	copy(payload[4:], irData)

	_, err = b.sendPacket(0x6a, payload)
	return err
}

// LearnIRCommand puts device in learning mode to capture IR signal
func (b *BroadlinkDevice) LearnIRCommand(timeout time.Duration) (string, error) {
	// Enter learning mode
	payload := []byte{0x03}
	_, err := b.sendPacket(0x6a, payload)
	if err != nil {
		return "", err
	}

	// Wait for IR signal
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		time.Sleep(1 * time.Second)

		// Check if learning is complete
		response, err := b.sendPacket(0x6a, []byte{0x04})
		if err != nil {
			continue
		}

		if len(response) > 4 && response[0] != 0 {
			// Learning complete, extract IR data
			irData := response[4:]
			return hex.EncodeToString(irData), nil
		}
	}

	return "", fmt.Errorf("learning timeout")
}

// sendPacket sends a packet to Broadlink device
func (b *BroadlinkDevice) sendPacket(command byte, payload []byte) ([]byte, error) {
	b.Count = (b.Count + 1) & 0xffff

	// Build packet
	packet := make([]byte, 0x38)

	// Header
	packet[0x00] = 0x5a
	packet[0x01] = 0xa5
	packet[0x02] = 0xaa
	packet[0x03] = 0x55
	packet[0x04] = 0x5a
	packet[0x05] = 0xa5
	packet[0x06] = 0xaa
	packet[0x07] = 0x55

	// Device ID
	copy(packet[0x08:0x0c], b.ID)

	// Packet count
	packet[0x28] = byte(b.Count & 0xff)
	packet[0x29] = byte(b.Count >> 8)

	// Command
	packet[0x26] = command

	// Checksum
	checksum := 0xbeaf
	for _, v := range payload {
		checksum += int(v)
	}
	packet[0x34] = byte(checksum & 0xff)
	packet[0x35] = byte(checksum >> 8)

	// Append encrypted payload
	// Note: In production, implement proper AES encryption
	packet = append(packet, payload...)

	// Send packet
	conn, err := net.DialTimeout("udp", fmt.Sprintf("%s:%d", b.IP, b.Port), 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	_, err = conn.Write(packet)
	if err != nil {
		return nil, fmt.Errorf("failed to send packet: %w", err)
	}

	// Read response
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	buffer := make([]byte, 2048)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if n < 0x38 {
		return nil, fmt.Errorf("invalid response length")
	}

	// Extract payload
	response := buffer[0x38:n]

	return response, nil
}

// CheckData checks if device has data to read
func (b *BroadlinkDevice) CheckData() ([]byte, error) {
	return b.sendPacket(0x6a, []byte{0x04})
}

// String returns string representation of device
func (b *BroadlinkDevice) String() string {
	return fmt.Sprintf("Broadlink Device (IP: %s, MAC: %s, Type: 0x%04x)", b.IP, b.MAC, b.DevType)
}

// Predefined IR commands for common devices
var (
	// AC Commands (example for common brands)
	ACCommands = map[string]map[string]string{
		"power_on": {
			"daikin":    "260050000001...",
			"panasonic": "260050000001...",
		},
		"power_off": {
			"daikin":    "260050000001...",
			"panasonic": "260050000001...",
		},
	}

	// TV Commands (example)
	TVCommands = map[string]string{
		"power":    "260050000001...",
		"vol_up":   "260050000001...",
		"vol_down": "260050000001...",
	}
)

// GetPredefinedCommand gets a predefined IR command
func GetPredefinedCommand(deviceType, brand, command string) (string, error) {
	switch deviceType {
	case "ac":
		if brandCommands, ok := ACCommands[command]; ok {
			if cmd, ok := brandCommands[brand]; ok {
				return cmd, nil
			}
		}
	case "tv":
		if cmd, ok := TVCommands[command]; ok {
			return cmd, nil
		}
	}

	return "", fmt.Errorf("command not found")
}
