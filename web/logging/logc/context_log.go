package logc

import (
	"context"
	mainLog "gitlab.id.vin/vincart/golib/log"
	"gitlab.id.vin/vincart/golib/web/constants"
	webContext "gitlab.id.vin/vincart/golib/web/context"
	"gitlab.id.vin/vincart/golib/web/logging"
)

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

func keysAndValuesFromContext(ctx context.Context) []interface{} {
	reqAttrCtxValue := ctx.Value(constants.ContextReqAttribute)
	if reqAttrCtxValue == nil {
		return []interface{}{}
	}
	requestAttributes, ok := reqAttrCtxValue.(*webContext.RequestAttributes)
	if !ok {
		return []interface{}{}
	}
	return []interface{}{constants.ContextReqMeta, buildLogContext(requestAttributes)}
}

func buildLogContext(requestAttributes *webContext.RequestAttributes) *logging.LogContext {
	return &logging.LogContext{
		DeviceId:          requestAttributes.DeviceId,
		DeviceSessionId:   requestAttributes.DeviceSessionId,
		CorrelationId:     requestAttributes.CorrelationId,
		UserId:            requestAttributes.SecurityAttributes.UserId,
		TechnicalUsername: requestAttributes.SecurityAttributes.TechnicalUsername,
	}
}
