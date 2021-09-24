package constant

const (
	HeaderCorrelationId     = "X-Request-ID"
	HeaderUserAgent         = "User-Agent"
	HeaderClientIpAddress   = "Client-IP-Address"
	HeaderServiceClientName = "Service-Client-Name"
	HeaderDeviceId          = "Device-ID"
	HeaderDeviceSessionId   = "Device-Session-ID"
	HeaderEventId           = "Event-ID"

	// Deprecated: Don't use this header anymore, keep it for backward compatible
	HeaderOldClientIpAddress = "client_ip_address"

	// Deprecated: Don't use this header anymore, keep it for backward compatible
	HeaderOldServiceClientName = "service_client_name"

	// Deprecated: Don't use this header anymore, keep it for backward compatible
	HeaderOldDeviceId = "device_id"

	// Deprecated: Don't use this header anymore, keep it for backward compatible
	HeaderOldDeviceSessionId = "device_session_id"
)
