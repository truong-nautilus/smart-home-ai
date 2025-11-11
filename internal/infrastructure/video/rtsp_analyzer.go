package video

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/truong-nautilus/smart-home-ai/internal/domain"
)

// RTSPAnalyzer triển khai VideoAnalyzer cho RTSP stream
type RTSPAnalyzer struct {
	rtspURL     string
	aiAssistant domain.AIAssistant
	logger      Logger
}

// Logger interface for logging
type Logger interface {
	Info(msg string)
	Error(msg string, err error)
}

// NewRTSPAnalyzer tạo RTSP analyzer mới
func NewRTSPAnalyzer(rtspURL string, aiAssistant domain.AIAssistant, logger Logger) *RTSPAnalyzer {
	return &RTSPAnalyzer{
		rtspURL:     rtspURL,
		aiAssistant: aiAssistant,
		logger:      logger,
	}
}

// CaptureFrame captures a single frame from RTSP stream
func (r *RTSPAnalyzer) CaptureFrame(ctx context.Context, outputPath string) error {
	cmd := exec.CommandContext(
		ctx,
		"ffmpeg",
		"-loglevel", "error", // Chỉ hiển thị lỗi, bỏ warnings
		"-rtsp_transport", "tcp", // Sử dụng TCP thay vì UDP để ổn định hơn
		"-i", r.rtspURL,
		"-frames:v", "1", // Chỉ lấy 1 frame
		"-q:v", "2", // Chất lượng cao
		"-vf", "scale=1280:720", // Resize về 720p
		"-y", // Ghi đè file output
		outputPath,
	)
	// Bỏ stderr để giảm log nhiễu
	// cmd.Stderr = os.Stderr
	return cmd.Run()
}

// StartContinuousAnalysis starts continuous video analysis loop
// Callback sẽ được gọi mỗi khi có mô tả mới về video
func (r *RTSPAnalyzer) StartContinuousAnalysis(ctx context.Context, intervalSec int, callback func(description string)) error {
	ticker := time.NewTicker(time.Duration(intervalSec) * time.Second)
	defer ticker.Stop()

	r.logger.Info(fmt.Sprintf("🎥 Bắt đầu phân tích video liên tục từ RTSP stream mỗi %d giây", intervalSec))
	r.logger.Info(fmt.Sprintf("📹 RTSP URL: %s", r.rtspURL))

	// Phân tích ngay lập tức lần đầu tiên
	if err := r.analyzeFrame(ctx, callback); err != nil {
		r.logger.Error("⚠️ Lỗi khi phân tích frame đầu tiên", err)
	}

	// Loop vô hạn, phân tích mỗi intervalSec giây
	for {
		select {
		case <-ctx.Done():
			r.logger.Info("🛑 Dừng phân tích video")
			return ctx.Err()
		case <-ticker.C:
			if err := r.analyzeFrame(ctx, callback); err != nil {
				r.logger.Error("⚠️ Lỗi khi phân tích frame", err)
				// Không return, tiếp tục thử frame tiếp theo
			}
		}
	}
}

// ShowVideoPreview hiển thị video stream trong cửa sổ preview (dùng ffplay)
func (r *RTSPAnalyzer) ShowVideoPreview(ctx context.Context) error {
	r.logger.Info("🖥️  Mở cửa sổ video preview...")
	cmd := exec.CommandContext(
		ctx,
		"ffplay",
		"-rtsp_transport", "tcp",
		"-i", r.rtspURL,
		"-window_title", "RTSP Video Preview",
		"-x", "960", // Chiều rộng cửa sổ
		"-y", "540", // Chiều cao cửa sổ
		"-left", "100", // Vị trí x
		"-top", "100", // Vị trí y
	)
	// Chạy ffplay và bỏ qua lỗi khi user đóng cửa sổ
	if err := cmd.Run(); err != nil {
		r.logger.Info("🛑 Đã đóng cửa sổ video preview")
	}
	return nil
}

// analyzeFrame captures và phân tích một frame từ RTSP stream
func (r *RTSPAnalyzer) analyzeFrame(ctx context.Context, callback func(description string)) error {
	// Tạo file tạm để lưu frame
	frameFile := fmt.Sprintf("frame_%d.jpg", time.Now().Unix())
	defer os.Remove(frameFile) // Cleanup

	// Bỏ log để giảm nhiễu
	// r.logger.Info("📸 Đang bắt frame từ RTSP stream...")

	// Capture frame
	if err := r.CaptureFrame(ctx, frameFile); err != nil {
		return fmt.Errorf("không thể capture frame: %w", err)
	}

	// Hiển thị frame vừa capture (mở bằng Preview app trên macOS)
	go func() {
		exec.Command("open", "-a", "Preview", frameFile).Run()
		time.Sleep(3 * time.Second) // Giữ Preview mở 3 giây
	}()

	// Bỏ log để giảm nhiễu
	// r.logger.Info("🧠 Đang phân tích nội dung video...")

	// Phân tích frame bằng AI
	prompt := "Mô tả ngắn gọn những gì bạn thấy trong video này. Hãy chỉ ra các đối tượng, hành động, và môi trường quan trọng."
	response, err := r.aiAssistant.AnalyzeMultimodal(ctx, prompt, frameFile)
	if err != nil {
		return fmt.Errorf("không thể phân tích frame: %w", err)
	}

	description := response.Text
	r.logger.Info(fmt.Sprintf("👁️  Phân tích: %s", description))

	// Gọi callback với kết quả
	if callback != nil {
		callback(description)
	}

	return nil
}
