package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"sync"
	"time"
)

// SecurityManager handles security features
type SecurityManager struct {
	allowedCommands map[string]bool
	rateLimit       *RateLimiter
	commandLog      []CommandLog
	mu              sync.RWMutex
}

// CommandLog logs executed commands
type CommandLog struct {
	Timestamp time.Time
	Command   string
	Device    string
	Success   bool
	Error     string
}

// RateLimiter limits command execution rate
type RateLimiter struct {
	maxRequests int
	window      time.Duration
	requests    []time.Time
	mu          sync.Mutex
}

// NewSecurityManager creates a new security manager
func NewSecurityManager() *SecurityManager {
	return &SecurityManager{
		allowedCommands: map[string]bool{
			"light.on":         true,
			"light.off":        true,
			"light.brightness": true,
			"light.color":      true,
			"light.color_temp": true,
			"switch.on":        true,
			"switch.off":       true,
			"switch.toggle":    true,
			"ac.on":            true,
			"ac.off":           true,
			"ac.set_temp":      true,
			"vacuum.start":     true,
			"vacuum.stop":      true,
			"vacuum.pause":     true,
			"vacuum.home":      true,
			"tv.power":         true,
			"tv.vol_up":        true,
			"tv.vol_down":      true,
		},
		rateLimit:  NewRateLimiter(10, 1*time.Minute),
		commandLog: make([]CommandLog, 0),
	}
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		maxRequests: maxRequests,
		window:      window,
		requests:    make([]time.Time, 0),
	}
}

// Allow checks if a request is allowed
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Remove old requests outside the window
	cutoff := now.Add(-rl.window)
	newRequests := make([]time.Time, 0)
	for _, req := range rl.requests {
		if req.After(cutoff) {
			newRequests = append(newRequests, req)
		}
	}
	rl.requests = newRequests

	// Check if we're under the limit
	if len(rl.requests) >= rl.maxRequests {
		return false
	}

	// Add this request
	rl.requests = append(rl.requests, now)
	return true
}

// ValidateCommand validates if a command is allowed
func (sm *SecurityManager) ValidateCommand(cmd *Command) error {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// Check rate limit
	if !sm.rateLimit.Allow() {
		return fmt.Errorf("rate limit exceeded")
	}

	// Check if command is in allowed list
	if !sm.allowedCommands[cmd.Action] {
		return fmt.Errorf("command not allowed: %s", cmd.Action)
	}

	// Additional validation based on device type
	if err := sm.validateDeviceCommand(cmd); err != nil {
		return err
	}

	return nil
}

// validateDeviceCommand validates device-specific commands
func (sm *SecurityManager) validateDeviceCommand(cmd *Command) error {
	// Validate brightness range
	if cmd.Action == "light.brightness" {
		if brightness, ok := cmd.Value.(float64); ok {
			if brightness < 0 || brightness > 100 {
				return fmt.Errorf("brightness must be between 0 and 100")
			}
		}
	}

	// Validate temperature range
	if cmd.Action == "ac.set_temp" {
		if temp, ok := cmd.Value.(float64); ok {
			if temp < 16 || temp > 30 {
				return fmt.Errorf("temperature must be between 16 and 30")
			}
		}
	}

	return nil
}

// LogCommand logs a command execution
func (sm *SecurityManager) LogCommand(cmd *Command, success bool, err error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	entry := CommandLog{
		Timestamp: time.Now(),
		Command:   cmd.Action,
		Device:    cmd.Device,
		Success:   success,
		Error:     errMsg,
	}

	sm.commandLog = append(sm.commandLog, entry)

	// Keep only last 1000 entries
	if len(sm.commandLog) > 1000 {
		sm.commandLog = sm.commandLog[len(sm.commandLog)-1000:]
	}

	// Log to console
	if success {
		log.Printf("[SECURITY] Command executed: %s on %s", cmd.Action, cmd.Device)
	} else {
		log.Printf("[SECURITY] Command failed: %s on %s - %s", cmd.Action, cmd.Device, errMsg)
	}
}

// GetCommandLog returns recent command logs
func (sm *SecurityManager) GetCommandLog(limit int) []CommandLog {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if limit <= 0 || limit > len(sm.commandLog) {
		limit = len(sm.commandLog)
	}

	result := make([]CommandLog, limit)
	copy(result, sm.commandLog[len(sm.commandLog)-limit:])

	return result
}

// HashPassword hashes a password using SHA256
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// VerifyPassword verifies a password against a hash
func VerifyPassword(password, hash string) bool {
	return HashPassword(password) == hash
}

// AddAllowedCommand adds a command to the allowed list
func (sm *SecurityManager) AddAllowedCommand(action string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.allowedCommands[action] = true
	log.Printf("[SECURITY] Added allowed command: %s", action)
}

// RemoveAllowedCommand removes a command from the allowed list
func (sm *SecurityManager) RemoveAllowedCommand(action string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.allowedCommands, action)
	log.Printf("[SECURITY] Removed allowed command: %s", action)
}

// IsCommandAllowed checks if a command is allowed
func (sm *SecurityManager) IsCommandAllowed(action string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.allowedCommands[action]
}
