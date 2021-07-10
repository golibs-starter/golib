package pubsub

type EventProducer interface {
	ProduceEvent() chan Event
}

type EventBus struct {
	producer      EventProducer
	eventMappings map[string][]Subscriber
}

func NewEventBus(eventProducer EventProducer) *EventBus {
	return &EventBus{
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
		//log.Info("Event ", event.GetName(), " was fired")
		for _, subscriber := range b.eventMappings[event.GetName()] {
			go subscriber.Handler(event)
		}
	}
}
