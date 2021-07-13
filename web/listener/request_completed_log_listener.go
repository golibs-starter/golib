package listener

import (
	"gitlab.id.vin/vincart/golib/log"
	"gitlab.id.vin/vincart/golib/pubsub"
	"gitlab.id.vin/vincart/golib/web/constants"
	"gitlab.id.vin/vincart/golib/web/event"
	"gitlab.id.vin/vincart/golib/web/logging"
)

type RequestCompletedLogListener struct {
}

func (r RequestCompletedLogListener) Handler(e pubsub.Event) {
	if e.GetName() != (event.RequestCompletedEvent{}).GetName() {
		return
	}
	payload, ok := e.GetPayload().(event.RequestCompletedPayload)
	if !ok {
		return
	}
	log.Infow([]interface{}{constants.ContextReqMeta, r.makeHttpRequestLog(&payload)}, "finish router")
}

func (r RequestCompletedLogListener) makeHttpRequestLog(message *event.RequestCompletedPayload) *logging.HttpRequestLog {
	return &logging.HttpRequestLog{
		LogContext: logging.LogContext{
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
