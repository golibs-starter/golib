package log

import (
	"context"
	mainLog "gitlab.id.vin/vincart/golib/log"
	"gitlab.id.vin/vincart/golib/web/constant"
	webContext "gitlab.id.vin/vincart/golib/web/context"
)

type LoggingContext struct {
	CorrelationId     string `json:"request_id"`
	UserId            string `json:"jwt_subject,omitempty"`
	DeviceId          string `json:"device_id,omitempty"`
	DeviceSessionId   string `json:"device_session_id,omitempty"`
	TechnicalUsername string `json:"technical_username,omitempty"`
}

func buildLoggingContext(requestAttributes *webContext.RequestAttributes) *LoggingContext {
	return &LoggingContext{
		DeviceId:          requestAttributes.DeviceId,
		DeviceSessionId:   requestAttributes.DeviceSessionId,
		CorrelationId:     requestAttributes.CorrelationId,
		UserId:            requestAttributes.SecurityAttributes.UserId,
		TechnicalUsername: requestAttributes.SecurityAttributes.TechnicalUsername,
	}
}

func keysAndValuesFromContext(ctx context.Context) []interface{} {
	if requestAttributes := webContext.GetRequestAttributes(ctx); requestAttributes != nil {
		return []interface{}{constant.ContextReqMeta, buildLoggingContext(requestAttributes)}
	}
	return nil
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
