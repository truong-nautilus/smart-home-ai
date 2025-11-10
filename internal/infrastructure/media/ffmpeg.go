package media

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

// FFmpegCapturer triển khai MediaCapturer sử dụng FFmpeg
type FFmpegCapturer struct{}

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
