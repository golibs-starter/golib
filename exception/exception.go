package exception

type Exception interface {

	// Code returns error code
	Code() uint

	// Message return the simple message
	Message() string

	// Error returns detailed error message
	Error() string
}
