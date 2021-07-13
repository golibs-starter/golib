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
	if _, exist := b.eventMappings[event.GetName()]; exist == false {
		b.eventMappings[event.GetName()] = make([]Subscriber, 0)
	}
	b.eventMappings[event.GetName()] = append(b.eventMappings[event.GetName()], subscriber...)
}

func (b *EventBus) Run() {
	for {
		event := <-b.producer.ProduceEvent()
		if b.logger != nil {
			b.logger.Debugf("Event [%s] was fired with payload [%s]", event.GetName(), event.String())
		}
		for _, subscriber := range b.eventMappings[event.GetName()] {
			go subscriber.Handler(event)
		}
	}
}
