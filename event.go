package golib

import (
	"gitlab.id.vin/vincart/golib/log"
	"gitlab.id.vin/vincart/golib/pubsub"
	"gitlab.id.vin/vincart/golib/utils"
	"gitlab.id.vin/vincart/golib/web/listener"
	"go.uber.org/fx"
)

type EventListener interface {
	pubsub.Subscriber

	// Events List of events that
	// this listener will subscribe
	Events() []pubsub.Event
}

type EventAutoConfigOut struct {
	fx.Out
	Bus              *pubsub.EventBus
	Publisher        pubsub.Publisher
	DefaultListeners []EventListener `group:"event_listener,flatten"`
}

func NewEventAutoConfig(logger log.Logger) EventAutoConfigOut {
	var debugLog pubsub.DebugLog = logger.Debugf
	publisher := pubsub.NewPublisher()
	bus := pubsub.NewEventBus(publisher, debugLog)
	return EventAutoConfigOut{
		Bus:       bus,
		Publisher: publisher,
		DefaultListeners: []EventListener{
			new(listener.RequestCompletedLogListener),
		},
	}
}

type RegisterEventAutoConfigIn struct {
	fx.In
	Logger    log.Logger
	Bus       *pubsub.EventBus
	Publisher pubsub.Publisher
	Listeners []EventListener `group:"event_listener"`
}

func RegisterEventAutoConfig(in RegisterEventAutoConfigIn) {
	pubsub.ReplaceGlobal(in.Publisher)

	// Subscribe events
	for _, eventListener := range in.Listeners {
		if subscribedEvents := eventListener.Events(); subscribedEvents != nil {
			for _, e := range subscribedEvents {
				in.Logger.Infof("[%s] is subscribed event [%s]", utils.GetTypeName(eventListener), e.Name())
				in.Bus.Subscribe(e, eventListener)
			}
		}
	}

	go in.Bus.Run()
}
