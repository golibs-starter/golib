package pubsub

import "context"

type Event interface {
	// Identifier returns the ID of event
	Identifier() string

	// Name returns event name of current event
	Name() string

	// Context returns the current context of event
	Context() context.Context

	// Payload returns event payload of current event
	Payload() interface{}

	// String convert event data to string
	String() string
}
