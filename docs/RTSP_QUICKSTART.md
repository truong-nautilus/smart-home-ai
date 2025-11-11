# 🚀 Quick Start: RTSP Video Analysis

## TL;DR

```bash
# 1. Cấu hình RTSP URL trong .env
echo 'RTSP_URL=rtsp://username:password@192.168.1.xxx:554/stream1' >> .env

# 2. Chạy
go run cmd/assistant/main.go

# 3. Enjoy!
# - Video được phân tích tự động mỗi 10 giây
# - Nhấn ENTER để voice chat với AI
```

---

## 📋 Chi Tiết

### Bước 1: Chuẩn Bị RTSP URL

**Tìm RTSP URL của camera:**
- Tapo Camera: `rtsp://username:password@ip:554/stream1`
- Other IP Camera: Xem manual hoặc app settings

**URL encode password:**
```
@ → %40
# → %23
! → %21
```

Ví dụ: `Tapo@2024` → `Tapo%402024`

### Bước 2: Cấu Hình

Tạo/sửa file `.env`:
```bash
# RTSP Stream
RTSP_URL=rtsp://obstinate:Tapo%402024@192.168.1.186:554/stream1

# AI Model (optional)
OLLAMA_MODEL=gemma2:2b

# ASR Model (optional)
ASR_MODEL=phowhisper
```

### Bước 3: Test RTSP Connection

```bash
# Test với ffmpeg
ffmpeg -rtsp_transport tcp -i "YOUR_RTSP_URL" -frames:v 1 test.jpg

# Hoặc test với script
./scripts/test-rtsp.sh
```

### Bước 4: Run

```bash
# Option 1: Run directly
go run cmd/assistant/main.go

# Option 2: Build first
go build -o smart-home-ai cmd/assistant/main.go
./smart-home-ai
```

### Bước 5: Sử Dụng

**Video Analysis (Automatic)**
- Tự động chạy mỗi 10 giây
- Xem output trong console:
  ```
  📺 Video: Căn phòng có ánh sáng tự nhiên...
  ```

**Voice Assistant (Manual)**
- Nhấn ENTER lần 1 → Ghi âm
- Nhấn ENTER lần 2 → Xử lý
- Nghe trả lời

---

## 🎯 Output Mẫu

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

[Chờ voice command...]
```

---

## ⚙️ Tuning

### Thay đổi interval

**File**: `cmd/assistant/main.go`, dòng ~112
```go
assistant.StartContinuousVideoAnalysis(ctx, 10)  // 10 giây
                                            ↑
                                    Đổi thành: 5, 15, 30...
```

### Thay đổi resolution

**File**: `internal/infrastructure/video/rtsp_analyzer.go`, dòng ~45
```go
"-vf", "scale=1280:720",  // HD
                  ↑
         Đổi thành: 640x480, 1920x1080...
```

### Thay đổi AI prompt

**File**: `internal/infrastructure/video/rtsp_analyzer.go`, dòng ~99
```go
prompt := "Mô tả ngắn gọn..."  // Tùy chỉnh prompt
```

---

## 🔧 Troubleshooting

### Cannot connect to RTSP
```bash
# Test kết nối
ping 192.168.1.186

# Test RTSP với VLC
vlc rtsp://...
```

### Video quá chậm
- Tăng interval: 10s → 30s
- Giảm resolution: 1280x720 → 640x480
- Dùng model nhẹ: `OLLAMA_MODEL=gemma2:2b`

### ffmpeg not found
```bash
brew install ffmpeg
```

---

## 📚 Docs

- [Chi tiết đầy đủ](./RTSP_VIDEO_ANALYSIS.md)
- [Tổng kết implementation](./RTSP_IMPLEMENTATION_COMPLETE.md)
- [Update summary](./RTSP_UPDATE_SUMMARY.md)

---

**Ready?** → `go run cmd/assistant/main.go` 🚀
