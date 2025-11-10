package ollama

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"regexp"
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

// AnalyzeMultimodal gọi Ollama với text (bỏ qua imagePath)
func (o *OllamaClient) AnalyzeMultimodal(ctx context.Context, text string, imagePath string) (*domain.AIResponse, error) {
	// imagePath bị bỏ qua - chỉ xử lý text thuần túy cho voice assistant

	// Tạo prompt với phong cách vui vẻ, hài hước, tự nhiên
	prompt := fmt.Sprintf(`Bạn là trợ lý AI thân thiện, vui vẻ, hài hước. Hãy trả lời câu hỏi sau:

"%s"

YÊU CẦU:
- Trả lời bằng tiếng Việt
- Phong cách: thân thiện, vui vẻ, hài hước, dễ thương
- Trả lời ngắn gọn (1-3 câu), tự nhiên như giọng nói
- KHÔNG dùng emoji, emoticon, ký tự đặc biệt
- Có thể đùa cợt nhẹ nhàng, tạo cảm giác gần gũi
- Nếu là yêu cầu điều khiển (bật/tắt đèn, quạt...), trả lời vui vẻ và xác nhận`, text)

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

	// Loại bỏ emoji và ký tự đặc biệt
	cleanResponse := removeEmoji(response)

	return &domain.AIResponse{
		Text:  cleanResponse,
		Model: o.Model,
	}, nil
}

// removeEmoji loại bỏ emoji và ký tự đặc biệt từ text
func removeEmoji(text string) string {
	// Regex để loại bỏ emoji và các ký tự Unicode đặc biệt
	emojiPattern := regexp.MustCompile(`[\x{1F600}-\x{1F64F}]|[\x{1F300}-\x{1F5FF}]|[\x{1F680}-\x{1F6FF}]|[\x{1F1E0}-\x{1F1FF}]|[\x{2600}-\x{26FF}]|[\x{2700}-\x{27BF}]|[\x{1F900}-\x{1F9FF}]|[\x{1FA70}-\x{1FAFF}]|[\x{FE00}-\x{FE0F}]|[\x{200D}]`)
	cleaned := emojiPattern.ReplaceAllString(text, "")

	// Loại bỏ markdown bold (**text**)
	cleaned = strings.ReplaceAll(cleaned, "**", "")

	// Loại bỏ khoảng trắng thừa
	cleaned = strings.TrimSpace(cleaned)
	cleaned = regexp.MustCompile(`\s+`).ReplaceAllString(cleaned, " ")

	return cleaned
}
