package logger

import (
	"context"
	"log/slog"
	"os"
)

type Slog struct {
	logger *slog.Logger
}

type Logger interface {
	Info(ctx context.Context, message string, attrs ...slog.Attr)
	Error(ctx context.Context, message string, err error, attrs ...slog.Attr)
}

func NewSlog() *Slog {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return &Slog{logger: logger}
}

func (s *Slog) Info(ctx context.Context, message string, attrs ...slog.Attr) {
	s.logger.InfoContext(ctx, message, attrs)
}
func (s *Slog) Error(ctx context.Context, message string, err error, attrs ...slog.Attr) {
	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
	}
	s.logger.ErrorContext(ctx, message, attrs)
}
