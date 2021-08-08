package golib

import (
	"gitlab.id.vin/vincart/golib/config"
	"gitlab.id.vin/vincart/golib/web/actuator"
	"go.uber.org/fx"
)

type ActuatorIn struct {
	fx.In
	Props     *config.AppProperties
	Checkers  []actuator.HealthChecker `group:"actuator_health_checker"`
	Informers []actuator.Informer      `group:"actuator_informer"`
}

type ActuatorOut struct {
	fx.Out
	Endpoint        *actuator.Endpoint
	HealthService   actuator.HealthService
	InformerService actuator.InfoService
}

func NewActuatorEndpointAutoConfig(in ActuatorIn) ActuatorOut {
	healthService := actuator.NewDefaultHealthService(in.Checkers)
	infoService := actuator.NewDefaultInfoService(in.Props, in.Informers)
	return ActuatorOut{
		Endpoint:        actuator.NewEndpoint(healthService, infoService),
		HealthService:   healthService,
		InformerService: infoService,
	}
}
