package golib

import (
	"gitlab.id.vin/vincart/golib/actuator"
	"gitlab.id.vin/vincart/golib/config"
	webActuator "gitlab.id.vin/vincart/golib/web/actuator"
	"go.uber.org/fx"
)

func ActuatorEndpointOpt() fx.Option {
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
	Endpoint        *webActuator.Endpoint
	HealthService   actuator.HealthService
	InformerService actuator.InfoService
}

// NewActuatorEndpoint Initiate actuator endpoint with
// health checker and informer automatically.
//
// ================= Health Checker ======================
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
// or using ProvideHealthChecker
//   func NewSampleHealthChecker() actuator.HealthChecker {
//   	return &SampleHealthChecker{}
//   }
//   ProvideHealthChecker(NewSampleHealthChecker)
//
//
// =================== Informer =========================
// Similar to Health Checker, an Informer also registered by produce an actuator.Informer.
// For example, a GitRevision provider produce the following result:
//   type GitRevisionOut struct {
//      fx.Out
//      Informer actuator.Informer `group:"actuator_informer"`
//   }
//   func NewGitRevision() (GitRevisionOut, error) {}
//
// or using ProvideInformer
//   func NewSampleInformer() actuator.Informer {
//   	return &SampleInformer{}
//   }
//   ProvideInformer(NewSampleInformer)
//
func NewActuatorEndpoint(in ActuatorIn) ActuatorOut {
	healthService := actuator.NewDefaultHealthService(in.Checkers)
	infoService := actuator.NewDefaultInfoService(in.Props, in.Informers)
	return ActuatorOut{
		Endpoint:        webActuator.NewEndpoint(healthService, infoService),
		HealthService:   healthService,
		InformerService: infoService,
	}
}

// ProvideHealthChecker A simple way to provide a health checker
func ProvideHealthChecker(healthCheckerConstructor interface{}) fx.Option {
	return fx.Provide(fx.Annotated{Group: "actuator_health_checker", Target: healthCheckerConstructor})
}

// ProvideInformer A simple way to provide an informer
func ProvideInformer(informerConstructor interface{}) fx.Option {
	return fx.Provide(fx.Annotated{Group: "actuator_informer", Target: informerConstructor})
}
