package example

// ==================================================
// ========= Example about declare an event =========
// ==================================================

import (
	"context"
	baseEvent "gitlab.com/golibs-starter/golib/event"
	"gitlab.com/golibs-starter/golib/web/event"
)

// NewSampleEvent In order to get more tracing content,
// inject the request context to your event.
// Then using pubsub.Publish(NewSampleEvent(ctx, &SampleEventMessage{})) to publish it
func NewSampleEvent(ctx context.Context, payload *SampleEventMessage) *SampleEvent {
	return &SampleEvent{event.NewAbstractEvent(ctx, "SampleEvent", baseEvent.WithPayload(payload))}
}

type SampleEvent struct {
	*event.AbstractEvent
}

func (a SampleEvent) String() string {
	return a.ToString(a)
}

type SampleEventMessage struct {
	Field1 string
	Field2 string
}
