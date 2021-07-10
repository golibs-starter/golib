package pubsub

var _publisher Publisher

func RegisterGlobal(publisher Publisher) {
	_publisher = publisher
}

func Publish(event Event) {
	_publisher.Publish(event)
}
