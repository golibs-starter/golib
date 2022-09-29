package client

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io"
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
	nativeHttpClient, err := NewNativeHttpClient(props)
	assert.NoError(t, err)
	client := NewDefaultHttpClient(nativeHttpClient)
	assert.IsType(t, &DefaultHttpClient{}, client)
	defaultClient := client.(*DefaultHttpClient)
	assert.NotNil(t, defaultClient.client)
	assert.Equal(t, props.Timeout, defaultClient.client.Timeout)
	assert.NotNil(t, defaultClient.client.Transport)
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
