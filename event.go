package golib

import (
	"context"
	"gitlab.com/golibs-starter/golib/event"
	"gitlab.com/golibs-starter/golib/pubsub"
	"gitlab.com/golibs-starter/golib/web/log"
	"go.uber.org/fx"
)

func EventOpt() fx.Option {
	return fx.Options(
		ProvideProps(event.NewProperties),

		SupplyEventBusOpt(pubsub.WithEventBusDebugLog(log.Debuge)),
		ProvideEventBusOpt(func(props *event.Properties) pubsub.EventBusOpt {
			return pubsub.WithEventChannelSize(props.ChannelSize)
		}),
		fx.Provide(NewDefaultEventBus),
		ProvideInformer(pubsub.NewDefaultBusInformer),

		SupplyEventPublisherOpt(pubsub.WithPublisherDebugLog(log.Debuge)),
		ProvideEventPublisherOpt(func(props *event.Properties) pubsub.PublisherOpt {
			return pubsub.WithPublisherNotLogPayload(props.Log.NotLogPayloadForEvents)
		}),
		fx.Provide(NewDefaultEventPublisher),

		fx.Invoke(RegisterEventPublisher),
		fx.Invoke(RunEventBus),
	)
}

func OnStopEventOpt() fx.Option {
	return fx.Invoke(func(lc fx.Lifecycle, bus pubsub.EventBus) {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				bus.Stop()
				return nil
			},
		})
	})
}

func ProvideEventListener(listener interface{}) fx.Option {
	return fx.Provide(fx.Annotated{Group: "event_listener", Target: listener})
}

func SupplyEventBusOpt(opt pubsub.EventBusOpt) fx.Option {
	return fx.Supply(fx.Annotated{Group: "event_bus_opt", Target: opt})
}

func ProvideEventBusOpt(optConstructor interface{}) fx.Option {
	return fx.Provide(fx.Annotated{Group: "event_bus_opt", Target: optConstructor})
}

type EventBusIn struct {
	fx.In
	Options []pubsub.EventBusOpt `group:"event_bus_opt"`
}

func NewDefaultEventBus(in EventBusIn) pubsub.EventBus {
	return pubsub.NewDefaultEventBus(in.Options...)
}

func SupplyEventPublisherOpt(opt pubsub.PublisherOpt) fx.Option {
	return fx.Supply(fx.Annotated{Group: "event_publisher_opt", Target: opt})
}

func ProvideEventPublisherOpt(optConstructor interface{}) fx.Option {
	return fx.Provide(fx.Annotated{Group: "event_publisher_opt", Target: optConstructor})
}

type EventPublisherIn struct {
	fx.In
	Bus     pubsub.EventBus
	Options []pubsub.PublisherOpt `group:"event_publisher_opt"`
}

func NewDefaultEventPublisher(in EventPublisherIn) pubsub.Publisher {
	return pubsub.NewDefaultPublisher(in.Bus, in.Options...)
}

type RegisterEventPublisherIn struct {
	fx.In
	Bus         pubsub.EventBus
	Publisher   pubsub.Publisher
	Subscribers []pubsub.Subscriber `group:"event_listener"`
}

func RegisterEventPublisher(in RegisterEventPublisherIn) {
	pubsub.ReplaceGlobal(in.Bus, in.Publisher)
	in.Bus.Register(in.Subscribers...)
}

func RunEventBus(bus pubsub.EventBus) {
	bus.Run()
}
