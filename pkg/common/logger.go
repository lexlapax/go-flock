// ABOUTME: This file provides a default logger implementation for the common.Logger interface.
// ABOUTME: It uses slog for compatibility with go-llms and supports debug mode.

package common

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
)

// slogLogger wraps slog.Logger to implement our Logger interface
type slogLogger struct {
	logger *slog.Logger
}

// Global logger instance
var (
	globalLogger Logger
	loggerOnce   sync.Once
)

// InitLogger initializes the global logger with debug mode setting
func InitLogger(debugMode bool) {
	loggerOnce.Do(func() {
		level := slog.LevelInfo
		if debugMode {
			level = slog.LevelDebug
		}

		handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: level,
		})
		globalLogger = &slogLogger{
			logger: slog.New(handler),
		}
	})
}

// GetLogger returns the global logger instance
func GetLogger() Logger {
	if globalLogger == nil {
		// Initialize with debug mode from environment if not already initialized
		debugMode := os.Getenv("FLOCK_DEBUG") == "true" || os.Getenv("FLOCK_DEBUG") == "1"
		InitLogger(debugMode)
	}
	return globalLogger
}

// Debug logs a debug message
func (l *slogLogger) Debug(ctx context.Context, msg string, args ...interface{}) {
	l.logger.DebugContext(ctx, msg, convertArgs(args)...)
}

// Info logs an info message
func (l *slogLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	l.logger.InfoContext(ctx, msg, convertArgs(args)...)
}

// Warn logs a warning message
func (l *slogLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	l.logger.WarnContext(ctx, msg, convertArgs(args)...)
}

// Error logs an error message
func (l *slogLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	l.logger.ErrorContext(ctx, msg, convertArgs(args)...)
}

// convertArgs converts variadic args to slog attributes
// It expects pairs of key-value arguments
func convertArgs(args []interface{}) []any {
	if len(args) == 0 {
		return nil
	}

	// If we have a single argument, treat it as a message
	if len(args) == 1 {
		return []any{slog.String("detail", fmt.Sprintf("%v", args[0]))}
	}

	// Convert pairs to slog attributes
	attrs := make([]any, 0, len(args))
	for i := 0; i < len(args)-1; i += 2 {
		key := fmt.Sprintf("%v", args[i])
		value := args[i+1]
		attrs = append(attrs, slog.Any(key, value))
	}

	// If we have an odd number of args, add the last one as "extra"
	if len(args)%2 != 0 {
		attrs = append(attrs, slog.Any("extra", args[len(args)-1]))
	}

	return attrs
}

// SetDebugMode allows changing debug mode at runtime
func SetDebugMode(enabled bool) {
	level := slog.LevelInfo
	if enabled {
		level = slog.LevelDebug
	}

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})
	globalLogger = &slogLogger{
		logger: slog.New(handler),
	}
}

// GetSlogger returns the underlying slog.Logger for compatibility with go-llms
func GetSlogger() *slog.Logger {
	if logger, ok := globalLogger.(*slogLogger); ok {
		return logger.logger
	}
	// Fallback to default slog logger
	return slog.Default()
}
