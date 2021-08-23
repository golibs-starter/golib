package golib

import (
	"gitlab.id.vin/vincart/golib/config"
	"gitlab.id.vin/vincart/golib/web/actuator"
	"go.uber.org/fx"
)

func ActuatorEndpointAutoConfig() fx.Option {
	return fx.Provide(NewActuatorEndpoint)
}

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

// NewActuatorEndpoint Initiate actuator endpoint with
// health checker and informer automatically.
//
// To register a Health Checker, your component have to
// produce an actuator.HealthChecker with group `actuator_health_checker`
// in the result of provider function.
// For example, a redis provider produce the following result:
//   type RedisOut struct {
//      fx.Out
//      Client        *redis.Client
//      HealthChecker actuator.HealthChecker `group:"actuator_health_checker"`
//   }
//   func NewRedis() (RedisOut, error) {}
//
// Similar to Health Checker, an Informer also registered by produce an actuator.Informer.
// For example, a GitRevision provider produce the following result:
//   type GitRevisionOut struct {
//      fx.Out
//      Informer actuator.Informer `group:"actuator_informer"`
//   }
//   func NewGitRevision() (GitRevisionOut, error) {}
func NewActuatorEndpoint(in ActuatorIn) ActuatorOut {
	healthService := actuator.NewDefaultHealthService(in.Checkers)
	infoService := actuator.NewDefaultInfoService(in.Props, in.Informers)
	return ActuatorOut{
		Endpoint:        actuator.NewEndpoint(healthService, infoService),
		HealthService:   healthService,
		InformerService: infoService,
	}
}