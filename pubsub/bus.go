package pubsub

type EventProducer interface {
	ProduceEvent() chan Event
}

type EventBus struct {
	logger        Logger
	producer      EventProducer
	eventMappings map[string][]Subscriber
}

func NewEventBus(eventProducer EventProducer, logger Logger) *EventBus {
	return &EventBus{
		logger:        logger,
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
		if b.logger != nil {
			b.logger.Debugf("Event [%s] was fired with payload [%s]", event.Name(), event.String())
		}
		for _, subscriber := range b.eventMappings[event.Name()] {
			go subscriber.Handler(event)
		}
	}
}
