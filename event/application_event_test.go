package event

import (
	assert "github.com/stretchr/testify/require"
	"testing"
)

func TestNewApplicationEvent_WhenNoOpts_ShouldInitCorrectly(t *testing.T) {
	eventName := "TestEvent"
	e := NewApplicationEvent(eventName)
	assert.Equal(t, eventName, e.Name())
	assert.Equal(t, eventName, e.Event)
	assert.NotEmpty(t, e.Id)
	assert.Equal(t, e.Id, e.Identifier())
	assert.Greater(t, e.Timestamp, int64(0))
	assert.Equal(t, DefaultEventSource, e.Source)
	assert.Empty(t, e.ServiceCode)
	assert.Nil(t, e.AdditionalData)
	assert.Nil(t, e.PayloadData)
	assert.Nil(t, e.Payload())
}

func TestNewApplicationEvent_WhenHasOpts_ShouldInitCorrectly(t *testing.T) {
	eventName := "TestEvent"
	payload := map[string]string{"a": "a"}
	e := NewApplicationEvent(eventName,
		WithId("test-id"),
		WithServiceCode("test-service-code"),
		WithSource("test-source"),
		WithAdditionalData(map[string]interface{}{
			"key1": "val1",
		}),
		WithPayload(payload),
	)
	assert.Equal(t, eventName, e.Name())
	assert.Equal(t, eventName, e.Event)
	assert.Equal(t, "test-id", e.Identifier())
	assert.Equal(t, "test-id", e.Id)
	assert.Greater(t, e.Timestamp, int64(0))
	assert.Equal(t, "test-source", e.Source)
	assert.Equal(t, "test-service-code", e.ServiceCode)
	assert.Equal(t, map[string]interface{}{
		"key1": "val1",
	}, e.AdditionalData)
	assert.Equal(t, payload, e.PayloadData)
	assert.Equal(t, payload, e.Payload())
}

func TestApplicationEvent_ToString(t *testing.T) {
	e1 := ApplicationEvent{
		Id:             "1",
		Event:          "TEST",
		Source:         "NOT_USED",
		ServiceCode:    "service-test",
		AdditionalData: map[string]interface{}{"a": "b"},
		PayloadData:    map[string]string{"x": "y"},
		Timestamp:      10,
	}
	assert.Equal(t, `{"id":"1","event":"TEST","source":"NOT_USED","service_code":"service-test","additional_data":{"a":"b"},"payload":{"x":"y"},"timestamp":10}`, e1.String())

	e2 := ApplicationEvent{
		Id:             "1",
		Event:          "TEST",
		Source:         "",
		ServiceCode:    "",
		AdditionalData: nil,
		PayloadData:    nil,
		Timestamp:      0,
	}
	assert.Equal(t, `{"id":"1","event":"TEST","source":"","service_code":"","payload":null,"timestamp":0}`, e2.String())
}

func TestApplicationEvent_WhenAddAdditionalData_ShouldAddCorrectly(t *testing.T) {
	eventName := "TestEvent"
	e := NewApplicationEvent(eventName,
		WithAdditionalData(map[string]interface{}{
			"key1": "val1",
		}),
	)
	assert.Equal(t, map[string]interface{}{
		"key1": "val1",
	}, e.AdditionalData)
	e.AddAdditionData("key2", "val2")
	assert.Equal(t, map[string]interface{}{
		"key1": "val1",
		"key2": "val2",
	}, e.AdditionalData)
}

func TestApplicationEvent_WhenDeleteAdditionalData_ShouldDeleteCorrectly(t *testing.T) {
	eventName := "TestEvent"
	e := NewApplicationEvent(eventName,
		WithAdditionalData(map[string]interface{}{
			"key1": "val1",
			"key2": "val2",
		}),
	)
	assert.Equal(t, map[string]interface{}{
		"key1": "val1",
		"key2": "val2",
	}, e.AdditionalData)
	e.DeleteAdditionData("key1")
	assert.Equal(t, map[string]interface{}{
		"key2": "val2",
	}, e.AdditionalData)
}
