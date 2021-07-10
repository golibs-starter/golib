package constants

const (
	HeaderUserAgent         = "User-Agent"
	HeaderClientIpAddress   = "Client-IP-Address"
	HeaderServiceClientName = "Service-Client-Name"
	HeaderCorrelationId     = "X-Request-ID"
	HeaderDeviceId          = "X-Device-ID"
	HeaderDeviceSessionId   = "X-Device-Session-ID"

	// Deprecated: Don't use this header any more, keep it for backward compatible
	HeaderOldClientIpAddress = "client_ip_address"

	// Deprecated: Don't use this header any more, keep it for backward compatible
	HeaderOldServiceClientName = "service_client_name"
)
