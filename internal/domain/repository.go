package domain

import "context"

// MediaCapturer handles camera and microphone capture
type MediaCapturer interface {
	CaptureImage(ctx context.Context, outputPath string) error
	RecordAudio(ctx context.Context, outputPath string, duration int) error
	PlayAudio(ctx context.Context, audioPath string) error
}

// SpeechRecognizer transcribes audio to text
type SpeechRecognizer interface {
	Transcribe(ctx context.Context, audioPath string) (*Transcription, error)
}

// AIAssistant provides multimodal AI capabilities
type AIAssistant interface {
	AnalyzeMultimodal(ctx context.Context, text string, imagePath string) (*AIResponse, error)
}

// SpeechSynthesizer converts text to speech
type SpeechSynthesizer interface {
	Synthesize(ctx context.Context, text string, outputPath string) (*SpeechOutput, error)
}
