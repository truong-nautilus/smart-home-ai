package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/truong-nautilus/smart-home-ai/internal/infrastructure/edgetts"
	"github.com/truong-nautilus/smart-home-ai/internal/infrastructure/gesture"
	"github.com/truong-nautilus/smart-home-ai/internal/infrastructure/media"
	"github.com/truong-nautilus/smart-home-ai/internal/infrastructure/ollama"
	"github.com/truong-nautilus/smart-home-ai/internal/infrastructure/whispercpp"
	"github.com/truong-nautilus/smart-home-ai/internal/usecase"
	"github.com/truong-nautilus/smart-home-ai/pkg/logger"
)

func main() {
	// Tải file .env nếu có
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Không tìm thấy file .env, sử dụng biến môi trường hệ thống")
	}

	// Khởi tạo các phụ thuộc (local)
	consoleLogger := logger.NewConsoleLogger()
	ffmpeg := media.NewFFmpegCapturer()

	// Whisper.cpp recognizer (local). Nếu cần, chỉnh path của binary và model qua .env
	whisperBin := os.Getenv("WHISPER_CPP_BIN") // optional
	whisperModel := os.Getenv("WHISPER_CPP_MODEL")
	recognizer := whispercpp.New(whisperBin, whisperModel)

	// Ollama local model
	ollamaModel := os.Getenv("OLLAMA_MODEL")
	if ollamaModel == "" {
		ollamaModel = "phi-3-mini"
	}
	aiClient := ollama.New(ollamaModel)

	// Edge TTS synthesizer (Microsoft neural TTS - giọng rất tự nhiên)
	edgeTTSVoice := os.Getenv("EDGE_TTS_VOICE") // vi-VN-HoaiMyNeural (nữ) hoặc vi-VN-NamMinhNeural (nam)
	edgeTTSBin := os.Getenv("EDGE_TTS_BIN")     // optional
	synthesizer := edgetts.New(edgeTTSVoice, edgeTTSBin)

	// Gesture detector (MediaPipe hand tracking)
	var gestureDetector usecase.GestureDetector

	// Nếu ENABLE_GESTURE=false, bỏ qua gesture detection
	enableGesture := os.Getenv("ENABLE_GESTURE")
	if enableGesture == "false" {
		gestureDetector = gesture.NewNoOpDetector() // Luôn trả về true
	} else {
		gestureDetector = gesture.NewDetector("./scripts/detect_gesture.py")
	}

	// Use case
	assistant := usecase.NewAssistantUseCase(
		gestureDetector, // gesture detector (MediaPipe)
		ffmpeg,          // media capturer
		recognizer,      // speech recognizer (whisper.cpp)
		aiClient,        // ai assistant (ollama)
		synthesizer,     // speech synthesizer (Edge TTS)
		consoleLogger,
	)

	// Thực thi vô hạn - chạy liên tục
	ctx := context.Background()
	consoleLogger.Info("🚀 Trợ lý AI đã sẵn sàng - chạy liên tục...")
	consoleLogger.Info("📌 Nhấn Ctrl+C để thoát")

	for {
		if err := assistant.Execute(ctx); err != nil {
			consoleLogger.Error("⚠️ Lỗi khi thực thi", err)
			// Không thoát, tiếp tục chạy
			consoleLogger.Info("🔄 Khởi động lại sau 2 giây...")
			time.Sleep(2 * time.Second)
			continue
		}

		// Nghỉ 1 giây trước khi chờ gesture tiếp theo
		consoleLogger.Info("✨ Hoàn tất! Sẵn sàng cho lần tiếp theo...")
		time.Sleep(1 * time.Second)
	}
}
