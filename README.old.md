# Trợ Lý Thông Minh Nhà Thông Minh# Smart Home AI Assistant



Một trợ lý AI cho macOS (Apple Silicon) sử dụng camera và microphone tích hợp của MacBook để tạo trải nghiệm AI đa phương thức.An AI assistant for macOS (Apple Silicon) that uses your MacBook's built-in camera and microphone to create a multimodal AI experience.



## 🎯 Chức Năng## 🎯 What It Does



1. **Bắt ảnh** từ camera FaceTime HD1. **Captures** a single frame from your FaceTime HD camera

2. **Thu âm** 5 giây từ microphone tích hợp2. **Records** 5 seconds of audio from your built-in microphone

3. **Chuyển đổi giọng nói** thành văn bản bằng OpenAI Whisper API3. **Transcribes** your speech using OpenAI Whisper API

4. **Phân tích** hình ảnh và văn bản cùng nhau bằng GPT-4o-mini4. **Analyzes** the image and text together using GPT-4o-mini

5. **Tổng hợp** phản hồi của AI thành giọng nói bằng TTS5. **Synthesizes** the AI's response to speech using TTS

6. **Phát** phản hồi âm thanh cho bạn6. **Plays** the audio response back to you



## 🧰 Yêu Cầu## 🧰 Prerequisites



### 1. Cài đặt Homebrew (nếu chưa có)### 1. Install Homebrew (if not already installed)

```bash```bash

/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

``````



### 2. Cài đặt ffmpeg và ffplay### 2. Install ffmpeg and ffplay

```bash```bash

brew install ffmpegbrew install ffmpeg

``````



### 3. Cài đặt Go 1.22+ (nếu chưa có)### 3. Install Go 1.22+ (if not already installed)

```bash```bash

brew install gobrew install go

``````



### 4. Lấy OpenAI API Key### 4. Get an OpenAI API Key

- Truy cập https://platform.openai.com/api-keys- Visit https://platform.openai.com/api-keys

- Tạo API key mới- Create a new API key

- Thiết lập biến môi trường:- Set it as an environment variable:



```bash```bash

export OPENAI_API_KEY="sk-your-api-key-here"export OPENAI_API_KEY="sk-your-api-key-here"

``````



**Mẹo:** Thêm lệnh export vào `~/.zshrc` để lưu vĩnh viễn:**Tip:** Add the export command to your `~/.zshrc` to persist it:

```bash```bash

echo 'export OPENAI_API_KEY="sk-your-api-key-here"' >> ~/.zshrcecho 'export OPENAI_API_KEY="sk-your-api-key-here"' >> ~/.zshrc

source ~/.zshrcsource ~/.zshrc

``````



Hoặc tạo file `.env`:## 🚀 Setup

```bash

echo 'OPENAI_API_KEY=sk-your-api-key-here' > .env1. **Clone or navigate to the project directory:**

``````bash

cd /Users/phamthetruong/github/smart-home-ai

## 🚀 Cài Đặt```



1. **Di chuyển đến thư mục dự án:**2. **Download dependencies:**

```bash```bash

cd /Users/phamthetruong/github/smart-home-aigo mod tidy

``````



2. **Tải các dependency:**## ▶️ Run

```bash

go mod tidy**Note:** The project now uses Clean Architecture. See [ARCHITECTURE.md](ARCHITECTURE.md) for details.

```

```bash

## ▶️ Chạy Ứng Dụng# Run directly

go run cmd/assistant/main.go

**Lưu ý:** Dự án sử dụng Clean Architecture. Xem [KY_THUAT.md](KY_THUAT.md) để biết chi tiết.

# Or build binary

```bashgo build -o smart-home-ai cmd/assistant/main.go

# Chạy trực tiếp./smart-home-ai

go run cmd/assistant/main.go```



# Hoặc build file thực thi### Expected Output:

go build -o tro-ly-thong-minh cmd/assistant/main.go```

./tro-ly-thong-minh[15:04:05] 🎥 Capturing image from FaceTime HD camera...

```[15:04:06] ✅ Image captured successfully

[15:04:06] 🎤 Recording audio from microphone (5 seconds)...

### Kết Quả Mong Đợi:[15:04:11] ✅ Audio recorded successfully

```[15:04:11] 🧠 Transcribing speech with Whisper API...

