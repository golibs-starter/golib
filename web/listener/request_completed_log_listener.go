package listener

import (
	"gitlab.com/golibs-starter/golib/config"
	"gitlab.com/golibs-starter/golib/log"
	"gitlab.com/golibs-starter/golib/log/field"
	"gitlab.com/golibs-starter/golib/pubsub"
	"gitlab.com/golibs-starter/golib/web/constant"
	"gitlab.com/golibs-starter/golib/web/event"
	webLog "gitlab.com/golibs-starter/golib/web/log"
	"gitlab.com/golibs-starter/golib/web/properties"
	"strings"
)

type RequestCompletedLogListener struct {
	appProps         *config.AppProperties
	httpRequestProps *properties.HttpRequestLogProperties
}

func NewRequestCompletedLogListener(
	appProps *config.AppProperties,
	httpRequestProps *properties.HttpRequestLogProperties,
) pubsub.Subscriber {
	return &RequestCompletedLogListener{
		appProps:         appProps,
		httpRequestProps: httpRequestProps,
	}
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
		// TODO Should remove context path in the highest filter
		if r.isDisabled(payload.Method, r.removeContextPath(payload.Uri, r.appProps.Path)) {
			return
		}
		log.WithField(field.Object(constant.ContextReqMeta, r.makeHttpRequestLog(payload))).Infof("finish router")
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

func (r RequestCompletedLogListener) removeContextPath(uri string, contextPath string) string {
	uri = strings.TrimPrefix(uri, contextPath)
	return "/" + strings.TrimLeft(uri, "/")
}

func (r RequestCompletedLogListener) makeHttpRequestLog(message *event.RequestCompletedMessage) *HttpRequestLog {
	return &HttpRequestLog{
		ContextAttributes: webLog.ContextAttributes{
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
