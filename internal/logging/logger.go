// Package logging provides structured logging functionality.
package logging

import (
	"fmt"
	"io"
	"os"

	"github.com/gumi-tsd/secret-env-manager/internal/formatting"
	"github.com/gumi-tsd/secret-env-manager/internal/functional"
)

// Level represents the severity level of a log message
type Level int

const (
	// DebugLevel represents debug level logs
	DebugLevel Level = iota
	// InfoLevel represents informational level logs
	InfoLevel
	// WarnLevel represents warning level logs
	WarnLevel
	// ErrorLevel represents error level logs
	ErrorLevel
	// SuccessLevel represents success level logs
	SuccessLevel
)

// Logger provides structured logging functionality
type Logger struct {
	writer    io.Writer
	minLevel  Level
	colorized bool
}

// DefaultLogger creates a new logger with default settings
func DefaultLogger() *Logger {
	return &Logger{
		writer:    os.Stdout,
		minLevel:  InfoLevel, // Default log level is Info
		colorized: true,
	}
}

// NewLogger creates a new logger with custom settings
func NewLogger(writer io.Writer, minLevel Level, colorized bool) *Logger {
	return &Logger{
		writer:    writer,
		minLevel:  minLevel,
		colorized: colorized,
	}
}

// WithLevel returns a copy of the logger with a new minimum log level
func (l *Logger) WithLevel(level Level) *Logger {
	return NewLogger(l.writer, level, l.colorized)
}

// WithColorized returns a copy of the logger with colorization setting
func (l *Logger) WithColorized(colorized bool) *Logger {
	return NewLogger(l.writer, l.minLevel, colorized)
}

// WithWriter returns a copy of the logger with a new writer
func (l *Logger) WithWriter(writer io.Writer) *Logger {
	return NewLogger(writer, l.minLevel, l.colorized)
}

// Debug logs a debug message if the level permits
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DebugLevel, format, args...)
}

// Info logs an info message if the level permits
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(InfoLevel, format, args...)
}

// Warn logs a warning message if the level permits
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WarnLevel, format, args...)
}

// Error logs an error message if the level permits
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ErrorLevel, format, args...)
}

// Success logs a success message if the level permits
func (l *Logger) Success(format string, args ...interface{}) {
	l.log(SuccessLevel, format, args...)
}

// log handles the actual message formatting and output
func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level < l.minLevel {
		return
	}

	var prefix string
	switch level {
	case DebugLevel:
		prefix = "[DEBUG] "
	case InfoLevel:
		prefix = "[INFO] "
	case WarnLevel:
		prefix = "[WARN] "
	case ErrorLevel:
		prefix = "[ERROR] "
	case SuccessLevel:
		prefix = "[SUCCESS] "
	}

	message := fmt.Sprintf(format, args...)
	if l.colorized {
		switch level {
		case DebugLevel:
			message = formatting.Info("%s", message)
		case InfoLevel:
			// Default color
		case WarnLevel:
			message = formatting.Warning("%s", message)
		case ErrorLevel:
			message = formatting.Error("%s", message)
		case SuccessLevel:
			message = formatting.Success("%s", message)
		}
		prefix = formatting.Info("%s", prefix)
	}

	fmt.Fprintf(l.writer, "%s%s\n", prefix, message)
}

// LoggingIO wraps logging functionality in an IO monad
type LoggingIO struct {
	logger *Logger
}

// NewLoggingIO creates a new logging IO with a default logger
func NewLoggingIO() LoggingIO {
	return LoggingIO{
		logger: DefaultLogger(),
	}
}

// WithLogger returns a new LoggingIO with the specified logger
func (l LoggingIO) WithLogger(logger *Logger) LoggingIO {
	return LoggingIO{
		logger: logger,
	}
}

// Debug returns an IO that logs a debug message when performed
func (l LoggingIO) Debug(format string, args ...interface{}) functional.IO[struct{}] {
	return functional.NewIO(func() struct{} {
		l.logger.Debug(format, args...)
		return struct{}{}
	})
}

// Info returns an IO that logs an info message when performed
func (l LoggingIO) Info(format string, args ...interface{}) functional.IO[struct{}] {
	return functional.NewIO(func() struct{} {
		l.logger.Info(format, args...)
		return struct{}{}
	})
}

// Warn returns an IO that logs a warning message when performed
func (l LoggingIO) Warn(format string, args ...interface{}) functional.IO[struct{}] {
	return functional.NewIO(func() struct{} {
		l.logger.Warn(format, args...)
		return struct{}{}
	})
}

// Error returns an IO that logs an error message when performed
func (l LoggingIO) Error(format string, args ...interface{}) functional.IO[struct{}] {
	return functional.NewIO(func() struct{} {
		l.logger.Error(format, args...)
		return struct{}{}
	})
}

// Success returns an IO that logs a success message when performed
func (l LoggingIO) Success(format string, args ...interface{}) functional.IO[struct{}] {
	return functional.NewIO(func() struct{} {
		l.logger.Success(format, args...)
		return struct{}{}
	})
}
