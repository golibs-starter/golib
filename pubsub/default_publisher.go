package pubsub

type DefaultPublisher struct {
	bus                    EventBus
	debugLog               DebugLog
	notLogPayloadForEvents map[string]bool
}

func NewDefaultPublisher(bus EventBus, opts ...PublisherOpt) *DefaultPublisher {
	pub := &DefaultPublisher{bus: bus}
	for _, opt := range opts {
		opt(pub)
	}
	if pub.debugLog == nil {
		pub.debugLog = defaultDebugLog
	}
	return pub
}

func (p *DefaultPublisher) Publish(event Event) {
	p.bus.Deliver(event)
	if p.notLogPayloadForEvents != nil && p.notLogPayloadForEvents[event.Name()] {
		p.debugLog(event, "Event [%s] was fired with id [%s]", event.Name(), event.Identifier())
	} else {
		p.debugLog(event, "Event [%s] was fired with id [%s], payload [%+v]",
			event.Name(), event.Identifier(), event.Payload())
	}
}
