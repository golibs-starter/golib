package listener

import (
	"gitlab.com/golibs-starter/golib/pubsub"
	"gitlab.com/golibs-starter/golib/web/constant"
	"gitlab.com/golibs-starter/golib/web/event"
	"gitlab.com/golibs-starter/golib/web/log"
	"gitlab.com/golibs-starter/golib/web/properties"
)

type RequestCompletedLogListener struct {
	httpRequestProps *properties.HttpRequestLogProperties
}

func NewRequestCompletedLogListener(httpRequestProps *properties.HttpRequestLogProperties) pubsub.Subscriber {
	return &RequestCompletedLogListener{httpRequestProps: httpRequestProps}
}

func (r RequestCompletedLogListener) Supports(e pubsub.Event) bool {
	_, ok := e.(*event.RequestCompletedEvent)
	return ok
}

func (r RequestCompletedLogListener) Handle(e pubsub.Event) {
	if r.httpRequestProps.Disabled {
		return
	}
	ev := e.(*event.RequestCompletedEvent)
	if payload, ok := ev.Payload().(*event.RequestCompletedMessage); ok {
		if r.isDisabled(payload.Method, payload.Uri) {
			return
		}
		log.Infow([]interface{}{constant.ContextReqMeta, r.makeHttpRequestLog(payload)}, "finish router")
	}
}

func (r RequestCompletedLogListener) isDisabled(method string, uri string) bool {
	for _, urlMatching := range r.httpRequestProps.AllDisabledUrls() {
		if urlMatching.Method != "" && urlMatching.Method != method {
			continue
		}
		if urlMatching.UrlRegexp() != nil && urlMatching.UrlRegexp().MatchString(uri) {
			return true
		}
	}
	return false
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
		ExecutionTime:  message.ExecutionTime.Milliseconds(),
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
