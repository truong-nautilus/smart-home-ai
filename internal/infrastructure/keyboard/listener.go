package keyboard

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

// Listener xử lý việc lắng nghe phím Space (hold/release)
// Sử dụng stdin đơn giản: nhấn Enter để bắt đầu, nhấn Enter lần nữa để dừng
type Listener struct {
	reader *bufio.Reader
}

// NewListener tạo instance mới của keyboard listener
func NewListener() *Listener {
	return &Listener{
		reader: bufio.NewReader(os.Stdin),
	}
}

// WaitForSpacePress chờ người dùng nhấn Enter lần 1 để bắt đầu ghi âm
func (l *Listener) WaitForSpacePress() error {
	fmt.Println("\n⏸️  Nhấn ENTER lần 1 để bắt đầu ghi âm...")

	_, err := l.reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("lỗi khi đọc input: %v", err)
	}

	fmt.Println("🔴 Đang ghi âm... (nhấn ENTER lần 2 để dừng và xử lý)")
	return nil
}

// WaitForSpaceRelease chờ người dùng nhấn Enter lần 2 để dừng ghi âm
func (l *Listener) WaitForSpaceRelease() error {
	// Thêm delay nhỏ để đảm bảo recording đã bắt đầu
	time.Sleep(100 * time.Millisecond)

	// Đọc line tiếp theo (người dùng nhấn Enter lần 2)
	_, err := l.reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("lỗi khi đọc input: %v", err)
	}

	fmt.Println("⏹️  Đã nhấn lần 2, đang xử lý...")
	return nil
}
