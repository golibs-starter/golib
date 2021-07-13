package middleware

import (
	mainContext "context"
	"errors"
	"gitlab.id.vin/vincart/golib/pubsub"
	"gitlab.id.vin/vincart/golib/web/constants"
	"gitlab.id.vin/vincart/golib/web/context"
	"gitlab.id.vin/vincart/golib/web/event"
	"gitlab.id.vin/vincart/golib/web/logging/logc"
	"net/http"
	"time"
)

func RequestContext() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestAttributes := getOrCreateRequestAttributes(r)
			next.ServeHTTP(w, r)
			if advancedResponseWriter, err := getAdvancedResponseWriter(w); err != nil {
				logc.Warn(r.Context(), "Cannot detect AdvancedResponseWriter with error [%s]", err.Error())
			} else {
				requestAttributes.StatusCode = advancedResponseWriter.Status()
			}
			requestAttributes.ExecutionTime = time.Now().Sub(start)
			publishEvent(r.Context(), requestAttributes)
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

func getOrCreateRequestAttributes(r *http.Request) *context.RequestAttributes {
	reqAttrCtxValue := r.Context().Value(constants.ContextReqAttribute)
	if reqAttrCtxValue == nil {
		return createNewRequestAttributes(r)
	}
	requestAttributes, ok := reqAttrCtxValue.(*context.RequestAttributes)
	if !ok {
		logc.Error(r.Context(), "Request attributes is not *RequestAttributes type")
		return createNewRequestAttributes(r)
	}
	return requestAttributes
}

func createNewRequestAttributes(r *http.Request) *context.RequestAttributes {
	requestAttributes := makeRequestAttributes(r)
	*r = *r.WithContext(mainContext.WithValue(r.Context(), constants.ContextReqAttribute, requestAttributes))
	return requestAttributes
}

func makeRequestAttributes(r *http.Request) *context.RequestAttributes {
	return &context.RequestAttributes{
		Method:             r.Method,
		Uri:                r.RequestURI,
		Query:              r.URL.RawQuery,
		Url:                r.URL.String(),
		UserAgent:          r.Header.Get(constants.HeaderUserAgent),
		ClientIpAddress:    getClientIpAddress(r),
		DeviceId:           r.Header.Get(constants.HeaderDeviceId),
		DeviceSessionId:    r.Header.Get(constants.HeaderDeviceSessionId),
		CallerId:           getServiceClientName(r),
		SecurityAttributes: context.SecurityAttributes{},
	}
}

func getClientIpAddress(r *http.Request) string {
	if clientIpAddress := r.Header.Get(constants.HeaderClientIpAddress); len(clientIpAddress) > 0 {
		return clientIpAddress
	}
	if clientIpAddress := r.Header.Get(constants.HeaderOldClientIpAddress); len(clientIpAddress) > 0 {
		return clientIpAddress
	}
	return r.RemoteAddr
}

func getServiceClientName(r *http.Request) string {
	serviceName := r.Header.Get(constants.HeaderServiceClientName)
	if len(serviceName) > 0 {
		return serviceName
	}
	return r.Header.Get(constants.HeaderOldServiceClientName)
}

func publishEvent(ctx mainContext.Context, requestAttributes *context.RequestAttributes) {
	pubsub.Publish(event.NewRequestCompletedEvent(ctx, event.RequestCompletedPayload{
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
