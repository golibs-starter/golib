package exception

// New returns a new error instance
func New(code uint, message string) Exception {
	return withCode{
		code:    code,
		message: message,
	}
}

// withCode represents a error with code and message
type withCode struct {
	code    uint
	message string
}

func (e withCode) Code() uint {
	return e.code
}

func (e withCode) Error() string {
	return e.message
}
