package opentelCore

import (
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	tp          *sdktrace.TracerProvider
	serviceName string
)

type Protocol uint
