package logger

import (
	"golang.org/x/exp/slog" // Import the slog package
	"letsquiz/config"
	"os"
)

var log *slog.Logger
var enabled bool

// InitLogger initializes the logger based on the configuration
func InitLogger(db bool) {
	log = slog.New(slog.NewTextHandler(os.Stdout, nil)) // Default to standard output

	if !db {
		enabled = config.AppConfig.LoggingEnabled
	} else {
		enabled = config.DbConfig.LoggingEnabled
	}

	if enabled && !db {
		SetLogFile(config.AppConfig.LogFile) // Set log file for application logs
	} else if enabled && db {
		SetLogFile(config.DbConfig.LogFile) // Set log file for database logs
	}
}

// Info logs an informational message with optional fields
func Info(msg string, keysAndValues ...interface{}) {
	if enabled {
		log.Info(msg, keysAndValues...)
	}
}

// Error logs an error message with optional fields
func Error(msg string, keysAndValues ...interface{}) {
	if enabled {
		log.Error(msg, keysAndValues...)
	}
}

// SetLogFile sets the log output to a file specified by 'filepath'
func SetLogFile(filepath string) {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Error("Failed to open log file", "error", err) // Log an error if file opening fails
		return
	}
	log = slog.New(slog.NewTextHandler(file, nil)) // Update logger to write to the file
}
