package example

// ==================================================
// ==== Example about how to bootstrap your app =====
// ==================================================

import (
	"github.com/golibs-starter/golib"
	"github.com/golibs-starter/golib/pubsub"
	"github.com/golibs-starter/golib/pubsub/executor"
	"go.uber.org/fx"
)

func All() fx.Option {
	return fx.Options(
		// Required
		golib.AppOpt(),
		golib.PropertiesOpt(),

		// When you want to use default logging strategy.
		golib.LoggingOpt(),

		// When you want to enable event publisher
		golib.EventOpt(),
		// When you want handle event in simple synchronous way
		golib.SupplyEventBusOpt(pubsub.WithEventExecutor(executor.NewSyncExecutor())),
		// Or want a custom executor, such as using worker pool
		fx.Provide(NewSampleEventExecutor),
		golib.ProvideEventBusOpt(func(executor *SampleEventExecutor) pubsub.EventBusOpt {
			return pubsub.WithEventExecutor(executor)
		}),

		// When you want to enable http request log
		golib.HttpRequestLogOpt(),

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

		// Graceful shutdown.
		// OnStop hooks will run in reverse order.
		golib.OnStopEventOpt(),
	)
}
