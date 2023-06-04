package listener

import (
	"gitlab.com/golibs-starter/golib/log"
	webLog "gitlab.com/golibs-starter/golib/web/log"
)

type HttpRequestLog struct {
	webLog.ContextAttributes
	Url            string `json:"url"`
	Status         int    `json:"status"`
	ExecutionTime  int64  `json:"execution_time"`
	RequestPattern string `json:"request_pattern,omitempty"`
	RequestPath    string `json:"request_path"`
	Method         string `json:"request_method"`
	Query          string `json:"query,omitempty"`
	RequestId      string `json:"request_id,omitempty"`
	CallerId       string `json:"caller_id,omitempty"`
	ClientIp       string `json:"client_ip,omitempty"`
	UserAgent      string `json:"user_agent,omitempty"`
}

func (h HttpRequestLog) MarshalLogObject(encoder log.ObjectEncoder) error {
	err := h.ContextAttributes.MarshalLogObject(encoder)
	if err != nil {
		return err
	}
	encoder.AddString("url", h.Url)
	encoder.AddInt("status", h.Status)
	encoder.AddInt64("execution_time", h.ExecutionTime)
	encoder.AddString("request_pattern", h.RequestPattern)
	encoder.AddString("request_path", h.RequestPath)
	encoder.AddString("request_method", h.Method)
	encoder.AddString("query", h.Query)
	encoder.AddString("request_id", h.RequestId)
	encoder.AddString("caller_id", h.CallerId)
	encoder.AddString("client_ip", h.ClientIp)
	encoder.AddString("user_agent", h.UserAgent)
	return nil
}
