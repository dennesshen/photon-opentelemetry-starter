package opentelStarter

import (
	"github.com/dennesshen/photon-core-starter/core"
	"github.com/dennesshen/photon-core-starter/log"
	"github.com/dennesshen/photon-opentelemetry-starter/opentelCore"
	"github.com/dennesshen/photon-opentelemetry-starter/opentelLog"
)

func init() {
	log.RegisterInitAction(opentelLog.Start)
	core.RegisterCoreDependency(opentelCore.Start)
	core.RegisterShutdownCoreDependency(opentelCore.Shutdown)
}
