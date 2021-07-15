package client

import "net/http"

type RequestOption func(r *http.Request)

func WithBasicAuth(username string, password string) RequestOption {
	return func(r *http.Request) {
		r.SetBasicAuth(username, password)
	}
}

func WithContentType(contentType string) RequestOption {
	return func(r *http.Request) {
		r.Header.Set("Content-Type", contentType)
	}
}

func WithHeader(key string, value string) RequestOption {
	return func(r *http.Request) {
		r.Header.Set(key, value)
	}
}
