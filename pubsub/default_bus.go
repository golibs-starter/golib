package pubsub

import (
	"context"
	"gitlab.com/golibs-starter/golib/pubsub/executor"
	"gitlab.com/golibs-starter/golib/utils"
	"sync"
)

type DefaultEventBus struct {
	debugLog    DebugLog
	subscribers map[string]Subscriber
	eventChSize int
	eventCh     chan Event
	stopCh      chan bool
	executor    Executor
	wg          sync.WaitGroup
}

func NewDefaultEventBus(opts ...EventBusOpt) *DefaultEventBus {
	bus := &DefaultEventBus{
		subscribers: make(map[string]Subscriber, 0),
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
	bus.stopCh = make(chan bool)
	return bus
}

func (b *DefaultEventBus) Register(subscribers ...Subscriber) {
	for _, subscriber := range subscribers {
		subscriberId := utils.GetStructFullname(subscriber)
		if _, exists := b.subscribers[subscriberId]; exists {
			b.debugLog(context.Background(), "Subscriber [%s] already registered", subscriberId)
			continue
		}
		b.subscribers[subscriberId] = subscriber
		b.debugLog(context.Background(), "Register subscriber [%s] successful", subscriberId)
	}
}

func (b *DefaultEventBus) Deliver(event Event) {
	b.eventCh <- event
}

func (b *DefaultEventBus) Run() {
	b.debugLog(context.Background(), "Default event bus is starting")
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		for {
			select {
			case event := <-b.eventCh:
				for _, subscriber := range b.subscribers {
					if subscriber.Supports(event) {
						subscriber := subscriber
						b.executor.Execute(func() {
							subscriber.Handle(event)
						})
					}
				}
				break
			case isStop := <-b.stopCh:
				if isStop {
					return
				}
				break
			}
		}
	}()
	b.debugLog(context.Background(), "Default event bus is started")
}

func (b *DefaultEventBus) Stop() {
	b.debugLog(context.Background(), "Default event bus is stopping")
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		for {
			if len(b.eventCh) == 0 {
				b.stopCh <- true
				return
			}
		}
	}()
	b.wg.Wait()
	b.debugLog(context.Background(), "Default event bus is stopped")
}
