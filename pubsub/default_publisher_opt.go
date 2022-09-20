package pubsub

type PublisherOpt func(pub *DefaultPublisher)

func WithPublisherDebugLog(debugLog DebugLog) PublisherOpt {
	return func(pub *DefaultPublisher) {
		pub.debugLog = debugLog
	}
}

func WithPublisherNotLogPayload(events []string) PublisherOpt {
	return func(pub *DefaultPublisher) {
		pub.notLogPayloadForEvents = make(map[string]bool)
		for _, e := range events {
			pub.notLogPayloadForEvents[e] = true
		}
	}
}
