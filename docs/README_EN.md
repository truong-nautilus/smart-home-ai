# Smart Home AI Assistant

An AI assistant for macOS (Apple Silicon) that uses your MacBook's built-in camera and microphone to create a multimodal AI experience.

## 🎯 What It Does

1. **Captures** a single frame from your FaceTime HD camera
2. **Records** 5 seconds of audio from your built-in microphone
3. **Transcribes** your speech using OpenAI Whisper API
4. **Analyzes** the image and text together using GPT-4o-mini
5. **Synthesizes** the AI's response to speech using TTS
6. **Plays** the audio response back to you

## 🧰 Prerequisites

### 1. Install Homebrew (if not already installed)
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

### 2. Install ffmpeg and ffplay
```bash
brew install ffmpeg
```

### 3. Install Go 1.22+ (if not already installed)
```bash
brew install go
```

### 4. Get an OpenAI API Key
- Visit https://platform.openai.com/api-keys
- Create a new API key
- Set it as an environment variable:

```bash
export OPENAI_API_KEY="sk-your-api-key-here"
```

**Tip:** Add the export command to your `~/.zshrc` to persist it:
```bash
echo 'export OPENAI_API_KEY="sk-your-api-key-here"' >> ~/.zshrc
source ~/.zshrc
```

## 🚀 Setup

1. **Clone or navigate to the project directory:**
```bash
cd /Users/phamthetruong/github/smart-home-ai
```

2. **Download dependencies:**
```bash
go mod tidy
```

## ▶️ Run

**Note:** The project now uses Clean Architecture. See [ARCHITECTURE.md](ARCHITECTURE.md) for details.

```bash
# Run directly
go run cmd/assistant/main.go

# Or build binary
go build -o smart-home-ai cmd/assistant/main.go
./smart-home-ai
```

### Expected Output:
```
[15:04:05] 🎥 Capturing image from FaceTime HD camera...
[15:04:06] ✅ Image captured successfully
[15:04:06] 🎤 Recording audio from microphone (5 seconds)...
[15:04:11] ✅ Audio recorded successfully
[15:04:11] 🧠 Transcribing speech with Whisper API...
[15:04:12] 📝 Transcription: "What do you see in front of me?"
[15:04:12] 🤖 GPT-4o-mini reasoning on text + image...
[15:04:14] 💬 GPT Response: "I can see you're sitting at a desk with a laptop..."
[15:04:14] 🔊 Converting response to speech with TTS...
[15:04:16] ✅ Speech generated successfully
[15:04:16] 🔈 Playing audio response...
[15:04:20] ✅ Done!
```

## 🔐 Permissions

On first run, macOS will ask for permissions:
- **Camera access** - Required to capture images
- **Microphone access** - Required to record audio

Click "Allow" when prompted.

## 🛠️ Troubleshooting

### Camera/Microphone Not Found
If you get an error about device not found, list available devices:
```bash
ffmpeg -f avfoundation -list_devices true -i ""
```

This will show something like:
```
[0] FaceTime HD Camera
[1] External Camera
[:0] Built-in Microphone
[:1] External Microphone
```

Update `main.go` if needed:
- For video: change `"-i", "0"` to your camera index
- For audio: change `"-i", ":0"` to your microphone index

### API Key Not Set
If you see `OPENAI_API_KEY environment variable not set`:
```bash
export OPENAI_API_KEY="your-key-here"
```

### FFmpeg Not Installed
If you get `executable file not found`:
```bash
brew install ffmpeg
```

## 📦 Dependencies

- **Go 1.22+**
- **github.com/sashabaranov/go-openai** v1.20.4 - OpenAI Go client
- **ffmpeg** - Audio/video capture and playback
- **OpenAI APIs:**
  - Whisper (speech-to-text)
  - GPT-4o-mini (multimodal reasoning)
  - TTS-1 (text-to-speech)

## 🧹 Cleanup

Temporary files (`frame.jpg`, `audio.wav`, `reply.mp3`) are automatically cleaned up after each run.

## 📝 Notes

- The assistant records **5 seconds** of audio by default (configurable in `main.go`)
- Camera captures a **640x480** frame (configurable)
- Audio is recorded in **mono at 16kHz** (optimal for Whisper)
- TTS uses the **Alloy voice** (configurable to: alloy, echo, fable, onyx, nova, shimmer)

## 🎨 Customization

Edit `main.go` to customize:
- `audioDuration` - Change recording length
- `Voice` - Change TTS voice in `textToSpeech()`
- `video_size` - Change camera resolution in `captureImage()`
- `Model` - Switch between GPT models

## 📄 License

MIT
