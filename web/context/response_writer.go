package context

import "net/http"

type WrappingResponseWriter interface {
	Writer() http.ResponseWriter
}

// AdvancedResponseWriter provide method to write response and
// effective ways to get response with more details
type AdvancedResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewAdvancedResponseWriter(w http.ResponseWriter) *AdvancedResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &AdvancedResponseWriter{w, http.StatusOK}
}

func (h *AdvancedResponseWriter) WriteHeader(code int) {
	h.statusCode = code
	h.ResponseWriter.WriteHeader(code)
}

func (h AdvancedResponseWriter) Status() int {
	return h.statusCode
}
