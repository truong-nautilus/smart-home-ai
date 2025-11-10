# 🎙️ Trợ Lý AI Giọng Nói Vui Vẻ (Local)# 🎙️ Trợ Lý AI Giọng Nói Vui Vẻ (Local)



Trợ lý AI voice-only **phong cách thân thiện, vui vẻ, hài hước** chạy hoàn toàn local trên macOS:Trợ lý AI voice-only **phong cách thân thiện, vui vẻ, hài hước** chạy hoàn toàn local trên macOS:

- 🎯 **PhoWhisper** (vinai/PhoWhisper-small 157M) - Nhận dạng tiếng Việt tối ưu- 🎯 **PhoWhisper** (vinai/PhoWhisper-small 157M) - Nhận dạng tiếng Việt tối ưu

- 🤖 **Ollama** (gemma2:2b) - AI trả lời vui nhộn, tự nhiên- 🤖 **Ollama** (gemma2:2b) - AI trả lời vui nhộn, tự nhiên

- 🗣️ **Edge TTS** (Microsoft Neural) - Giọng nói như người thật- 🗣️ **Edge TTS** (Microsoft Neural) - Giọng nói như người thật



## ✨ Tính Năng## ✨ Tính Năng



### 🎤 Cách Sử Dụng### 🎤 Cách Sử Dụng

1. **Nhấn ENTER** → Bắt đầu ghi âm1. **Nhấn ENTER** → Bắt đầu ghi âm

2. **Nói bất cứ gì** → Chat tự nhiên, hỏi chuyện, điều khiển nhà thông minh, kể chuyện cười...2. **Nói bất cứ gì** → Chat tự nhiên, hỏi chuyện, điều khiển nhà thông minh, kể chuyện cười...

3. **Nhấn ENTER** lần nữa → Dừng và xử lý3. **Nhấn ENTER** lần nữa → Dừng và xử lý

4. **Nghe phản hồi** → AI trả lời vui vẻ, hài hước4. **Nghe phản hồi** → AI trả lời vui vẻ, hài hước



### 🎯 Phong Cách Giao Tiếp### 🎯 Phong Cách Giao Tiếp

- ✅ **Thân thiện, gần gũi** - Như bạn bè tâm sự- ✅ **Thân thiện, gần gũi** - Như bạn bè tâm sự

- ✅ **Vui vẻ, hài hước** - Có thể đùa cợt nhẹ nhàng- ✅ **Vui vẻ, hài hước** - Có thể đùa cợt nhẹ nhàng

- ✅ **Tự nhiên** - Không cứng nhắc, không formal- ✅ **Tự nhiên** - Không cứng nhắc, không formal

- ✅ **Đa chủ đề** - Smart home, thời tiết, tâm sự, giải trí...- ✅ **Đa chủ đề** - Smart home, thời tiết, tâm sự, giải trí...



### ⚙️ Quy Trình Xử Lý### ⚙️ Quy Trình Xử Lý

1. **STT** → PhoWhisper (độ chính xác cao nhất cho tiếng Việt)1. **STT** → PhoWhisper (độ chính xác cao nhất cho tiếng Việt)

2. **AI** → Ollama phân tích và trả lời vui vẻ2. **AI** → Ollama phân tích và trả lời vui vẻ

3. **TTS** → Edge TTS tổng hợp giọng tự nhiên3. **TTS** → Edge TTS tổng hợp giọng tự nhiên

4. **Play** → Phát âm thanh phản hồi4. **Play** → Phát âm thanh phản hồi



## 🛠️ Stack Công Nghệ (100% Local)## 🛠️ Stack Công Nghệ (100% Local)



- **Language**: Go 1.22+- **Language**: Go 1.22+

- **Audio I/O**: FFmpeg AVFoundation- **Audio I/O**: FFmpeg AVFoundation

- **STT**: PhoWhisper **vinai/PhoWhisper-small** (157M) - tối ưu cho tiếng Việt- **STT**: PhoWhisper **vinai/PhoWhisper-small** (157M) - tối ưu cho tiếng Việt

