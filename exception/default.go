package exception

import "net/http"

var (
	BadRequest = New(http.StatusBadRequest, "Bad Request")
	Forbidden  = New(http.StatusForbidden, "Forbidden")
	NotFound   = New(http.StatusNotFound, "Resource Not Found")
)
