package middleware

import (
	assert "github.com/stretchr/testify/require"
	"gitlab.id.vin/vincart/golib/web/context"
	"net/http"
	"testing"
)

func TestAdvancedResponseWriter_ShouldReplaceDefaultWriter(t *testing.T) {
	next := dummyTestRequestContextHandler{responseStatus: http.StatusOK}
	handler := AdvancedResponseWriter()
	assert.NotNil(t, handler)

	internalHandler := handler(&next)
	assert.NotNil(t, internalHandler)

	handlerFunc, ok := internalHandler.(http.HandlerFunc)
	assert.True(t, ok)

	r, _ := http.NewRequest("GET", "/test", nil)
	handlerFunc(&mockResponseWriter{}, r)

	assert.IsType(t, &context.AdvancedResponseWriter{}, next.writer)
}
