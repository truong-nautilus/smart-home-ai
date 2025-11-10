# Hướng Dẫn Nhanh

## 🚀 Chạy Ứng Dụng

```bash
# Đảm bảo đã cài ffmpeg
brew install ffmpeg

# Thiết lập OpenAI API key trong file .env
echo 'OPENAI_API_KEY=sk-your-key-here' > .env

# Chạy ứng dụng
go run cmd/assistant/main.go
```

## 📂 Cấu Trúc Dự Án (Clean Architecture)

```
smart-home-ai/
├── cmd/assistant/main.go          # 🚪 Điểm vào
├── internal/
│   ├── domain/                    # 🎯 Logic nghiệp vụ cốt lõi
│   │   ├── entity.go             # Các model dữ liệu
│   │   └── repository.go         # Interface (ports)
│   ├── usecase/                   # 🔄 Logic ứng dụng
│   │   └── assistant.go          # Điều phối quy trình
│   └── infrastructure/            # 🔌 Adapter bên ngoài
│       ├── openai/               # OpenAI API client
│       └── media/                # FFmpeg wrapper
└── pkg/logger/                    # 📝 Tiện ích chia sẻ
```

## 🎯 Cách Hoạt Động

1. **Bắt ảnh** → Chụp ảnh từ camera
2. **Thu âm** → Thu 5 giây âm thanh
3. **Chuyển đổi** → Whisper API chuyển giọng nói thành văn bản
4. **Phân tích** → GPT-4o-mini hiểu hình ảnh + văn bản
5. **Tổng hợp** → TTS chuyển phản hồi thành giọng nói
6. **Phát** → Phát phản hồi âm thanh

## 🧩 Các Mẫu Thiết Kế Chính

- **Clean Architecture**: Tách biệt các mối quan tâm
- **Dependency Injection**: Kết nối trong `cmd/assistant/main.go`
- **Repository Pattern**: Interface trong `domain/repository.go`
- **Adapter Pattern**: Triển khai trong `infrastructure/`

## 🔧 Ví Dụ Tùy Chỉnh

### Thay Đổi Giọng TTS
Chỉnh sửa `internal/infrastructure/openai/client.go`:
```go
Voice: openai.VoiceNova, // hoặc: Alloy, Echo, Fable, Onyx, Shimmer
```

### Thay Đổi Thời Lượng Thu Âm
Chỉnh sửa `internal/usecase/assistant.go`:
```go
const thoiLuongThuAm = 10 // 10 giây thay vì 5
```

### Thêm Nhà Cung Cấp AI Mới
1. Tạo `internal/infrastructure/nhacungcap/client.go`
2. Triển khai `domain.BoNhanDienGiongNoi`, `domain.TroLyAI`, v.v.
3. Kết nối trong `cmd/assistant/main.go`

## 📚 Tài Liệu

- **[README.md](README.md)** - Hướng dẫn cài đặt đầy đủ
- **[KY_THUAT.md](KY_THUAT.md)** - Hướng dẫn kiến trúc chi tiết
- **[SO_DO_KY_THUAT.md](SO_DO_KY_THUAT.md)** - Sơ đồ kiến trúc trực quan
- **[TOM_TAT_TAI_CAU_TRUC.md](TOM_TAT_TAI_CAU_TRUC.md)** - Những gì đã thay đổi và tại sao

## ❓ Khắc Phục Sự Cố

### Không Tìm Thấy Camera/Mic
```bash
# Liệt kê các thiết bị có sẵn
ffmpeg -f avfoundation -list_devices true -i ""
```

### Lỗi Quyền Bị Từ Chối
Cho phép truy cập Camera & Microphone trong Cài đặt Hệ thống → Quyền riêng tư & Bảo mật

### Lỗi API Key
Đảm bảo file `.env` tồn tại và chứa:
```
OPENAI_API_KEY=sk-your-actual-key
```

## 🧪 Kiểm Thử (Tương Lai)

```bash
# Unit tests (sẽ được thêm)
go test ./internal/...

# Integration tests (sẽ được thêm)
go test ./test/integration/...
```

## 🎨 Lợi Ích So Với Code Cũ

| Tính Năng | Cũ (main.go) | Mới (Clean Arch) |
|---------|--------------|------------------|
| **Số dòng code** | 235 trong 1 file | Tổ chức theo layer |
| **Khả năng kiểm thử** | ❌ Khó | ✅ Dễ với mock |
| **Khả năng mở rộng** | ❌ Phải sửa main | ✅ Chỉ thêm adapter |
| **Khả năng đọc** | ⚠️ Lẫn lộn | ✅ Tách biệt rõ ràng |
| **Khả năng bảo trì** | ⚠️ Liên kết chặt | ✅ Liên kết lỏng |

---

Cần trợ giúp? Xem các file tài liệu ở trên! 📖
