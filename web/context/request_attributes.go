package context

import "time"

type RequestAttributes struct {
	ServiceCode        string             `json:"service_code"`
	StatusCode         int                `json:"status"`
	ExecutionTime      time.Duration      `json:"duration_ms"`
	Uri                string             `json:"uri"`
	Query              string             `json:"query"`
	Mapping            string             `json:"mapping"`
	Url                string             `json:"url"`
	Method             string             `json:"method"`
	CallerId           string             `json:"caller_id"`
	DeviceId           string             `json:"device_id"`
	DeviceSessionId    string             `json:"device_session_id"`
	CorrelationId      string             `json:"correlation_id"`
	ClientIpAddress    string             `json:"client_ip_address"`
	UserAgent          string             `json:"user_agent"`
	SecurityAttributes SecurityAttributes `json:"security_attributes"`
}
