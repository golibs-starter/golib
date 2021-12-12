package client

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestNewDefaultHttpClient(t *testing.T) {
	props := &HttpClientProperties{
		Timeout:             10 * time.Second,
		MaxIdleConns:        20,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     15,
		Proxy: ProxyProperties{
			Url: "http://localhost:8080",
			AppliedUris: []string{
				"https://example.com/",
			},
		},
	}
	client, err := NewDefaultHttpClient(&http.Client{}, props)
	assert.Nil(t, err)
	assert.IsType(t, &DefaultHttpClient{}, client)
	defaultClient := client.(*DefaultHttpClient)
	assert.Equal(t, defaultClient.properties, props)
	assert.NotNil(t, defaultClient.client)
	assert.Equal(t, props.Timeout, defaultClient.client.Timeout)
	assert.NotNil(t, defaultClient.client.Transport)
}

func TestNewDefaultHttpClient_WhenTransportIsError_ShouldReturnError(t *testing.T) {
	props := &HttpClientProperties{
		Proxy: ProxyProperties{
			Url: "", //make transport error
			AppliedUris: []string{
				"https://example.com/",
			},
		},
	}
	_, err := NewDefaultHttpClient(&http.Client{}, props)
	assert.NotNil(t, err)
}

func TestSetupHttpTransport_WhenProvideValidProps_ShouldReturnCorrectValues(t *testing.T) {
	props := &HttpClientProperties{
		MaxIdleConns:        20,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     15,
		Proxy: ProxyProperties{
			Url: "http://localhost:8080",
			AppliedUris: []string{
				"https://example.com/",
			},
		},
	}
	transport, err := setupHttpTransport(props)
	assert.Nil(t, err)
	assert.Equal(t, props.MaxIdleConns, transport.MaxIdleConns)
	assert.Equal(t, props.MaxIdleConnsPerHost, transport.MaxIdleConnsPerHost)
	assert.Equal(t, props.MaxConnsPerHost, transport.MaxConnsPerHost)
	assert.NotNil(t, transport.Proxy)
}

func TestSetupHttpTransportWithProxy_WhenNotEnabledProxy_ShouldReturnNil(t *testing.T) {
	transport := http.Transport{}
	proxyProps := ProxyProperties{
		Url: "http://localhost:8080",
	}
	err := setupHttpTransportWithProxy(&transport, &proxyProps)
	assert.Nil(t, err)
	assert.Nil(t, transport.Proxy)
}

func TestSetupHttpTransportWithProxy_WhenEnabledProxyAndProxyUrlIsEmpty_ShouldReturnError(t *testing.T) {
	transport := http.Transport{}
	err := setupHttpTransportWithProxy(&transport, &ProxyProperties{
		AppliedUris: []string{
			"https://example.com/",
		},
	})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "proxy url must be defined")
}

func TestSetupHttpTransportWithProxy_WhenEnabledProxyAndProxyUrlIsInvalid_ShouldReturnError(t *testing.T) {
	transport := http.Transport{}
	err := setupHttpTransportWithProxy(&transport, &ProxyProperties{
		Url: "https://abc:zyz/",
		AppliedUris: []string{
			"https://example.com/",
		},
	})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "proxy url is not valid")
	assert.NotNil(t, errors.Cause(err))
}

func TestSetupHttpTransportWithProxy_WhenExecuteWithRequestUrlIsMatched_ShouldReturnProxyUrl(t *testing.T) {
	transport := http.Transport{}
	proxyProps := ProxyProperties{
		Url: "http://localhost:8080",
		AppliedUris: []string{
			"https://example.com/",
		},
	}
	err := setupHttpTransportWithProxy(&transport, &proxyProps)
	reqUrl, _ := url.Parse("https://example.com/path/")
	proxyUrl, err := transport.Proxy(&http.Request{URL: reqUrl})
	assert.Nil(t, err)
	assert.Equal(t, proxyProps.Url, proxyUrl.String())
}

func TestSetupHttpTransportWithProxy_WhenExecuteWithRequestUrlIsNotMatched_ShouldReturnNil(t *testing.T) {
	transport := http.Transport{}
	proxyProps := ProxyProperties{
		Url: "http://localhost:8080",
		AppliedUris: []string{
			"https://example.com/",
		},
	}
	err := setupHttpTransportWithProxy(&transport, &proxyProps)
	reqUrl, _ := url.Parse("https://abc.com/path/")
	proxyUrl, err := transport.Proxy(&http.Request{URL: reqUrl})
	assert.Nil(t, err)
	assert.Nil(t, proxyUrl)
}

