package exception

import "net/http"

var (
	BadRequest   = New(http.StatusBadRequest, "Bad Request")
	Unauthorized = New(http.StatusUnauthorized, "Unauthorized")
	Forbidden    = New(http.StatusForbidden, "Forbidden")
	NotFound     = New(http.StatusNotFound, "Resource Not Found")
	SystemError  = New(http.StatusInternalServerError, "System Error")
)
