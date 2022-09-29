package client

import (
	"github.com/pkg/errors"
	assert "github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestNewNativeHttpClient(t *testing.T) {
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
	assert.Equal(t, props.Timeout, nativeHttpClient.Timeout)
	assert.NotNil(t, nativeHttpClient.Transport)
}

func TestNewNativeHttpClient_WhenTransportIsError_ShouldReturnError(t *testing.T) {
	props := &HttpClientProperties{
		Proxy: ProxyProperties{
			Url: "", //make transport error
			AppliedUris: []string{
				"https://example.com/",
			},
		},
	}
	_, err := NewNativeHttpClient(props)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "proxy url must be defined")
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
	assert.ErrorContains(t, err, "proxy url must be defined")
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
	assert.ErrorContains(t, err, "proxy url is not valid")
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
