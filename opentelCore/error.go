package opentelCore

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func recordErrorMsg(span trace.Span, err error) {
	span.SetAttributes(attribute.String("error.msg", err.Error()))
	span.RecordError(err, trace.WithAttributes(attribute.String("error", err.Error())))
}

func recordError(span trace.Span, msg string, err error) {
	if err != nil {
		recordErrorMsg(span, err)
	}
	span.SetStatus(codes.Error, msg)
}

func recordSpanError(span trace.Span, err error) {
	recordErrorMsg(span, err)
	span.SetStatus(codes.Error, err.Error())
}

func ErrorWithSpan(span trace.Span, err error) {
	if err == nil {
		return
	}
	recordSpanError(span, err)
}

func ErrorWithCtx(ctx context.Context, err error) {
	if err == nil {
		return
	}
	span := trace.SpanFromContext(ctx)
	recordSpanError(span, err)
}
