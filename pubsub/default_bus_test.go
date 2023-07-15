package pubsub

import (
	"context"
	assert "github.com/stretchr/testify/require"
	"gitlab.com/golibs-starter/golib/pubsub/executor"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestNewDefaultEventBus_WhenNoOpts_ShouldUseDefaultValue(t *testing.T) {
	bus := NewDefaultEventBus()
	assert.NotNil(t, bus.subscribers)
	assert.Len(t, bus.subscribers, 0)
	assert.NotNil(t, bus.subscribers)
	assert.Len(t, bus.subscribers, 0)
	assert.Equal(t, bus.eventChSize, 0)
	assert.NotNil(t, bus.eventCh)
	assert.Len(t, bus.eventCh, 0)
	assert.Equal(t, cap(bus.eventCh), 0)
	assert.Equal(t, reflect.ValueOf(defaultDebugLog).Pointer(), reflect.ValueOf(bus.debugLog).Pointer())
	assert.IsType(t, new(executor.AsyncExecutor), bus.executor)
}

func TestNewDefaultEventBus_WhenUseOpts_ShouldSetOptCorrectly(t *testing.T) {
	var logger DebugLog = func(ctx context.Context, msgFormat string, args ...interface{}) {
	}
	syncExecutor := executor.NewSyncExecutor()
	bus := NewDefaultEventBus(
		WithEventBusDebugLog(logger),
		WithEventExecutor(syncExecutor),
		WithEventChannelSize(12),
	)
	assert.NotNil(t, bus.subscribers)
	assert.Len(t, bus.subscribers, 0)
	assert.NotNil(t, bus.subscribers)
	assert.Len(t, bus.subscribers, 0)
	assert.Equal(t, bus.eventChSize, 12)
	assert.NotNil(t, bus.eventCh)
	assert.Len(t, bus.eventCh, 0)
	assert.Equal(t, cap(bus.eventCh), 12)
	assert.Equal(t, reflect.ValueOf(logger).Pointer(), reflect.ValueOf(bus.debugLog).Pointer())
	assert.Equal(t, syncExecutor, bus.executor)
}

type DummySubscriber1 struct {
	eventRun          map[string]bool
	eventRunMu        sync.RWMutex
	orderedEventRun   []string
	orderedEventRunMu sync.RWMutex
}

func (d *DummySubscriber1) Supports(event Event) bool {
	return true
}

func (d *DummySubscriber1) Handle(event Event) {
	d.eventRunMu.Lock()
	if d.eventRun == nil {
		d.eventRun = make(map[string]bool)
	}
	d.eventRun[event.Name()] = true
	d.eventRunMu.Unlock()

	d.orderedEventRunMu.Lock()
	if d.orderedEventRun == nil {
		d.orderedEventRun = make([]string, 0)
	}
	d.orderedEventRun = append(d.orderedEventRun, event.Name())
	d.orderedEventRunMu.Unlock()
}

func (d *DummySubscriber1) eventRan(eventName string) bool {
	d.eventRunMu.RLock()
	defer d.eventRunMu.RUnlock()
	return d.eventRun[eventName]
}

func (d *DummySubscriber1) numberOfEventRan() int {
	d.eventRunMu.RLock()
	defer d.eventRunMu.RUnlock()
	return len(d.eventRun)
}

func (d *DummySubscriber1) numberOfOrderedEventRun() int {
	d.orderedEventRunMu.RLock()
	defer d.orderedEventRunMu.RUnlock()
	return len(d.orderedEventRun)
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
	ctx  context.Context
}

func (d DummyEvent) Identifier() string {
	return "dummy-event-id"
}

func (d DummyEvent) Context() context.Context {
	return d.ctx
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
	assert.Len(t, bus.subscribers, 2)
	sub1, s1Exists := bus.subscribers["pubsub.DummySubscriber1"]
	sub2, s2Exists := bus.subscribers["pubsub.DummySubscriber2"]
	assert.True(t, s1Exists)
	assert.True(t, s2Exists)
	assert.Len(t, bus.subscribers, 2)
	assert.Equal(t, &s1, sub1)
	assert.Equal(t, &s2, sub2)
}

func TestDefaultEventBus_GivenAsyncExecutor_WhenDeliverEvent_ShouldRunCorrectly(t *testing.T) {
	bus := NewDefaultEventBus()
	s1 := DummySubscriber1{}
	bus.Register(&s1)
	bus.Run()
	defer bus.Stop()
	bus.Deliver(&DummyEvent{name: "event-1"})
	bus.Deliver(&DummyEvent{name: "event-2"})
	bus.Deliver(&DummyEvent{name: "event-3"})
	bus.Deliver(&DummyEvent{name: "event-4"})
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 4, s1.numberOfEventRan())
	assert.True(t, s1.eventRan("event-1"))
	assert.True(t, s1.eventRan("event-2"))
	assert.True(t, s1.eventRan("event-3"))
	assert.True(t, s1.eventRan("event-4"))
	assert.Equal(t, 4, s1.numberOfOrderedEventRun())
	assert.Contains(t, s1.orderedEventRun, "event-1")
	assert.Contains(t, s1.orderedEventRun, "event-2")
	assert.Contains(t, s1.orderedEventRun, "event-3")
	assert.Contains(t, s1.orderedEventRun, "event-4")
}

func TestDefaultEventBus_GivenSyncExecutor_WhenDeliverEvent_ShouldRunCorrectly(t *testing.T) {
	bus := NewDefaultEventBus(WithEventExecutor(executor.NewSyncExecutor()))
	s1 := DummySubscriber1{}
	bus.Register(&s1)
	bus.Run()
	bus.Deliver(&DummyEvent{name: "event-1"})
	bus.Deliver(&DummyEvent{name: "event-2"})
	bus.Deliver(&DummyEvent{name: "event-3"})
	bus.Deliver(&DummyEvent{name: "event-4"})
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 4, s1.numberOfEventRan())
	assert.True(t, s1.eventRan("event-1"))
	assert.True(t, s1.eventRan("event-2"))
	assert.True(t, s1.eventRan("event-3"))
	assert.True(t, s1.eventRan("event-4"))
	assert.Equal(t, 4, s1.numberOfOrderedEventRun())
	assert.Equal(t, "event-1", s1.orderedEventRun[0])
	assert.Equal(t, "event-2", s1.orderedEventRun[1])
	assert.Equal(t, "event-3", s1.orderedEventRun[2])
	assert.Equal(t, "event-4", s1.orderedEventRun[3])
}
