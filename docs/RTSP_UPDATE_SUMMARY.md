# 🎥 RTSP Video Analysis Update

## Tóm Tắt Thay Đổi

Hệ thống đã được nâng cấp để **thay thế camera tích hợp** bằng **RTSP stream** và thêm tính năng **phân tích video liên tục**.

## ✨ Tính Năng Mới

### 1. RTSP Stream Support
- Kết nối với camera IP qua RTSP protocol
- Hỗ trợ authentication (username/password)
- URL mặc định: `rtsp://obstinate:Tapo%402024@192.168.1.186:554/stream1`

### 2. Continuous Video Analysis
- Tự động phân tích video **mỗi 10 giây**
- AI mô tả những gì đang xảy ra trong video
- Chạy song song với voice assistant

### 3. Dual Mode Operation
- **Voice Assistant**: Interactive voice commands (nhấn ENTER để ghi âm)
- **Video Monitor**: Automatic video analysis mỗi 10 giây (background)

## 📁 Files Đã Thay Đổi

### Mới
- `internal/infrastructure/video/rtsp_analyzer.go` - RTSP video analyzer
- `docs/RTSP_VIDEO_ANALYSIS.md` - Hướng dẫn chi tiết

### Cập nhật
- `internal/domain/repository.go` - Thêm VideoAnalyzer interface
- `internal/usecase/assistant.go` - Thêm continuous video analysis
- `cmd/assistant/main.go` - Khởi tạo RTSP analyzer
- `.env` - Thêm RTSP_URL configuration

## 🚀 Cách Sử Dụng

### 1. Cấu hình RTSP URL

Sửa file `.env`:
```bash
RTSP_URL=rtsp://username:password@ip:port/path
```

**Lưu ý**: URL encode password nếu có ký tự đặc biệt:
- `@` → `%40`
- `#` → `%23`
- Ví dụ: `Tapo@2024` → `Tapo%402024`

### 2. Chạy chương trình

```bash
go run cmd/assistant/main.go
```

### 3. Output mẫu

```
🚀 Trợ lý AI đã sẵn sàng!
📌 Cách dùng: Nhấn ENTER lần 1 → ghi âm → nhấn ENTER lần 2 → xử lý
🎥 Video RTSP sẽ được phân tích liên tục mỗi 10 giây
🛑 Nhấn Ctrl+C để thoát

🎥 Bắt đầu phân tích video liên tục từ RTSP stream mỗi 10 giây
📹 RTSP URL: rtsp://obstinate:Tapo%402024@192.168.1.186:554/stream1
📸 Đang bắt frame từ RTSP stream...
🧠 Đang phân tích nội dung video...
👁️  Phân tích: Căn phòng với ánh sáng tự nhiên, bàn làm việc ở góc phải
📺 Video: Căn phòng với ánh sáng tự nhiên, bàn làm việc ở góc phải

[Voice assistant ready for commands]
[Nhấn ENTER để bắt đầu ghi âm...]
```

## ⚙️ Tùy Chỉnh

### Thay đổi interval (từ 10 giây)

File: `cmd/assistant/main.go`
```go
assistant.StartContinuousVideoAnalysis(ctx, 10) // Đổi 10 → 5, 15, 30, etc.
```

### Thay đổi prompt phân tích

File: `internal/infrastructure/video/rtsp_analyzer.go`
```go
prompt := "Mô tả ngắn gọn những gì bạn thấy..." // Tùy chỉnh prompt
```

### Thay đổi resolution

File: `internal/infrastructure/video/rtsp_analyzer.go`
```go
"-vf", "scale=1280:720", // 1280x720 → 640x480, 1920x1080, etc.
```

## 🏗️ Architecture

```
┌─────────────────────────────────────────┐
│           Main Application              │
│         (cmd/assistant/main.go)         │
└────────────┬──────────────┬─────────────┘
             │              │
             │              │
    ┌────────▼─────┐   ┌───▼──────────────┐
    │   Voice      │   │  Video Analyzer  │
    │  Assistant   │   │  (RTSP Stream)   │
    │  (Keyboard)  │   │  Every 10s       │
    └──────────────┘   └──────────────────┘
             │                   │
             │                   │
             └──────┬────────────┘
                    │
             ┌──────▼─────────┐
             │   AI Assistant │
             │    (Ollama)    │
             └────────────────┘
```

## 🔧 Technical Details

### RTSP Capture
- Protocol: RTSP over TCP (more stable)
- Frame capture: ffmpeg single frame extraction
- Format: JPEG, 720p
- Quality: High (q:v = 2)

### Video Analysis Loop
- Timer: time.Ticker (10 seconds interval)
- Async: Goroutine for non-blocking
- Error handling: Continue on error (don't crash)
- Cleanup: Automatic temp file removal

### AI Integration
- Model: Ollama (multimodal)
- Prompt: Customizable Vietnamese prompt
- Response: Natural language description
- Performance: Depends on model size

## 📖 Documentation

Chi tiết đầy đủ: [RTSP_VIDEO_ANALYSIS.md](./RTSP_VIDEO_ANALYSIS.md)

## 🐛 Troubleshooting

### Cannot connect to RTSP
1. Test với VLC: Media → Open Network Stream
2. Kiểm tra ping: `ping [camera_ip]`
3. Kiểm tra username/password
4. URL encode password đúng cách

### Video analysis quá chậm
1. Tăng interval (10s → 30s)
2. Giảm resolution (1280x720 → 640x480)
3. Dùng model nhẹ hơn (gemma2:2b thay vì phi-3-mini)

### ffmpeg errors
1. Cài đặt: `brew install ffmpeg`
2. Kiểm tra version: `ffmpeg -version`
3. Test RTSP: `ffmpeg -i [rtsp_url] -frames:v 1 test.jpg`

## 📝 Notes

- **Camera tích hợp**: Không còn được sử dụng (đã thay bằng RTSP)
- **Gesture detection**: Vẫn có thể dùng nếu cần (ENABLE_GESTURE=true)
- **Privacy**: Tất cả xử lý local, không gửi cloud
- **Performance**: CPU-intensive, tối ưu interval theo máy

## 🎯 Use Cases

- Giám sát an ninh nhà thông minh
- Theo dõi trẻ em/thú cưng
- Phát hiện chuyển động bất thường
- Mô tả môi trường tự động
- Smart home automation triggers
