package log

type HttpRequestLog struct {
	LoggingContext
	Status         int    `json:"status"`
	ExecutionTime  int64  `json:"execution_time"`
	RequestPattern string `json:"request_pattern"`
	RequestPath    string `json:"request_path"`
	Method         string `json:"request_method"`
	Query          string `json:"query,omitempty"`
	Url            string `json:"url"`
	RequestId      string `json:"request_id"`
	CallerId       string `json:"caller_id"`
	ClientIp       string `json:"client_ip"`
	UserAgent      string `json:"user_agent"`
}
