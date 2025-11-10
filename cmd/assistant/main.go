package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/truong-nautilus/smart-home-ai/internal/infrastructure/edgetts"
	"github.com/truong-nautilus/smart-home-ai/internal/infrastructure/keyboard"
	"github.com/truong-nautilus/smart-home-ai/internal/infrastructure/media"
	"github.com/truong-nautilus/smart-home-ai/internal/infrastructure/ollama"
	"github.com/truong-nautilus/smart-home-ai/internal/infrastructure/phowhisper"
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

	// PhoWhisper recognizer (vinai/PhoWhisper-small - tối ưu cho tiếng Việt)
	phowhisperScript := os.Getenv("PHOWHISPER_SCRIPT")
	if phowhisperScript == "" {
		// Sử dụng đường dẫn tuyệt đối
		phowhisperScript = "/Users/phamthetruong/github/smart-home-ai/scripts/phowhisper_transcribe.py"
	}
	recognizer := phowhisper.NewPhoWhisperRecognizer(phowhisperScript)

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

	// Keyboard listener (Space key để ghi âm)
	keyboardListener := keyboard.NewListener()

	// Use case (với keyboard listener)
	assistant := usecase.NewAssistantUseCase(
		ffmpeg,           // media capturer
		recognizer,       // speech recognizer (PhoWhisper)
		aiClient,         // ai assistant (ollama)
		synthesizer,      // speech synthesizer (Edge TTS)
		keyboardListener, // keyboard listener (Enter key)
		consoleLogger,
	)

	// Thực thi vô hạn - chế độ press twice
	ctx := context.Background()
	consoleLogger.Info("🚀 Trợ lý AI đã sẵn sàng!")
	consoleLogger.Info("📌 Cách dùng: Nhấn ENTER lần 1 → ghi âm → nhấn ENTER lần 2 → xử lý")
	consoleLogger.Info("🛑 Nhấn Ctrl+C để thoát")

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
