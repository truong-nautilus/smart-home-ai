package devices

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TapoConfig holds Tapo device configuration
type TapoConfig struct {
	Email    string
	Password string
}

// TapoDevice represents a Tapo smart device
type TapoDevice struct {
	IP     string
	Model  string
	Token  string
	Cookie string
	client *http.Client
	config TapoConfig
}

// TapoRequest represents a request to Tapo device
type TapoRequest struct {
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params,omitempty"`
}

// TapoResponse represents a response from Tapo device
type TapoResponse struct {
	ErrorCode int                    `json:"error_code"`
	Result    map[string]interface{} `json:"result,omitempty"`
}

// NewTapoDevice creates a new Tapo device controller
func NewTapoDevice(ip, model string, config TapoConfig) *TapoDevice {
	return &TapoDevice{
		IP:     ip,
		Model:  model,
		config: config,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Handshake performs initial handshake with Tapo device
func (t *TapoDevice) Handshake() error {
	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return fmt.Errorf("failed to generate RSA key: %w", err)
	}

	// Encode public key
	publicKeyPEM := fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----",
		base64.StdEncoding.EncodeToString(privateKey.PublicKey.N.Bytes()))

	// Send handshake request
	handshakeReq := map[string]interface{}{
		"method": "handshake",
		"params": map[string]interface{}{
			"key": publicKeyPEM,
		},
	}

	resp, err := t.sendRequest("/app", handshakeReq, false)
	if err != nil {
		return fmt.Errorf("handshake failed: %w", err)
	}

	if resp.ErrorCode != 0 {
		return fmt.Errorf("handshake error code: %d", resp.ErrorCode)
	}

	// Extract token and cookie
	if key, ok := resp.Result["key"].(string); ok {
		t.Token = key
	}

	return nil
}

// Login logs in to the Tapo device
func (t *TapoDevice) Login() error {
	// Create credentials hash
	credentials := map[string]interface{}{
		"username": base64.StdEncoding.EncodeToString([]byte(t.config.Email)),
		"password": base64.StdEncoding.EncodeToString([]byte(t.config.Password)),
	}

	loginReq := TapoRequest{
		Method: "login_device",
		Params: credentials,
	}

	resp, err := t.sendSecureRequest(loginReq)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	if resp.ErrorCode != 0 {
		return fmt.Errorf("login error code: %d", resp.ErrorCode)
	}

	if token, ok := resp.Result["token"].(string); ok {
		t.Token = token
	}

	return nil
}

// TurnOn turns on the device
func (t *TapoDevice) TurnOn() error {
	req := TapoRequest{
		Method: "set_device_info",
		Params: map[string]interface{}{
			"device_on": true,
		},
	}

	resp, err := t.sendSecureRequest(req)
	if err != nil {
		return err
	}

	if resp.ErrorCode != 0 {
		return fmt.Errorf("error code: %d", resp.ErrorCode)
	}

	return nil
}

// TurnOff turns off the device
func (t *TapoDevice) TurnOff() error {
	req := TapoRequest{
		Method: "set_device_info",
		Params: map[string]interface{}{
			"device_on": false,
		},
	}

	resp, err := t.sendSecureRequest(req)
	if err != nil {
		return err
	}

	if resp.ErrorCode != 0 {
		return fmt.Errorf("error code: %d", resp.ErrorCode)
	}

	return nil
}

// SetBrightness sets the brightness (1-100) for L530
func (t *TapoDevice) SetBrightness(brightness int) error {
	if brightness < 1 || brightness > 100 {
		return fmt.Errorf("brightness must be between 1 and 100")
	}

	req := TapoRequest{
		Method: "set_device_info",
		Params: map[string]interface{}{
			"device_on":  true,
			"brightness": brightness,
		},
	}

	resp, err := t.sendSecureRequest(req)
	if err != nil {
		return err
	}

	if resp.ErrorCode != 0 {
		return fmt.Errorf("error code: %d", resp.ErrorCode)
	}

	return nil
}

