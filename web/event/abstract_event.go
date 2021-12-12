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
	RequestId         string `json:"request_id"`
	UserId            string `json:"user_id"`
	TechnicalUsername string `json:"technical_username"`
	ctx               context.Context
}

func NewAbstractEvent(ctx context.Context, eventName string) *AbstractEvent {
	absEvent := AbstractEvent{
		ApplicationEvent: event.NewApplicationEvent(eventName),
	}
	requestAttributes := webContext.GetRequestAttributes(ctx)
	if requestAttributes == nil {
		return &absEvent
	}
	absEvent.RequestId = requestAttributes.CorrelationId
	absEvent.UserId = requestAttributes.SecurityAttributes.UserId
	absEvent.TechnicalUsername = requestAttributes.SecurityAttributes.TechnicalUsername
	absEvent.AdditionalData = map[string]interface{}{
		constant.HeaderClientIpAddress:    requestAttributes.ClientIpAddress,
		constant.HeaderDeviceId:           requestAttributes.DeviceId,
		constant.HeaderDeviceSessionId:    requestAttributes.DeviceSessionId,
		constant.HeaderOldDeviceId:        requestAttributes.DeviceId,
		constant.HeaderOldDeviceSessionId: requestAttributes.DeviceSessionId,
	}
	absEvent.generateContext(ctx)
	return &absEvent
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
