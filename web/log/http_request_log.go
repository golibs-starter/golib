package log

type HttpRequestLog struct {
	LoggingContext
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
