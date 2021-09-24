package golib

import (
	"gitlab.id.vin/vincart/golib/event"
	"gitlab.id.vin/vincart/golib/log"
	"gitlab.id.vin/vincart/golib/pubsub"
	"gitlab.id.vin/vincart/golib/web/listener"
	"go.uber.org/fx"
)

func EventOpt() fx.Option {
	return fx.Options(
		ProvideEventListener(listener.NewRequestCompletedLogListener),
		ProvideProps(event.NewProperties),
		fx.Provide(NewEventPublisher),
		fx.Invoke(RegisterEventPublisher),
		fx.Invoke(RunEventBus),
	)
}

func ProvideEventListener(listener interface{}) fx.Option {
	return fx.Provide(fx.Annotated{Group: "event_listener", Target: listener})
}

type EventPublisherOut struct {
	fx.Out
	Publisher pubsub.Publisher
	EventBus  *pubsub.EventBus
}

func NewEventPublisher(logger log.Logger, props *event.Properties) EventPublisherOut {
	publisher := pubsub.NewPublisher(
		pubsub.WithPublisherDebugLog(logger.Debugf),
		pubsub.WithPublisherNotLogPayload(props.Log.NotLogPayloadForEvents),
	)
	eventBus := pubsub.NewEventBus(publisher, pubsub.WithEventBusDebugLog(logger.Debugf))
	return EventPublisherOut{
		EventBus:  eventBus,
		Publisher: publisher,
	}
}

type RegisterEventAutoConfigIn struct {
	fx.In
	EventBus    *pubsub.EventBus
	Publisher   pubsub.Publisher
	Subscribers []pubsub.Subscriber `group:"event_listener"`
}

func RegisterEventPublisher(in RegisterEventAutoConfigIn) {
	pubsub.ReplaceGlobal(in.Publisher)
	in.EventBus.Register(in.Subscribers...)
}

func RunEventBus(eventBus *pubsub.EventBus) {
	go eventBus.Run()
}
