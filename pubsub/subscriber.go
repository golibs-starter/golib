package pubsub

type Subscriber interface {

	// Supports indicates whether an event is supported by this Subscriber or not.
	Supports(event Event) bool

	// Handle a supported Event that indicated by the Supports function.
	Handle(event Event)
}
