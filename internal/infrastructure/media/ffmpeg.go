package media

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

// FFmpegCapturer triển khai MediaCapturer sử dụng FFmpeg
type FFmpegCapturer struct {
	recordingCmd *exec.Cmd
	recordingMu  sync.Mutex
}

// NewFFmpegCapturer tạo FFmpeg capturer mới
func NewFFmpegCapturer() *FFmpegCapturer {
	return &FFmpegCapturer{}
}

// CaptureImage chụp một khung hình từ camera
func (f *FFmpegCapturer) CaptureImage(ctx context.Context, outputPath string) error {
	cmd := exec.CommandContext(
		ctx,
		"ffmpeg",
		"-f", "avfoundation",
		"-framerate", "30",
		"-video_size", "640x480",
		"-i", "0", // 0 là camera tích hợp
		"-frames:v", "1",
		"-y", // Ghi đè file output
		outputPath,
	)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RecordAudio ghi âm từ microphone
func (f *FFmpegCapturer) RecordAudio(ctx context.Context, outputPath string, duration int) error {
	cmd := exec.CommandContext(
		ctx,
		"ffmpeg",
		"-f", "avfoundation",
		"-i", ":0", // :0 là microphone tích hợp
		"-t", fmt.Sprintf("%d", duration),
		"-ac", "1", // Mono
		"-ar", "16000", // 16kHz sample rate
		"-y", // Ghi đè file output
		outputPath,
	)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// StartRecording bắt đầu ghi âm trong background
// Ghi âm sẽ tiếp tục cho đến khi StopRecording() được gọi hoặc context bị cancel
func (f *FFmpegCapturer) StartRecording(ctx context.Context, outputPath string) (context.CancelFunc, error) {
	f.recordingMu.Lock()
	defer f.recordingMu.Unlock()

	if f.recordingCmd != nil {
		return nil, fmt.Errorf("đang có recording đang chạy")
	}

	// Tạo context có thể cancel
	recordCtx, cancel := context.WithCancel(ctx)

	// Tạo command ghi âm không giới hạn thời gian
	cmd := exec.CommandContext(
		recordCtx,
		"ffmpeg",
		"-f", "avfoundation",
		"-i", ":0", // :0 là microphone tích hợp
		"-ac", "1",     // Mono
		"-ar", "16000", // 16kHz sample rate
		"-y",           // Ghi đè file output
		outputPath,
	)
	cmd.Stderr = os.Stderr

	// Start recording
	if err := cmd.Start(); err != nil {
		cancel()
		return nil, fmt.Errorf("không thể bắt đầu ghi âm: %w", err)
	}

	f.recordingCmd = cmd

	// Chạy goroutine để chờ command hoàn thành
	go func() {
		cmd.Wait()
		f.recordingMu.Lock()
		f.recordingCmd = nil
		f.recordingMu.Unlock()
	}()

	return cancel, nil
}

// StopRecording dừng recording đang chạy
func (f *FFmpegCapturer) StopRecording() error {
	f.recordingMu.Lock()
	defer f.recordingMu.Unlock()

	if f.recordingCmd == nil {
		return fmt.Errorf("không có recording nào đang chạy")
	}

	// Gửi signal SIGTERM để FFmpeg kết thúc gracefully
	if err := f.recordingCmd.Process.Signal(os.Interrupt); err != nil {
		// Nếu SIGTERM fail, thử SIGKILL
		f.recordingCmd.Process.Kill()
	}

	// Đợi process kết thúc
	f.recordingCmd.Wait()
	f.recordingCmd = nil

	return nil
}

// PlayAudio phát file audio sử dụng ffplay
func (f *FFmpegCapturer) PlayAudio(ctx context.Context, audioPath string) error {
	cmd := exec.CommandContext(
		ctx,
		"ffplay",
		"-nodisp",   // Không hiển thị video
		"-autoexit", // Thoát khi phát xong
		"-loglevel", "quiet", // Tắt log ffplay
		audioPath,
	)
	return cmd.Run()
}
