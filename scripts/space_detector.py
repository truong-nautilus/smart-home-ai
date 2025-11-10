#!/usr/bin/env python3
"""
Script để phát hiện Space key press và release
Output: "PRESSED" khi nhấn Space, "RELEASED" khi nhả Space
"""

import sys
import termios
import tty
import select

def detect_space_key():
    """
    Phát hiện Space key press/release trong raw mode
    """
    # Lưu cài đặt terminal hiện tại
    old_settings = termios.tcgetattr(sys.stdin)
    
    try:
        # Chuyển sang raw mode để đọc từng key press
        tty.setraw(sys.stdin.fileno())
        
        print("⏸️  Giữ SPACE để ghi âm, nhả SPACE để xử lý...", file=sys.stderr, flush=True)
        
        while True:
            # Chờ input
            if select.select([sys.stdin], [], [], 0.1)[0]:
                char = sys.stdin.read(1)
                
                # Kiểm tra nếu là Space (ASCII 32 hay ' ')
                if char == ' ':
                    print("PRESSED", flush=True)
                    
                    # Chờ người dùng nhả phím
                    # Trong terminal raw mode, không có sự kiện key release
                    # Nên ta chỉ detect khi Space được nhấn lần nữa
                    while True:
                        if select.select([sys.stdin], [], [], 0.1)[0]:
                            next_char = sys.stdin.read(1)
                            if next_char == ' ':
                                # Coi như nhả rồi (do nhấn lần nữa)
                                continue
                            else:
                                # Nhấn phím khác - bỏ qua
                                pass
                        else:
                            # Timeout - giả sử đã nhả phím
                            break
                    
                    print("RELEASED", flush=True)
                
                # Ctrl+C để thoát
                elif char == '\x03':
                    break
    
    finally:
        # Khôi phục cài đặt terminal
        termios.tcsetattr(sys.stdin, termios.TCSADRAIN, old_settings)

if __name__ == "__main__":
    try:
        detect_space_key()
    except KeyboardInterrupt:
        print("\nThoát...", file=sys.stderr)
        sys.exit(0)
