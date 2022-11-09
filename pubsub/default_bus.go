package pubsub

import (
	"gitlab.com/golibs-starter/golib/pubsub/executor"
	"gitlab.com/golibs-starter/golib/utils"
)

type DefaultEventBus struct {
	debugLog       DebugLog
	subscribers    []Subscriber
	mapSubscribers map[string]bool
	eventChSize    int
	eventCh        chan Event
	executor       Executor
}

func NewDefaultEventBus(opts ...EventBusOpt) *DefaultEventBus {
	bus := &DefaultEventBus{
		subscribers:    make([]Subscriber, 0),
		mapSubscribers: make(map[string]bool),
	}
	for _, opt := range opts {
		opt(bus)
	}
	if bus.debugLog == nil {
		bus.debugLog = defaultDebugLog
	}
	if bus.eventChSize < 0 {
		bus.eventChSize = 0
	}
	if bus.eventCh == nil {
		bus.eventCh = make(chan Event, bus.eventChSize)
	}
	if bus.executor == nil {
		bus.executor = executor.NewAsyncExecutor()
	}
	return bus
}

func (b *DefaultEventBus) Register(subscribers ...Subscriber) {
	for _, subscriber := range subscribers {
		subscriberId := utils.GetStructFullname(subscriber)
		if _, exists := b.mapSubscribers[subscriberId]; exists {
			b.debugLog(nil, "Subscriber [%s] already registered", subscriberId)
			continue
		}
		b.mapSubscribers[subscriberId] = true
		b.subscribers = append(b.subscribers, subscriber)
		b.debugLog(nil, "Register subscriber [%s] successful", subscriberId)
	}
}

func (b *DefaultEventBus) Deliver(event Event) {
	b.eventCh <- event
}

func (b *DefaultEventBus) Run() {
	for {
		event := <-b.eventCh
		for _, subscriber := range b.subscribers {
			if subscriber.Supports(event) {
				subscriber := subscriber
				b.executor.Execute(func() {
					subscriber.Handle(event)
				})
			}
		}
	}
}
