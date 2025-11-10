# Quick Start Guide

## 🚀 Run the Application

```bash
# Make sure ffmpeg is installed
brew install ffmpeg

# Set your OpenAI API key in .env file
echo 'OPENAI_API_KEY=sk-your-key-here' > .env

# Run the application
go run cmd/assistant/main.go
```

## 📂 Project Structure (Clean Architecture)

```
smart-home-ai/
├── cmd/assistant/main.go          # 🚪 Entry point
├── internal/
│   ├── domain/                    # 🎯 Core business logic
│   │   ├── entity.go             # Data models
│   │   └── repository.go         # Interfaces (ports)
│   ├── usecase/                   # 🔄 Application logic
│   │   └── assistant.go          # Workflow orchestration
│   └── infrastructure/            # 🔌 External adapters
│       ├── openai/               # OpenAI API client
│       └── media/                # FFmpeg wrapper
└── pkg/logger/                    # 📝 Shared utilities
```

## 🎯 How It Works

1. **Capture** → Takes photo from camera
2. **Record** → Records 5 seconds of audio
3. **Transcribe** → Whisper API converts speech to text
4. **Analyze** → GPT-4o-mini understands image + text
5. **Synthesize** → TTS converts response to speech
6. **Play** → Plays audio response

## 🧩 Key Design Patterns

- **Clean Architecture**: Separation of concerns
- **Dependency Injection**: Wired in `cmd/assistant/main.go`
- **Repository Pattern**: Interfaces in `domain/repository.go`
- **Adapter Pattern**: Implementations in `infrastructure/`

## 🔧 Customization Examples

### Change TTS Voice
Edit `internal/infrastructure/openai/client.go`:
```go
Voice: openai.VoiceNova, // or: Alloy, Echo, Fable, Onyx, Shimmer
```

### Change Recording Duration
Edit `internal/usecase/assistant.go`:
```go
const audioDuration = 10 // 10 seconds instead of 5
```

### Add New AI Provider
1. Create `internal/infrastructure/yourprovider/client.go`
2. Implement `domain.SpeechRecognizer`, `domain.AIAssistant`, etc.
3. Wire it in `cmd/assistant/main.go`

## 📚 Documentation

- **[ARCHITECTURE.md](ARCHITECTURE.md)** - Detailed architecture guide
- **[REFACTORING_SUMMARY.md](REFACTORING_SUMMARY.md)** - What changed and why
- **[README.md](README.md)** - Full setup instructions

## ❓ Troubleshooting

### Camera/Mic Not Found
```bash
# List available devices
ffmpeg -f avfoundation -list_devices true -i ""
```

### Permission Denied
Allow Camera & Microphone access in System Settings → Privacy & Security

### API Key Error
Make sure `.env` file exists and contains:
```
OPENAI_API_KEY=sk-your-actual-key
```

## 🧪 Testing (Future)

```bash
# Unit tests (to be added)
go test ./internal/...

# Integration tests (to be added)
go test ./test/integration/...
```

## 🎨 Benefits Over Old Code

| Feature | Old (main.go) | New (Clean Arch) |
|---------|--------------|------------------|
| **Lines of code** | 235 in one file | Organized across layers |
| **Testability** | ❌ Hard | ✅ Easy with mocks |
| **Extensibility** | ❌ Requires editing main | ✅ Just add adapters |
| **Readability** | ⚠️ Mixed concerns | ✅ Clear separation |
| **Maintainability** | ⚠️ Coupled | ✅ Loosely coupled |

---

Need help? Check the documentation files above! 📖