- **AI**: Ollama gemma2:2b (1.6GB) - phong cách vui vẻ- **AI**: Ollama gemma2:2b (1.6GB) - phong cách vui vẻ

- **TTS**: Edge TTS Neural Voices (vi-VN-NamMinhNeural)- **TTS**: Edge TTS Neural Voices (vi-VN-NamMinhNeural)



## Yêu Cầu Hệ Thống## Yêu Cầu Hệ Thống



- macOS (Apple Silicon hoặc Intel)- macOS (Apple Silicon hoặc Intel)

- Go 1.22+- Go 1.22+

- FFmpeg với hỗ trợ AVFoundation- FFmpeg với hỗ trợ AVFoundation

- Ollama CLI- Ollama CLI

- Python 3.9+ với PhoWhisper- Python 3.9+ với PhoWhisper

- Edge TTS (Python package)- Edge TTS (Python package)



## 📦 Cài Đặt



### 1. Cài đặt Homebrew (nếu chưa có)## Cài đặt



```bash### 1. Cài đặt Homebrew (nếu chưa có)### 1. Install Homebrew (if not already installed)

/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

```### 1. Cài đặt Go



### 2. Cài đặt Go và FFmpeg```bash```bash



```bash```bash

brew install go ffmpeg

```brew install go/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"



### 3. Cài đặt Ollama```



```bash``````

# Cài đặt Ollama CLI

brew install ollama### 2. Cài đặt FFmpeg



# Khởi động Ollama service

ollama serve &

```bash

# Pull model gemma2:2b (hoặc model nhẹ khác)

ollama pull gemma2:2bbrew install ffmpeg### 2. Cài đặt ffmpeg và ffplay### 2. Install ffmpeg and ffplay

```

```

**Các model nhẹ được đề xuất:**

- `gemma2:2b` (2B parameters) - nhỏ gọn, hiệu suất tốt```bash```bash

- `phi-3-mini` (3.8B parameters) - cân bằng tốc độ và chất lượng

- `tinyllama` (1.1B parameters) - rất nhẹ, phù hợp máy yếu### 3. Cài đặt Ollama



### 4. Cài đặt PhoWhisperbrew install ffmpegbrew install ffmpeg



```bash```bash

# Tạo virtual environment cho Python

python3 -m venv phowhisper-env# Cài đặt Ollama CLI``````

source phowhisper-env/bin/activate

brew install ollama

# Cài đặt dependencies

pip install -r requirements.txt



# Model vinai/PhoWhisper-small (~157MB) sẽ tự động download khi chạy lần đầu# Khởi động Ollama service

```

ollama serve &### 3. Cài đặt Go 1.22+ (nếu chưa có)### 3. Install Go 1.22+ (if not already installed)

**Lưu ý quan trọng:** Cập nhật shebang trong `scripts/phowhisper_transcribe.py` với đường dẫn Python của bạn:

```python

#!/path/to/your/phowhisper-env/bin/python3

```# Pull model phi-3-mini (hoặc model nhẹ khác)```bash```bash



Ví dụ:ollama pull phi-3-mini

```python

#!/Users/phamthetruong/phowhisper-env/bin/python3```brew install gobrew install go

```



### 5. Cài đặt Edge TTS

Các model nhẹ được đề xuất:``````

```bash

# Kích hoạt Python virtual environment (nếu chưa)- `phi-3-mini` (3.8B parameters) - cân bằng tốc độ và chất lượng

source phowhisper-env/bin/activate

- `tinyllama` (1.1B parameters) - rất nhẹ, phù hợp máy yếu

# Cài đặt Edge TTS

pip install edge-tts- `gemma:2b` (2B parameters) - nhỏ gọn, hiệu suất tốt

