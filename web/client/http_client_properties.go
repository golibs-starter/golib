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
}

func (h HttpClientProperties) Prefix() string {
	return "vinid.httpclient"
}
