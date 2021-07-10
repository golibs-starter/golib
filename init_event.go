package golib

import (
	"gitlab.id.vin/vincart/golib/pubsub"
	"gitlab.id.vin/vincart/golib/web/event"
	"gitlab.id.vin/vincart/golib/web/listeners"
)

func InitEventBus(eventMapping map[pubsub.Event][]pubsub.Subscriber) {
	publisher := pubsub.NewPublisher()
	pubsub.RegisterGlobal(publisher)
	bus := pubsub.NewEventBus(publisher)
	subscribeEvents(bus, eventMapping)
	go bus.Run()
}

func subscribeEvents(bus *pubsub.EventBus, eventMapping map[pubsub.Event][]pubsub.Subscriber) {
	if eventMapping != nil {
		for e, subscribers := range eventMapping {
			bus.Subscribe(e, subscribers...)
		}
	}
	bus.Subscribe(new(event.RequestCompletedEvent), new(listeners.RequestCompletedLogListener))
}