```

### 4. Lấy OpenAI API Key### 4. Get an OpenAI API Key

### 6. Cấp quyền cho Terminal

### 4. Build whisper.cpp

Cấp quyền cho Terminal/iTerm2 truy cập microphone trong **System Preferences → Privacy & Security → Microphone**.

- Truy cập https://platform.openai.com/api-keys- Visit https://platform.openai.com/api-keys

### 7. Clone và build project

```bash

```bash

# Clone repository# Clone repository- Tạo API key mới- Create a new API key

git clone https://github.com/truong-nautilus/smart-home-ai.git

cd smart-home-aigit clone https://github.com/ggerganov/whisper.cpp.git



# Tải dependenciescd whisper.cpp- Thiết lập biến môi trường:- Set it as an environment variable:

go mod download



# Build

go build -o smart-home-ai cmd/assistant/main.go# Build

```

make

## ⚙️ Cấu Hình

```bash```bash

Tạo file `.env` trong thư mục gốc của project:

# Download model nhỏ (base hoặc tiny)

```env

# PhoWhisper configurationbash ./models/download-ggml-model.sh baseexport OPENAI_API_KEY="sk-your-api-key-here"export OPENAI_API_KEY="sk-your-api-key-here"

PHOWHISPER_SCRIPT=scripts/phowhisper_transcribe.py

# Hoặc model nhỏ hơn:

# Ollama configuration

OLLAMA_MODEL=gemma2:2b# bash ./models/download-ggml-model.sh tiny``````



# Edge TTS configuration

EDGE_TTS_VOICE=vi-VN-NamMinhNeural

# Giọng khác: vi-VN-HoaiMyNeural (nữ)# Lưu đường dẫn binary và model

```

# Binary: ./main

## ▶️ Chạy Ứng Dụng

# Model: ./models/ggml-base.bin**Mẹo:** Thêm lệnh export vào `~/.zshrc` để lưu vĩnh viễn:**Tip:** Add the export command to your `~/.zshrc` to persist it:

```bash

# Chạy trực tiếp```

./smart-home-ai

```bash```bash

# Hoặc dùng go run

go run cmd/assistant/main.go### 5. Clone và build project

```

echo 'export OPENAI_API_KEY="sk-your-api-key-here"' >> ~/.zshrcecho 'export OPENAI_API_KEY="sk-your-api-key-here"' >> ~/.zshrc

### Kết Quả Mong Đợi:

```bash

```

🚀 Trợ lý AI đã sẵn sàng!# Clone repositorysource ~/.zshrcsource ~/.zshrc

📌 Cách dùng: Nhấn ENTER lần 1 → ghi âm → nhấn ENTER lần 2 → xử lý

🛑 Nhấn Ctrl+C để thoátgit clone <repository-url>

🎤 Đang ghi âm... (nhấn ENTER để dừng)

✅ Đã ghi âm xongcd smart-home-ai``````

🧠 Đang chuyển giọng nói thành văn bản...

[🔍 PhoWhisper output: "Hôm nay thời tiết thế nào"]

📝 Câu hỏi: "Hôm nay thời tiết thế nào"

🤖 Đang xử lý câu hỏi...# Tải dependencies

💬 Trả lời: "Chào bạn! Hôm nay thời tiết đẹp lắm nhé..."

🔊 Đang tổng hợp giọng nói...go mod download

✅ Tổng hợp giọng nói thành công

🔈 Đang phát âm thanh phản hồi...Hoặc tạo file `.env`:## 🚀 Setup

✅ Hoàn thành!

✨ Hoàn tất! Sẵn sàng cho lần tiếp theo...# Build

```

go build -o smart-home-ai cmd/assistant/main.go```bash

## 🛠️ Troubleshooting

```

### Lỗi: "ffmpeg: No such file or directory"

echo 'OPENAI_API_KEY=sk-your-api-key-here' > .env1. **Clone or navigate to the project directory:**

```bash

# Kiểm tra ffmpeg đã cài đặt chưa## Cấu hình

which ffmpeg

``````bash

# Cài đặt nếu chưa có

brew install ffmpegTạo file `.env` trong thư mục gốc của project:

```

