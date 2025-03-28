package server

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
)

func init() {
	// Create logs directory if it doesn't exist
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		log.Fatal("Could not create logs directory:", err)
	}

	// Create or append to log file with timestamp in name
	currentTime := time.Now()
	logFileName := filepath.Join("logs", fmt.Sprintf("shandris_%s.log", currentTime.Format("2006-01-02")))
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Could not open log file:", err)
	}

	// Create loggers with different prefixes
	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime)
	DebugLogger = log.New(file, "DEBUG: ", log.Ldate|log.Ltime)

	// Also write to stdout for development
	InfoLogger.SetOutput(os.Stdout)
	ErrorLogger.SetOutput(os.Stdout)
	DebugLogger.SetOutput(os.Stdout)
}

// LogOperation logs the start and end of an operation with its details
func LogOperation(operation string, details map[string]interface{}) func(error) {
	startTime := time.Now()
	_, file, line, _ := runtime.Caller(1)
	file = filepath.Base(file)

	// Log operation start
	var detailsStr []string
	for k, v := range details {
		detailsStr = append(detailsStr, fmt.Sprintf("%s: %v", k, v))
	}
	InfoLogger.Printf("🔵 Starting %s at %s:%d [%s]", operation, file, line, strings.Join(detailsStr, ", "))

	// Return function to be called when operation ends
	return func(err error) {
		duration := time.Since(startTime)
		if err != nil {
			ErrorLogger.Printf("❌ %s failed after %v: %v", operation, duration, err)
		} else {
			InfoLogger.Printf("✅ %s completed in %v", operation, duration)
		}
	}
}

// LogMemoryOperation logs memory-related operations
func LogMemoryOperation(operation string, sessionID string, key string, value interface{}) {
	_, file, line, _ := runtime.Caller(1)
	file = filepath.Base(file)

	details := fmt.Sprintf("SessionID: %s, Key: %s, Value: %v", sessionID, key, value)
	InfoLogger.Printf("💾 %s at %s:%d [%s]", operation, file, line, details)
}

// LogProfileOperation logs profile-related operations
func LogProfileOperation(operation string, sessionID string, profile interface{}) {
	_, file, line, _ := runtime.Caller(1)
	file = filepath.Base(file)

	InfoLogger.Printf("👤 %s at %s:%d [SessionID: %s, Profile: %v]", operation, file, line, sessionID, profile)
}

// LogChatOperation logs chat-related operations
func LogChatOperation(operation string, sessionID string, message string, topic string) {
	_, file, line, _ := runtime.Caller(1)
	file = filepath.Base(file)

	details := fmt.Sprintf("SessionID: %s, Topic: %s, Message: %s", sessionID, topic, message)
	InfoLogger.Printf("💬 %s at %s:%d [%s]", operation, file, line, details)
}

// LogError logs error details with stack trace
func LogError(err error, context string) {
	_, file, line, _ := runtime.Caller(1)
	file = filepath.Base(file)

	ErrorLogger.Printf("❌ Error in %s at %s:%d: %v", context, file, line, err)
}
