package claude

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/truong-nautilus/smart-home-ai/core"
)

// ClaudeConfig holds Claude API configuration
type ClaudeConfig struct {
	APIKey       string
	Model        string
	WebSocketURL string
	SystemPrompt string
	MaxTokens    int
	Temperature  float64
}

// RealtimeClient handles communication with Claude Realtime API
type RealtimeClient struct {
	config         ClaudeConfig
	conn           *websocket.Conn
	audioInChan    chan []byte
	commandOutChan chan *core.Command
	isConnected    bool
	mu             sync.Mutex
	stopChan       chan struct{}
}

// Command is imported from core package

// RealtimeEvent represents events sent to/from Claude
type RealtimeEvent struct {
	Type  string                 `json:"type"`
	Event map[string]interface{} `json:"event,omitempty"`
}

// InputAudioBufferAppend event structure
type InputAudioBufferAppendEvent struct {
	Type  string `json:"type"`
	Audio string `json:"audio"` // base64 encoded PCM16
}

// NewRealtimeClient creates a new Claude Realtime client
func NewRealtimeClient(config ClaudeConfig) *RealtimeClient {
	return &RealtimeClient{
		config:         config,
		audioInChan:    make(chan []byte, 100),
		commandOutChan: make(chan *core.Command, 10),
		isConnected:    false,
		stopChan:       make(chan struct{}),
	}
}

// Connect establishes WebSocket connection to Claude API
func (c *RealtimeClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isConnected {
		return fmt.Errorf("already connected")
	}

	// Note: Claude API v1 uses standard HTTPS endpoints with streaming
	// For realtime, we'll use the messages streaming endpoint
	url := "wss://api.anthropic.com/v1/messages"

	headers := make(map[string][]string)
	headers["anthropic-version"] = []string{"2023-06-01"}
	headers["x-api-key"] = []string{c.config.APIKey}
	headers["content-type"] = []string{"application/json"}

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = 10 * time.Second

	conn, _, err := dialer.Dial(url, headers)
	if err != nil {
		return fmt.Errorf("failed to connect to Claude API: %w", err)
	}

	c.conn = conn
	c.isConnected = true

	log.Println("Connected to Claude Realtime API")

	// Send initial session configuration
	if err := c.sendSessionConfig(); err != nil {
		return fmt.Errorf("failed to send session config: %w", err)
	}

	// Start listening for responses
	go c.listenForResponses()

	// Start processing audio input
	go c.processAudioInput()

	return nil
}

// sendSessionConfig sends initial session configuration
func (c *RealtimeClient) sendSessionConfig() error {
	config := map[string]interface{}{
		"type": "session.update",
		"session": map[string]interface{}{
			"model":               c.config.Model,
			"modalities":          []string{"text", "audio"},
			"instructions":        c.config.SystemPrompt,
			"voice":               "alloy",
			"input_audio_format":  "pcm16",
			"output_audio_format": "pcm16",
			"input_audio_transcription": map[string]interface{}{
				"enabled": true,
				"model":   "whisper-1",
			},
			"turn_detection": map[string]interface{}{
				"type":                "server_vad",
				"threshold":           0.5,
				"prefix_padding_ms":   300,
				"silence_duration_ms": 500,
			},
			"temperature":                c.config.Temperature,
			"max_response_output_tokens": c.config.MaxTokens,
		},
	}

	return c.sendJSON(config)
}

// SendAudio sends audio data to Claude
func (c *RealtimeClient) SendAudio(audioData []byte) error {
	select {
	case c.audioInChan <- audioData:
		return nil
	default:
		return fmt.Errorf("audio input channel is full")
	}
}

// processAudioInput processes audio data and sends to Claude
func (c *RealtimeClient) processAudioInput() {
	for {
		select {
		case <-c.stopChan:
			return
		case audioData := <-c.audioInChan:
			if err := c.sendAudioBuffer(audioData); err != nil {
				log.Printf("Error sending audio buffer: %v", err)
			}
		}
	}
}

// sendAudioBuffer sends audio buffer to Claude
func (c *RealtimeClient) sendAudioBuffer(audioData []byte) error {
	// Encode audio to base64
	encoded := base64.StdEncoding.EncodeToString(audioData)

	event := InputAudioBufferAppendEvent{
		Type:  "input_audio_buffer.append",
		Audio: encoded,
	}

	return c.sendJSON(event)
}

// CommitAudioBuffer commits the audio buffer to trigger processing
func (c *RealtimeClient) CommitAudioBuffer() error {
	event := map[string]interface{}{
		"type": "input_audio_buffer.commit",
	}
	return c.sendJSON(event)
}

// listenForResponses listens for responses from Claude
func (c *RealtimeClient) listenForResponses() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in listenForResponses: %v", r)
		}
	}()

	for {
		select {
		case <-c.stopChan:
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v", err)
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					return
				}
				continue
			}

			c.handleResponse(message)
		}
	}
}

