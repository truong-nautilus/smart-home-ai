package domain

import "time"

// MediaCapture represents captured media data
type MediaCapture struct {
	ImagePath string
	AudioPath string
	Timestamp time.Time
}

// Transcription represents speech-to-text result
type Transcription struct {
	Text      string
	Language  string
	Duration  float64
	Timestamp time.Time
}

// AIResponse represents the AI assistant's response
type AIResponse struct {
	Text      string
	Timestamp time.Time
	Model     string
}

// SpeechOutput represents synthesized speech
type SpeechOutput struct {
	AudioPath string
	Text      string
	Voice     string
	Timestamp time.Time
}
