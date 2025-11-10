package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/truong-nautilus/smart-home-ai/internal/infrastructure/edgetts"
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

	// Use case
	assistant := usecase.NewAssistantUseCase(
		ffmpeg,      // media capturer
		recognizer,  // speech recognizer (whisper.cpp)
		aiClient,    // ai assistant (ollama)
		synthesizer, // speech synthesizer (MacTTS)
		consoleLogger,
	)

	// Thực thi
	ctx := context.Background()
	if err := assistant.Execute(ctx); err != nil {
		consoleLogger.Error("Không thể thực thi trợ lý", err)
		os.Exit(1)
	}
}
