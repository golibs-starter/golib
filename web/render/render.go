package render

import (
	"errors"
	"gitlab.com/golibs-starter/golib/log"
	"net/http"
	"syscall"
)

// Renderer interface is to be implemented by JSON, XML, HTML, YAML and so on.
type Renderer interface {
	// Render writes data with custom ContentType.
	Render(http.ResponseWriter) error

	// WriteContentType writes custom ContentType.
	WriteContentType(w http.ResponseWriter)
}

// Render writes data with custom http status code
func Render(w http.ResponseWriter, httpStatus int, r Renderer) {
	w.WriteHeader(httpStatus)
	if err := r.Render(w); err != nil &&
		!errors.Is(err, syscall.EPIPE) &&
		!errors.Is(err, syscall.ECONNRESET) {
		log.Errorf("Cannot render response with error [%v]", err)
	}
}

// writeContentType writes content type to a writer
func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
