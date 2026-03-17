package lib

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type logEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

var (
	InfoLogger  = log.New(os.Stdout, "", 0) // no prefix/flags, we'll handle our own
	ErrorLogger = log.New(os.Stderr, "", 0)
	WarnLogger  = log.New(os.Stdout, "", 0)
)

// Helper to write a JSON log entry
func writeLog(logger *log.Logger, level, msg string) {
	entry := logEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level,
		Message:   msg,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		// Fallback in case JSON marshaling fails
		logger.Printf(`{"level":"ERROR","msg":"failed to marshal log entry: %v"}`, err)
		return
	}
	logger.Println(string(data))
}

// Public logging functions
func Info(msg string) {
	writeLog(InfoLogger, "INFO", msg)
}

func Error(msg string) {
	writeLog(ErrorLogger, "ERROR", msg)
}

func Warn(msg string) {
	writeLog(WarnLogger, "WARN", msg)
}
