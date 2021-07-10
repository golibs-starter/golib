package pubsub

type Subscriber interface {
	Handler(event Event)
}
