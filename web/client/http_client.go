package client

import (
	"bytes"
	"encoding/json"
	"gitlab.id.vin/vincart/golib/log"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type HeaderAttributes struct {
	RequestID string
}

// HttpClient ...
type HttpClient struct {
	client     *http.Client
	properties *HttpClientProperties

	proxyURL          *url.URL
	basicAuthUsername string
	basicAuthPassword string
}

// HttpResponse ...
type HttpResponse struct {
	Status     string
	StatusCode int
}

// NewHttpClient create new http client
func NewHttpClient(properties *HttpClientProperties) *HttpClient {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = properties.MaxIdleConns
	transport.MaxIdleConnsPerHost = properties.MaxIdleConnsPerHost
	transport.MaxConnsPerHost = properties.MaxConnsPerHost
	//transport.Proxy = http.ProxyURL(hc.proxyURL)
	return &HttpClient{
		client: &http.Client{
			Timeout:   properties.Timeout,
			Transport: transport,
		},
		properties: properties,
	}
}

// Request ...
func (hc *HttpClient) Request(method string, url string, headers map[string]string, parameters interface{}, result interface{}) (*HttpResponse, error) {
	request, err := hc.MakeRequest(method, url, headers, parameters)
	if err != nil {
		return nil, err
	}

	response, err := hc.client.Do(request)
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

	logMsg := ""
	if res.StatusCode != http.StatusOK {
		var buf bytes.Buffer
		tee := io.TeeReader(response.Body, &buf)
		str, _ := ioutil.ReadAll(tee)
		logMsg = string(str)
		response.Body = ioutil.NopCloser(bytes.NewBuffer(str))
	}

	if result != nil {
		if err := json.NewDecoder(response.Body).Decode(result); err != nil {
			log.Warnf("[HttpRequestDebug] Decode fail, detail: %+v", logMsg)
			return res, err
		}
	}

	return res, nil
}

// MakeRequest ...
func (hc *HttpClient) MakeRequest(method string, url string, headers map[string]string, parameters interface{}) (*http.Request, error) {
	buf := new(bytes.Buffer)
	if parameters != nil {
		if err := json.NewEncoder(buf).Encode(parameters); err != nil {
			return nil, err
		}
	}
	request, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}
	if hc.basicAuthPassword != "" && hc.basicAuthUsername != "" {
		request.SetBasicAuth(hc.basicAuthUsername, hc.basicAuthPassword)
	}
	request.Header.Set("Content-Type", "application/json")
	if headers != nil {
		for key, value := range headers {
			request.Header.Set(key, value)
		}
	}
	return request, nil
}
