package exception

type Exception interface {

	// Code returns error code
	Code() uint

	// Error returns error message
	Error() string
}
