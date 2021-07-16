package client

import (
	"context"
	"gitlab.id.vin/vincart/golib/web/constant"
	webContext "gitlab.id.vin/vincart/golib/web/context"
	"net/http"
)

type TraceableHttpClient struct {
	client HttpClient
}

func NewTraceableHttpClient(client HttpClient) *TraceableHttpClient {
	return &TraceableHttpClient{client: client}
}

func (t *TraceableHttpClient) Get(ctx context.Context, url string, result interface{}, options ...RequestOption) (*HttpResponse, error) {
	return t.Request(ctx, http.MethodGet, url, nil, result, options...)
}

func (t *TraceableHttpClient) Post(ctx context.Context, url string, body interface{}, result interface{}, options ...RequestOption) (*HttpResponse, error) {
	return t.Request(ctx, http.MethodPost, url, body, result, options...)
}

func (t *TraceableHttpClient) Put(ctx context.Context, url string, body interface{}, result interface{}, options ...RequestOption) (*HttpResponse, error) {
	return t.Request(ctx, http.MethodPut, url, body, result, options...)
}

func (t *TraceableHttpClient) Patch(ctx context.Context, url string, body interface{}, result interface{}, options ...RequestOption) (*HttpResponse, error) {
	return t.Request(ctx, http.MethodPatch, url, body, result, options...)
}

func (t *TraceableHttpClient) Delete(ctx context.Context, url string, body interface{}, result interface{}, options ...RequestOption) (*HttpResponse, error) {
	return t.Request(ctx, http.MethodDelete, url, body, result, options...)
}

func (t *TraceableHttpClient) Request(ctx context.Context, method string, url string, body interface{}, result interface{}, options ...RequestOption) (*HttpResponse, error) {
	httpOpts := make([]RequestOption, 0)
	if reqAttrs := webContext.GetRequestAttributes(ctx); reqAttrs != nil {
		httpOpts = append(httpOpts,
			WithHeader(constant.HeaderCorrelationId, reqAttrs.CorrelationId),
			WithHeader(constant.HeaderDeviceId, reqAttrs.DeviceId),
			// TODO add more header
		)
	}
	httpOpts = append(httpOpts, options...)
	return t.client.Request(method, url, body, result, httpOpts...)
}