cd /Users/phamthetruong/github/smart-home-ai

### Lỗi: "ollama: command not found"

```env

```bash

# Cài đặt Ollama# Whisper.cpp configuration## 🚀 Cài Đặt```

brew install ollama

WHISPER_CPP_BIN=/path/to/whisper.cpp/main

# Khởi động service

ollama serve &WHISPER_CPP_MODEL=/path/to/whisper.cpp/models/ggml-base.bin



# Pull model

ollama pull gemma2:2b

```# Ollama configuration1. **Di chuyển đến thư mục dự án:**2. **Download dependencies:**



### Lỗi: "PhoWhisper không nhận diện được văn bản"OLLAMA_MODEL=phi-3-mini



Kiểm tra lại:```bash```bash

1. Đường dẫn shebang trong `scripts/phowhisper_transcribe.py`

2. Virtual environment đã cài đủ dependencies: `pip install -r requirements.txt`# MacTTS configuration (optional)

3. Quyền thực thi: `chmod +x scripts/phowhisper_transcribe.py`

MACTTS_VOICE=Alexcd /Users/phamthetruong/github/smart-home-aigo mod tidy

### Lỗi: "avfoundation: Cannot find device"

```

Cấp quyền cho Terminal/iTerm2 truy cập microphone trong **System Preferences → Privacy & Security**.

``````

Liệt kê các thiết bị có sẵn:

**Các giọng nói macOS khả dụng:**

```bash

ffmpeg -f avfoundation -list_devices true -i ""- `Alex` (Mỹ, nam)

```

- `Samantha` (Mỹ, nữ)

Kết quả sẽ hiển thị:

```- `Daniel` (Anh, nam)2. **Tải các dependency:**## ▶️ Run

[:0] Built-in Microphone

[:1] External Microphone- `Victoria` (Anh, nữ)

```

- Xem danh sách đầy đủ: `say -v ?````bash

## 📦 Dependencies



- **Go 1.22+**

- **github.com/sashabaranov/go-openai** v1.20.4 - OpenAI Go client (cho tương thích)## Chạy ứng dụnggo mod tidy**Note:** The project now uses Clean Architecture. See [ARCHITECTURE.md](ARCHITECTURE.md) for details.

- **github.com/joho/godotenv** v1.5.1 - Tải biến môi trường từ .env

- **ffmpeg** - Bắt và phát âm thanh/video

- **Python Libraries:**

  - torch >= 2.0.0```bash```

  - transformers >= 4.30.0

  - librosa >= 0.10.0# Chạy trực tiếp

  - edge-tts

./smart-home-ai```bash

## 📁 Cấu Trúc Dự Án



```

smart-home-ai/# Hoặc dùng go run## ▶️ Chạy Ứng Dụng# Run directly

├── cmd/assistant/              # Điểm vào ứng dụng

│   └── main.go                 # Dependency injectiongo run cmd/assistant/main.go

├── internal/

│   ├── domain/                 # Logic nghiệp vụ cốt lõi```go run cmd/assistant/main.go

│   │   ├── entity.go           # Các entity nghiệp vụ

│   │   └── repository.go       # Interface (ports)

│   ├── usecase/                # Logic ứng dụng

│   │   └── assistant.go        # Điều phối quy trình## Luồng hoạt động**Lưu ý:** Dự án sử dụng Clean Architecture. Xem [KY_THUAT.md](KY_THUAT.md) để biết chi tiết.

│   └── infrastructure/         # Adapter dịch vụ bên ngoài

│       ├── phowhisper/         # PhoWhisper recognizer

│       │   └── recognizer.go

│       ├── ollama/             # Ollama AI client1. **Chụp ảnh từ camera** - FFmpeg capture một frame từ camera tích hợp# Or build binary

│       │   └── client.go

│       ├── edgetts/            # Edge TTS synthesizer2. **Ghi âm câu hỏi** - FFmpeg ghi âm 5 giây từ microphone

