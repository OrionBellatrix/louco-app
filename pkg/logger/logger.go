package logger

import (
	"os"
	"time"

	"github.com/louco-event/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
	*zerolog.Logger
}

func New(cfg config.LoggerConfig) *Logger {
	// Set global log level
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Configure output format
	var logger zerolog.Logger
	if cfg.Format == "console" {
		logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}).With().Timestamp().Logger()
	} else {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	return &Logger{Logger: &logger}
}

func (l *Logger) WithRequestID(requestID string) *zerolog.Logger {
	logger := l.Logger.With().Str("request_id", requestID).Logger()
	return &logger
}

func (l *Logger) WithUserID(userID string) *zerolog.Logger {
	logger := l.Logger.With().Str("user_id", userID).Logger()
	return &logger
}

func (l *Logger) WithContext(ctx map[string]interface{}) *zerolog.Logger {
	event := l.Logger.With()
	for key, value := range ctx {
		event = event.Interface(key, value)
	}
	logger := event.Logger()
	return &logger
}

// Global logger functions for convenience
func Info() *zerolog.Event {
	return log.Info()
}

func Error() *zerolog.Event {
	return log.Error()
}

func Debug() *zerolog.Event {
	return log.Debug()
}

func Warn() *zerolog.Event {
	return log.Warn()
}

func Fatal() *zerolog.Event {
	return log.Fatal()
}
