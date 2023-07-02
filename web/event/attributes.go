package event

import (
	"context"
	"gitlab.com/golibs-starter/golib/web/constant"
)

type Attributes struct {
	CorrelationId     string `json:"request_id,omitempty"`
	UserId            string `json:"jwt_subject,omitempty"`
	DeviceId          string `json:"device_id,omitempty"`
	DeviceSessionId   string `json:"device_session_id,omitempty"`
	TechnicalUsername string `json:"technical_username,omitempty"`
}

func MakeAttributes(e *AbstractEvent) *Attributes {
	deviceId, _ := e.AdditionalData[constant.HeaderDeviceId].(string)
	deviceSessionId, _ := e.AdditionalData[constant.HeaderDeviceSessionId].(string)
	return &Attributes{
		CorrelationId:     e.RequestId,
		UserId:            e.UserId,
		DeviceId:          deviceId,
		DeviceSessionId:   deviceSessionId,
		TechnicalUsername: e.TechnicalUsername,
	}
}

func GetAttributes(ctx context.Context) *Attributes {
	eventAttrCtxValue := ctx.Value(constant.ContextEventAttributes)
	if eventAttrCtxValue == nil {
		return nil
	}
	attrs, ok := eventAttrCtxValue.(*Attributes)
	if !ok {
		return nil
	}
	return attrs
}