// SetColor sets the color (hue, saturation) for L530
func (t *TapoDevice) SetColor(hue, saturation int) error {
	if hue < 0 || hue > 360 {
		return fmt.Errorf("hue must be between 0 and 360")
	}
	if saturation < 0 || saturation > 100 {
		return fmt.Errorf("saturation must be between 0 and 100")
	}

	req := TapoRequest{
		Method: "set_device_info",
		Params: map[string]interface{}{
			"device_on":  true,
			"hue":        hue,
			"saturation": saturation,
		},
	}

	resp, err := t.sendSecureRequest(req)
	if err != nil {
		return err
	}

	if resp.ErrorCode != 0 {
		return fmt.Errorf("error code: %d", resp.ErrorCode)
	}

	return nil
}

// SetColorTemp sets the color temperature (2500-6500K) for L530
func (t *TapoDevice) SetColorTemp(temp int) error {
	if temp < 2500 || temp > 6500 {
		return fmt.Errorf("color temperature must be between 2500 and 6500")
	}

	req := TapoRequest{
		Method: "set_device_info",
		Params: map[string]interface{}{
			"device_on":  true,
			"color_temp": temp,
		},
	}

	resp, err := t.sendSecureRequest(req)
	if err != nil {
		return err
	}

	if resp.ErrorCode != 0 {
		return fmt.Errorf("error code: %d", resp.ErrorCode)
	}

	return nil
}

// GetDeviceInfo gets device information
func (t *TapoDevice) GetDeviceInfo() (map[string]interface{}, error) {
	req := TapoRequest{
		Method: "get_device_info",
	}

	resp, err := t.sendSecureRequest(req)
	if err != nil {
		return nil, err
	}

	if resp.ErrorCode != 0 {
		return nil, fmt.Errorf("error code: %d", resp.ErrorCode)
	}

	return resp.Result, nil
}

// sendSecureRequest sends an encrypted request to Tapo device
func (t *TapoDevice) sendSecureRequest(req TapoRequest) (*TapoResponse, error) {
	// Marshal request
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// Encrypt data
	encrypted, err := t.encrypt(jsonData)
	if err != nil {
		return nil, err
	}

	secureReq := map[string]interface{}{
		"method": "securePassthrough",
		"params": map[string]interface{}{
			"request": encrypted,
		},
	}

	return t.sendRequest("/app?token="+t.Token, secureReq, true)
}

// sendRequest sends a request to Tapo device
func (t *TapoDevice) sendRequest(path string, data interface{}, withCookie bool) (*TapoResponse, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("http://%s%s", t.IP, path)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if withCookie && t.Cookie != "" {
		req.Header.Set("Cookie", t.Cookie)
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Save cookie
	if cookies := resp.Cookies(); len(cookies) > 0 {
		t.Cookie = cookies[0].String()
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tapoResp TapoResponse
	if err := json.Unmarshal(body, &tapoResp); err != nil {
		return nil, err
	}

	return &tapoResp, nil
}

// encrypt encrypts data using AES
func (t *TapoDevice) encrypt(data []byte) (string, error) {
	// Use SHA256 of token as key
	hash := sha256.Sum256([]byte(t.Token))
	key := hash[:16] // Use first 16 bytes for AES-128

	// Use SHA1 of token as IV
	hashIV := sha1.Sum([]byte(t.Token))
	iv := hashIV[:16]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// PKCS7 padding
	padding := aes.BlockSize - len(data)%aes.BlockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	data = append(data, padText...)

	ciphertext := make([]byte, len(data))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, data)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decrypt decrypts data using AES
func (t *TapoDevice) decrypt(encryptedData string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}

	// Use SHA256 of token as key
	hash := sha256.Sum256([]byte(t.Token))
	key := hash[:16]

	// Use SHA1 of token as IV
	hashIV := sha1.Sum([]byte(t.Token))
	iv := hashIV[:16]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of block size")
	}

	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	// Remove PKCS7 padding
	padding := int(plaintext[len(plaintext)-1])
	plaintext = plaintext[:len(plaintext)-padding]

	return plaintext, nil
}
