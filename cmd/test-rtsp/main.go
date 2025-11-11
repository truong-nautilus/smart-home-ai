package testrtsp
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/truong-nautilus/smart-home-ai/internal/infrastructure/ollama"
	"github.com/truong-nautilus/smart-home-ai/internal/infrastructure/video"
	"github.com/truong-nautilus/smart-home-ai/pkg/logger"
)

// Test RTSP video analyzer
func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Không tìm thấy file .env")
	}

	// Get RTSP URL
	rtspURL := os.Getenv("RTSP_URL")
	if rtspURL == "" {
		rtspURL = "rtsp://obstinate:Tapo%402024@192.168.1.186:554/stream1"
	}

	// Initialize components
	consoleLogger := logger.NewConsoleLogger()
	ollamaModel := os.Getenv("OLLAMA_MODEL")
	if ollamaModel == "" {
		ollamaModel = "gemma2:2b"
	}
	aiClient := ollama.New(ollamaModel)

	// Create video analyzer
	videoAnalyzer := video.NewRTSPAnalyzer(rtspURL, aiClient, consoleLogger)

	consoleLogger.Info("🎬 Test RTSP Video Analyzer")
	consoleLogger.Info(fmt.Sprintf("📹 RTSP URL: %s", rtspURL))
	consoleLogger.Info("🔄 Sẽ phân tích 3 frame (mỗi 5 giây)")
	consoleLogger.Info("🛑 Nhấn Ctrl+C để dừng sớm")

	ctx := context.Background()

	// Test 1: Capture single frame
	consoleLogger.Info("\n--- Test 1: Capture Frame ---")
	testFile := "test_frame.jpg"
	if err := videoAnalyzer.CaptureFrame(ctx, testFile); err != nil {
		consoleLogger.Error("❌ Lỗi capture frame", err)
		return
	}
	consoleLogger.Info(fmt.Sprintf("✅ Đã lưu frame: %s", testFile))

	// Test 2: Analyze 3 frames with 5 second interval
	consoleLogger.Info("\n--- Test 2: Continuous Analysis (3 frames) ---")
	
	count := 0
	maxCount := 3
	
	// Create context with timeout
	testCtx, cancel := context.WithTimeout(ctx, time.Duration(maxCount*5+5)*time.Second)
	defer cancel()

	// Callback to count frames
	callback := func(description string) {
		count++
		consoleLogger.Info(fmt.Sprintf("📊 Frame %d/%d đã phân tích", count, maxCount))
		if count >= maxCount {
			consoleLogger.Info("\n✅ Test hoàn tất! Dừng...")
			cancel()
		}
	}

	// Start continuous analysis
	if err := videoAnalyzer.StartContinuousAnalysis(testCtx, 5, callback); err != nil {
		if err == context.Canceled {
			consoleLogger.Info("✅ Test đã hoàn thành thành công!")
		} else {
			consoleLogger.Error("❌ Lỗi continuous analysis", err)
		}
	}

	consoleLogger.Info("\n🎉 Test kết thúc!")
	consoleLogger.Info(fmt.Sprintf("📸 Kiểm tra file: %s", testFile))
}
