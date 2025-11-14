package audio

import "testing"

func TestAudioConfig(t *testing.T) {
	config := AudioConfig{
		SampleRate: 16000,
		Channels:   1,
		Format:     2, // S16
		BufferSize: 3200,
	}

	if config.SampleRate != 16000 {
		t.Errorf("Expected sample rate 16000, got %d", config.SampleRate)
	}

	if config.Channels != 1 {
		t.Errorf("Expected 1 channel, got %d", config.Channels)
	}
}

func TestConvertToPCM16(t *testing.T) {
	// Test data: 4 bytes = 2 samples in 16-bit PCM
	data := []byte{0x00, 0x01, 0xFF, 0x7F}
	samples := ConvertToPCM16(data)

	if len(samples) != 2 {
		t.Errorf("Expected 2 samples, got %d", len(samples))
	}

	expectedFirst := int16(0x0100)
	if samples[0] != expectedFirst {
		t.Errorf("Expected first sample %d, got %d", expectedFirst, samples[0])
	}
}
