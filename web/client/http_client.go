package client

type HttpClient interface {
	Get(url string, result interface{}, options ...RequestOption) (*HttpResponse, error)
	Post(url string, body interface{}, result interface{}, options ...RequestOption) (*HttpResponse, error)
	Put(url string, body interface{}, result interface{}, options ...RequestOption) (*HttpResponse, error)
	Patch(url string, body interface{}, result interface{}, options ...RequestOption) (*HttpResponse, error)
	Delete(url string, body interface{}, result interface{}, options ...RequestOption) (*HttpResponse, error)
	Request(method string, url string, body interface{}, result interface{}, options ...RequestOption) (*HttpResponse, error)
}

type HttpResponse struct {
	Status     string
	StatusCode int
}
