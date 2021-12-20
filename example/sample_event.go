package example

// ==================================================
// ========= Example about declare an event =========
// ==================================================

import (
	"context"
	"gitlab.com/golibs-starter/golib/web/event"
)

// NewSampleEvent In order to get more tracing content,
// inject the request context to your event.
// Then using pubsub.Publish(NewSampleEvent(ctx, &SampleEventMessage{})) to publish it
func NewSampleEvent(ctx context.Context, payload *SampleEventMessage) *SampleEvent {
	return &SampleEvent{
		AbstractEvent: event.NewAbstractEvent(ctx, "SampleEvent"),
		PayloadData:   payload,
	}
}

type SampleEvent struct {
	*event.AbstractEvent
	PayloadData *SampleEventMessage `json:"payload"`
}

func (a SampleEvent) Payload() interface{} {
	return a.PayloadData
}

func (a SampleEvent) String() string {
	return a.ToString(a)
}

type SampleEventMessage struct {
	Field1 string
	Field2 string
}
