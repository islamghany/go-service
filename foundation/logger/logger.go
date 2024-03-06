package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"log/slog"
)

// TracerIDFn is a function that returns a string to be used as the tracer ID.
type TracerIDFn func(ctx context.Context) string

// Logger represents a logger that can be used to log messages.
// The interface is left to be implemented by the user.
type Logger struct {
	handler   slog.Handler
	traceIDFn TracerIDFn // private
}

// New creates a new Logger instance.
func New(w io.Writer, minLevel Level, serviceName string, traceIDFn TracerIDFn) *Logger {
	return new(w, minLevel, serviceName, traceIDFn, Events{})
}

// NewWithEvents creates a new Logger instance with custom event functions.
func NewWithEvents(w io.Writer, minLevel Level, serviceName string, traceIDFn TracerIDFn, events Events) *Logger {
	return new(w, minLevel, serviceName, traceIDFn, events)
}

// NewWithHandler creates a new Logger instance with a custom handler.
func NewWithHandler(handler slog.Handler) *Logger {
	return &Logger{
		handler: handler,
	}
}

// NewStdLogger returns a standard library Logger that wraps the slog Logger.
func NewStdLogger(logger *Logger, level Level) *log.Logger {
	return slog.NewLogLogger(logger.handler, slog.Level(level))
}

// Debug logs a debug message.
func (l *Logger) Debug(ctx context.Context, message string, attrs ...any) {
	l.write(ctx, LevelDebug, 3, message, attrs...)
}

// DebugCtx logs a debug message with a the specified call stack position.
func (l *Logger) DebugCtx(ctx context.Context, callDepth int, message string, attrs ...any) {
	l.write(ctx, LevelDebug, callDepth, message, attrs...)
}

// Info logs an info message.
func (l *Logger) Info(ctx context.Context, message string, attrs ...any) {
	l.write(ctx, LevelInfo, 3, message, attrs...)
}

// InfoCtx logs an info message with a the specified call stack position.
func (l *Logger) InfoCtx(ctx context.Context, callDepth int, message string, attrs ...any) {
	l.write(ctx, LevelInfo, callDepth, message, attrs...)
}

// Warn logs a warning message.
func (l *Logger) Warn(ctx context.Context, message string, attrs ...any) {
	l.write(ctx, LevelWarn, 3, message, attrs...)
}

// WarnCtx logs a warning message with a the specified call stack position.
func (l *Logger) WarnCtx(ctx context.Context, callDepth int, message string, attrs ...any) {
	l.write(ctx, LevelWarn, callDepth, message, attrs...)
}

// Error logs an error message.
func (l *Logger) Error(ctx context.Context, message string, attrs ...any) {
	l.write(ctx, LevelError, 3, message, attrs...)
}

// ErrorCtx logs an error message with a the specified call stack position.
func (l *Logger) ErrorCtx(ctx context.Context, callDepth int, message string, attrs ...any) {
	l.write(ctx, LevelError, callDepth, message, attrs...)
}

func (l *Logger) write(ctx context.Context, level Level, callDepth int, message string, attrs ...any) {
	slogLevel := slog.Level(level)
	// Check if the log level is enabled, this is determined when you pass the minimum log level to the slog.NewLogger function.
	if !l.handler.Enabled(ctx, slogLevel) {
		return
	}

	// pcs is the program counter slice, it's used to get the call stack.
	var pcs [1]uintptr
	runtime.Callers(callDepth, pcs[:])
	// Create a new log record.
	r := slog.NewRecord(time.Now(), slogLevel, message, pcs[0])
	// if the traceIDFn is not nil, then we add the trace ID to the log record.
	if l.traceIDFn != nil {
		attrs = append(attrs, "trace_id", l.traceIDFn(ctx))
	}
	// Add the attributes to the log record.
	r.Add(attrs...)

	// Call the handler to handle the log record.
	l.handler.Handle(ctx, r)
}

func new(w io.Writer, minLevel Level, serviceName string, traceIDFn TracerIDFn, events Events) *Logger {
	// convert the file name to just the name.ext format.
	fn := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			if source, ok := a.Value.Any().(*slog.Source); ok {
				v := fmt.Sprintf("%s:%d", filepath.Base(source.File), source.Line)
				return slog.Attr{Key: "file", Value: slog.StringValue(v)}
			}

		}
		return a
	}
	handler := slog.Handler(slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level:       slog.Level(minLevel),
		AddSource:   true,
		ReplaceAttr: fn,
	}))
	// If events are to be processed, wrap the JSON handler around the custom
	// log handler.
	if events.Debug != nil || events.Info != nil || events.Warn != nil || events.Error != nil {
		handler = newLogHandler(handler, events)
	}
	// Attributes to add to every log.
	attrs := []slog.Attr{
		{Key: "service", Value: slog.StringValue(serviceName)},
	}

	// Add those attributes and capture the final handler.
	handler = handler.WithAttrs(attrs)

	return &Logger{
		handler:   handler,
		traceIDFn: traceIDFn,
	}
}
