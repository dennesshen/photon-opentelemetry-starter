package opentelLog

import (
	"context"
	"log/slog"
	
	"github.com/dennesshen/photon-core-starter/configuration"
	"github.com/dennesshen/photon-core-starter/log"
)

func Start(ctx context.Context) (log.CoreLogger, error) {
	slog.Info("init otel log")
	
	// log 設定
	config, err := configuration.Get[Config](ctx)
	if err != nil {
		return nil, err
	}
	logLevel := log.TransferLogLevel(config.Log.LogLevel)
	
	return NewOpentelLogger(logLevel), nil
}
