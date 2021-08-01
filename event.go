package golib

import (
	"gitlab.id.vin/vincart/golib/pubsub"
	"gitlab.id.vin/vincart/golib/web/listener"
)

type EventListener interface {
	pubsub.Subscriber

	// Events List of events that
	// this listener will be subscribed on
	Events() []pubsub.Event
}

func WithEventAutoConfig(listeners ...EventListener) Module {
	return func(app *App) {
		var debugLog pubsub.DebugLog
		if app.Logger != nil {
			debugLog = app.Logger.Debugf
		}

		publisher := pubsub.NewPublisher()
		pubsub.RegisterGlobal(publisher)
		app.Publisher = publisher

		bus := pubsub.NewEventBus(publisher, debugLog)
		registerListeners(bus, listeners)
		go bus.Run()
	}
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
