package example

import (
	"context"
	"gitlab.com/golibs-starter/golib/pubsub"
	"gitlab.com/golibs-starter/golib/web/client"
	"gitlab.com/golibs-starter/golib/web/log"
)

// ==================================================
// ======== Example about inject dependencies =======
// ==================================================

// NewSampleService In this case Contextual Http Client is required
// Use fx.Provide(NewSampleService) to register a service
func NewSampleService(httpClient client.ContextualHttpClient) *SampleService {
	return &SampleService{httpClient: httpClient}
}

type SampleService struct {
	httpClient client.ContextualHttpClient
}

func (s SampleService) DoSomething(ctx context.Context) error {
	// You can write log with current context
	log.Info(ctx, "Write something to log with context")

	// Then pass the context to ContextualHttpClient's call
	var result struct{}
	_, err := s.httpClient.Get(ctx, "https://example", &result)
	if err != nil {
		log.Error(ctx, "Http client call with error [%s]", err)
		return err
	}

	// Even pass the context to an event
	pubsub.Publish(NewSampleEvent(ctx, &SampleEventMessage{
		Field1: "val1",
		Field2: "val2",
	}))
	return nil
}
