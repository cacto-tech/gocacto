package logger

import (
	"context"
	"log/slog"
	"os"
)

// Logger wraps slog.Logger with additional methods
type Logger struct {
	*slog.Logger
}

// New creates a new logger instance
func New(level slog.Level) *Logger {
	opts := &slog.HandlerOptions{
		Level: level,
		AddSource: true,
	}

	var handler slog.Handler = slog.NewJSONHandler(os.Stdout, opts)
	
	// In development, use text handler for better readability
	if os.Getenv("ENV") == "development" || os.Getenv("ENV") == "dev" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	
	return &Logger{Logger: logger}
}

// Default creates a logger with Info level
func Default() *Logger {
	return New(slog.LevelInfo)
}

// WithContext returns a logger with context
func (l *Logger) WithContext(ctx context.Context) *Logger {
	return &Logger{Logger: l.Logger.With("request_id", getRequestID(ctx))}
}

// WithError adds error to logger
func (l *Logger) WithError(err error) *Logger {
	return &Logger{Logger: l.Logger.With("error", err.Error())}
}

// WithField adds a field to logger
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{Logger: l.Logger.With(key, value)}
}

// WithFields adds multiple fields to logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &Logger{Logger: l.Logger.With(args...)}
}

// getRequestID extracts request ID from context
func getRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if id, ok := ctx.Value("request_id").(string); ok {
		return id
	}
	return ""
}

// Global logger instance
var defaultLogger = Default()

// SetDefault sets the default logger
func SetDefault(logger *Logger) {
	defaultLogger = logger
}

// GetDefault returns the default logger
func GetDefault() *Logger {
	return defaultLogger
}

// Helper functions for global logger
func Info(msg string, args ...interface{}) {
	defaultLogger.Info(msg, args...)
}

func Error(msg string, args ...interface{}) {
	defaultLogger.Error(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	defaultLogger.Warn(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	defaultLogger.Debug(msg, args...)
}
