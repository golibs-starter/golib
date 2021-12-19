package example

import (
	"context"
	"gitlab.id.vin/vincart/golib/pubsub"
	"gitlab.id.vin/vincart/golib/web/client"
)

// ==================================================
// ======== Example about inject dependencies =======
// ==================================================

// NewSampleService In this case Contextual Http Client is required
func NewSampleService(httpClient client.ContextualHttpClient) *SampleService {
	return &SampleService{httpClient: httpClient}
}

type SampleService struct {
	httpClient client.ContextualHttpClient
}

func (s SampleService) DoSomething(ctx context.Context) {

	// Sample for publish an event
	pubsub.Publish(NewSampleEvent(ctx, &SampleEventMessage{
		Field1: "val1",
		Field2: "val2",
	}))
}
