package logger

import (
	"context"
	"time"

	"log/slog"
)

// Level represents the log level for different log messages.
type Level slog.Level

// Log levels.
const (
	LevelDebug = Level(slog.LevelDebug)
	LevelInfo  = Level(slog.LevelInfo)
	LevelWarn  = Level(slog.LevelWarn)
	LevelError = Level(slog.LevelError)
)

// Record represents a log record.
type Record struct {
	Level      Level
	Time       time.Time
	Message    string
	Attributes map[string]any
}

// toRecord converts a slog.Record to a Record of our choice.
func toRecord(r slog.Record) Record {
	atts := make(map[string]any, r.NumAttrs())
	fn := func(attr slog.Attr) bool {
		atts[attr.Key] = attr.Value
		return true
	}
	r.Attrs(fn)
	return Record{
		Level:      Level(r.Level),
		Time:       r.Time,
		Message:    r.Message,
		Attributes: make(map[string]any),
	}
}

// EventFn is a function that logs an event, that will be attached to the log level.
type EventFn func(ctx context.Context, r Record)

// Events contains an assignment of an event function to a log level.
type Events struct {
	Debug EventFn
	Info  EventFn
	Warn  EventFn
	Error EventFn
}
