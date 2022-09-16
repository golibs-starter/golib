package pubsub

import (
	assert "github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestReplaceGlobal(t *testing.T) {
	assert.NotNil(t, _bus)
	assert.NotNil(t, _publisher)
	assert.Equal(t, _bus, GetEventBus())
	assert.Equal(t, _publisher, GetPublisher())

	newBus := NewDefaultEventBus()
	newPub := NewDefaultPublisher(newBus)
	ReplaceGlobal(newBus, newPub)
	assert.Equal(t, newBus, _bus)
	assert.Equal(t, newPub, _publisher)
	assert.Equal(t, newBus, GetEventBus())
	assert.Equal(t, newPub, GetPublisher())
}

func TestGlobalRegister(t *testing.T) {
	s1 := DummySubscriber1{}
	s2 := DummySubscriber2{}
	Register(&s1, &s2)
	assert.Len(t, _bus.(*DefaultEventBus).mapSubscribers, 2)
	assert.True(t, _bus.(*DefaultEventBus).mapSubscribers["pubsub.DummySubscriber1"])
	assert.True(t, _bus.(*DefaultEventBus).mapSubscribers["pubsub.DummySubscriber2"])
	assert.Len(t, _bus.(*DefaultEventBus).subscribers, 2)
	assert.Equal(t, &s1, _bus.(*DefaultEventBus).subscribers[0])
	assert.Equal(t, &s2, _bus.(*DefaultEventBus).subscribers[1])
}

func TestGlobalPublish(t *testing.T) {
	s1 := DummySubscriber1{}
	Register(&s1)
	go Run()

	Publish(&DummyEvent{name: "event-1"})
	Publish(&DummyEvent{name: "event-2"})

	time.Sleep(200 * time.Millisecond)

	assert.Len(t, s1.eventRun, 2)
	assert.True(t, s1.eventRun["event-1"])
	assert.True(t, s1.eventRun["event-2"])

	assert.Len(t, s1.orderedEventRun, 2)
	assert.Contains(t, s1.orderedEventRun, "event-1")
	assert.Contains(t, s1.orderedEventRun, "event-2")
}
