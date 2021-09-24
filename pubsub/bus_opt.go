package pubsub

type EventBusOpt func(bus *EventBus)

func WithEventBusDebugLog(debugLog DebugLog) EventBusOpt {
	return func(bus *EventBus) {
		bus.debugLog = debugLog
	}
}