[15:04:05] 🎥 Đang bắt ảnh từ camera FaceTime HD...[15:04:12] 📝 Transcription: "What do you see in front of me?"

[15:04:06] ✅ Đã bắt ảnh thành công[15:04:12] 🤖 GPT-4o-mini reasoning on text + image...

[15:04:06] 🎤 Đang thu âm từ microphone (5 giây)...[15:04:14] 💬 GPT Response: "I can see you're sitting at a desk with a laptop..."

[15:04:11] ✅ Đã thu âm thành công[15:04:14] 🔊 Converting response to speech with TTS...

[15:04:11] 🧠 Đang chuyển giọng nói thành văn bản với Whisper API...[15:04:16] ✅ Speech generated successfully

[15:04:12] 📝 Văn bản: "Bạn thấy gì trước mặt tôi?"[15:04:16] 🔈 Playing audio response...

[15:04:12] 🤖 GPT-4o-mini đang phân tích văn bản + hình ảnh...[15:04:20] ✅ Done!

[15:04:14] 💬 Phản hồi GPT: "Tôi thấy bạn đang ngồi ở bàn làm việc với laptop..."```

[15:04:14] 🔊 Đang chuyển phản hồi thành giọng nói với TTS...

[15:04:16] ✅ Đã tạo giọng nói thành công## 🔐 Permissions

[15:04:16] 🔈 Đang phát âm thanh phản hồi...

[15:04:20] ✅ Hoàn thành!On first run, macOS will ask for permissions:

```- **Camera access** - Required to capture images

- **Microphone access** - Required to record audio

## 🔐 Quyền Truy Cập

Click "Allow" when prompted.

Lần chạy đầu tiên, macOS sẽ yêu cầu quyền:

- **Truy cập Camera** - Cần thiết để bắt ảnh## 🛠️ Troubleshooting

- **Truy cập Microphone** - Cần thiết để thu âm

### Camera/Microphone Not Found

Nhấp "Cho phép" khi được nhắc.If you get an error about device not found, list available devices:

```bash

## 🛠️ Khắc Phục Sự Cốffmpeg -f avfoundation -list_devices true -i ""

```

### Không Tìm Thấy Camera/Microphone

Nếu gặp lỗi về thiết bị không tìm thấy, liệt kê các thiết bị có sẵn:This will show something like:

```bash```

ffmpeg -f avfoundation -list_devices true -i ""[0] FaceTime HD Camera

```[1] External Camera

[:0] Built-in Microphone

Kết quả sẽ hiển thị:[:1] External Microphone

``````

[0] FaceTime HD Camera

[1] External CameraUpdate `main.go` if needed:

[:0] Built-in Microphone- For video: change `"-i", "0"` to your camera index

[:1] External Microphone- For audio: change `"-i", ":0"` to your microphone index

```

### API Key Not Set

Cập nhật `internal/infrastructure/media/ffmpeg.go` nếu cần:If you see `OPENAI_API_KEY environment variable not set`:

- Với video: thay đổi `"-i", "0"` thành chỉ số camera của bạn```bash

- Với audio: thay đổi `"-i", ":0"` thành chỉ số microphone của bạnexport OPENAI_API_KEY="your-key-here"

```

### Chưa Thiết Lập API Key

Nếu thấy `Chưa thiết lập biến môi trường OPENAI_API_KEY`:### FFmpeg Not Installed

```bashIf you get `executable file not found`:

export OPENAI_API_KEY="your-key-here"```bash

```brew install ffmpeg

```

Hoặc tạo file `.env`:

```bash## 📦 Dependencies

echo 'OPENAI_API_KEY=sk-your-key-here' > .env

```- **Go 1.22+**

- **github.com/sashabaranov/go-openai** v1.20.4 - OpenAI Go client

### Chưa Cài FFmpeg- **ffmpeg** - Audio/video capture and playback

Nếu gặp lỗi `executable file not found`:- **OpenAI APIs:**

```bash  - Whisper (speech-to-text)

brew install ffmpeg  - GPT-4o-mini (multimodal reasoning)

```  - TTS-1 (text-to-speech)



## 📦 Các Dependency## 🧹 Cleanup



- **Go 1.22+**Temporary files (`frame.jpg`, `audio.wav`, `reply.mp3`) are automatically cleaned up after each run.

- **github.com/sashabaranov/go-openai** v1.20.4 - OpenAI Go client

