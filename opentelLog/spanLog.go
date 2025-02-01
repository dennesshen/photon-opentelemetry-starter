package opentelLog

import (
	"context"
	
	"github.com/dennesshen/photon-opentelemetry-starter/opentelCore"
	
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type spanLog struct {
	span trace.Span
}

func StartSpanLog(ctx context.Context, traceName string, attributes ...attribute.KeyValue) *spanLog {
	_, span := opentelCore.NewTrace(ctx, traceName)
	span.SetAttributes(attributes...)
	return &spanLog{span: span}
}

func (sl *spanLog) SetAttributes(attributes ...attribute.KeyValue) {
	if sl.span != nil {
		sl.span.SetAttributes(attributes...)
	}
}

func (sl *spanLog) Finish(attributes ...attribute.KeyValue) {
	if sl.span != nil {
		sl.span.SetAttributes(attributes...)
		sl.span.End()
	}
}
