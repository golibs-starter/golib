package pubsub

type EventBusOpt func(bus *DefaultEventBus)

func WithEventBusDebugLog(debugLog DebugLog) EventBusOpt {
	return func(bus *DefaultEventBus) {
		bus.debugLog = debugLog
	}
}

func WithEventExecutor(executor Executor) EventBusOpt {
	return func(bus *DefaultEventBus) {
		bus.executor = executor
	}
}
