package golib

import (
	"gitlab.com/golibs-starter/golib/event"
	"gitlab.com/golibs-starter/golib/pubsub"
	"gitlab.com/golibs-starter/golib/web/log"
	"go.uber.org/fx"
)

func EventOpt() fx.Option {
	return fx.Options(
		ProvideProps(event.NewProperties),

		SupplyEventBusOpt(pubsub.WithEventBusDebugLog(log.Debuge)),
		fx.Provide(NewDefaultEventBus),

		SupplyEventPublisherOpt(pubsub.WithPublisherDebugLog(log.Debuge)),
		ProvideEventPublisherOpt(func(props *event.Properties) pubsub.PublisherOpt {
			return pubsub.WithPublisherNotLogPayload(props.Log.NotLogPayloadForEvents)
		}),
		fx.Provide(NewDefaultEventPublisher),

		fx.Invoke(RegisterEventPublisher),
		fx.Invoke(RunEventBus),
	)
}

func ProvideEventListener(listener interface{}) fx.Option {
	return fx.Provide(fx.Annotated{Group: "event_listener", Target: listener})
}

func SupplyEventBusOpt(opt pubsub.EventBusOpt) fx.Option {
	return fx.Provide(fx.Annotated{Group: "event_bus_opt", Target: fx.Supply(opt)})
}

type EventBusIn struct {
	Options []pubsub.EventBusOpt `group:"event_bus_opt"`
}

func NewDefaultEventBus(in EventBusIn) pubsub.EventBus {
	return pubsub.NewDefaultEventBus(in.Options...)
}

func SupplyEventPublisherOpt(opt pubsub.PublisherOpt) fx.Option {
	return fx.Provide(fx.Annotated{Group: "event_publisher_opt", Target: fx.Supply(opt)})
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
	go bus.Run()
}
