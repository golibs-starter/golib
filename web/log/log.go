package log

import (
	"context"
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

func BuildLoggingContextFromEventAttr(attributes *event.Attributes) *LoggingContext {
	return &LoggingContext{
		DeviceId:          attributes.DeviceId,
		DeviceSessionId:   attributes.DeviceSessionId,
		CorrelationId:     attributes.CorrelationId,
		UserId:            attributes.UserId,
		TechnicalUsername: attributes.TechnicalUsername,
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
	if eventAttributes := event.GetAttributes(ctx); eventAttributes != nil {
		return []interface{}{constant.ContextReqMeta, BuildLoggingContextFromEventAttr(eventAttributes)}
	}
	return nil
}

func keysAndValuesFromEvent(e pubsub.Event) []interface{} {
	var logContext = make([]interface{}, 0)
	if e == nil {
		return logContext
	}
	if we, ok := e.(event.AbstractEventWrapper); ok {
		logContext = []interface{}{constant.ContextReqMeta, BuildLoggingContextFromEvent(we.GetAbstractEvent())}
	}
	return logContext
}

func Debug(ctx context.Context, msgFormat string, args ...interface{}) {
	global.Debugw(keysAndValuesFromContext(ctx), msgFormat, args...)
}

func Info(ctx context.Context, msgFormat string, args ...interface{}) {
	global.Infow(keysAndValuesFromContext(ctx), msgFormat, args...)
}

func Warn(ctx context.Context, msgFormat string, args ...interface{}) {
	global.Warnw(keysAndValuesFromContext(ctx), msgFormat, args...)
}

func Error(ctx context.Context, msgFormat string, args ...interface{}) {
	global.Errorw(keysAndValuesFromContext(ctx), msgFormat, args...)
}

func Fatal(ctx context.Context, msgFormat string, args ...interface{}) {
	global.Fatalw(keysAndValuesFromContext(ctx), msgFormat, args...)
}

func Debuge(e pubsub.Event, msgFormat string, args ...interface{}) {
	global.Debugw(keysAndValuesFromEvent(e), msgFormat, args...)
}

func Infoe(e pubsub.Event, msgFormat string, args ...interface{}) {
	global.Infow(keysAndValuesFromEvent(e), msgFormat, args...)
}

func Warne(e pubsub.Event, msgFormat string, args ...interface{}) {
	global.Warnw(keysAndValuesFromEvent(e), msgFormat, args...)
}

func Errore(e pubsub.Event, msgFormat string, args ...interface{}) {
	global.Errorw(keysAndValuesFromEvent(e), msgFormat, args...)
}

func Fatale(e pubsub.Event, msgFormat string, args ...interface{}) {
	global.Fatalw(keysAndValuesFromEvent(e), msgFormat, args...)
}

func Debugf(msgFormat string, args ...interface{}) {
	global.Debugf(msgFormat, args...)
}

func Infof(msgFormat string, args ...interface{}) {
	global.Infof(msgFormat, args...)
}

func Warnf(msgFormat string, args ...interface{}) {
	global.Warnf(msgFormat, args...)
}

func Errorf(msgFormat string, args ...interface{}) {
	global.Errorf(msgFormat, args...)
}

func Fatalf(msgFormat string, args ...interface{}) {
	global.Fatalf(msgFormat, args...)
}

func Debugw(keysAndValues []interface{}, msgFormat string, args ...interface{}) {
	global.Debugw(keysAndValues, msgFormat, args...)
}

func Infow(keysAndValues []interface{}, msgFormat string, args ...interface{}) {
	global.Infow(keysAndValues, msgFormat, args...)
}

func Warnw(keysAndValues []interface{}, msgFormat string, args ...interface{}) {
	global.Warnw(keysAndValues, msgFormat, args...)
}

func Errorw(keysAndValues []interface{}, msgFormat string, args ...interface{}) {
	global.Errorw(keysAndValues, msgFormat, args...)
}

func Fatalw(keysAndValues []interface{}, msgFormat string, args ...interface{}) {
	global.Fatalw(keysAndValues, msgFormat, args...)
}
