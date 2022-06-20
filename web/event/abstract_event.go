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
	RequestId         string `json:"request_id,omitempty"`
	UserId            string `json:"user_id,omitempty"`
	TechnicalUsername string `json:"technical_username,omitempty"`
	ctx               context.Context
}

func NewAbstractEvent(ctx context.Context, name string, options ...event.AppEventOpt) *AbstractEvent {
	e := AbstractEvent{
		ApplicationEvent: event.NewApplicationEvent(name, options...),
	}
	requestAttributes := webContext.GetRequestAttributes(ctx)
	if requestAttributes == nil {
		return &e
	}
	e.RequestId = requestAttributes.CorrelationId
	e.UserId = requestAttributes.SecurityAttributes.UserId
	e.TechnicalUsername = requestAttributes.SecurityAttributes.TechnicalUsername
	e.AddAdditionData(constant.HeaderClientIpAddress, requestAttributes.ClientIpAddress)
	e.AddAdditionData(constant.HeaderDeviceId, requestAttributes.DeviceId)
	e.AddAdditionData(constant.HeaderDeviceSessionId, requestAttributes.DeviceSessionId)
	e.AddAdditionData(constant.HeaderOldDeviceId, requestAttributes.DeviceId)
	e.AddAdditionData(constant.HeaderOldDeviceSessionId, requestAttributes.DeviceSessionId)
	e.generateContext(ctx)
	return &e
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
