package client

import (
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strings"
)

func NewNativeHttpClient(properties *HttpClientProperties) (*http.Client, error) {
	base := &http.Client{}
	transport, err := setupHttpTransport(properties)
	if err != nil {
		return nil, err
	}
	base.Timeout = properties.Timeout
	base.Transport = transport
	return base, nil
}

func setupHttpTransport(props *HttpClientProperties) (*http.Transport, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = props.MaxIdleConns
	transport.MaxIdleConnsPerHost = props.MaxIdleConnsPerHost
	transport.MaxConnsPerHost = props.MaxConnsPerHost
	if err := setupHttpTransportWithProxy(transport, &props.Proxy); err != nil {
		return nil, errors.WithMessage(err, "cannot setup http transport proxy")
	}
	return transport, nil
}

func setupHttpTransportWithProxy(t *http.Transport, proxyProps *ProxyProperties) error {
	var enabledProxy = false
	if len(proxyProps.AppliedUris) > 0 {
		if len(proxyProps.Url) == 0 {
			return errors.New("proxy url must be defined")
		}
		enabledProxy = true
	}
	if !enabledProxy {
		return nil
	}
	proxyUrl, err := url.Parse(proxyProps.Url)
	if err != nil {
		return errors.WithMessage(err, "proxy url is not valid")
	}
	appliedUrls := proxyProps.AppliedUris
	t.Proxy = func(r *http.Request) (*url.URL, error) {
		for _, appliedUrl := range appliedUrls {
			if strings.HasPrefix(r.URL.String(), appliedUrl) {
				return proxyUrl, nil
			}
		}
		// No proxy is used
		return nil, nil
	}
	return nil
}