// handleResponse processes responses from Claude
func (c *RealtimeClient) handleResponse(message []byte) {
	var event map[string]interface{}
	if err := json.Unmarshal(message, &event); err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return
	}

	eventType, ok := event["type"].(string)
	if !ok {
		log.Println("Response missing type field")
		return
	}

	log.Printf("Received event type: %s", eventType)

	switch eventType {
	case "session.created":
		log.Println("Session created successfully")

	case "session.updated":
		log.Println("Session updated successfully")

	case "input_audio_buffer.speech_started":
		log.Println("Speech detected")

	case "input_audio_buffer.speech_stopped":
		log.Println("Speech ended")

	case "conversation.item.created":
		log.Println("Conversation item created")

	case "response.done":
		log.Println("Response completed")
		c.handleResponseDone(event)

	case "response.text.delta":
		c.handleTextDelta(event)

	case "response.audio.delta":
		c.handleAudioDelta(event)

	case "response.function_call_arguments.done":
		c.handleFunctionCall(event)

	case "error":
		c.handleError(event)
	}
}

// handleResponseDone handles completed responses
func (c *RealtimeClient) handleResponseDone(event map[string]interface{}) {
	response, ok := event["response"].(map[string]interface{})
	if !ok {
		return
	}

	output, ok := response["output"].([]interface{})
	if !ok || len(output) == 0 {
		return
	}

	for _, item := range output {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		if itemType, ok := itemMap["type"].(string); ok && itemType == "message" {
			if content, ok := itemMap["content"].([]interface{}); ok {
				for _, contentItem := range content {
					contentMap, ok := contentItem.(map[string]interface{})
					if !ok {
						continue
					}

					if text, ok := contentMap["text"].(string); ok {
						c.parseCommand(text)
					}
				}
			}
		}
	}
}

// handleTextDelta handles streaming text responses
func (c *RealtimeClient) handleTextDelta(event map[string]interface{}) {
	if delta, ok := event["delta"].(map[string]interface{}); ok {
		if text, ok := delta["text"].(string); ok {
			log.Printf("Claude: %s", text)
		}
	}
}

// handleAudioDelta handles streaming audio responses
func (c *RealtimeClient) handleAudioDelta(event map[string]interface{}) {
	if delta, ok := event["delta"].(map[string]interface{}); ok {
		if audioBase64, ok := delta["audio"].(string); ok {
			// Decode and play audio (implementation depends on audio output system)
			audioData, err := base64.StdEncoding.DecodeString(audioBase64)
			if err != nil {
				log.Printf("Error decoding audio: %v", err)
				return
			}
			log.Printf("Received audio chunk: %d bytes", len(audioData))
			// TODO: Play audio through speakers
		}
	}
}

// handleFunctionCall handles function call responses
func (c *RealtimeClient) handleFunctionCall(event map[string]interface{}) {
	if args, ok := event["arguments"].(string); ok {
		c.parseCommand(args)
	}
}

// handleError handles error events
func (c *RealtimeClient) handleError(event map[string]interface{}) {
	if errData, ok := event["error"].(map[string]interface{}); ok {
		log.Printf("Claude API Error: %v", errData)
	}
}

// parseCommand parses command from Claude's text response
func (c *RealtimeClient) parseCommand(text string) {
	// Try to parse as JSON command
	var command core.Command
	if err := json.Unmarshal([]byte(text), &command); err != nil {
		// Not a valid JSON command, might be regular text
		log.Printf("Not a command, regular text: %s", text)
		return
	}

	if command.Action == "" {
		return
	}

	log.Printf("Parsed command: action=%s, device=%s, value=%v", command.Action, command.Device, command.Value)

	select {
	case c.commandOutChan <- &command:
		log.Println("Command sent to execution queue")
	default:
		log.Println("Warning: Command output channel is full")
	}
}

// GetCommandChannel returns the channel for receiving commands
func (c *RealtimeClient) GetCommandChannel() <-chan *core.Command {
	return c.commandOutChan
}

// sendJSON sends a JSON message to Claude
func (c *RealtimeClient) sendJSON(data interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isConnected {
		return fmt.Errorf("not connected")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// Disconnect closes the WebSocket connection
func (c *RealtimeClient) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isConnected {
		return nil
	}

	close(c.stopChan)

	if c.conn != nil {
		err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Printf("Error sending close message: %v", err)
		}
		c.conn.Close()
	}

	c.isConnected = false
	close(c.audioInChan)
	close(c.commandOutChan)

	log.Println("Disconnected from Claude Realtime API")

	return nil
}

// IsConnected returns whether the client is connected
func (c *RealtimeClient) IsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.isConnected
}
