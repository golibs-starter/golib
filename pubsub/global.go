package pubsub

var _bus EventBus = NewDefaultEventBus()
var _publisher Publisher = NewDefaultPublisher(_bus)

func GetEventBus() EventBus {
	return _bus
}

func GetPublisher() Publisher {
	return _publisher
}

func Register(subscribers ...Subscriber) {
	_bus.Register(subscribers...)
}

func Run() {
	_bus.Run()
}

func Publish(event Event) {
	_publisher.Publish(event)
}

func ReplaceGlobal(bus EventBus, publisher Publisher) {
	_bus = bus
	_publisher = publisher
}
