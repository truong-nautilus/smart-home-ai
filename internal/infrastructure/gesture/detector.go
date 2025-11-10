package gesture

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// Detector nhận diện cử chỉ từ camera
type Detector struct {
	ScriptPath string
}

// NewDetector tạo gesture detector mới
func NewDetector(scriptPath string) *Detector {
	if scriptPath == "" {
		scriptPath = "./scripts/detect_gesture.py"
	}
	return &Detector{ScriptPath: scriptPath}
}

// WaitForTwoFingers chờ vô hạn cho đến khi phát hiện 2 ngón tay từ camera
func (d *Detector) WaitForTwoFingers(ctx context.Context) (bool, error) {
	// Không dùng context timeout để chờ vô hạn
	cmd := exec.Command("python3", d.ScriptPath)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return false, fmt.Errorf("gesture: không tạo được stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return false, fmt.Errorf("gesture: không tạo được stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return false, fmt.Errorf("gesture: không khởi động được script: %w", err)
	}

	// Đọc stderr để log debug info
	go func() {
		stderrScanner := bufio.NewScanner(stderr)
		for stderrScanner.Scan() {
			// Log stderr để debug (nếu cần)
			_ = stderrScanner.Text()
		}
	}()

	scanner := bufio.NewScanner(stdout)
	detected := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "DETECTED_2_FINGERS" {
			detected = true
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("gesture: đọc output thất bại: %w", err)
	}

	cmd.Wait()

	return detected, nil
}

// NoOpDetector bỏ qua gesture detection (luôn trả về true)
type NoOpDetector struct{}

// NewNoOpDetector tạo no-op detector
func NewNoOpDetector() *NoOpDetector {
	return &NoOpDetector{}
}

// WaitForTwoFingers luôn trả về true ngay lập tức
func (d *NoOpDetector) WaitForTwoFingers(ctx context.Context) (bool, error) {
	return true, nil
}
