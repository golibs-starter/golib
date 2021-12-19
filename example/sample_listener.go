package example

import (
	"gitlab.id.vin/vincart/golib/pubsub"
)

// ==================================================
// === Example about declare listener (subscriber) ==
// ==================================================

// NewSampleListener
// Use golib.ProvideEventListener(NewSampleListener) to declare a listener
func NewSampleListener() pubsub.Subscriber {
	return &SampleListener{}
}

type SampleListener struct {
}

func (s SampleListener) Supports(e pubsub.Event) bool {
	_, ok := e.(*SampleEvent)
	return ok
}

func (s SampleListener) Handle(e pubsub.Event) {
	// Handle when receive event
}
