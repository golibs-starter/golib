package logging

type LogContext struct {
	CorrelationId     string `json:"request_id"`
	UserId            string `json:"jwt_subject,omitempty"`
	DeviceId          string `json:"device_id,omitempty"`
	DeviceSessionId   string `json:"device_session_id,omitempty"`
	TechnicalUsername string `json:"technical_username,omitempty"`
}
