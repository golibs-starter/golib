package example

// ==================================================
// ==== Example about how to bootstrap your app =====
// ==================================================

import (
	"gitlab.id.vin/vincart/golib"
	"go.uber.org/fx"
)

func All() []fx.Option {
	return []fx.Option{
		// Required
		golib.AppOpt(),

		// Required
		golib.PropertiesOpt(),

		// When you want to use default logging strategy.
		golib.LoggingOpt(),

		// When you want to enable event publisher
		golib.EventOpt(),

		// When you want to enable actuator endpoints.
		// By default, we provide HealthService and InfoService.
		golib.ActuatorEndpointOpt(),
		// When you want to provide build info to above InfoService.
		golib.BuildInfoOpt(Version, CommitHash, BuildTime),
		// When you want to provide custom health checker and informer
		golib.ProvideHealthChecker(NewSampleHealthChecker),
		golib.ProvideInformer(NewSampleInformer),

		// When you want to enable http client auto config with contextual client by default
		golib.HttpClientOpt(),

		// When you want to tell GoLib to load your properties.
		golib.ProvideProps(NewSampleProperties),

		// When you want to declare a service
		fx.Provide(NewSampleService),

		// When you want to register your event listener.
		golib.ProvideEventListener(NewSampleListener),
	}
}
