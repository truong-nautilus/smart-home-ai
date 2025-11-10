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
	"github.com/truong-nautilus/smart-home-ai/internal/infrastructure/wav2vec2"
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

	// Speech recognizer - chọn giữa PhoWhisper hoặc Wav2Vec2
	asrModel := os.Getenv("ASR_MODEL")
	if asrModel == "" {
		asrModel = "phowhisper" // mặc định
	}

	var recognizer *phowhisper.PhoWhisperRecognizer
	var wav2vec2Recognizer *wav2vec2.Wav2Vec2Recognizer

	if asrModel == "wav2vec2" {
		// Wav2Vec2 recognizer (fast Vietnamese CTC model)
		wav2vec2Script := "/Users/phamthetruong/github/smart-home-ai/scripts/wav2vec2_transcribe.py"
		wav2vec2Recognizer = wav2vec2.NewWav2Vec2Recognizer(wav2vec2Script)
		consoleLogger.Info("🎤 Sử dụng Wav2Vec2 ASR")
	} else {
		// PhoWhisper recognizer (vinai/PhoWhisper - tối ưu cho tiếng Việt)
		phowhisperScript := os.Getenv("PHOWHISPER_SCRIPT")
		if phowhisperScript == "" {
			phowhisperScript = "/Users/phamthetruong/github/smart-home-ai/scripts/phowhisper_transcribe.py"
		}
		recognizer = phowhisper.NewPhoWhisperRecognizer(phowhisperScript)
		consoleLogger.Info("🎤 Sử dụng PhoWhisper ASR")
	}

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
	var assistant *usecase.AssistantUseCase
	if wav2vec2Recognizer != nil {
		assistant = usecase.NewAssistantUseCase(
			ffmpeg,             // media capturer
			wav2vec2Recognizer, // speech recognizer (Wav2Vec2)
			aiClient,           // ai assistant (ollama)
			synthesizer,        // speech synthesizer (Edge TTS)
			keyboardListener,   // keyboard listener (Enter key)
			consoleLogger,
		)
	} else {
		assistant = usecase.NewAssistantUseCase(
			ffmpeg,           // media capturer
			recognizer,       // speech recognizer (PhoWhisper)
			aiClient,         // ai assistant (ollama)
			synthesizer,      // speech synthesizer (Edge TTS)
			keyboardListener, // keyboard listener (Enter key)
			consoleLogger,
		)
	}

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
