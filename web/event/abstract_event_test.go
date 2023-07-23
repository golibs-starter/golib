package event

import (
	"context"
	assert "github.com/stretchr/testify/require"
	"gitlab.com/golibs-starter/golib/event"
	"gitlab.com/golibs-starter/golib/web/constant"
	context2 "gitlab.com/golibs-starter/golib/web/context"
	"testing"
)

func TestNewAbstractEvent_GivenAName_ShouldInitCorrectly(t *testing.T) {
	eventName := "TestEvent"
	e := NewAbstractEvent(context.Background(), eventName)
	assert.NotNil(t, e.GetAbstractEvent())
	assert.NotNil(t, e.ApplicationEvent)
	assert.Equal(t, eventName, e.Name())
	assert.Equal(t, eventName, e.Event)
	assert.NotEmpty(t, e.Id)
	assert.Equal(t, e.Id, e.Identifier())
	assert.Greater(t, e.Timestamp, int64(0))
	assert.Equal(t, event.DefaultEventSource, e.Source)
	assert.Empty(t, e.ServiceCode)
	assert.Empty(t, e.RequestId)
	assert.Empty(t, e.UserId)
	assert.Empty(t, e.TechnicalUsername)
	assert.Nil(t, e.AdditionalData)
}

func TestNewAbstractEvent_GivenANameAndRequestAttribute_ShouldInitCorrectly(t *testing.T) {
	eventName := "TestEvent"
	attr := context2.RequestAttributes{
		ServiceCode:     "test-service-code1",
		CorrelationId:   "test-id1",
		ClientIpAddress: "test-client-ip1",
		DeviceId:        "test-device-id1",
		DeviceSessionId: "test-device-session-id1",
		SecurityAttributes: context2.SecurityAttributes{
			UserId:            "test-uid1",
			TechnicalUsername: "test-username1",
		},
	}
	ctx := context.WithValue(context.Background(), constant.ContextReqAttribute, &attr)
	e := NewAbstractEvent(ctx, eventName)
	assert.NotNil(t, e.GetAbstractEvent())
	assert.NotNil(t, e.ApplicationEvent)
	assert.Equal(t, eventName, e.Name())
	assert.Equal(t, eventName, e.Event)
	assert.NotEmpty(t, e.Id)
	assert.Equal(t, e.Id, e.Identifier())
	assert.Greater(t, e.Timestamp, int64(0))
	assert.Equal(t, event.DefaultEventSource, e.Source)
	assert.Equal(t, attr.ServiceCode, e.ServiceCode)
	assert.Equal(t, attr.CorrelationId, e.RequestId)
	assert.Equal(t, attr.SecurityAttributes.UserId, e.UserId)
	assert.Equal(t, attr.SecurityAttributes.TechnicalUsername, e.TechnicalUsername)
	assert.Equal(t, map[string]interface{}{
		constant.HeaderClientIpAddress: attr.ClientIpAddress,
		constant.HeaderDeviceId:        attr.DeviceId,
		constant.HeaderDeviceSessionId: attr.DeviceSessionId,
	}, e.AdditionalData)
}

func TestNewAbstractEvent_GivenANameAndOptions_ShouldRunOptionsCorrectly(t *testing.T) {
	eventName := "TestEvent"
	payload := map[string]string{"a": "a"}
	e := NewAbstractEvent(context.Background(), eventName,
		event.WithId("test-id"),
		event.WithServiceCode("test-service-code"),
		event.WithSource("test-source"),
		event.WithAdditionalData(map[string]interface{}{
			"key1": "val1",
		}),
		event.WithPayload(payload),
	)
	assert.NotNil(t, e.GetAbstractEvent())
	assert.NotNil(t, e.ApplicationEvent)
	assert.Equal(t, eventName, e.Name())
	assert.Equal(t, eventName, e.Event)
	assert.Equal(t, "test-id", e.Id)
	assert.Equal(t, "test-id", e.Identifier())
	assert.Greater(t, e.Timestamp, int64(0))
	assert.Equal(t, "test-source", e.Source)
	assert.Equal(t, "test-service-code", e.ServiceCode)
	assert.Empty(t, e.RequestId)
	assert.Empty(t, e.UserId)
	assert.Empty(t, e.TechnicalUsername)
	assert.Equal(t, map[string]interface{}{
		"key1": "val1",
	}, e.AdditionalData)
	assert.Equal(t, payload, e.PayloadData)
	assert.Equal(t, payload, e.Payload())
}

func TestNewAbstractEvent_GivenANameAndCustomAdditionalData_ShouldMergeAdditionalDataCorrectly(t *testing.T) {
	eventName := "TestEvent"
	attr := context2.RequestAttributes{
		ClientIpAddress: "test-client-ip1",
	}
	ctx := context.WithValue(context.Background(), constant.ContextReqAttribute, &attr)
	e := NewAbstractEvent(ctx, eventName, event.WithAdditionalData(map[string]interface{}{
		"key1": "val1",
	}))
	assert.NotNil(t, e.GetAbstractEvent())
	assert.NotNil(t, e.ApplicationEvent)
	assert.Equal(t, eventName, e.Name())
	assert.Equal(t, eventName, e.Event)
	assert.NotEmpty(t, e.Id)
	assert.Equal(t, e.Id, e.Identifier())
	assert.Greater(t, e.Timestamp, int64(0))
	assert.Equal(t, event.DefaultEventSource, e.Source)
	assert.Empty(t, e.ServiceCode)
	assert.Empty(t, e.RequestId)
	assert.Empty(t, e.UserId)
	assert.Empty(t, e.TechnicalUsername)
	assert.Equal(t, map[string]interface{}{
		constant.HeaderClientIpAddress: "test-client-ip1",
		"key1":                         "val1",
	}, e.AdditionalData)
}

func TestAbstractEvent_ToString(t *testing.T) {
	e1 := AbstractEvent{
		ApplicationEvent: &event.ApplicationEvent{
			Id:             "1",
			Event:          "TEST",
			Source:         "NOT_USED",
			ServiceCode:    "service-test",
			AdditionalData: map[string]interface{}{"a": "b"},
			PayloadData:    map[string]string{"x": "y"},
			Timestamp:      10,
		},
		RequestId:         "test-request-id-1",
		UserId:            "test-user-id-1",
		TechnicalUsername: "test-technical-username-1",
	}
	assert.Equal(t, `{"id":"1","event":"TEST","source":"NOT_USED","service_code":"service-test","additional_data":{"a":"b"},"payload":{"x":"y"},"timestamp":10,"request_id":"test-request-id-1","user_id":"test-user-id-1","technical_username":"test-technical-username-1"}`, e1.String())
}
