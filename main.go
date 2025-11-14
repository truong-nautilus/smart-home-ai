package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/truong-nautilus/smart-home-ai/audio"
	"github.com/truong-nautilus/smart-home-ai/claude"
	"github.com/truong-nautilus/smart-home-ai/core"
	"github.com/truong-nautilus/smart-home-ai/devices"
)

const (
	configFile = "config.json"
	envFile    = ".env"
)

func main() {
	log.Println("=== Jarvis AI Smart Home System ===")
	log.Println("Initializing...")

	// Load environment variables
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Load configuration
	config, err := core.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize security manager
	security := core.NewSecurityManager()
	log.Println("Security manager initialized")

	// Initialize command router
	router := core.NewCommandRouter(config)

	// Setup device configurations
	tapoConfig := devices.TapoConfig{
		Email:    os.Getenv("TAPO_USER"),
		Password: os.Getenv("TAPO_PASS"),
	}

	mqttConfig := devices.MQTTConfig{
		Host:     os.Getenv("MQTT_HOST"),
		Port:     1883,
		Username: os.Getenv("MQTT_USER"),
		Password: os.Getenv("MQTT_PASS"),
		ClientID: "jarvis-ai-" + time.Now().Format("20060102150405"),
	}

	// Initialize devices
	if err := router.Initialize(tapoConfig, mqttConfig); err != nil {
		log.Printf("Warning: Some devices failed to initialize: %v", err)
	}

	// Initialize Claude Realtime client
	claudeConfig := claude.ClaudeConfig{
		APIKey:       os.Getenv("CLAUDE_API_KEY"),
		Model:        config.Claude.Model,
		WebSocketURL: config.Claude.WebSocketURL,
		SystemPrompt: config.Claude.SystemPrompt,
		MaxTokens:    config.Claude.MaxTokens,
		Temperature:  config.Claude.Temperature,
	}

	claudeClient := claude.NewRealtimeClient(claudeConfig)
	log.Println("Connecting to Claude Realtime API...")

	if err := claudeClient.Connect(); err != nil {
		log.Fatalf("Failed to connect to Claude: %v", err)
	}
	defer claudeClient.Disconnect()

	log.Println("Connected to Claude Realtime API")

	// Initialize audio recorder
	log.Println("Initializing audio recorder...")
	recorder, err := audio.NewRecorder(
		uint32(config.Audio.SampleRate),
		uint32(config.Audio.Channels),
		uint32(config.Audio.BufferSize),
	)
	if err != nil {
		log.Fatalf("Failed to create audio recorder: %v", err)
	}
	defer recorder.Close()

	if err := recorder.Start(); err != nil {
		log.Fatalf("Failed to start audio recorder: %v", err)
	}
	defer recorder.Stop()

	log.Println("Audio recorder started")

	// Setup command processing goroutine
	go func() {
		commandChan := claudeClient.GetCommandChannel()
		for cmd := range commandChan {
			log.Printf("Received command: %+v", cmd)

			// Validate command
			if err := security.ValidateCommand(cmd); err != nil {
				log.Printf("Command validation failed: %v", err)
				security.LogCommand(cmd, false, err)
				continue
			}

			// Execute command
			if err := router.ExecuteCommand(cmd); err != nil {
				log.Printf("Command execution failed: %v", err)
				security.LogCommand(cmd, false, err)
			} else {
				log.Printf("Command executed successfully: %s on %s", cmd.Action, cmd.Device)
				security.LogCommand(cmd, true, nil)
			}
		}
	}()

	// Setup audio streaming goroutine
	go func() {
		audioChan := recorder.GetAudioChannel()
		lastCommit := time.Now()

		for audioData := range audioChan {
			// Send audio to Claude
			if err := claudeClient.SendAudio(audioData); err != nil {
				log.Printf("Error sending audio: %v", err)
				continue
			}

			// Commit audio buffer every 500ms to trigger processing
			if time.Since(lastCommit) > 500*time.Millisecond {
				if err := claudeClient.CommitAudioBuffer(); err != nil {
					log.Printf("Error committing audio buffer: %v", err)
				}
				lastCommit = time.Now()
			}
		}
	}()

	// Display startup message
	log.Println("\n============================================================")
	log.Println("Jarvis AI Smart Home is running!")
	log.Println("Speak into your microphone to control your smart home devices")
	log.Println("Example: 'Turn on the living room light'")
	log.Println("         'Set air conditioner to 26 degrees'")
	log.Println("         'Start the vacuum robot'")
	log.Println("Press Ctrl+C to stop")
	log.Println("============================================================")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down gracefully...")
	router.Close()
	log.Println("Goodbye!")
}
