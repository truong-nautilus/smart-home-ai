package usecase

import (
	"context"
	"fmt"
	"os"

	"github.com/truong-nautilus/smart-home-ai/internal/domain"
)

// GestureDetector phát hiện cử chỉ từ camera
type GestureDetector interface {
	WaitForTwoFingers(ctx context.Context) (bool, error)
}

// AssistantUseCase orchestrates the AI assistant workflow
type AssistantUseCase struct {
	gestureDetector   GestureDetector
	mediaCapturer     domain.MediaCapturer
	speechRecognizer  domain.SpeechRecognizer
	aiAssistant       domain.AIAssistant
	speechSynthesizer domain.SpeechSynthesizer
	logger            Logger
}

// Logger interface for logging
type Logger interface {
	Info(msg string)
	Error(msg string, err error)
}

// NewAssistantUseCase creates a new assistant use case
func NewAssistantUseCase(
	gestureDetector GestureDetector,
	mediaCapturer domain.MediaCapturer,
	speechRecognizer domain.SpeechRecognizer,
	aiAssistant domain.AIAssistant,
	speechSynthesizer domain.SpeechSynthesizer,
	logger Logger,
) *AssistantUseCase {
	return &AssistantUseCase{
		gestureDetector:   gestureDetector,
		mediaCapturer:     mediaCapturer,
		speechRecognizer:  speechRecognizer,
		aiAssistant:       aiAssistant,
		speechSynthesizer: speechSynthesizer,
		logger:            logger,
	}
}

// Execute runs the complete AI assistant workflow
func (uc *AssistantUseCase) Execute(ctx context.Context) error {
	const (
		imageFile     = "frame.jpg"
		audioFile     = "audio.wav"
		replyFile     = "reply.mp3"
		audioDuration = 5
	)

	// Cleanup temp files on exit
	defer uc.cleanup(imageFile, audioFile, replyFile)

	// Step 1: Wait for gesture trigger (chờ vô hạn)
	uc.logger.Info("👋 Hãy giơ 2 ngón tay trước camera để bắt đầu (đang chờ...)...")
	detected, err := uc.gestureDetector.WaitForTwoFingers(ctx)
	if err != nil {
		return fmt.Errorf("không thể phát hiện cử chỉ: %w", err)
	}
	if !detected {
		return fmt.Errorf("không phát hiện được cử chỉ 2 ngón tay")
	}
	uc.logger.Info("✅ Đã phát hiện cử chỉ 2 ngón tay!")

	// Step 2: Capture image
	uc.logger.Info("🎥 Đang chụp ảnh từ camera...")
	if err := uc.mediaCapturer.CaptureImage(ctx, imageFile); err != nil {
		return fmt.Errorf("không thể chụp ảnh: %w", err)
	}
	uc.logger.Info("✅ Chụp ảnh thành công")

	// Step 3: Record audio
	uc.logger.Info("🎤 Đang ghi âm từ microphone (5 giây)...")
	if err := uc.mediaCapturer.RecordAudio(ctx, audioFile, audioDuration); err != nil {
		return fmt.Errorf("không thể ghi âm: %w", err)
	}
	uc.logger.Info("✅ Ghi âm thành công")

	// Step 4: Transcribe audio
	uc.logger.Info("🧠 Đang chuyển giọng nói thành văn bản (whisper.cpp)...")
	transcription, err := uc.speechRecognizer.Transcribe(ctx, audioFile)
	if err != nil {
		return fmt.Errorf("không thể chuyển giọng nói: %w", err)
	}
	uc.logger.Info(fmt.Sprintf("📝 Văn bản: \"%s\"", transcription.Text))

	// Step 5: Analyze with AI
	uc.logger.Info("🤖 Đang phân tích (ollama, mô hình local)...")
	response, err := uc.aiAssistant.AnalyzeMultimodal(ctx, transcription.Text, imageFile)
	if err != nil {
		return fmt.Errorf("không thể nhận phản hồi từ AI: %w", err)
	}
	uc.logger.Info(fmt.Sprintf("💬 Phản hồi AI: \"%s\"", response.Text))

	// Step 6: Synthesize speech
	uc.logger.Info("🔊 Đang tổng hợp giọng nói (sử dụng 'say' trên macOS)...")
	if _, err := uc.speechSynthesizer.Synthesize(ctx, response.Text, replyFile); err != nil {
		return fmt.Errorf("không thể tổng hợp giọng nói: %w", err)
	}
	uc.logger.Info("✅ Tổng hợp giọng nói thành công")

	// Step 7: Play audio
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
