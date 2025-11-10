package logger

import (
	"fmt"
	"time"
)

// ConsoleLogger is a simple console logger
type ConsoleLogger struct{}

// NewConsoleLogger creates a new console logger
func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

// Info logs an informational message
func (l *ConsoleLogger) Info(msg string) {
	fmt.Printf("[%s] %s\n", time.Now().Format("15:04:05"), msg)
}

// Error logs an error message
func (l *ConsoleLogger) Error(msg string, err error) {
	fmt.Printf("[%s] ❌ %s: %v\n", time.Now().Format("15:04:05"), msg, err)
}
