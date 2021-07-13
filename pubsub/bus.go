package pubsub

import (
	"fmt"
)

type EventProducer interface {
	ProduceEvent() chan Event
}

type EventBus struct {
	debugLog      DebugLog
	producer      EventProducer
	eventMappings map[string][]Subscriber
}

func NewEventBus(eventProducer EventProducer, debugLog DebugLog) *EventBus {
	if debugLog == nil {
		debugLog = func(msgFormat string, args ...interface{}) {
			_, _ = fmt.Printf(msgFormat+"\n", args...)
		}
	}
	return &EventBus{
		debugLog:      debugLog,
		producer:      eventProducer,
		eventMappings: make(map[string][]Subscriber),
	}
}

func (b *EventBus) Subscribe(event Event, subscriber ...Subscriber) {
	if _, exist := b.eventMappings[event.Name()]; exist == false {
		b.eventMappings[event.Name()] = make([]Subscriber, 0)
	}
	b.eventMappings[event.Name()] = append(b.eventMappings[event.Name()], subscriber...)
}

func (b *EventBus) Run() {
	for {
		event := <-b.producer.ProduceEvent()
		if b.debugLog != nil {
			b.debugLog("Event [%s] was fired with payload [%s]", event.Name(), event.String())
		}
		for _, subscriber := range b.eventMappings[event.Name()] {
			go subscriber.Handler(event)
		}
	}
}