func TestDefaultHttpClientMakeRequest_WhenBodyIsRequestReader_ShouldReturnSuccess(t *testing.T) {
	client := &DefaultHttpClient{}
	request, err := client.makeRequest("GET", "https://example.com",
		&testStructWithRequestReader{}, WithHeader("Test1", "Val1"))
	assert.NoError(t, err)
	assert.Equal(t, "GET", request.Method)
	assert.Equal(t, "https://example.com", request.URL.String())
	assert.Equal(t, "Val1", request.Header.Get("Test1"))
	actualRequestBody := new(strings.Builder)
	_, err = io.Copy(actualRequestBody, request.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"field1": "abc"}`, actualRequestBody.String())
}

func TestDefaultHttpClientMakeRequest_WhenBodyIsRequestReaderWithError_ShouldReturnError(t *testing.T) {
	client := &DefaultHttpClient{}
	request, err := client.makeRequest("GET", "https://example.com", &testStructWithRequestReaderError{})
	assert.Nil(t, request)
	assert.Error(t, err)
	assert.Equal(t, "error message 1", err.Error())
}

func TestDefaultHttpClientMakeRequest_WhenBodyIsIoReader_ShouldReturnSuccess(t *testing.T) {
	client := &DefaultHttpClient{}
	request, err := client.makeRequest("GET", "https://example.com",
		strings.NewReader(`{"field1": "abc"}`), WithHeader("Test1", "Val1"))
	assert.NoError(t, err)
	assert.Equal(t, "GET", request.Method)
	assert.Equal(t, "https://example.com", request.URL.String())
	assert.Equal(t, "Val1", request.Header.Get("Test1"))
	actualRequestBody := new(strings.Builder)
	_, err = io.Copy(actualRequestBody, request.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"field1": "abc"}`, actualRequestBody.String())
}

func TestDefaultHttpClientMakeRequest_WhenBodyIsUrlValues_ShouldReturnSuccess(t *testing.T) {
	data := url.Values{}
	data.Set("foo", "bar")
	client := &DefaultHttpClient{}
	request, err := client.makeRequest("GET", "https://example.com",
		data, WithHeader("Test1", "Val1"))
	assert.NoError(t, err)
	assert.Equal(t, "GET", request.Method)
	assert.Equal(t, "https://example.com", request.URL.String())
	assert.Equal(t, "application/x-www-form-urlencoded", request.Header.Get("Content-Type"))
	assert.Equal(t, "7", request.Header.Get("Content-Length"))
	assert.Equal(t, "Val1", request.Header.Get("Test1"))
	actualRequestBody := new(strings.Builder)
	_, err = io.Copy(actualRequestBody, request.Body)
	assert.NoError(t, err)
	assert.Equal(t, `foo=bar`, actualRequestBody.String())
}

func TestDefaultHttpClientMakeRequest_WhenBodyIsUrlValuesAndOverwriteHeader_ShouldReturnSuccessWithNewHeader(t *testing.T) {
	data := url.Values{}
	data.Set("foo", "bar")
	client := &DefaultHttpClient{}
	request, err := client.makeRequest("GET", "https://example.com",
		data, WithContentType("application/url-form-urlencoded"))
	assert.NoError(t, err)
	assert.Equal(t, "GET", request.Method)
	assert.Equal(t, "https://example.com", request.URL.String())
	assert.Equal(t, "application/url-form-urlencoded", request.Header.Get("Content-Type"))
	assert.Equal(t, "7", request.Header.Get("Content-Length"))
	actualRequestBody := new(strings.Builder)
	_, err = io.Copy(actualRequestBody, request.Body)
	assert.NoError(t, err)
	assert.Equal(t, `foo=bar`, actualRequestBody.String())
}

func TestDefaultHttpClientMakeRequest_WhenBodyIsSimpleStruct_ShouldReturnSuccess(t *testing.T) {
	data := testStruct{Field1: "xyz"}
	client := &DefaultHttpClient{}
	request, err := client.makeRequest("GET", "https://example.com", data, WithHeader("Test1", "Val1"))
	assert.NoError(t, err)
	assert.Equal(t, "GET", request.Method)
	assert.Equal(t, "https://example.com", request.URL.String())
	assert.Equal(t, "application/json", request.Header.Get("Content-Type"))
	assert.Equal(t, "Val1", request.Header.Get("Test1"))
	actualRequestBody := new(strings.Builder)
	_, err = io.Copy(actualRequestBody, request.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"field_1":"xyz"}`, strings.TrimSpace(actualRequestBody.String()))
}

type testStruct struct {
	Field1 string `json:"field_1"`
}

type testStructWithRequestReader struct {
	Field1 string
}

func (testStructWithRequestReader) Read() (io.Reader, error) {
	return strings.NewReader(`{"field1": "abc"}`), nil
}

type testStructWithRequestReaderError struct {
	Field1 string
}

func (testStructWithRequestReaderError) Read() (io.Reader, error) {
	return nil, errors.New("error message 1")
}
