package opentelCore

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func NewTrace(ctx context.Context, traceName string) (context.Context, trace.Span) {
	tracer := otel.Tracer(serviceName)
	return tracer.Start(ctx, traceName)
}

func NewServerTrace(ctx context.Context, traceName string) (context.Context, trace.Span) {
	tracer := otel.Tracer(serviceName)
	return tracer.Start(ctx, traceName, trace.WithSpanKind(trace.SpanKindServer))
}

func NewTraceFromHTTP(r *http.Request, traceName string) (*http.Request, trace.Span) {
	ctx, span := NewTrace(r.Context(), traceName)
	return r.WithContext(ctx), span
}
