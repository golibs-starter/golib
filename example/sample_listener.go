package example

import (
	"github.com/golibs-starter/golib/log"
	"github.com/golibs-starter/golib/pubsub"
)

// ==================================================
// === Example about declare listener (subscriber) ==
// ==================================================

// NewSampleListener
// Use golib.ProvideEventListener(NewSampleListener) to declare a listener
func NewSampleListener(service *SampleService) pubsub.Subscriber {
	return &SampleListener{service: service}
}

type SampleListener struct {
	service *SampleService
}

func (s SampleListener) Supports(e pubsub.Event) bool {
	_, ok := e.(*SampleEvent)
	return ok
}

func (s SampleListener) Handle(e pubsub.Event) {
	// You can use event as a log context
	// Note that the context only appear when your event embeds web AbstractEvent
	log.WithCtx(e.Context()).Info("A log with context")

	// Cast to concrete event
	sampleEvent := e.(*SampleEvent)

	// You can get context in the web abstract event directly
	log.WithCtx(sampleEvent.Context()).Info("Another log with context")

	// Then pass the context to the next call
	_ = s.service.DoSomething(sampleEvent.Context())
}
