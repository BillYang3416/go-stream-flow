package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// Logger -.
type Logger interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

// logger -.
type logger struct {
	logger *zerolog.Logger
}

// New -.
func New(level string) Logger {
	var l zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)

	// skipFrameCount is used to skip frames from the stacktrace
	// Debug() => foramtMessage() => log()
	// 3 frames should be skipped
	skipFrameCount := 3
	var z zerolog.Logger
	if l == zerolog.DebugLevel {
		z = zerolog.New(os.Stdout).Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
			With().
			Timestamp().
			CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).
			Logger()
	} else {
		z = zerolog.New(os.Stdout).
			With().
			Timestamp().
			CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).
			Logger()
	}

	return &logger{
		logger: &z,
	}
}

func (l *logger) formatMessage(message interface{}) string {
	switch t := message.(type) {
	case error:
		return t.Error()
	case string:
		return t
	default:
		return fmt.Sprintf("Unknown type %v", message)
	}
}

// Debug -.
func (l *logger) Debug(message interface{}, args ...interface{}) {
	mf := l.formatMessage(message)
	l.log(l.logger.Debug(), mf, args...)
}

// Info -.
func (l *logger) Info(message string, args ...interface{}) {
	mf := l.formatMessage(message)
	l.log(l.logger.Info(), mf, args...)
}

// Warn -.
func (l *logger) Warn(message string, args ...interface{}) {
	mf := l.formatMessage(message)
	l.log(l.logger.Warn(), mf, args...)
}

// Error -.
func (l *logger) Error(message interface{}, args ...interface{}) {
	mf := l.formatMessage(message)
	l.log(l.logger.Error(), mf, args...)
}

// Fatal -.
func (l *logger) Fatal(message interface{}, args ...interface{}) {
	mf := l.formatMessage(message)
	l.log(l.logger.Fatal(), mf, args...)

	os.Exit(1)
}

func (l *logger) log(e *zerolog.Event, m string, args ...interface{}) {
	if len(args) == 0 {
		e.Msg(m)
	} else {
		e.Msgf(m, args...)
	}
}
