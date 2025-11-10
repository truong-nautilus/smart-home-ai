package wav2vec2

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/truong-nautilus/smart-home-ai/internal/domain"
)

// Wav2Vec2Recognizer triển khai SpeechRecognizer sử dụng Wav2Vec2-Base-Vietnamese
type Wav2Vec2Recognizer struct {
	scriptPath string
}

// NewWav2Vec2Recognizer tạo recognizer mới
func NewWav2Vec2Recognizer(scriptPath string) *Wav2Vec2Recognizer {
	return &Wav2Vec2Recognizer{
		scriptPath: scriptPath,
	}
}

// Transcribe chuyển đổi file audio thành text sử dụng Wav2Vec2
func (r *Wav2Vec2Recognizer) Transcribe(ctx context.Context, audioPath string) (*domain.Transcription, error) {
	// Chuyển sang absolute path
	absPath, err := filepath.Abs(audioPath)
	if err != nil {
		return nil, fmt.Errorf("không thể resolve audio path: %w", err)
	}

	// Gọi Python script với Python interpreter từ virtual environment
	pythonPath := "/Users/phamthetruong/phowhisper-env/bin/python3"
	cmd := exec.CommandContext(ctx, pythonPath, r.scriptPath, absPath)

	// Chỉ lấy stdout, stderr để riêng (tránh logging lẫn với kết quả)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Wav2Vec2 error: %w", err)
	}

	text := strings.TrimSpace(string(output))
	if text == "" {
		return nil, fmt.Errorf("Wav2Vec2 không nhận diện được văn bản")
	}

	fmt.Printf("[🔍 Wav2Vec2 output: \"%s\"]\n", text)

	return &domain.Transcription{
		Text:     text,
		Language: "vi", // Tiếng Việt
	}, nil
}
