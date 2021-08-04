package resource

import (
	"gitlab.id.vin/vincart/golib/exception"
	"net/http"
)

type Resource struct {
	Meta Meta        `json:"meta,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func New(code int, message string, data interface{}) Resource {
	return Resource{
		Meta: Meta{
			Code:    code,
			Message: message,
		},
		Data: data,
	}
}

func Ok(data interface{}) Resource {
	return New(http.StatusOK, "Successful", data)
}

func Created(data interface{}) Resource {
	return New(http.StatusCreated, "Resource created", data)
}

func Error(err error) Resource {
	code := http.StatusInternalServerError
	message := "Internal Server Error"
	if e, ok := err.(exception.Exception); ok {
		code = int(e.Code())
		message = err.Error()
	}
	return Resource{
		Meta: Meta{
			Code:    code,
			Message: message,
		},
	}
}
