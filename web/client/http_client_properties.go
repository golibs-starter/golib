package client

import "time"

type HttpClientProperties struct {
	// Timeout specifies a time limit for requests made by this client.
	// The timeout includes connection time, any redirects, and reading the response body.
	// A Timeout of zero means no timeout.
	Timeout time.Duration `mapstructure:"timeout" default:"60s"`

	// MaxIdleConns is the connection pool size,
	// and this is the maximum connection that can be open;
	// its default value is 100 connections, zero means no limit.
	MaxIdleConns int `mapstructure:"max_idle_conns" default:"100"`

	// MaxIdleConnsPerHost is the number of connection can be allowed to open per host basic.
	// If zero, http.Transport DefaultMaxIdleConnsPerHost is used.
	MaxIdleConnsPerHost int `mapstructure:"max_idle_conns_per_host" default:"100"`

	// MaxConnsPerHost optionally limits the total number of
	// connections per host, including connections in the dialing,
	// active, and idle states. On limit violation, dials will block.
	// Zero means no limit.
	MaxConnsPerHost int `mapstructure:"max_conns_per_host" default:"100"`

	// Proxy specify proxy url and urls will apply proxy and
	// the requests with these urls will be requested under proxy
	Proxy ProxyProperties `mapstructure:"proxy"`
}

func (h HttpClientProperties) Prefix() string {
	return "vinid.httpclient"
}

type ProxyProperties struct {
	// Url is proxy url. Example: http://localhost:8080
	Url string `mapstructure:"url"`

	// AppliedUris is the list of URIs, which will be requested under above proxy
	// Example:
	// 		https://example.com/path/
	//		All URL starts with https://example.com/path/ will be request under proxy
	AppliedUris []string `mapstructure:"applied_uris"`
}
