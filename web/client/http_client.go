package client

import "context"

type HttpClient interface {
	Get(url string, result interface{}, options ...RequestOption) (*HttpResponse, error)
	Post(url string, body interface{}, result interface{}, options ...RequestOption) (*HttpResponse, error)
	Put(url string, body interface{}, result interface{}, options ...RequestOption) (*HttpResponse, error)
	Patch(url string, body interface{}, result interface{}, options ...RequestOption) (*HttpResponse, error)
	Delete(url string, body interface{}, result interface{}, options ...RequestOption) (*HttpResponse, error)
	Request(method string, url string, body interface{}, result interface{}, options ...RequestOption) (*HttpResponse, error)
}

type ContextualHttpClient interface {
	Get(ctx context.Context, url string, result interface{}, options ...RequestOption) (*HttpResponse, error)
	Post(ctx context.Context, url string, body interface{}, result interface{}, options ...RequestOption) (*HttpResponse, error)
	Put(ctx context.Context, url string, body interface{}, result interface{}, options ...RequestOption) (*HttpResponse, error)
	Patch(ctx context.Context, url string, body interface{}, result interface{}, options ...RequestOption) (*HttpResponse, error)
	Delete(ctx context.Context, url string, body interface{}, result interface{}, options ...RequestOption) (*HttpResponse, error)
	Request(ctx context.Context, method string, url string, body interface{}, result interface{}, options ...RequestOption) (*HttpResponse, error)
}

type HttpResponse struct {
	Status     string
	StatusCode int
}
