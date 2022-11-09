package pubsub

import (
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
	var logger DebugLog = func(e Event, msgFormat string, args ...interface{}) {
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
	var logger DebugLog = func(e Event, msgFormat string, args ...interface{}) {
		logMsgs[e.Name()] = msgFormat
	}

	bus := NewDefaultEventBus()
	s1 := DummySubscriber1{}
	bus.Register(&s1)
	bus.Run()

	pub := NewDefaultPublisher(bus,
		WithPublisherDebugLog(logger),
		WithPublisherNotLogPayload([]string{"event-2"}),
	)
	pub.Publish(&DummyEvent{name: "event-1"})
	pub.Publish(&DummyEvent{name: "event-2"})

	time.Sleep(100 * time.Millisecond)

	assert.Len(t, s1.eventRun, 2)
	assert.True(t, s1.eventRun["event-1"])
	assert.True(t, s1.eventRun["event-2"])

	assert.Len(t, s1.orderedEventRun, 2)
	assert.Contains(t, s1.orderedEventRun, "event-1")
	assert.Contains(t, s1.orderedEventRun, "event-2")

	// Test log message
	assert.Len(t, logMsgs, 2)
	assert.Contains(t, logMsgs["event-1"], "payload")
	assert.NotContains(t, logMsgs["event-2"], "payload")
}
