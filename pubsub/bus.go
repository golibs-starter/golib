package pubsub

type EventBus interface {

	// Register subscriber(s) with the bus
	Register(subscribers ...Subscriber)

	// Deliver an event
	Deliver(event Event)

	// Run the bus
	Run()

	// Stop the bus
	Stop()
}
