package mactts

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/truong-nautilus/smart-home-ai/internal/domain"
)

// MacTTSSynthesizer triển khai domain.SpeechSynthesizer sử dụng lệnh `say` của macOS
type MacTTSSynthesizer struct {
	Voice string // Tùy chọn giọng nói, ví dụ: "Alex", "Samantha"
}

// New tạo MacTTS synthesizer mới
func New(voice string) *MacTTSSynthesizer {
	return &MacTTSSynthesizer{Voice: voice}
}

// Synthesize chuyển văn bản thành giọng nói và phát trực tiếp (outputPath bị bỏ qua vì `say` phát luôn)
func (m *MacTTSSynthesizer) Synthesize(ctx context.Context, text string, outputPath string) (*domain.SpeechOutput, error) {
	args := []string{}
	if m.Voice != "" {
		args = append(args, "-v", m.Voice)
	}
	args = append(args, text)

	cmd := exec.CommandContext(ctx, "say", args...)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("mactts: không thực thi được lệnh say: %w", err)
	}

	// Không tạo file audio, chỉ phát trực tiếp
	return &domain.SpeechOutput{
		Text:      text,
		AudioPath: "", // Không có file audio
		Voice:     m.Voice,
		Timestamp: time.Now(),
	}, nil
}
