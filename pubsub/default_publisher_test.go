package pubsub

import (
	"context"
	assert "github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

func TestNewDefaultPublisher_WhenNoOpts_ShouldUseDefaultValue(t *testing.T) {
	bus := NewDefaultEventBus()
	pub := NewDefaultPublisher(bus)
	assert.Equal(t, bus, pub.bus)
	assert.Equal(t, reflect.ValueOf(defaultDebugLog).Pointer(), reflect.ValueOf(pub.debugLog).Pointer())
	assert.Nil(t, pub.notLogPayloadForEvents)
}

func TestNewDefaultPublisher_WhenUseOpts_ShouldSetOptCorrectly(t *testing.T) {
	var logger DebugLog = func(ctx context.Context, msgFormat string, args ...interface{}) {
	}
	var eventNotLogPayloads = []string{"event-test"}
	bus := NewDefaultEventBus()
	pub := NewDefaultPublisher(
		bus,
		WithPublisherDebugLog(logger),
		WithPublisherNotLogPayload(eventNotLogPayloads),
	)
	assert.Equal(t, bus, pub.bus)
	assert.Equal(t, reflect.ValueOf(logger).Pointer(), reflect.ValueOf(pub.debugLog).Pointer())
	assert.Equal(t, map[string]bool{"event-test": true}, pub.notLogPayloadForEvents)
}

func TestDefaultPublisher_WhenPublishEvent_ShouldPublishCorrectly(t *testing.T) {
	logMsgs := make(map[string]string)
	var logger DebugLog = func(ctx context.Context, msgFormat string, args ...interface{}) {
		event := ctx.Value("event").(string)
		logMsgs[event] = msgFormat
	}

	bus := NewDefaultEventBus()
	s1 := DummySubscriber1{}
	bus.Register(&s1)
	bus.Run()

	pub := NewDefaultPublisher(
		bus,
		WithPublisherDebugLog(logger),
		WithPublisherNotLogPayload([]string{"event-2"}),
	)
	pub.Publish(&DummyEvent{name: "event-1", ctx: context.WithValue(context.Background(), "event", "event-1")})
	pub.Publish(&DummyEvent{name: "event-2", ctx: context.WithValue(context.Background(), "event", "event-2")})

	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, 2, s1.numberOfEventRan())
	assert.True(t, s1.eventRan("event-1"))
	assert.True(t, s1.eventRan("event-2"))

	assert.Equal(t, 2, s1.numberOfOrderedEventRun())
	assert.Contains(t, s1.orderedEventRun, "event-1")
	assert.Contains(t, s1.orderedEventRun, "event-2")

	// Test log message
	assert.Len(t, logMsgs, 2)
	assert.Contains(t, logMsgs["event-1"], "payload")
	assert.NotContains(t, logMsgs["event-2"], "payload")
}
