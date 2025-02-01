package opentelCore

import (
	"context"
	"encoding/json"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type header struct {
	md metadata.MD
}

func (h *header) Get(key string) string {
	if len(h.md) == 0 {
		return ""
	}
	values := h.md[key]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (h *header) Set(key, val string) {
	h.md.Set(key, val)
}

func (h *header) Keys() []string {
	keys := make([]string, 0, len(h.md))
	for k := range h.md {
		keys = append(keys, k)
	}
	return keys
}

func UnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			ctx = otel.GetTextMapPropagator().Extract(ctx, &header{md})
		}
		ctx, span := NewTrace(ctx, "api start")
		defer span.End()
		if req != nil {
			reqBytes, _ := json.Marshal(req)
			span.SetAttributes(attribute.String("request.body", string(reqBytes)))
		}
		p, _ := peer.FromContext(ctx)
		if p != nil {
			span.SetAttributes(attribute.String("client.ip", p.Addr.String()))
		}
		result, err := handler(ctx, req)
		if err != nil {
			recordError(span, "response error", err)
		}
		if result != nil {
			resultBytes, err := json.Marshal(result)
			if err != nil {
				recordError(span, "decode response error", err)
			}
			span.SetAttributes(attribute.String("response.body", string(resultBytes)))
		}
		return result, err
	})
}

func GrpcServerHandler() grpc.ServerOption {
	return grpc.StatsHandler(otelgrpc.NewServerHandler())
}

func GrpcClientHandler() grpc.DialOption {
	return grpc.WithStatsHandler(otelgrpc.NewClientHandler())
}
