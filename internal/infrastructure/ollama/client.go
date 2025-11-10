package ollama

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/truong-nautilus/smart-home-ai/internal/domain"
)

// OllamaClient triển khai domain.AIAssistant sử dụng Ollama CLI
type OllamaClient struct {
	Model string
}

// New tạo Ollama client mới
func New(model string) *OllamaClient {
	if model == "" {
		model = "phi-3-mini"
	}
	return &OllamaClient{Model: model}
}

// AnalyzeMultimodal gọi Ollama với text + ảnh (phi3:mini không hỗ trợ multimodal, chỉ xử lý text)
func (o *OllamaClient) AnalyzeMultimodal(ctx context.Context, text string, imagePath string) (*domain.AIResponse, error) {
	// Note: phi3:mini không hỗ trợ image input, chỉ xử lý text
	// Để sử dụng multimodal, cần model như llava hoặc bakllava
	
	// Tạo prompt đơn giản chỉ với text
	prompt := fmt.Sprintf("Câu hỏi của người dùng: %s\n\nHãy trả lời ngắn gọn và hữu ích.", text)

	// Sử dụng "ollama run" thay vì "ollama generate"
	cmd := exec.CommandContext(ctx, "ollama", "run", o.Model, prompt)
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("ollama: stdout pipe: %w", err)
	}
	
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("ollama: stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("ollama: không khởi động được: %w", err)
	}

	// Đọc stderr để debug
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			// Bỏ qua stderr, không log
		}
	}()

	scanner := bufio.NewScanner(out)
	var sb strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ollama: đọc output thất bại: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("ollama: command thất bại: %w", err)
	}

	response := strings.TrimSpace(sb.String())
	if response == "" {
		return nil, fmt.Errorf("ollama: không nhận được phản hồi")
	}

	return &domain.AIResponse{
		Text:  response,
		Model: o.Model,
	}, nil
}
