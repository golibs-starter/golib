package pubsub

type Publisher interface {
	Publish(event Event)
}

type publisher struct {
	eventCh                chan Event
	debugLog               DebugLog
	notLogPayloadForEvents map[string]bool
}

func NewPublisher(opts ...PublisherOpt) *publisher {
	pub := &publisher{eventCh: make(chan Event)}
	for _, opt := range opts {
		opt(pub)
	}
	if pub.debugLog == nil {
		pub.debugLog = defaultDebugLog
	}
	return pub
}

func (p *publisher) Publish(event Event) {
	p.eventCh <- event
	if p.notLogPayloadForEvents != nil && p.notLogPayloadForEvents[event.Name()] {
		p.debugLog(event, "Event [%s] was fired with id [%s]", event.Name(), event.Identifier())
	} else {
		p.debugLog(event, "Event [%s] was fired with id [%s], payload [%s]",
			event.Name(), event.Identifier(), event.String())
	}
}

func (p *publisher) ProduceEvent() chan Event {
	return p.eventCh
}
