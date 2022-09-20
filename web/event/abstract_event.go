package event

import (
	"context"
	"gitlab.com/golibs-starter/golib/event"
	"gitlab.com/golibs-starter/golib/web/constant"
	webContext "gitlab.com/golibs-starter/golib/web/context"
)

type AbstractEventWrapper interface {
	GetAbstractEvent() *AbstractEvent
}

type AbstractEvent struct {
	*event.ApplicationEvent
	ctx               context.Context
	RequestId         string `json:"request_id,omitempty"`
	UserId            string `json:"user_id,omitempty"`
	TechnicalUsername string `json:"technical_username,omitempty"`
}

func NewAbstractEvent(ctx context.Context, name string, options ...event.AppEventOpt) *AbstractEvent {
	evt := AbstractEvent{}
	evt.ApplicationEvent = event.NewApplicationEvent(name, options...)
	requestAttributes := webContext.GetRequestAttributes(ctx)
	if requestAttributes != nil {
		evt.ServiceCode = requestAttributes.ServiceCode
		evt.RequestId = requestAttributes.CorrelationId
		evt.UserId = requestAttributes.SecurityAttributes.UserId
		evt.TechnicalUsername = requestAttributes.SecurityAttributes.TechnicalUsername
		if len(requestAttributes.ClientIpAddress) > 0 {
			evt.AddAdditionData(constant.HeaderClientIpAddress, requestAttributes.ClientIpAddress)
		}
		if len(requestAttributes.DeviceId) > 0 {
			evt.AddAdditionData(constant.HeaderDeviceId, requestAttributes.DeviceId)
		}
		if len(requestAttributes.DeviceSessionId) > 0 {
			evt.AddAdditionData(constant.HeaderDeviceSessionId, requestAttributes.DeviceSessionId)
		}
	}
	evt.generateContext(ctx)
	return &evt
}

func (a *AbstractEvent) generateContext(parent context.Context) {
	if parent == nil {
		parent = context.Background()
	}
	a.ctx = context.WithValue(parent, constant.ContextEventAttributes, MakeAttributes(a))
}

func (a *AbstractEvent) Context() context.Context {
	if a.ctx == nil {
		a.generateContext(nil)
	}
	return a.ctx
}

func (a AbstractEvent) String() string {
	return a.ToString(a)
}

func (a *AbstractEvent) GetAbstractEvent() *AbstractEvent {
	return a
}
