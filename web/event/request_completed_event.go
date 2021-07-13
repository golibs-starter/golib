package event

import (
	"context"
	"time"
)

type RequestCompletedEvent struct {
	AbstractEvent
}

type RequestCompletedPayload struct {
	Status            int           `json:"status"`
	ExecutionTime     time.Duration `json:"duration_ms"`
	Uri               string        `json:"uri"`
	Query             string        `json:"query"`
	Mapping           string        `json:"mapping"`
	Url               string        `json:"url"`
	Method            string        `json:"method"`
	CorrelationId     string        `json:"correlation_id"`
	CallerId          string        `json:"caller_id"`
	ClientIpAddress   string        `json:"client_ip_address"`
	Locale            string        `json:"locale"`
	UserAgent         string        `json:"user_agent"`
	UserId            string        `json:"user_id"`
	DeviceId          string        `json:"device_id"`
	DeviceSessionId   string        `json:"device_session_id"`
	TechnicalUsername string        `json:"technical_username"`
}

func NewRequestCompletedEvent(ctx context.Context, payload RequestCompletedPayload) *RequestCompletedEvent {
	event := RequestCompletedEvent{}
	event.AbstractEvent = NewAbstractEvent(ctx, event.GetName(), payload)
	return &event
}

func (r RequestCompletedEvent) GetName() string {
	return "RequestCompletedEvent"
}

func (r RequestCompletedEvent) GetPayload() interface{} {
	return r.Payload
}
