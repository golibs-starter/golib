package log

import (
	"github.com/golibs-starter/golib/log/field"
	webContext "github.com/golibs-starter/golib/web/context"
	"github.com/golibs-starter/golib/web/event"
)

type ContextAttributes struct {
	CorrelationId     string `json:"request_id,omitempty"`
	UserId            string `json:"jwt_subject,omitempty"`
	DeviceId          string `json:"device_id,omitempty"`
	DeviceSessionId   string `json:"device_session_id,omitempty"`
	TechnicalUsername string `json:"technical_username,omitempty"`
}

func (c ContextAttributes) MarshalLogObject(encoder field.ObjectEncoder) error {
	if c.CorrelationId != "" {
		encoder.AddString("request_id", c.CorrelationId)
	}
	if c.UserId != "" {
		encoder.AddString("user_id", c.UserId)
	}
	if c.DeviceId != "" {
		encoder.AddString("device_id", c.DeviceId)
	}
	if c.DeviceSessionId != "" {
		encoder.AddString("device_session_id", c.DeviceSessionId)
	}
	if c.TechnicalUsername != "" {
		encoder.AddString("technical_username", c.TechnicalUsername)
	}
	return nil
}

func NewContextAttributesFromReqAttr(requestAttributes *webContext.RequestAttributes) *ContextAttributes {
	return &ContextAttributes{
		DeviceId:          requestAttributes.DeviceId,
		DeviceSessionId:   requestAttributes.DeviceSessionId,
		CorrelationId:     requestAttributes.CorrelationId,
		UserId:            requestAttributes.SecurityAttributes.UserId,
		TechnicalUsername: requestAttributes.SecurityAttributes.TechnicalUsername,
	}
}

func NewContextAttributesFromEventAttrs(attributes *event.Attributes) *ContextAttributes {
	return &ContextAttributes{
		DeviceId:          attributes.DeviceId,
		DeviceSessionId:   attributes.DeviceSessionId,
		CorrelationId:     attributes.CorrelationId,
		UserId:            attributes.UserId,
		TechnicalUsername: attributes.TechnicalUsername,
	}
}
