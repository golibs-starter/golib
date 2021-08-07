package golib

import (
	"gitlab.id.vin/vincart/golib/log"
	"gitlab.id.vin/vincart/golib/pubsub"
	"gitlab.id.vin/vincart/golib/web/listener"
)

type EventListener interface {
	pubsub.Subscriber

	// Events List of events that
	// this listener will be subscribed on
	Events() []pubsub.Event
}

func NewEventAutoConfig(logger log.Logger, listeners ...EventListener) (*pubsub.EventBus, pubsub.Publisher) {
	var debugLog pubsub.DebugLog = logger.Debugf
	publisher := pubsub.NewPublisher()
	bus := pubsub.NewEventBus(publisher, debugLog)
	registerListeners(bus, listeners)
	return bus, publisher
}

func RegisterEventAutoConfig(bus *pubsub.EventBus, publisher pubsub.Publisher) {
	pubsub.ReplaceGlobal(publisher)
	go bus.Run()
}

func registerListeners(bus *pubsub.EventBus, listeners []EventListener) {
	if listeners == nil {
		listeners = make([]EventListener, 0)
	}
	listeners = appendPredefinedListeners(listeners)
	for _, eventListener := range listeners {
		if subscribedEvents := eventListener.Events(); subscribedEvents != nil {
			for _, e := range subscribedEvents {
				bus.Subscribe(e, eventListener)
			}
		}
	}
}

func appendPredefinedListeners(listeners []EventListener) []EventListener {
	return append(listeners, new(listener.RequestCompletedLogListener))
}
