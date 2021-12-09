package client

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
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

func Test_setupHttpTransport(t *testing.T) {
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

func Test_setupHttpTransportWithProxy_WhenNotEnabledProxy_ShouldReturnNil(t *testing.T) {
	transport := http.Transport{}
	proxyProps := ProxyProperties{
		Url: "http://localhost:8080",
	}
	err := setupHttpTransportWithProxy(&transport, &proxyProps)
	assert.Nil(t, err)
	assert.Nil(t, transport.Proxy)
}

func Test_setupHttpTransportWithProxy_WhenEnabledProxyAndProxyUrlIsEmpty_ShouldReturnError(t *testing.T) {
	transport := http.Transport{}
	err := setupHttpTransportWithProxy(&transport, &ProxyProperties{
		AppliedUris: []string{
			"https://example.com/",
		},
	})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "proxy url must be defined")
}

func Test_setupHttpTransportWithProxy_WhenEnabledProxyAndProxyUrlIsInvalid_ShouldReturnError(t *testing.T) {
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

func Test_setupHttpTransportWithProxy_WhenExecuteWithRequestUrlIsMatched_ShouldReturnProxyUrl(t *testing.T) {
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

func Test_setupHttpTransportWithProxy_WhenExecuteWithRequestUrlIsNotMatched_ShouldReturnNil(t *testing.T) {
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
