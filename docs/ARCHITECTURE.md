# Smart Home AI - Clean Architecture

## 📁 Project Structure

```
smart-home-ai/
├── cmd/
│   └── assistant/          # Application entry points
│       └── main.go         # Main CLI application
├── internal/
│   ├── domain/             # Business entities & interfaces (core)
│   │   ├── entity.go       # Core business models
│   │   └── repository.go   # Port interfaces
│   ├── usecase/            # Application business logic
│   │   └── assistant.go    # AI assistant orchestration
│   └── infrastructure/     # External adapters (implementations)
│       ├── openai/         # OpenAI API client
│       │   └── client.go
│       └── media/          # FFmpeg media capture
│           └── ffmpeg.go
├── pkg/                    # Public shared libraries
│   └── logger/
│       └── console.go      # Console logging utility
├── .env                    # Environment variables
├── .gitignore
├── go.mod
├── go.sum
├── main.go                 # Legacy entry (redirects to new structure)
└── README.md
```

## 🏗️ Clean Architecture Layers

### 1. **Domain Layer** (`internal/domain/`)
- **Purpose**: Core business logic, entities, and interfaces
- **Dependencies**: None (innermost layer)
- **Files**:
  - `entity.go`: Business models (MediaCapture, Transcription, AIResponse, SpeechOutput)
  - `repository.go`: Port interfaces (MediaCapturer, SpeechRecognizer, AIAssistant, SpeechSynthesizer)

### 2. **Use Case Layer** (`internal/usecase/`)
- **Purpose**: Application-specific business rules and orchestration
- **Dependencies**: Only depends on domain layer
- **Files**:
  - `assistant.go`: Orchestrates the complete AI assistant workflow

### 3. **Infrastructure Layer** (`internal/infrastructure/`)
- **Purpose**: External service adapters and implementations
- **Dependencies**: Implements domain interfaces
- **Packages**:
  - `openai/`: OpenAI API client (implements SpeechRecognizer, AIAssistant, SpeechSynthesizer)
  - `media/`: FFmpeg wrapper (implements MediaCapturer)

### 4. **Delivery Layer** (`cmd/`)
- **Purpose**: Entry points and dependency injection
- **Dependencies**: Wires all layers together
- **Files**:
  - `cmd/assistant/main.go`: CLI application entry point

### 5. **Shared Packages** (`pkg/`)
- **Purpose**: Reusable utilities accessible by any layer
- **Files**:
  - `logger/console.go`: Console logging implementation

## 🚀 Running the Application

### Quick Start
```bash
# Run directly
go run cmd/assistant/main.go

# Or build binary
go build -o smart-home-ai cmd/assistant/main.go
./smart-home-ai
```

### Environment Setup
```bash
# 1. Install dependencies
go mod tidy

# 2. Install FFmpeg
brew install ffmpeg

# 3. Set OpenAI API key in .env
echo 'OPENAI_API_KEY=sk-your-key-here' > .env
```

## 🧪 Benefits of This Architecture

### ✅ **Separation of Concerns**
- Each layer has a single responsibility
- Business logic isolated from external dependencies

### ✅ **Testability**
- Easy to mock interfaces for unit testing
- Domain logic testable without external services

### ✅ **Flexibility**
- Swap implementations without changing business logic
- Example: Replace OpenAI with local LLM by implementing domain interfaces

### ✅ **Maintainability**
- Clear structure makes codebase navigable
- Changes isolated to specific layers

### ✅ **Dependency Rule**
- Dependencies point inward (outer → inner)
- Domain layer has no external dependencies

## 🔄 Data Flow

```
User Input
    ↓
cmd/assistant/main.go (Dependency Injection)
    ↓
usecase/assistant.go (Orchestration)
    ↓
domain/repository.go (Interfaces)
    ↓
infrastructure/* (Implementations)
    ↓
External Services (OpenAI, FFmpeg)
```

## 🧩 Adding New Features

### Example: Add Google Cloud Speech-to-Text

1. **No changes to domain** - interface already exists
2. **Create new adapter**:
   ```go
   // internal/infrastructure/google/speech.go
   type GoogleSpeechRecognizer struct {}
   
   func (g *GoogleSpeechRecognizer) Transcribe(ctx context.Context, audioPath string) (*domain.Transcription, error) {
       // Google Cloud implementation
   }
   ```

3. **Update dependency injection** in `cmd/assistant/main.go`:
   ```go
   // Old: openaiClient
   // New: googleRecognizer
   assistantUseCase := usecase.NewAssistantUseCase(
       mediaCapturer,
       googleRecognizer, // <- Changed
       openaiClient,
       openaiClient,
       consoleLogger,
   )
   ```

## 📚 Design Patterns Used

- **Dependency Injection**: Manual DI in main.go
- **Repository Pattern**: domain/repository.go interfaces
- **Adapter Pattern**: infrastructure/* implementations
- **Use Case Pattern**: usecase/* orchestration

## 🔐 Environment Variables

```bash
OPENAI_API_KEY=sk-...  # Required: OpenAI API authentication
```

## 📖 Further Reading

- [The Clean Architecture by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
