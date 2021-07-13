package client

import "time"

type HttpClientConfig struct {
	ConnectionRequestTimeout time.Duration `mapstructure:"connection_request_timeout"`
	ConnectTimeout           time.Duration `mapstructure:"connect_timeout"`
}

func (h HttpClientConfig) Prefix() string {
	return "vinid.httpclient"
}
