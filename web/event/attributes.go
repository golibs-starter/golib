package event

import (
    "context"
    "gitlab.com/golibs-starter/golib/pubsub"
    "gitlab.com/golibs-starter/golib/web/constant"
)

type Attributes struct {
    CorrelationId     string `json:"request_id"`
    UserId            string `json:"jwt_subject,omitempty"`
    DeviceId          string `json:"device_id,omitempty"`
    DeviceSessionId   string `json:"device_session_id,omitempty"`
    TechnicalUsername string `json:"technical_username,omitempty"`
}

func MakeAttributes(e pubsub.Event) *Attributes {
    if we, ok := e.(AbstractEventWrapper); ok {
        absEvent := we.GetAbstractEvent()
        deviceId, _ := absEvent.AdditionalData[constant.HeaderDeviceId].(string)
        deviceSessionId, _ := absEvent.AdditionalData[constant.HeaderDeviceSessionId].(string)
        return &Attributes{
            CorrelationId:     absEvent.RequestId,
            UserId:            absEvent.UserId,
            DeviceId:          deviceId,
            DeviceSessionId:   deviceSessionId,
            TechnicalUsername: absEvent.TechnicalUsername,
        }
    }
    return nil
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