│       │   └── synthesizer.go

│       ├── media/              # FFmpeg wrapper3. **Chuyển giọng nói thành text** - whisper.cpp transcribe file audio```bashgo build -o smart-home-ai cmd/assistant/main.go

│       │   └── ffmpeg.go

│       └── keyboard/           # Keyboard listener4. **Phân tích đa phương thức** - Ollama xử lý text + ảnh

│           └── listener.go

├── scripts/5. **Phát giọng nói phản hồi** - macOS `say` đọc câu trả lời# Chạy trực tiếp./smart-home-ai

│   └── phowhisper_transcribe.py  # Python script cho PhoWhisper

├── pkg/logger/                 # Tiện ích chia sẻ

│   └── console.go

├── requirements.txt            # Python dependencies## Cấu trúc thư mụcgo run cmd/assistant/main.go```

├── .env                        # Biến môi trường

├── go.mod                      # Định nghĩa Go module

└── README.md

``````



## 🧹 Dọn Dẹpsmart-home-ai/



Các file tạm (`audio.wav`, `reply.mp3`) được tự động dọn dẹp sau mỗi lần chạy.├── cmd/# Hoặc build file thực thi### Expected Output:



## 📝 Ghi Chú│   └── assistant/



- Thời lượng thu âm: **Nhấn và giữ ENTER** để bắt đầu, **nhả ENTER** để dừng│       └── main.go              # Entry pointgo build -o tro-ly-thong-minh cmd/assistant/main.go```

- Âm thanh thu ở **mono với tần số 16kHz** (tối ưu cho PhoWhisper)

- Edge TTS sử dụng giọng **vi-VN-NamMinhNeural** (nam) hoặc **vi-VN-HoaiMyNeural** (nữ)├── internal/



## 🎨 Tùy Chỉnh│   ├── domain/./tro-ly-thong-minh[15:04:05] 🎥 Capturing image from FaceTime HD camera...



Chỉnh sửa các file để tùy chỉnh:│   │   ├── entity.go            # Domain entities

- `cmd/assistant/main.go` - Chọn model Ollama khác

- `.env` - Thay đổi giọng Edge TTS (vi-VN-HoaiMyNeural cho giọng nữ)│   │   └── repository.go        # Repository interfaces (ports)```[15:04:06] ✅ Image captured successfully

- `scripts/phowhisper_transcribe.py` - Thay đổi batch_size, chunk_length_s để tối ưu performance

│   ├── usecase/

### So sánh PhoWhisper vs Whisper.cpp

│   │   └── assistant.go         # Business logic orchestration[15:04:06] 🎤 Recording audio from microphone (5 seconds)...

| Tiêu chí | PhoWhisper | Whisper.cpp |

|----------|-----------|-------------|│   └── infrastructure/

| **Model Size** | 157MB | 1.5GB (medium) |

| **Độ chính xác (Tiếng Việt)** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |│       ├── media/### Kết Quả Mong Đợi:[15:04:11] ✅ Audio recorded successfully

| **Tốc độ** | Trung bình | Nhanh |

| **Tối ưu cho** | Tiếng Việt | Đa ngôn ngữ |│       │   └── ffmpeg.go        # FFmpeg adapter

| **Setup** | Dễ (pip install) | Phức tạp (build cmake) |

