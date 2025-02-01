package opentelCore

import (
	"context"
	"log/slog"
	"strings"
	"time"
	
	"github.com/dennesshen/photon-core-starter/configuration"
	
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 初始化 opentelmetery
func Start(startContext context.Context) (err error) {
	slog.Info("init otel")
	
	// TODO 那個defer cancel()是幹嘛的 ?
	ctx, cancel := context.WithTimeout(startContext, time.Second*5)
	defer cancel()
	
	// log 設定
	config, err := configuration.Get[Config](ctx)
	if err != nil {
		return
	}
	
	// 初始化 trace exporter
	initTraceExporter(ctx, &config)
	
	// 啟動 metric service
	startMetricService(&config)
	
	return
}

func initTraceExporter(ctx context.Context, config *Config) {
	var err error
	var traceExporter *otlptrace.Exporter
	if strings.HasPrefix(config.OpenTel.URL, "grpc://") {
		traceExporter, err = connectGrpc(ctx, config.OpenTel.URL[7:])
	} else {
		traceExporter, err = otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(config.OpenTel.URL))
	}
	if err != nil {
		slog.Error("[otel]dial to otel collect server failed", "error", err, "url", config.OpenTel.URL)
	}
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	serviceSemName := semconv.ServiceNameKey.String(config.OpenTel.ServiceName)
	tp = sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(bsp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			serviceSemName,
		)),
	)
	otel.SetTracerProvider(tp)
	
	otel.SetTextMapPropagator(propagation.TraceContext{})
}

// 取得 gRpc 連線
func connectGrpc(ctx context.Context, host string) (*otlptrace.Exporter, error) {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("[otel]dial to otel grpc server failed", "error", err, "host", host)
	}
	
	return otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
}
