# Hướng Dẫn Phân Tích Video RTSP Liên Tục

## 📹 Tổng Quan

Hệ thống đã được nâng cấp để hỗ trợ **phân tích video liên tục** từ RTSP stream. Thay vì chụp ảnh từ camera tích hợp, hệ thống giờ đây sẽ:

- Kết nối với camera RTSP (ví dụ: camera IP, Tapo, v.v.)
- Tự động phân tích video **mỗi 10 giây**
- Mô tả những gì đang xảy ra trong video
- Chạy song song với tính năng voice assistant

## 🔧 Cấu Hình RTSP URL

### Mặc định
URL RTSP mặc định đã được cấu hình trong code:
```
rtsp://obstinate:Tapo%402024@192.168.1.186:554/stream1
```

### Tùy chỉnh qua biến môi trường
Tạo file `.env` và thêm:
```bash
RTSP_URL=rtsp://username:password@192.168.1.xxx:554/stream1
```

### Format RTSP URL
```
rtsp://[username]:[password]@[ip]:[port]/[path]
```

**Lưu ý**: Ký tự đặc biệt trong password cần được URL encode:
- `@` → `%40`
- `#` → `%23`
- `!` → `%21`
- Ví dụ: `Tapo@2024` → `Tapo%402024`

## 🚀 Cách Sử Dụng

### 1. Chạy chương trình
```bash
go run cmd/assistant/main.go
```

### 2. Tính năng hoạt động song song

Chương trình sẽ chạy **2 tác vụ song song**:

#### A. Voice Assistant (interactive)
- Nhấn **ENTER** lần 1 → Bắt đầu ghi âm
- Nhấn **ENTER** lần 2 → Dừng ghi âm và xử lý
- Hệ thống sẽ trả lời bằng giọng nói

#### B. Video Analysis (automatic)
- Tự động phân tích video **mỗi 10 giây**
- In ra màn hình mô tả về nội dung video
- Chạy ngầm trong background, không cần tương tác

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
👁️  Phân tích: Căn phòng có ánh sáng tự nhiên từ cửa sổ, một chiếc bàn làm việc ở góc phải với laptop và cốc cà phê. Không có người trong khung hình.
📺 Video: Căn phòng có ánh sáng tự nhiên từ cửa sổ, một chiếc bàn làm việc ở góc phải với laptop và cốc cà phê. Không có người trong khung hình.

[Sau 10 giây]
📸 Đang bắt frame từ RTSP stream...
🧠 Đang phân tích nội dung video...
👁️  Phân tích: Một người đang bước vào phòng từ cửa bên trái.
📺 Video: Một người đang bước vào phòng từ cửa bên trái.
```

## ⚙️ Tùy Chỉnh

### Thay đổi interval phân tích

Trong file `cmd/assistant/main.go`, tìm dòng:
```go
go func() {
    if err := assistant.StartContinuousVideoAnalysis(ctx, 10); err != nil {
        consoleLogger.Error("❌ Lỗi video analysis", err)
    }
}()
```

Thay `10` bằng số giây khác (ví dụ: `5` cho 5 giây, `30` cho 30 giây)

### Thay đổi prompt phân tích

Trong file `internal/infrastructure/video/rtsp_analyzer.go`, tìm dòng:
```go
prompt := "Mô tả ngắn gọn những gì bạn thấy trong video này. Hãy chỉ ra các đối tượng, hành động, và môi trường quan trọng."
```

Tùy chỉnh prompt theo nhu cầu của bạn.

### Thay đổi chất lượng/kích thước frame

Trong file `internal/infrastructure/video/rtsp_analyzer.go`, sửa:
```go
"-vf", "scale=1280:720", // Thay đổi resolution
"-q:v", "2",             // 1 (tốt nhất) đến 5 (thấp nhất)
```

## 🐛 Troubleshooting

### Lỗi: Cannot connect to RTSP stream

**Nguyên nhân**: 
- RTSP URL sai
- Camera không hỗ trợ RTSP
- Network không thể kết nối

**Giải pháp**:
1. Kiểm tra RTSP URL bằng VLC Media Player:
   - Mở VLC → Media → Open Network Stream
   - Dán RTSP URL và test

2. Kiểm tra camera có bật RTSP stream chưa
3. Ping địa chỉ IP camera: `ping 192.168.1.186`

### Lỗi: ffmpeg command not found

**Giải pháp**: Cài đặt ffmpeg
```bash
brew install ffmpeg
```

### Video analysis quá chậm

**Nguyên nhân**: Model AI xử lý chậm

**Giải pháp**:
- Tăng interval (từ 10s → 30s hoặc cao hơn)
- Giảm resolution frame (từ 1280x720 → 640x480)
- Sử dụng model AI nhẹ hơn (phi-3-mini thay vì llama3)

### Lỗi: Authentication failed

**Nguyên nhân**: Username/password sai hoặc chưa URL encode

**Giải pháp**:
- Kiểm tra username/password camera
- URL encode ký tự đặc biệt trong password
- Test với VLC trước

## 📝 Lưu Ý

1. **Performance**: Video analysis tiêu tốn tài nguyên CPU/GPU. Trên máy yếu, hãy tăng interval.

2. **Network**: RTSP stream cần kết nối mạng ổn định. Sử dụng dây LAN nếu có thể.

3. **Privacy**: Hệ thống chỉ xử lý local, không gửi video lên cloud.

4. **File cleanup**: Các frame tạm được tự động xóa sau khi phân tích.

## 🎯 Use Cases

- **Giám sát an ninh**: Phát hiện người lạ xuất hiện
- **Smart home**: Theo dõi hoạt động trong nhà
- **Giám sát trẻ em**: Theo dõi trẻ từ xa
- **Pet monitoring**: Theo dõi thú cưng
- **Văn phòng**: Theo dõi không gian làm việc

## 🔗 Tích Hợp Với Voice Assistant

Bạn có thể hỏi về video qua voice:
- "Có ai trong phòng không?"
- "Hiện tại trong video đang có gì?"
- "Mô tả cho tôi những gì bạn thấy"

Hệ thống sẽ capture frame mới nhất và phân tích khi bạn hỏi.
