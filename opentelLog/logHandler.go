package opentelLog

import (
	"context"
	"log/slog"
)

type OpentelLogHandler struct {
	Level slog.Level
}

func NewOpentelLogHandler(level slog.Level) *OpentelLogHandler {
	return &OpentelLogHandler{Level: level}
}

func (h *OpentelLogHandler) Enabled(context context.Context, level slog.Level) bool {
	return h.Level <= level
}

func (h *OpentelLogHandler) Handle(context context.Context, record slog.Record) error {
	return nil
}

func (h *OpentelLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *OpentelLogHandler) WithGroup(name string) slog.Handler {
	return h
}
