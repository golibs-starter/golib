package client

import (
	"context"
	"gitlab.com/golibs-starter/golib/config"
	"gitlab.com/golibs-starter/golib/web/constant"
	webContext "gitlab.com/golibs-starter/golib/web/context"
	"net/http"
)

type TraceableHttpClient struct {
	client   HttpClient
	appProps *config.AppProperties
}

func NewTraceableHttpClient(client HttpClient, appProps *config.AppProperties) ContextualHttpClient {
	return &TraceableHttpClient{
		client:   client,
		appProps: appProps,
	}
}

func (t *TraceableHttpClient) Get(ctx context.Context, url string, result interface{},
	options ...RequestOption) (*HttpResponse, error) {
	return t.Request(ctx, http.MethodGet, url, nil, result, options...)
}

func (t *TraceableHttpClient) Post(ctx context.Context, url string, body interface{}, result interface{},
	options ...RequestOption) (*HttpResponse, error) {
	return t.Request(ctx, http.MethodPost, url, body, result, options...)
}

func (t *TraceableHttpClient) Put(ctx context.Context, url string, body interface{}, result interface{},
	options ...RequestOption) (*HttpResponse, error) {
	return t.Request(ctx, http.MethodPut, url, body, result, options...)
}

func (t *TraceableHttpClient) Patch(ctx context.Context, url string, body interface{}, result interface{},
	options ...RequestOption) (*HttpResponse, error) {
	return t.Request(ctx, http.MethodPatch, url, body, result, options...)
}

func (t *TraceableHttpClient) Delete(ctx context.Context, url string, body interface{}, result interface{},
	options ...RequestOption) (*HttpResponse, error) {
	return t.Request(ctx, http.MethodDelete, url, body, result, options...)
}

func (t *TraceableHttpClient) Request(ctx context.Context, method string, url string, body interface{},
	result interface{}, options ...RequestOption) (*HttpResponse, error) {
	httpOpts := []RequestOption{
		WithHeader(constant.HeaderServiceClientName, t.appProps.Name),

		// Set header with old format for backward compatible
		WithHeader(constant.HeaderOldServiceClientName, t.appProps.Name),
	}
	if reqAttrs := webContext.GetRequestAttributes(ctx); reqAttrs != nil {
		httpOpts = append(httpOpts,
			WithHeader(constant.HeaderCorrelationId, reqAttrs.CorrelationId),
			WithHeader(constant.HeaderDeviceId, reqAttrs.DeviceId),
			WithHeader(constant.HeaderDeviceSessionId, reqAttrs.DeviceSessionId),
			WithHeader(constant.HeaderClientIpAddress, reqAttrs.ClientIpAddress),

			// Set header with old format for backward compatible
			WithHeader(constant.HeaderOldDeviceId, reqAttrs.DeviceId),
			WithHeader(constant.HeaderOldDeviceSessionId, reqAttrs.DeviceSessionId),
			WithHeader(constant.HeaderOldClientIpAddress, reqAttrs.ClientIpAddress),
		)
	}
	httpOpts = append(httpOpts, options...)
	return t.client.Request(method, url, body, result, httpOpts...)
}
