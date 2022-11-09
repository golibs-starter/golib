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
	newBus := NewDefaultEventBus()
	newPub := NewDefaultPublisher(newBus)
	ReplaceGlobal(newBus, newPub)

	s1 := DummySubscriber1{}
	s2 := DummySubscriber2{}
	Register(&s1, &s2)
	sub1, sub1e := _bus.(*DefaultEventBus).subscribers["pubsub.DummySubscriber1"]
	sub2, sub2e := _bus.(*DefaultEventBus).subscribers["pubsub.DummySubscriber2"]
	assert.Len(t, _bus.(*DefaultEventBus).subscribers, 2)
	assert.True(t, sub1e)
	assert.True(t, sub2e)
	assert.Len(t, _bus.(*DefaultEventBus).subscribers, 2)
	assert.Equal(t, &s1, sub1)
	assert.Equal(t, &s2, sub2)
}

func TestGlobalPublish(t *testing.T) {
	newBus := NewDefaultEventBus()
	newPub := NewDefaultPublisher(newBus)
	ReplaceGlobal(newBus, newPub)

	s1 := DummySubscriber1{}
	Register(&s1)
	Run()

	Publish(&DummyEvent{name: "event-1"})
	Publish(&DummyEvent{name: "event-2"})

	time.Sleep(100 * time.Millisecond)

	assert.Len(t, s1.eventRun, 2)
	assert.True(t, s1.eventRun["event-1"])
	assert.True(t, s1.eventRun["event-2"])

	assert.Len(t, s1.orderedEventRun, 2)
	assert.Contains(t, s1.orderedEventRun, "event-1")
	assert.Contains(t, s1.orderedEventRun, "event-2")
}
