package pubsub

type PublisherOpt func(pub *publisher)

func WithPublisherDebugLog(debugLog DebugLog) PublisherOpt {
	return func(pub *publisher) {
		pub.debugLog = debugLog
	}
}

func WithPublisherNotLogPayload(events []string) PublisherOpt {
	return func(pub *publisher) {
		pub.notLogPayloadForEvents = make(map[string]bool)
		for _, e := range events {
			pub.notLogPayloadForEvents[e] = true
		}
	}
}
