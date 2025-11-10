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
	// Tối ưu cho large-v3 model (mạnh nhất)
	args = append(args, "--temperature", "0.0")    // Giảm randomness, tăng consistency
	args = append(args, "--best-of", "5")          // Thử 5 beam và chọn tốt nhất
	args = append(args, "--beam-size", "5")        // Beam search width
	args = append(args, "--entropy-thold", "2.4")  // Entropy threshold cho quality
	args = append(args, "--logprob-thold", "-1.0") // Log probability threshold
	args = append(args, "--max-len", "0")          // Không giới hạn độ dài
	args = append(args, "--word-thold", "0.01")    // Threshold thấp cho word detection
	// Prompt tự nhiên cho giao tiếp
	args = append(args, "--prompt", "Đây là cuộc trò chuyện bằng tiếng Việt")
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

	rawText := strings.TrimSpace(sb.String())
	if rawText == "" {
		return nil, fmt.Errorf("whisper: no transcription produced")
	}

	// Extract text content (loại bỏ timestamp nếu có)
	text := extractTextContent(rawText)

	// Không sử dụng post-processing - tin tưởng vào medium model

	return &domain.Transcription{
		Text:      text,
		Language:  "unknown",
		Duration:  0,
		Timestamp: time.Now(),
	}, nil
}

// extractTextContent lấy nội dung text từ output của Whisper (loại bỏ timestamp)
func extractTextContent(raw string) string {
	// Format: "[00:00:00.000 --> 00:00:02.000]   Text content"
	// Tìm dấu "]" cuối cùng và lấy phần sau
	idx := strings.LastIndex(raw, "]")
	if idx == -1 {
		return raw // Không có timestamp, trả về nguyên bản
	}
	return strings.TrimSpace(raw[idx+1:])
}
