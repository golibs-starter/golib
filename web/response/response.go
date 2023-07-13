package response

import (
	"github.com/pkg/errors"
	"gitlab.com/golibs-starter/golib/exception"
	"net/http"
)

type Response[T any] struct {
	Meta Meta        `json:"meta,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func New[T any](code int, message string, data T) Response[T] {
	return Response[T]{
		Meta: Meta{
			Code:    code,
			Message: message,
		},
		Data: data,
	}
}

func Ok[T any](data T) Response[T] {
	return New(http.StatusOK, "Successful", data)
}

func Created[T any](data T) Response[T] {
	return New(http.StatusCreated, "Resource created", data)
}

func Error(err error) Response[error] {
	code := http.StatusInternalServerError
	message := "Internal Server Error"
	switch e := errors.Cause(err).(type) {
	case exception.Exception:
		code = int(e.Code())
		message = e.Message()
	}
	return Response[error]{
		Meta: Meta{
			Code:    code,
			Message: message,
		},
	}
}
