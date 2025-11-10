package phowhisper

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/truong-nautilus/smart-home-ai/internal/domain"
)

// PhoWhisperRecognizer triển khai SpeechRecognizer sử dụng PhoWhisper model
type PhoWhisperRecognizer struct {
	scriptPath string
}

// NewPhoWhisperRecognizer tạo recognizer mới
func NewPhoWhisperRecognizer(scriptPath string) *PhoWhisperRecognizer {
	return &PhoWhisperRecognizer{
		scriptPath: scriptPath,
	}
}

// Transcribe chuyển đổi file audio thành text sử dụng PhoWhisper
func (r *PhoWhisperRecognizer) Transcribe(ctx context.Context, audioPath string) (*domain.Transcription, error) {
	// Chuyển sang absolute path
	absPath, err := filepath.Abs(audioPath)
	if err != nil {
		return nil, fmt.Errorf("không thể resolve audio path: %w", err)
	}

	// Gọi Python script thông qua Python interpreter
	// Sử dụng python3 từ shebang của script
	cmd := exec.CommandContext(ctx, "/Users/phamthetruong/phowhisper-env/bin/python3", r.scriptPath, absPath)

	// Chỉ lấy stdout, bỏ qua stderr (logging)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("PhoWhisper error: %w", err)
	}

	text := strings.TrimSpace(string(output))
	if text == "" {
		return nil, fmt.Errorf("PhoWhisper không nhận diện được văn bản")
	}

	fmt.Printf("[🔍 PhoWhisper output: \"%s\"]\n", text)

	return &domain.Transcription{
		Text:     text,
		Language: "vi", // Tiếng Việt
	}, nil
}
