package golib

import (
	"gitlab.id.vin/vincart/common/golib/event"
	"gitlab.id.vin/vincart/common/golib/pubsub"
	"gitlab.id.vin/vincart/common/golib/web/listener"
	"gitlab.id.vin/vincart/common/golib/web/log"
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

func NewEventPublisher(props *event.Properties) EventPublisherOut {
	publisher := pubsub.NewPublisher(
		pubsub.WithPublisherDebugLog(log.Debuge),
		pubsub.WithPublisherNotLogPayload(props.Log.NotLogPayloadForEvents),
	)
	eventBus := pubsub.NewEventBus(publisher, pubsub.WithEventBusDebugLog(log.Debuge))
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
