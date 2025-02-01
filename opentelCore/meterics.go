package opentelCore

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"
	
	"github.com/dennesshen/photon-core-starter/utils/convert"
	
	"github.com/prometheus/client_golang/prometheus/promhttp"
	runtimeMetrics "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	metr "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

type Int64UpDownCounter struct {
	metr.Int64UpDownCounter
}

type Int64UpCounter struct {
	metr.Int64Counter
}

type Int64Guage struct {
	metr.Int64ObservableGauge
}

var (
	provider       *metric.MeterProvider
	exporter       *prometheus.Exporter
	upDownCounters = make(map[string]Int64UpDownCounter)
	upCounters     = make(map[string]Int64UpCounter)
	upDownMutex    sync.Mutex
	upMutex        sync.Mutex
)

func startMetricService(config *Config) {
	if config.OpenTel.MetricPath == nil && config.OpenTel.MetricPort == nil {
		return
	}
	slog.Info("[otel]start metric service")
	if err := runtimeMetrics.Start(); err != nil {
		slog.Error("[otel]start otel metric service error", "error", err)
	}
	exporter = configureMetrics()
	http.Handle(*config.OpenTel.MetricPath, promhttp.Handler())
	slog.Info("[otel]metric listening", "port", *config.OpenTel.MetricPort, "path", *config.OpenTel.MetricPath)
	go func() {
		server := &http.Server{
			Addr:              ":" + *config.OpenTel.MetricPort,
			ReadHeaderTimeout: 3 * time.Second,
		}
		if err := server.ListenAndServe(); err != nil {
			slog.Error("[otel]start otel metric http service error", "error", err)
		}
	}()
}

func configureMetrics() *prometheus.Exporter {
	exporter, err := prometheus.New()
	if err != nil {
		slog.Error("set metric service error", "error", err)
	}
	
	provider = metric.NewMeterProvider(metric.WithReader(exporter))
	
	return exporter
}

/*
Add meter up/down count
*/
func AddMeterUpDownCount(action, description, unit string, value int64, attributes ...attribute.KeyValue) {
	upDownMutex.Lock()
	defer upDownMutex.Unlock()
	counter, ok := upDownCounters[action]
	if !ok {
		meter := provider.Meter(serviceName)
		c, err := meter.Int64UpDownCounter(
			action,
			metr.WithDescription(description),
			metr.WithUnit(unit),
		)
		if err != nil {
			slog.Error("[otel]add meter counter error", "error", err)
		}
		counter.Int64UpDownCounter = c
		upDownCounters[action] = counter
	}
	ctx := context.Background()
	counter.addMeterUpDownCount(ctx, value, metr.WithAttributes(attributes...))
}

func (c *Int64UpDownCounter) addMeterUpDownCount(ctx context.Context, value int64, options ...metr.AddOption) {
	c.Int64UpDownCounter.Add(ctx, value, options...)
}

func AddMeterCount(action, description, unit string, value uint64, attributes ...attribute.KeyValue) {
	upMutex.Lock()
	defer upMutex.Unlock()
	counter, ok := upCounters[action]
	if !ok {
		meter := provider.Meter(serviceName)
		c, err := meter.Int64Counter(
			action,
			metr.WithDescription(description),
			metr.WithUnit(unit),
		)
		if err != nil {
			slog.Error("[otel]add meter counter error", "error", err)
		}
		counter.Int64Counter = c
		upCounters[action] = counter
	}
	ctx := context.Background()
	counter.addMeterCount(ctx, value, metr.WithAttributes(attributes...))
}

func (c *Int64UpCounter) addMeterCount(ctx context.Context, value uint64, options ...metr.AddOption) {
	c.Int64Counter.Add(ctx, convert.Signed[uint64, int64](value), options...)
}
