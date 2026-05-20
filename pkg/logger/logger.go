package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Level represents the log level
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var levelNames = map[Level]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
}

type Logger struct {
	level  Level
	logger *log.Logger
}

var defaultLogger *Logger

func init() {
	defaultLogger = New(INFO)
}

// New creates a new logger instance
func New(level Level) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...interface{}) {
	if l.level <= DEBUG {
		l.log("DEBUG", msg, fields...)
	}
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...interface{}) {
	if l.level <= INFO {
		l.log("INFO", msg, fields...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...interface{}) {
	if l.level <= WARN {
		l.log("WARN", msg, fields...)
	}
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...interface{}) {
	if l.level <= ERROR {
		l.log("ERROR", msg, fields...)
	}
}

func (l *Logger) log(level string, msg string, fields ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fieldStr := ""

	// Format fields as key=value pairs
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			fieldStr += fmt.Sprintf(" %v=%v", fields[i], fields[i+1])
		}
	}

	logMsg := fmt.Sprintf("[%s] %s: %s%s", timestamp, level, msg, fieldStr)
	l.logger.Println(logMsg)
}

// Global convenience functions
func Debug(msg string, fields ...interface{}) {
	defaultLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...interface{}) {
	defaultLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...interface{}) {
	defaultLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...interface{}) {
	defaultLogger.Error(msg, fields...)
}

func SetLevel(level Level) {
	defaultLogger.SetLevel(level)
}