│       ├── whispercpp/```[15:04:11] 🧠 Transcribing speech with Whisper API...

## 🏗️ Kiến Trúc

│       │   └── recognizer.go    # Whisper.cpp adapter

Dự án sử dụng **Clean Architecture** với:

- **Domain Layer**: Entity và interface cốt lõi│       ├── ollama/[15:04:05] 🎥 Đang bắt ảnh từ camera FaceTime HD...[15:04:12] 📝 Transcription: "What do you see in front of me?"

- **Use Case Layer**: Logic điều phối nghiệp vụ

- **Infrastructure Layer**: Triển khai adapter│       │   └── client.go        # Ollama adapter

- **Delivery Layer**: Điểm vào và dependency injection

│       └── mactts/[15:04:06] ✅ Đã bắt ảnh thành công[15:04:12] 🤖 GPT-4o-mini reasoning on text + image...

## 📚 Tài Liệu Bổ Sung

│           └── synthesizer.go   # MacTTS adapter

- **[ARCHITECTURE.md](docs/ARCHITECTURE.md)** - Kiến trúc chi tiết

- **[QUICKSTART.md](docs/QUICKSTART.md)** - Hướng dẫn nhanh├── pkg/[15:04:06] 🎤 Đang thu âm từ microphone (5 giây)...[15:04:14] 💬 GPT Response: "I can see you're sitting at a desk with a laptop..."



## 📄 Giấy Phép│   └── logger/



MIT│       └── console.go           # Console logger[15:04:11] ✅ Đã thu âm thành công[15:04:14] 🔊 Converting response to speech with TTS...



---├── go.mod



**Chúc bạn lập trình vui vẻ!** 🎯├── go.sum[15:04:11] 🧠 Đang chuyển giọng nói thành văn bản với Whisper API...[15:04:16] ✅ Speech generated successfully


├── .env                         # Configuration (tạo mới)

└── README.md[15:04:12] 📝 Văn bản: "Bạn thấy gì trước mặt tôi?"[15:04:16] 🔈 Playing audio response...

```

[15:04:12] 🤖 GPT-4o-mini đang phân tích văn bản + hình ảnh...[15:04:20] ✅ Done!

## Troubleshooting

[15:04:14] 💬 Phản hồi GPT: "Tôi thấy bạn đang ngồi ở bàn làm việc với laptop..."```

### Lỗi: "ffmpeg: No such file or directory"

[15:04:14] 🔊 Đang chuyển phản hồi thành giọng nói với TTS...

```bash

# Kiểm tra ffmpeg đã cài đặt chưa[15:04:16] ✅ Đã tạo giọng nói thành công## 🔐 Permissions

which ffmpeg

[15:04:16] 🔈 Đang phát âm thanh phản hồi...

# Cài đặt nếu chưa có

brew install ffmpeg[15:04:20] ✅ Hoàn thành!On first run, macOS will ask for permissions:

```

```- **Camera access** - Required to capture images

### Lỗi: "ollama: command not found"

- **Microphone access** - Required to record audio

```bash

# Cài đặt Ollama## 🔐 Quyền Truy Cập

brew install ollama

Click "Allow" when prompted.

# Khởi động service

ollama serve &Lần chạy đầu tiên, macOS sẽ yêu cầu quyền:



# Pull model- **Truy cập Camera** - Cần thiết để bắt ảnh## 🛠️ Troubleshooting

ollama pull phi-3-mini

```- **Truy cập Microphone** - Cần thiết để thu âm



### Lỗi: "whisper.cpp: không tìm thấy binary"### Camera/Microphone Not Found



Kiểm tra lại đường dẫn trong file `.env`:Nhấp "Cho phép" khi được nhắc.If you get an error about device not found, list available devices:



```env```bash

WHISPER_CPP_BIN=/Users/yourusername/whisper.cpp/main

WHISPER_CPP_MODEL=/Users/yourusername/whisper.cpp/models/ggml-base.bin## 🛠️ Khắc Phục Sự Cốffmpeg -f avfoundation -list_devices true -i ""

```

```

### Lỗi: "avfoundation: Cannot find device"

### Không Tìm Thấy Camera/Microphone

Cấp quyền cho Terminal/iTerm2 truy cập camera và microphone trong **System Preferences → Privacy & Security**.

Nếu gặp lỗi về thiết bị không tìm thấy, liệt kê các thiết bị có sẵn:This will show something like:

## License

```bash```

MIT

ffmpeg -f avfoundation -list_devices true -i ""[0] FaceTime HD Camera

## Contributing

```[1] External Camera

Pull requests are welcome!

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
