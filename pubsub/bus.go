package pubsub

import (
	"reflect"
)

type EventProducer interface {
	ProduceEvent() chan Event
}

type EventBus struct {
	debugLog       DebugLog
	producer       EventProducer
	subscribers    []Subscriber
	mapSubscribers map[string]bool
}

func NewEventBus(eventProducer EventProducer, opts ...EventBusOpt) *EventBus {
	bus := &EventBus{
		producer:       eventProducer,
		subscribers:    make([]Subscriber, 0),
		mapSubscribers: make(map[string]bool),
	}
	for _, opt := range opts {
		opt(bus)
	}
	if bus.debugLog == nil {
		bus.debugLog = defaultDebugLog
	}
	return bus
}

func (b *EventBus) Register(subscribers ...Subscriber) {
	for _, subscriber := range subscribers {
		subscriberId := reflect.TypeOf(subscriber).String()
		if _, exists := b.mapSubscribers[subscriberId]; exists {
			b.debugLog("Subscriber [%s] already registered", subscriberId)
			continue
		}
		b.mapSubscribers[subscriberId] = true
		b.subscribers = append(b.subscribers, subscriber)
		b.debugLog("Register subscriber [%s] successful", subscriberId)
	}
}

func (b *EventBus) Run() {
	for {
		event := <-b.producer.ProduceEvent()
		for _, subscriber := range b.subscribers {
			if subscriber.Supports(event) {
				go subscriber.Handle(event)
			}
		}
	}
}
