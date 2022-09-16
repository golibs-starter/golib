package pubsub

import (
	assert "github.com/stretchr/testify/require"
	"gitlab.com/golibs-starter/golib/pubsub/executor"
	"reflect"
	"testing"
	"time"
)

func TestNewDefaultEventBus_WhenNoOpts_ShouldUseDefaultValue(t *testing.T) {
	bus := NewDefaultEventBus()
	assert.NotNil(t, bus.subscribers)
	assert.Len(t, bus.subscribers, 0)
	assert.NotNil(t, bus.mapSubscribers)
	assert.Len(t, bus.mapSubscribers, 0)
	assert.NotNil(t, bus.eventCh)
	assert.Len(t, bus.eventCh, 0)
	assert.Equal(t, reflect.ValueOf(defaultDebugLog).Pointer(), reflect.ValueOf(bus.debugLog).Pointer())
	assert.IsType(t, new(executor.AsyncExecutor), bus.executor)
}

func TestNewDefaultEventBus_WhenUseOpts_ShouldSetOptCorrectly(t *testing.T) {
	var logger DebugLog = func(e Event, msgFormat string, args ...interface{}) {
	}
	syncExecutor := executor.NewSyncExecutor()
	bus := NewDefaultEventBus(WithEventBusDebugLog(logger), WithEventExecutor(syncExecutor))
	assert.NotNil(t, bus.subscribers)
	assert.Len(t, bus.subscribers, 0)
	assert.NotNil(t, bus.mapSubscribers)
	assert.Len(t, bus.mapSubscribers, 0)
	assert.NotNil(t, bus.eventCh)
	assert.Len(t, bus.eventCh, 0)
	assert.Equal(t, reflect.ValueOf(logger).Pointer(), reflect.ValueOf(bus.debugLog).Pointer())
	assert.Equal(t, syncExecutor, bus.executor)
}

type DummySubscriber1 struct {
	eventRun        map[string]bool
	orderedEventRun []string
}

func (d DummySubscriber1) Supports(event Event) bool {
	return true
}

func (d *DummySubscriber1) Handle(event Event) {
	if d.eventRun == nil {
		d.eventRun = make(map[string]bool)
	}
	if d.orderedEventRun == nil {
		d.orderedEventRun = make([]string, 0)
	}
	d.eventRun[event.Name()] = true
	d.orderedEventRun = append(d.orderedEventRun, event.Name())
}

type DummySubscriber2 struct {
}

func (d DummySubscriber2) Supports(event Event) bool {
	return true
}

func (d DummySubscriber2) Handle(event Event) {
}

type DummyEvent struct {
	name string
}

func (d DummyEvent) Identifier() string {
	return "dummy-event-id"
}

func (d DummyEvent) Name() string {
	return d.name
}

func (d DummyEvent) Payload() interface{} {
	return map[string]string{
		"a": "b",
	}
}

func (d DummyEvent) String() string {
	return "dummy-event-string"
}

func TestDefaultEventBus_WhenRegisterSubscribers_ShouldRegisterCorrectly(t *testing.T) {
	bus := NewDefaultEventBus()
	s1 := DummySubscriber1{}
	s2 := DummySubscriber2{}
	bus.Register(&s1, &s2)
	assert.Len(t, bus.mapSubscribers, 2)
	assert.True(t, bus.mapSubscribers["pubsub.DummySubscriber1"])
	assert.True(t, bus.mapSubscribers["pubsub.DummySubscriber2"])
	assert.Len(t, bus.subscribers, 2)
	assert.Equal(t, &s1, bus.subscribers[0])
	assert.Equal(t, &s2, bus.subscribers[1])
}

func TestDefaultEventBus_GivenAsyncExecutor_WhenDeliverEvent_ShouldRunCorrectly(t *testing.T) {
	bus := NewDefaultEventBus()
	s1 := DummySubscriber1{}
	bus.Register(&s1)
	go bus.Run()
	bus.Deliver(&DummyEvent{name: "event-1"})
	bus.Deliver(&DummyEvent{name: "event-2"})
	bus.Deliver(&DummyEvent{name: "event-3"})
	bus.Deliver(&DummyEvent{name: "event-4"})
	time.Sleep(20 * time.Millisecond)
	assert.Len(t, s1.eventRun, 4)
	assert.True(t, s1.eventRun["event-1"])
	assert.True(t, s1.eventRun["event-2"])
	assert.True(t, s1.eventRun["event-3"])
	assert.True(t, s1.eventRun["event-4"])
	assert.Len(t, s1.orderedEventRun, 4)
	assert.Contains(t, s1.orderedEventRun, "event-1")
	assert.Contains(t, s1.orderedEventRun, "event-2")
	assert.Contains(t, s1.orderedEventRun, "event-3")
	assert.Contains(t, s1.orderedEventRun, "event-4")
}

func TestDefaultEventBus_GivenSyncExecutor_WhenDeliverEvent_ShouldRunCorrectly(t *testing.T) {
	bus := NewDefaultEventBus(WithEventExecutor(executor.NewSyncExecutor()))
	s1 := DummySubscriber1{}
	bus.Register(&s1)
	go bus.Run()
	bus.Deliver(&DummyEvent{name: "event-1"})
	bus.Deliver(&DummyEvent{name: "event-2"})
	bus.Deliver(&DummyEvent{name: "event-3"})
	bus.Deliver(&DummyEvent{name: "event-4"})
	time.Sleep(20 * time.Millisecond)
	assert.Len(t, s1.eventRun, 4)
	assert.True(t, s1.eventRun["event-1"])
	assert.True(t, s1.eventRun["event-2"])
	assert.True(t, s1.eventRun["event-3"])
	assert.True(t, s1.eventRun["event-4"])
	assert.Len(t, s1.orderedEventRun, 4)
	assert.Equal(t, "event-1", s1.orderedEventRun[0])
	assert.Equal(t, "event-2", s1.orderedEventRun[1])
	assert.Equal(t, "event-3", s1.orderedEventRun[2])
	assert.Equal(t, "event-4", s1.orderedEventRun[3])
}
