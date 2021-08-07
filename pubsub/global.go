package pubsub

var _publisher Publisher = NewPublisher()

func ReplaceGlobal(publisher Publisher) {
	_publisher = publisher
}

func Publish(event Event) {
	_publisher.Publish(event)
}
