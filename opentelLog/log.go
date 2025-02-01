package opentelLog

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	
	"github.com/dennesshen/photon-opentelemetry-starter/opentelCore"
	
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type OpentelLogger struct {
	logLevel slog.Level
	handler  slog.Handler
}

func NewOpentelLogger(logLevel slog.Level) *OpentelLogger {
	return &OpentelLogger{
		logLevel: logLevel,
		handler:  NewOpentelLogHandler(logLevel),
	}
}

func (l *OpentelLogger) log(ctx context.Context, level string, callerSkip int, msg string, args ...any) {
	if !l.handler.Enabled(ctx, l.logLevel) {
		return
	}
	
	_, span := opentelCore.NewTrace(ctx, "log")
	defer span.End()
	span.SetAttributes(attribute.String("log.type", level))
	span.SetAttributes(attribute.String("log.message", msg))
	for i := 0; i < len(args)-1; i = i + 2 {
		span.SetAttributes(attribute.String(fmt.Sprintf("log.params.%s", args[i]), fmt.Sprintf("%+v", args[i+1])))
	}
	
	_, file, line, ok := runtime.Caller(callerSkip)
	if ok {
		span.SetAttributes(attribute.String("caller.file", file))
		span.SetAttributes(attribute.Int("caller.line", line))
	}
	
	if level == "ERROR" {
		stake := ""
		for i := 0; i < 10; i++ {
			_, file, line, ok := runtime.Caller(callerSkip + i)
			if ok {
				stake = stake + fmt.Sprintf("%s:%d\n", file, line)
			} else {
				break
			}
		}
		span.SetAttributes(attribute.String("caller.stake", stake[:len(stake)-1]))
		span.SetStatus(codes.Error, msg)
	}
}

func (l *OpentelLogger) Info(ctx context.Context, msg string, args ...any) {
	l.log(ctx, "INFO", 2, msg, args...)
	slog.Info(msg, args...)
}

func (l *OpentelLogger) Error(ctx context.Context, msg string, args ...any) {
	l.log(ctx, "ERROR", 2, msg, args...)
	slog.Error(msg, args...)
}

func (l *OpentelLogger) Debug(ctx context.Context, msg string, args ...any) {
	l.log(ctx, "DEBUG", 2, msg, args...)
	slog.Debug(msg, args...)
}

func (l *OpentelLogger) Warn(ctx context.Context, msg string, args ...any) {
	l.log(ctx, "WARN", 2, msg, args...)
	slog.Warn(msg, args...)
}

func (l *OpentelLogger) DebugContext(msg string, args ...any) {
	l.Debug(context.Background(), msg, args...)
}

func (l *OpentelLogger) InfoContext(msg string, args ...any) {
	l.Info(context.Background(), msg, args...)
}

func (l *OpentelLogger) WarnContext(msg string, args ...any) {
	l.Warn(context.Background(), msg, args...)
}

func (l *OpentelLogger) ErrorContext(msg string, args ...any) {
	l.Error(context.Background(), msg, args...)
}
