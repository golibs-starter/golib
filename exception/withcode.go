package exception

// New returns a new error instance
func New(code uint, message string) Exception {
	return withCode{
		code:    code,
		message: message,
		details: message,
	}
}

func NewWithCause(cause Exception, message string) Exception {
	return withCode{
		code:    cause.Code(),
		message: cause.Message(),
		details: cause.Error() + ": " + message,
	}
}

// withCode represents an error with code and message
type withCode struct {
	code    uint
	message string
	details string
}

func (e withCode) Code() uint {
	return e.code
}

func (e withCode) Message() string {
	return e.message
}

func (e withCode) Error() string {
	return e.details
}
