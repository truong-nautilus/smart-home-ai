package usecase

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/truong-nautilus/smart-home-ai/internal/domain"
)

// KeyboardListener interface cho việc lắng nghe phím bấm (hold/release)
type KeyboardListener interface {
	WaitForSpacePress() error   // Chờ Space được nhấn
	WaitForSpaceRelease() error // Chờ Space được nhả
}

// AssistantUseCase orchestrates the AI assistant workflow
type AssistantUseCase struct {
	mediaCapturer     domain.MediaCapturer
	speechRecognizer  domain.SpeechRecognizer
	aiAssistant       domain.AIAssistant
	speechSynthesizer domain.SpeechSynthesizer
	keyboardListener  KeyboardListener
	logger            Logger
}

// Logger interface for logging
type Logger interface {
	Info(msg string)
	Error(msg string, err error)
}

// NewAssistantUseCase creates a new assistant use case
func NewAssistantUseCase(
	mediaCapturer domain.MediaCapturer,
	speechRecognizer domain.SpeechRecognizer,
	aiAssistant domain.AIAssistant,
	speechSynthesizer domain.SpeechSynthesizer,
	keyboardListener KeyboardListener,
	logger Logger,
) *AssistantUseCase {
	return &AssistantUseCase{
		mediaCapturer:     mediaCapturer,
		speechRecognizer:  speechRecognizer,
		aiAssistant:       aiAssistant,
		speechSynthesizer: speechSynthesizer,
		keyboardListener:  keyboardListener,
		logger:            logger,
	}
}

// Execute runs the complete AI assistant workflow (hold-space voice mode)
func (uc *AssistantUseCase) Execute(ctx context.Context) error {
	const (
		audioFile = "audio.wav"
		replyFile = "reply.mp3"
	)

	// Cleanup temp files on exit
	defer uc.cleanup(audioFile, replyFile)

	// Step 1: Chờ người dùng nhấn Space
	if err := uc.keyboardListener.WaitForSpacePress(); err != nil {
		uc.logger.Error("❌ Lỗi khi đọc phím", err)
		return fmt.Errorf("không thể đọc phím bấm: %w", err)
	}

	// Step 2: Bắt đầu ghi âm trong background
	cancelRecording, err := uc.mediaCapturer.StartRecording(ctx, audioFile)
	if err != nil {
		uc.logger.Error("❌ Lỗi bắt đầu ghi âm", err)
		return fmt.Errorf("không thể bắt đầu ghi âm: %w", err)
	}
	defer cancelRecording() // Đảm bảo cancel nếu có lỗi

	// Step 3: Chờ người dùng nhả Space
	if err := uc.keyboardListener.WaitForSpaceRelease(); err != nil {
		uc.logger.Error("❌ Lỗi khi đọc phím", err)
		return fmt.Errorf("không thể đọc phím nhả: %w", err)
	}

	// Step 4: Dừng ghi âm
	if err := uc.mediaCapturer.StopRecording(); err != nil {
		uc.logger.Error("❌ Lỗi dừng ghi âm", err)
		return fmt.Errorf("không thể dừng ghi âm: %w", err)
	}
	uc.logger.Info("✅ Đã ghi âm xong")

	// Step 5: Transcribe audio
	uc.logger.Info("🧠 Đang chuyển giọng nói thành văn bản...")
	transcription, err := uc.speechRecognizer.Transcribe(ctx, audioFile)
	if err != nil {
		uc.logger.Error("❌ Lỗi transcribe", err)
		return fmt.Errorf("không thể chuyển giọng nói: %w", err)
	}

	// Log để debug
	text := transcription.Text
	uc.logger.Info(fmt.Sprintf("🔍 Whisper output: \"%s\"", text))

	// Kiểm tra xem có nội dung thực sự không (bỏ qua blank audio, music, noise)
	if text == "" ||
		strings.Contains(text, "[BLANK_AUDIO]") ||
		strings.Contains(text, "[Music]") ||
		strings.Contains(text, "[Silence]") ||
		strings.Contains(text, "(electronic beeping)") ||
		len(strings.TrimSpace(text)) < 3 {
		uc.logger.Info("⚠️ Không phát hiện giọng nói rõ ràng, tiếp tục lắng nghe...")
		return nil // Không lỗi, chỉ là không có giọng nói
	}

	uc.logger.Info(fmt.Sprintf("📝 Câu hỏi: \"%s\"", text))

	// Step 6: Analyze with AI (không cần hình ảnh)
	uc.logger.Info("🤖 Đang xử lý câu hỏi...")
	response, err := uc.aiAssistant.AnalyzeMultimodal(ctx, text, "")
	if err != nil {
		return fmt.Errorf("không thể nhận phản hồi từ AI: %w", err)
	}
	uc.logger.Info(fmt.Sprintf("💬 Trả lời: \"%s\"", response.Text))

	// Step 7: Synthesize speech
	uc.logger.Info("🔊 Đang tổng hợp giọng nói (sử dụng 'say' trên macOS)...")
	if _, err := uc.speechSynthesizer.Synthesize(ctx, response.Text, replyFile); err != nil {
		return fmt.Errorf("không thể tổng hợp giọng nói: %w", err)
	}
	uc.logger.Info("✅ Tổng hợp giọng nói thành công")

	// Step 8: Play audio
	uc.logger.Info("🔈 Đang phát âm thanh phản hồi...")
	if err := uc.mediaCapturer.PlayAudio(ctx, replyFile); err != nil {
		return fmt.Errorf("không thể phát âm thanh: %w", err)
	}

	uc.logger.Info("✅ Hoàn thành!")
	return nil
}

func (uc *AssistantUseCase) cleanup(files ...string) {
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			os.Remove(file)
		}
	}
}
