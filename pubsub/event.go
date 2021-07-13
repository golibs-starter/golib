package pubsub

type Event interface {

	// GetName returns event name of current event
	GetName() string

	// GetPayload returns event payload of current event
	GetPayload() interface{}

	// String convert event data to string
	String() string
}
