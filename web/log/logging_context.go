package log

import (
	"context"
	mainLog "gitlab.id.vin/vincart/golib/log"
	"gitlab.id.vin/vincart/golib/pubsub"
	"gitlab.id.vin/vincart/golib/web/constant"
	webContext "gitlab.id.vin/vincart/golib/web/context"
	"gitlab.id.vin/vincart/golib/web/event"
)

type LoggingContext struct {
	CorrelationId     string `json:"request_id"`
	UserId            string `json:"jwt_subject,omitempty"`
	DeviceId          string `json:"device_id,omitempty"`
	DeviceSessionId   string `json:"device_session_id,omitempty"`
	TechnicalUsername string `json:"technical_username,omitempty"`
}

func BuildLoggingContextFromReqAttr(requestAttributes *webContext.RequestAttributes) *LoggingContext {
	return &LoggingContext{
		DeviceId:          requestAttributes.DeviceId,
		DeviceSessionId:   requestAttributes.DeviceSessionId,
		CorrelationId:     requestAttributes.CorrelationId,
		UserId:            requestAttributes.SecurityAttributes.UserId,
		TechnicalUsername: requestAttributes.SecurityAttributes.TechnicalUsername,
	}
}

func BuildLoggingContextFromEvent(e *event.AbstractEvent) *LoggingContext {
	deviceId, _ := e.AdditionalData[constant.HeaderDeviceId].(string)
	deviceSessionId, _ := e.AdditionalData[constant.HeaderDeviceSessionId].(string)
	return &LoggingContext{
		UserId:            e.UserId,
		DeviceId:          deviceId,
		DeviceSessionId:   deviceSessionId,
		CorrelationId:     e.RequestId,
		TechnicalUsername: e.TechnicalUsername,
	}
}

func keysAndValuesFromContext(ctx context.Context) []interface{} {
	if requestAttributes := webContext.GetRequestAttributes(ctx); requestAttributes != nil {
		return []interface{}{constant.ContextReqMeta, BuildLoggingContextFromReqAttr(requestAttributes)}
	}
	return nil
}

func keysAndValuesFromEvent(e pubsub.Event) []interface{} {
	var logContext = make([]interface{}, 0)
	if we, ok := e.(event.AbstractEventWrapper); ok {
		logContext = []interface{}{constant.ContextReqMeta, BuildLoggingContextFromEvent(we.GetAbstractEvent())}
	}
	return logContext
}

func Debug(ctx context.Context, msgFormat string, args ...interface{}) {
	mainLog.Debugw(keysAndValuesFromContext(ctx), msgFormat, args...)
}

func Info(ctx context.Context, msgFormat string, args ...interface{}) {
	mainLog.Infow(keysAndValuesFromContext(ctx), msgFormat, args...)
}

func Warn(ctx context.Context, msgFormat string, args ...interface{}) {
	mainLog.Warnw(keysAndValuesFromContext(ctx), msgFormat, args...)
}

func Error(ctx context.Context, msgFormat string, args ...interface{}) {
	mainLog.Errorw(keysAndValuesFromContext(ctx), msgFormat, args...)
}

func Fatal(ctx context.Context, msgFormat string, args ...interface{}) {
	mainLog.Fatalw(keysAndValuesFromContext(ctx), msgFormat, args...)
}

func Debuge(e pubsub.Event, msgFormat string, args ...interface{}) {
	mainLog.Debugw(keysAndValuesFromEvent(e), msgFormat, args...)
}

func Infoe(e pubsub.Event, msgFormat string, args ...interface{}) {
	mainLog.Infow(keysAndValuesFromEvent(e), msgFormat, args...)
}

func Warne(e pubsub.Event, msgFormat string, args ...interface{}) {
	mainLog.Warnw(keysAndValuesFromEvent(e), msgFormat, args...)
}

func Errore(e pubsub.Event, msgFormat string, args ...interface{}) {
	mainLog.Errorw(keysAndValuesFromEvent(e), msgFormat, args...)
}

func Fatale(e pubsub.Event, msgFormat string, args ...interface{}) {
	mainLog.Fatalw(keysAndValuesFromEvent(e), msgFormat, args...)
}
