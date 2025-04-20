package internal

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Logger struct {
	startTime time.Time
	logger    *log.Logger
	level     slog.Level
}

func NewLogger(level slog.Level, output *os.File) *Logger {
	return &Logger{
		startTime: time.Now(),
		logger:    log.New(output, "", 0),
		level:     level,
	}
}

// Enabled determines if the log level is enabled.
func (l *Logger) Enabled(ctx context.Context, level slog.Level) bool {
	return l.level <= level
}

// Handle formats and writes the log record.
func (l *Logger) Handle(ctx context.Context, record slog.Record) error {
	elapsed := strconv.FormatFloat(time.Since(l.startTime).Seconds(), 'f', 2, 64)
	message := record.Message

	// Format attributes into a string
	attrs := ""
	record.Attrs(func(attr slog.Attr) bool {
		attrs += fmt.Sprintf(" %s=%v", attr.Key, attr.Value)
		return true
	})

	// Write the log message
	l.logger.Printf("[%10s] %s%s", elapsed, message, attrs)
	return nil
}

func (l *Logger) WithAttrs(attrs []slog.Attr) slog.Handler {
	// For simplicity, return the same logger. Extend if needed.
	return l
}

func (l *Logger) WithGroup(name string) slog.Handler {
	// For simplicity, return the same logger. Extend if needed.
	return l
}
