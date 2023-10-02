package event

import (
	"context"
	"github.com/golibs-starter/golib/event"
	"github.com/golibs-starter/golib/web/constant"
	webContext "github.com/golibs-starter/golib/web/context"
)

type AbstractEventWrapper interface {
	GetAbstractEvent() *AbstractEvent
}

type AbstractEvent struct {
	*event.ApplicationEvent
	RequestId         string `json:"request_id,omitempty"`
	UserId            string `json:"user_id,omitempty"`
	TechnicalUsername string `json:"technical_username,omitempty"`
}

func NewAbstractEvent(ctx context.Context, name string, options ...event.AppEventOpt) *AbstractEvent {
	evt := AbstractEvent{}
	evt.ApplicationEvent = event.NewApplicationEvent(ctx, name, options...)
	reqAttrs := webContext.GetRequestAttributes(ctx)
	if reqAttrs != nil {
		evt.ServiceCode = reqAttrs.ServiceCode
		evt.RequestId = reqAttrs.CorrelationId
		evt.UserId = reqAttrs.SecurityAttributes.UserId
		evt.TechnicalUsername = reqAttrs.SecurityAttributes.TechnicalUsername
		if len(reqAttrs.ClientIpAddress) > 0 {
			evt.AddAdditionData(constant.HeaderClientIpAddress, reqAttrs.ClientIpAddress)
		}
		if len(reqAttrs.DeviceId) > 0 {
			evt.AddAdditionData(constant.HeaderDeviceId, reqAttrs.DeviceId)
		}
		if len(reqAttrs.DeviceSessionId) > 0 {
			evt.AddAdditionData(constant.HeaderDeviceSessionId, reqAttrs.DeviceSessionId)
		}
	}
	if attrs := MakeAttributes(&evt); attrs != nil {
		evt.Ctx = context.WithValue(ctx, constant.ContextEventAttributes, MakeAttributes(&evt))
	}
	return &evt
}

func (a *AbstractEvent) String() string {
	return a.ToString(a)
}

func (a *AbstractEvent) GetAbstractEvent() *AbstractEvent {
	return a
}