- **github.com/joho/godotenv** v1.5.1 - Tải biến môi trường từ .env## 📝 Notes

- **ffmpeg** - Bắt và phát âm thanh/video

- **OpenAI APIs:**- The assistant records **5 seconds** of audio by default (configurable in `main.go`)

  - Whisper (chuyển giọng nói thành văn bản)- Camera captures a **640x480** frame (configurable)

  - GPT-4o-mini (phân tích đa phương thức)- Audio is recorded in **mono at 16kHz** (optimal for Whisper)

  - TTS-1 (chuyển văn bản thành giọng nói)- TTS uses the **Alloy voice** (configurable to: alloy, echo, fable, onyx, nova, shimmer)



## 📁 Cấu Trúc Dự Án## 🎨 Customization



```Edit `main.go` to customize:

smart-home-ai/- `audioDuration` - Change recording length

├── cmd/assistant/              # Điểm vào ứng dụng- `Voice` - Change TTS voice in `textToSpeech()`

│   └── main.go                 # Dependency injection- `video_size` - Change camera resolution in `captureImage()`

├── internal/- `Model` - Switch between GPT models

│   ├── domain/                 # Logic nghiệp vụ cốt lõi

│   │   ├── entity.go           # Các entity nghiệp vụ## 📄 License

│   │   └── repository.go       # Interface (ports)

│   ├── usecase/                # Logic ứng dụngMIT

│   │   └── assistant.go        # Điều phối quy trình
│   └── infrastructure/         # Adapter dịch vụ bên ngoài
│       ├── openai/             # OpenAI API client
│       │   └── client.go
│       └── media/              # FFmpeg wrapper
│           └── ffmpeg.go
├── pkg/logger/                 # Tiện ích chia sẻ
│   └── console.go
├── .env                        # Biến môi trường
├── .gitignore                  # Quy tắc ignore Git
├── go.mod                      # Định nghĩa Go module
└── Các file tài liệu
```

## 🧹 Dọn Dẹp

Các file tạm (`hinh-anh.jpg`, `am-thanh.wav`, `tra-loi.mp3`) được tự động dọn dẹp sau mỗi lần chạy.

## 📝 Ghi Chú

- Trợ lý thu âm **5 giây** mặc định (có thể cấu hình trong `internal/usecase/assistant.go`)
- Camera bắt ảnh **640x480** (có thể cấu hình)
- Âm thanh thu ở **mono với tần số 16kHz** (tối ưu cho Whisper)
- TTS sử dụng giọng **Alloy** (có thể chọn: alloy, echo, fable, onyx, nova, shimmer)

## 🎨 Tùy Chỉnh

Chỉnh sửa các file để tùy chỉnh:
- `internal/usecase/assistant.go` - Thay đổi thời lượng thu âm (`thoiLuongThuAm`)
- `internal/infrastructure/openai/client.go` - Thay đổi giọng TTS trong phương thức `TongHop()`
- `internal/infrastructure/media/ffmpeg.go` - Thay đổi độ phân giải camera trong `BatAnh()`
- Chuyển đổi giữa các model GPT khác nhau

## 🏗️ Kiến Trúc

Dự án sử dụng **Clean Architecture** với:
- **Domain Layer**: Entity và interface cốt lõi
- **Use Case Layer**: Logic điều phối nghiệp vụ
- **Infrastructure Layer**: Triển khai adapter
- **Delivery Layer**: Điểm vào và dependency injection

Xem [KY_THUAT.md](KY_THUAT.md) để biết chi tiết về kiến trúc.

## 📚 Tài Liệu Bổ Sung

- **[HUONG_DAN_NHANH.md](HUONG_DAN_NHANH.md)** - Hướng dẫn nhanh 2 phút
- **[KY_THUAT.md](KY_THUAT.md)** - Tìm hiểu sâu về thiết kế
- **[SO_DO_KY_THUAT.md](SO_DO_KY_THUAT.md)** - Hướng dẫn trực quan
- **[TOM_TAT_TAI_CAU_TRUC.md](TOM_TAT_TAI_CAU_TRUC.md)** - Những gì đã thay đổi
- **[README_EN.md](README_EN.md)** - English version

## 📄 Giấy Phép

MIT

---

**Chúc bạn lập trình vui vẻ!** 🎯
