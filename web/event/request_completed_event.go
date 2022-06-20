package event

import (
	"context"
	"time"
)

func NewRequestCompletedEvent(ctx context.Context, payload *RequestCompletedMessage) *RequestCompletedEvent {
	return &RequestCompletedEvent{
		AbstractEvent: NewAbstractEvent(ctx, "RequestCompletedEvent"),
		PayloadData:   payload,
	}
}

type RequestCompletedEvent struct {
	*AbstractEvent
	PayloadData *RequestCompletedMessage `json:"payload"`
}

func (a RequestCompletedEvent) Payload() interface{} {
	return a.PayloadData
}

func (a RequestCompletedEvent) String() string {
	return a.ToString(a)
}

type RequestCompletedMessage struct {
	Status            int           `json:"status"`
	ExecutionTime     time.Duration `json:"duration_ms"`
	Uri               string        `json:"uri"`
	Query             string        `json:"query,omitempty"`
	Mapping           string        `json:"mapping"`
	Url               string        `json:"url"`
	Method            string        `json:"method"`
	CorrelationId     string        `json:"correlation_id,omitempty"`
	CallerId          string        `json:"caller_id,omitempty"`
	ClientIpAddress   string        `json:"client_ip_address,omitempty"`
	Locale            string        `json:"locale,omitempty"`
	UserAgent         string        `json:"user_agent,omitempty"`
	UserId            string        `json:"user_id,omitempty"`
	DeviceId          string        `json:"device_id,omitempty"`
	DeviceSessionId   string        `json:"device_session_id,omitempty"`
	TechnicalUsername string        `json:"technical_username,omitempty"`
}
