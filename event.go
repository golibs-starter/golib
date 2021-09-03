package golib

import (
	"gitlab.id.vin/vincart/golib/event"
	"gitlab.id.vin/vincart/golib/log"
	"gitlab.id.vin/vincart/golib/pubsub"
	"gitlab.id.vin/vincart/golib/utils"
	"gitlab.id.vin/vincart/golib/web/listener"
	"go.uber.org/fx"
)

func EventOpt() fx.Option {
	return fx.Options(
		ProvideEventListener(listener.NewRequestCompletedLogListener),
		fx.Provide(NewEventPublisher),
		fx.Invoke(RegisterEventPublisher),
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

func NewEventPublisher(logger log.Logger) EventPublisherOut {
	publisher := pubsub.NewPublisher()
	eventBus := pubsub.NewEventBus(publisher, logger.Debugf)
	return EventPublisherOut{
		EventBus:  eventBus,
		Publisher: publisher,
	}
}

type RegisterEventAutoConfigIn struct {
	fx.In
	Logger    log.Logger
	EventBus  *pubsub.EventBus
	Publisher pubsub.Publisher
	Listeners []event.Listener `group:"event_listener"`
}

func RegisterEventPublisher(in RegisterEventAutoConfigIn) {
	pubsub.ReplaceGlobal(in.Publisher)

	// Subscribe events
	for _, eventListener := range in.Listeners {
		if subscribedEvents := eventListener.Events(); subscribedEvents != nil {
			for _, e := range subscribedEvents {
				in.Logger.Infof("[%s] is subscribed event [%s]", utils.GetTypeName(eventListener), e.Name())
				in.EventBus.Subscribe(e, eventListener)
			}
		}
	}

	go in.EventBus.Run()
}
