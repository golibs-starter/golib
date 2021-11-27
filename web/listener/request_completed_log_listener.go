package listener

import (
	"gitlab.com/golibs-starter/golib/pubsub"
	"gitlab.com/golibs-starter/golib/web/constant"
	"gitlab.com/golibs-starter/golib/web/event"
	"gitlab.com/golibs-starter/golib/web/log"
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
	if payload, ok := ev.Payload().(*event.RequestCompletedMessage); ok {
		log.Infow([]interface{}{constant.ContextReqMeta, r.makeHttpRequestLog(payload)}, "finish router")
	}
}

func (r RequestCompletedLogListener) makeHttpRequestLog(message *event.RequestCompletedMessage) *log.HttpRequestLog {
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
