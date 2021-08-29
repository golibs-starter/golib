package context

import (
	"context"
	"gitlab.id.vin/vincart/golib/web/constant"
	"net/http"
)

func GetRequestAttributes(ctx context.Context) *RequestAttributes {
	reqAttrCtxValue := ctx.Value(constant.ContextReqAttribute)
	if reqAttrCtxValue == nil {
		return nil
	}
	requestAttributes, ok := reqAttrCtxValue.(*RequestAttributes)
	if !ok {
		return nil
	}
	return requestAttributes
}

func GetOrCreateRequestAttributes(r *http.Request) *RequestAttributes {
	requestAttributes := GetRequestAttributes(r.Context())
	if requestAttributes == nil {
		return replaceNewRequestAttributes(r)
	}
	return requestAttributes
}

func replaceNewRequestAttributes(r *http.Request) *RequestAttributes {
	requestAttributes := makeRequestAttributes(r)
	*r = *r.WithContext(context.WithValue(r.Context(), constant.ContextReqAttribute, requestAttributes))
	return requestAttributes
}

func makeRequestAttributes(r *http.Request) *RequestAttributes {
	return &RequestAttributes{
		Method:             r.Method,
		Uri:                r.URL.Path,
		Query:              r.URL.RawQuery,
		Url:                r.URL.String(),
		UserAgent:          r.Header.Get(constant.HeaderUserAgent),
		ClientIpAddress:    getClientIpAddress(r),
		DeviceId:           r.Header.Get(constant.HeaderDeviceId),
		DeviceSessionId:    r.Header.Get(constant.HeaderDeviceSessionId),
		CallerId:           getServiceClientName(r),
		SecurityAttributes: SecurityAttributes{},
	}
}

func getClientIpAddress(r *http.Request) string {
	if clientIpAddress := r.Header.Get(constant.HeaderClientIpAddress); len(clientIpAddress) > 0 {
		return clientIpAddress
	}
	if clientIpAddress := r.Header.Get(constant.HeaderOldClientIpAddress); len(clientIpAddress) > 0 {
		return clientIpAddress
	}
	return r.RemoteAddr
}

func getServiceClientName(r *http.Request) string {
	serviceName := r.Header.Get(constant.HeaderServiceClientName)
	if len(serviceName) > 0 {
		return serviceName
	}
	return r.Header.Get(constant.HeaderOldServiceClientName)
}
