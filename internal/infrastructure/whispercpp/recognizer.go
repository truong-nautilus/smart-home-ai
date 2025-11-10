package whispercpp

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/truong-nautilus/smart-home-ai/internal/domain"
)

// WhisperCPPRecognizer triển khai domain.SpeechRecognizer sử dụng binary whisper.cpp
type WhisperCPPRecognizer struct {
	Binary string
	Model  string
}

// New tạo recognizer mới
func New(binary, model string) *WhisperCPPRecognizer {
	if binary == "" {
		binary = "main"
	}
	return &WhisperCPPRecognizer{Binary: binary, Model: model}
}

// Transcribe thực thi whisper.cpp và trả về domain.Transcription
func (w *WhisperCPPRecognizer) Transcribe(ctx context.Context, audioPath string) (*domain.Transcription, error) {
	args := []string{}
	if w.Model != "" {
		args = append(args, "-m", w.Model)
	}
	// Thêm language hint cho tiếng Việt
	args = append(args, "-l", "vi")
	args = append(args, "-f", audioPath)

	cmd := exec.CommandContext(ctx, w.Binary, args...)
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("whisper: stdout pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("whisper: failed to start: %w", err)
	}

	scanner := bufio.NewScanner(out)
	var sb strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		sb.WriteString(line)
		sb.WriteString("\n")
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("whisper: reading stdout: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("whisper: command failed: %w", err)
	}

	text := strings.TrimSpace(sb.String())
	if text == "" {
		return nil, fmt.Errorf("whisper: no transcription produced")
	}

	return &domain.Transcription{
		Text:      text,
		Language:  "unknown",
		Duration:  0,
		Timestamp: time.Now(),
	}, nil
}
