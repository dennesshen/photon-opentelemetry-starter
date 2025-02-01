package opentelCore

import (
	"context"
	"log/slog"
)

func Shutdown(ctx context.Context) error {
	slog.Info("[otel]shutdown otel")

	if exporter != nil {
		if err := exporter.Shutdown(ctx); err != nil {
			slog.Error("[otel]shutdown metric exporter error", "error", err)
		}
	}
	if tp == nil {
		return nil
	}
	if err := tp.Shutdown(ctx); err != nil {
		slog.Error("[otel]shutdown trace provider error", "error", err)
	}
	return nil
}
