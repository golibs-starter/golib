package client

import (
	"bytes"
	"encoding/json"
	"github.com/golibs-starter/golib/log"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type DefaultHttpClient struct {
	client *http.Client
}

func NewDefaultHttpClient(base *http.Client) HttpClient {
	return &DefaultHttpClient{client: base}
}

func (d *DefaultHttpClient) Get(url string, result interface{}, options ...RequestOption) (*HttpResponse, error) {
	return d.Request(http.MethodGet, url, nil, result, options...)
}

func (d *DefaultHttpClient) Post(url string, body interface{}, result interface{}, options ...RequestOption) (*HttpResponse, error) {
	return d.Request(http.MethodPost, url, body, result, options...)
}

func (d *DefaultHttpClient) Put(url string, body interface{}, result interface{}, options ...RequestOption) (*HttpResponse, error) {
	return d.Request(http.MethodPut, url, body, result, options...)
}

func (d *DefaultHttpClient) Patch(url string, body interface{}, result interface{}, options ...RequestOption) (*HttpResponse, error) {
	return d.Request(http.MethodPatch, url, body, result, options...)
}

func (d *DefaultHttpClient) Delete(url string, body interface{}, result interface{}, options ...RequestOption) (*HttpResponse, error) {
	return d.Request(http.MethodDelete, url, body, result, options...)
}

func (d *DefaultHttpClient) Request(method string, url string, body interface{}, result interface{}, options ...RequestOption) (*HttpResponse, error) {
	request, err := d.makeRequest(method, url, body, options...)
	if err != nil {
		return nil, err
	}

	response, err := d.client.Do(request)
	if err != nil {
		return nil, err
	}
	if response != nil {
		defer func() {
			_ = response.Body.Close()
		}()
	}

	res := new(HttpResponse)
	res.Status = response.Status
	res.StatusCode = response.StatusCode

	bodyWhenError := ""
	if NewHttpSeries(res.StatusCode).IsError() {
		var buf bytes.Buffer
		tee := io.TeeReader(response.Body, &buf)
		str, _ := io.ReadAll(tee)
		bodyWhenError = string(str)
		response.Body = io.NopCloser(bytes.NewBuffer(str))
	}

	if result != nil {
		if err := json.NewDecoder(response.Body).Decode(result); err != nil {
			log.Warnf("[HttpRequestDebug] Decode fail, detail: [%s]", bodyWhenError)
			return res, err
		}
	}

	return res, nil
}

func (d *DefaultHttpClient) makeRequest(method string, reqUrl string, body interface{}, options ...RequestOption) (*http.Request, error) {
	var err error
	var request *http.Request
	if body == nil {
		if request, err = http.NewRequest(method, reqUrl, new(bytes.Buffer)); err != nil {
			return nil, err
		}
	} else {
		switch bodyCasted := body.(type) {
		case io.Reader:
			if request, err = http.NewRequest(method, reqUrl, bodyCasted); err != nil {
				return nil, err
			}
			break
		case url.Values:
			bodyEncoded := bodyCasted.Encode()
			if request, err = http.NewRequest(method, reqUrl, strings.NewReader(bodyEncoded)); err != nil {
				return nil, err
			}
			WithContentType("application/x-www-form-urlencoded")(request)
			WithContentLength(len(bodyEncoded))(request)
			break
		case RequestReader:
			var buf io.Reader
			if buf, err = bodyCasted.Read(); err != nil {
				return nil, err
			}
			if request, err = http.NewRequest(method, reqUrl, buf); err != nil {
				return nil, err
			}
			break
		default:
			var buf = new(bytes.Buffer)
			if err = json.NewEncoder(buf).Encode(body); err != nil {
				return nil, err
			}
			if request, err = http.NewRequest(method, reqUrl, buf); err != nil {
				return nil, err
			}
			WithContentType("application/json")(request)
			break
		}
	}
	if options != nil && len(options) > 0 {
		for _, option := range options {
			option(request)
		}
	}
	return request, nil
}
