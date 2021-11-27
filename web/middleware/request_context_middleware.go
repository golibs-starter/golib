package middleware

import (
	mainContext "context"
	"errors"
	"gitlab.com/golibs-starter/golib/pubsub"
	"gitlab.com/golibs-starter/golib/web/context"
	"gitlab.com/golibs-starter/golib/web/event"
	"gitlab.com/golibs-starter/golib/web/log"
	"net/http"
	"time"
)

// RequestContext middleware responsible to inject attributes to the request's context.
// This middleware should be run as soon as possible to
// create a uniform context for the request.
func RequestContext() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestAttributes := context.GetOrCreateRequestAttributes(r)
			next.ServeHTTP(w, r)
			if advancedResponseWriter, err := getAdvancedResponseWriter(w); err != nil {
				log.Warn(r.Context(), "Cannot detect AdvancedResponseWriter with error [%s]", err.Error())
			} else {
				requestAttributes.StatusCode = advancedResponseWriter.Status()
			}
			requestAttributes.ExecutionTime = time.Now().Sub(start)
			publishRequestCompletedEvent(r.Context(), requestAttributes)
		})
	}
}

func getAdvancedResponseWriter(w http.ResponseWriter) (*context.AdvancedResponseWriter, error) {
	if advancedResponseWriter, ok := w.(*context.AdvancedResponseWriter); ok {
		return advancedResponseWriter, nil
	}
	if wrappingWriter, ok := w.(context.WrappingResponseWriter); ok {
		if advancedResponseWriter, ok := wrappingWriter.Writer().(*context.AdvancedResponseWriter); ok {
			return advancedResponseWriter, nil
		}
		return nil, errors.New("ResponseWriter is wrapped by more than two level")
	}
	return nil, errors.New("your ResponseWriter is not implement context.WrappingResponseWriter")
}

func publishRequestCompletedEvent(ctx mainContext.Context, requestAttributes *context.RequestAttributes) {
	pubsub.Publish(event.NewRequestCompletedEvent(ctx, &event.RequestCompletedMessage{
		Status:            requestAttributes.StatusCode,
		ExecutionTime:     requestAttributes.ExecutionTime,
		Uri:               requestAttributes.Uri,
		Query:             requestAttributes.Query,
		Mapping:           requestAttributes.Mapping,
		Url:               requestAttributes.Url,
		Method:            requestAttributes.Method,
		CallerId:          requestAttributes.CallerId,
		DeviceId:          requestAttributes.DeviceId,
		DeviceSessionId:   requestAttributes.DeviceSessionId,
		CorrelationId:     requestAttributes.CorrelationId,
		ClientIpAddress:   requestAttributes.ClientIpAddress,
		UserAgent:         requestAttributes.UserAgent,
		UserId:            requestAttributes.SecurityAttributes.UserId,
		TechnicalUsername: requestAttributes.SecurityAttributes.TechnicalUsername,
	}))
}
