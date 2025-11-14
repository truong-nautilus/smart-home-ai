package audio

import (
	"encoding/binary"
	"fmt"
	"log"
	"sync"

	"github.com/gen2brain/malgo"
)

// AudioConfig holds the audio configuration
type AudioConfig struct {
	SampleRate uint32
	Channels   uint32
	Format     malgo.FormatType
	BufferSize uint32
}

// Recorder handles microphone input
type Recorder struct {
	ctx        *malgo.AllocatedContext
	device     *malgo.Device
	config     AudioConfig
	audioQueue chan []byte
	stopChan   chan struct{}
	isRunning  bool
	mu         sync.Mutex
}

// NewRecorder creates a new audio recorder
func NewRecorder(sampleRate, channels, bufferSize uint32) (*Recorder, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize malgo context: %w", err)
	}

	recorder := &Recorder{
		ctx: ctx,
		config: AudioConfig{
			SampleRate: sampleRate,
			Channels:   channels,
			Format:     malgo.FormatS16, // 16-bit PCM
			BufferSize: bufferSize,
		},
		audioQueue: make(chan []byte, 100),
		stopChan:   make(chan struct{}),
		isRunning:  false,
	}

	return recorder, nil
}

// Start begins recording audio from the microphone
func (r *Recorder) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.isRunning {
		return fmt.Errorf("recorder is already running")
	}

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = r.config.Format
	deviceConfig.Capture.Channels = r.config.Channels
	deviceConfig.SampleRate = r.config.SampleRate
	deviceConfig.Alsa.NoMMap = 1

	// Callback for audio data
	onRecvFrames := func(pOutputSample, pInputSamples []byte, framecount uint32) {
		if len(pInputSamples) > 0 {
			// Create a copy of the buffer
			data := make([]byte, len(pInputSamples))
			copy(data, pInputSamples)

			// Send to queue (non-blocking)
			select {
			case r.audioQueue <- data:
			default:
				// Queue is full, skip this frame
				log.Println("Warning: Audio queue is full, dropping frame")
			}
		}
	}

	device, err := malgo.InitDevice(r.ctx.Context, deviceConfig, malgo.DeviceCallbacks{
		Data: onRecvFrames,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize capture device: %w", err)
	}

	r.device = device

	if err := r.device.Start(); err != nil {
		return fmt.Errorf("failed to start device: %w", err)
	}

	r.isRunning = true
	log.Println("Audio recorder started successfully")
	log.Printf("Sample Rate: %d Hz, Channels: %d, Format: 16-bit PCM", r.config.SampleRate, r.config.Channels)

	return nil
}

// Stop stops the audio recording
func (r *Recorder) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.isRunning {
		return nil
	}

	close(r.stopChan)

	if r.device != nil {
		r.device.Uninit()
	}

	r.isRunning = false
	log.Println("Audio recorder stopped")

	return nil
}

// GetAudioChannel returns the channel for receiving audio data
func (r *Recorder) GetAudioChannel() <-chan []byte {
	return r.audioQueue
}

// Close closes the recorder and releases resources
func (r *Recorder) Close() error {
	if err := r.Stop(); err != nil {
		return err
	}

	if r.ctx != nil {
		_ = r.ctx.Uninit()
		r.ctx.Free()
	}

	close(r.audioQueue)

	return nil
}

// ConvertToPCM16 converts raw audio bytes to PCM16 format
func ConvertToPCM16(data []byte) []int16 {
	samples := make([]int16, len(data)/2)
	for i := 0; i < len(samples); i++ {
		samples[i] = int16(binary.LittleEndian.Uint16(data[i*2 : i*2+2]))
	}
	return samples
}

// IsRunning returns whether the recorder is currently running
func (r *Recorder) IsRunning() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.isRunning
}
