package golib

import (
	"gitlab.id.vin/vincart/golib/pubsub"
	"gitlab.id.vin/vincart/golib/web/event"
	"gitlab.id.vin/vincart/golib/web/listener"
)

func WithEventBusAutoConfig(eventMapping map[pubsub.Event][]pubsub.Subscriber) Module {
	return func(app *App) {
		var debugLog pubsub.DebugLog
		if app.Logger != nil {
			debugLog = app.Logger.Debugf
		}

		publisher := pubsub.NewPublisher()
		pubsub.RegisterGlobal(publisher)
		app.Publisher = publisher

		bus := pubsub.NewEventBus(publisher, debugLog)
		subscribeEvents(bus, eventMapping)
		go bus.Run()
	}
}

func subscribeEvents(bus *pubsub.EventBus, eventMapping map[pubsub.Event][]pubsub.Subscriber) {
	if eventMapping != nil {
		for e, subscribers := range eventMapping {
			bus.Subscribe(e, subscribers...)
		}
	}
	bus.Subscribe(new(event.RequestCompletedEvent), new(listener.RequestCompletedLogListener))
}
