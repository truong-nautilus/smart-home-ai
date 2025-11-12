package ollama

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
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

// AnalyzeMultimodal gọi Ollama với text và image (nếu có)
func (o *OllamaClient) AnalyzeMultimodal(ctx context.Context, text string, imagePath string) (*domain.AIResponse, error) {
	var prompt string
	var cmd *exec.Cmd

	// Nếu có imagePath, sử dụng vision model với ảnh
	if imagePath != "" {
		// Đọc và encode ảnh thành base64
		imageData, err := os.ReadFile(imagePath)
		if err != nil {
			return nil, fmt.Errorf("không thể đọc ảnh: %w", err)
		}
		base64Image := base64.StdEncoding.EncodeToString(imageData)

		// Tạo JSON request cho Ollama API với vision
		type Message struct {
			Role    string   `json:"role"`
			Content string   `json:"content"`
			Images  []string `json:"images,omitempty"`
		}
		type Request struct {
			Model    string    `json:"model"`
			Messages []Message `json:"messages"`
			Stream   bool      `json:"stream"`
		}

		reqData := Request{
			Model: "gemma2:2b", // Model hỗ trợ vision
			Messages: []Message{
				{
					Role:    "user",
					Content: text,
					Images:  []string{base64Image},
				},
			},
			Stream: false,
		}

		jsonData, err := json.Marshal(reqData)
		if err != nil {
			return nil, fmt.Errorf("không thể tạo JSON request: %w", err)
		}

		// Gọi Ollama API qua curl
		cmd = exec.CommandContext(ctx, "curl", "-s", "http://localhost:11434/api/chat",
			"-d", string(jsonData))
	} else {
		// Không có ảnh - xử lý text thuần túy cho voice assistant
		prompt = fmt.Sprintf(`Bạn là trợ lý AI thân thiện, vui vẻ, hài hước. Hãy trả lời câu hỏi sau:

"%s"

YÊU CẦU:
- Trả lời bằng tiếng Việt
- Phong cách: thân thiện, vui vẻ, hài hước, dễ thương
- Trả lời ngắn gọn (1-3 câu), tự nhiên như giọng nói
- KHÔNG dùng emoji, emoticon, ký tự đặc biệt
- Có thể đùa cợt nhẹ nhàng, tạo cảm giác gần gũi
- Nếu là yêu cầu điều khiển (bật/tắt đèn, quạt...), trả lời vui vẻ và xác nhận`, text)

		cmd = exec.CommandContext(ctx, "ollama", "run", o.Model, prompt)
	}

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

	// Nếu có ảnh, parse JSON response từ API
	if imagePath != "" {
		type Response struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		}
		var apiResp Response
		if err := json.Unmarshal([]byte(response), &apiResp); err != nil {
			return nil, fmt.Errorf("không thể parse JSON response: %w", err)
		}
		response = apiResp.Message.Content
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
