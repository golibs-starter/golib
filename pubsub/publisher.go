package pubsub

type Publisher interface {

	// Publish an event
	Publish(event Event)
}
