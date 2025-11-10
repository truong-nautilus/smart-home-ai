package edgetts

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/truong-nautilus/smart-home-ai/internal/domain"
)

// EdgeTTSSynthesizer triển khai domain.SpeechSynthesizer sử dụng Microsoft Edge TTS
// Giọng nói neural chất lượng cao, giống người thật nhất
type EdgeTTSSynthesizer struct {
	Voice      string
	EdgeTTSBin string
}

// New tạo EdgeTTS synthesizer mới
func New(voice, edgeTTSBin string) *EdgeTTSSynthesizer {
	if voice == "" {
		voice = "vi-VN-NamMinhNeural"
	}
	if edgeTTSBin == "" {
		homeDir, _ := os.UserHomeDir()
		edgeTTSBin = filepath.Join(homeDir, "Library", "Python", "3.9", "bin", "edge-tts")
	}
	return &EdgeTTSSynthesizer{
		Voice:      voice,
		EdgeTTSBin: edgeTTSBin,
	}
}

// Synthesize chuyển văn bản thành giọng nói và phát luôn
func (e *EdgeTTSSynthesizer) Synthesize(ctx context.Context, text string, outputPath string) (*domain.SpeechOutput, error) {
	tempFile := "/tmp/edge_tts_output.mp3"

	cmd := exec.CommandContext(ctx, e.EdgeTTSBin,
		"--voice", e.Voice,
		"--text", text,
		"--write-media", tempFile,
	)

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("edge-tts: không thể tạo audio: %w", err)
	}

	playCmd := exec.CommandContext(ctx, "afplay", tempFile)
	if err := playCmd.Run(); err != nil {
		return nil, fmt.Errorf("edge-tts: không thể phát audio: %w", err)
	}

	os.Remove(tempFile)

	return &domain.SpeechOutput{
		Text:      text,
		AudioPath: "",
		Voice:     e.Voice,
		Timestamp: time.Now(),
	}, nil
}
