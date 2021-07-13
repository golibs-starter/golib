package client

import "time"

type HttpClientProperties struct {
	ConnectionRequestTimeout time.Duration `mapstructure:"connection_request_timeout"`
	ConnectTimeout           time.Duration `mapstructure:"connect_timeout"`
}

func (h HttpClientProperties) Prefix() string {
	return "vinid.httpclient"
}
