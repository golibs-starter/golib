package event

import "gitlab.id.vin/vincart/golib/pubsub"

type Listener interface {
	pubsub.Subscriber

	// Events List of events that
	// this listener will subscribe
	Events() []pubsub.Event
}
