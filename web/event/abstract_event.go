package event

import (
	"context"
	"encoding/json"
	"gitlab.id.vin/vincart/golib/event"
	"gitlab.id.vin/vincart/golib/web/constant"
	webContext "gitlab.id.vin/vincart/golib/web/context"
)

type AbstractEvent struct {
	event.ApplicationEvent
	RequestId string `json:"request_id"`
	UserId    string `json:"user_id"`
}

func NewAbstractEvent(ctx context.Context, eventName string, payload interface{}) AbstractEvent {
	absEvent := AbstractEvent{
		// TODO add service code here
		ApplicationEvent: event.NewApplicationEvent("", eventName, payload),
	}
	requestAttributes := webContext.GetRequestAttributes(ctx)
	if requestAttributes == nil {
		return absEvent
	}
	absEvent.RequestId = requestAttributes.CorrelationId
	absEvent.UserId = requestAttributes.SecurityAttributes.UserId
	absEvent.AdditionalData = map[string]interface{}{
		constant.HeaderDeviceId:           requestAttributes.DeviceId,
		constant.HeaderOldDeviceId:        requestAttributes.DeviceId,
		constant.HeaderDeviceSessionId:    requestAttributes.DeviceSessionId,
		constant.HeaderOldDeviceSessionId: requestAttributes.DeviceSessionId,
		constant.HeaderClientIpAddress:    requestAttributes.ClientIpAddress,
	}
	return absEvent
}

func (a AbstractEvent) String() string {
	data, _ := json.Marshal(a)
	return string(data)
}
