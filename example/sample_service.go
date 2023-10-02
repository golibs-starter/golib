package example

import (
	"context"
	"github.com/golibs-starter/golib/log"
	"github.com/golibs-starter/golib/log/field"
	"github.com/golibs-starter/golib/pubsub"
	"github.com/golibs-starter/golib/web/client"
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
	log.WithCtx(ctx).Info("Write something to log with context")

	// Then pass the context to ContextualHttpClient's call
	var result struct{}
	_, err := s.httpClient.Get(ctx, "https://example", &result)
	if err != nil {
		log.WithCtx(ctx).WithErrors(err).Error("Http client call failed")
		return err
	}

	log.WithCtx(ctx).WithField(field.String("field1", "value1")).
		Infof("This is a log with single field and message format: %s", "example-value")

	log.WithCtx(ctx).WithField(
		field.String("field1", "value1"),
		field.Any("field2", []string{"value1", "value2"}),
	).Info("This is log with multiple fields")

	log.WithCtx(ctx).WithAny("field1", "value1").
		Info("This is log with simple key value field")

	// Even pass the context to an event
	pubsub.Publish(NewSampleEvent(ctx, &SampleEventMessage{
		Field1: "val1",
		Field2: "val2",
	}))
	return nil
}
