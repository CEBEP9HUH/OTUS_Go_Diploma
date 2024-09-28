package loggerwrapper

import (
	"fmt"
	"log/slog"
	"os"
)

type stdLogger struct {
	logger *slog.Logger
	level  string
}

func NewStdLogger(name, level string) (Logger, error) {
	var l slog.Level
	if err := l.UnmarshalText([]byte(level)); err != nil {
		return nil, fmt.Errorf("logger creation: %w", err)
	}
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: l,
	})
	return stdLogger{
		logger: slog.New(h.WithAttrs([]slog.Attr{
			slog.String("name", name),
		})),
		level: level,
	}, nil
}

func (l stdLogger) Debug(msg string, values ...interface{}) {
	l.logger.Debug(fmt.Sprintf(msg, values...))
}

func (l stdLogger) Info(msg string, values ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, values...))
}

func (l stdLogger) Warn(msg string, values ...interface{}) {
	l.logger.Warn(fmt.Sprintf(msg, values...))
}

func (l stdLogger) Error(msg string, values ...interface{}) {
	l.logger.Error(fmt.Sprintf(msg, values...))
}

func (l stdLogger) Fatal(msg string, values ...interface{}) {
	l.logger.Error(fmt.Sprintf(msg, values...))
	os.Exit(1)
}

func (l stdLogger) Level() string {
	return l.level
}
