package pubsub

type Publisher interface {
	Publish(event Event)
}

type publisher struct {
	eventCh chan Event
}

func NewPublisher() *publisher {
	return &publisher{eventCh: make(chan Event)}
}

func (p *publisher) Publish(event Event) {
	p.eventCh <- event
}

func (p *publisher) ProduceEvent() chan Event {
	return p.eventCh
}
