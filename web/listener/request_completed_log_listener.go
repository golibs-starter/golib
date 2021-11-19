package listener

import (
	"gitlab.id.vin/vincart/common/golib/pubsub"
	"gitlab.id.vin/vincart/common/golib/web/constant"
	"gitlab.id.vin/vincart/common/golib/web/event"
	"gitlab.id.vin/vincart/common/golib/web/log"
)

type RequestCompletedLogListener struct {
}

func NewRequestCompletedLogListener() pubsub.Subscriber {
	return &RequestCompletedLogListener{}
}

func (r RequestCompletedLogListener) Supports(e pubsub.Event) bool {
	_, ok := e.(*event.RequestCompletedEvent)
	return ok
}

func (r RequestCompletedLogListener) Handle(e pubsub.Event) {
	ev := e.(*event.RequestCompletedEvent)
	if payload, ok := ev.Payload().(event.RequestCompletedPayload); ok {
		log.Infow([]interface{}{constant.ContextReqMeta, r.makeHttpRequestLog(&payload)}, "finish router")
	}
}

func (r RequestCompletedLogListener) makeHttpRequestLog(message *event.RequestCompletedPayload) *log.HttpRequestLog {
	return &log.HttpRequestLog{
		LoggingContext: log.LoggingContext{
			UserId:            message.UserId,
			DeviceId:          message.DeviceId,
			DeviceSessionId:   message.DeviceSessionId,
			CorrelationId:     message.CorrelationId,
			TechnicalUsername: message.TechnicalUsername,
		},
		Status:         message.Status,
		ExecutionTime:  message.ExecutionTime,
		RequestPattern: message.Mapping,
		RequestPath:    message.Uri,
		Method:         message.Method,
		Query:          message.Query,
		Url:            message.Url,
		RequestId:      message.CorrelationId,
		CallerId:       message.CallerId,
		ClientIp:       message.ClientIpAddress,
		UserAgent:      message.UserAgent,
	}
}
