package middleware

import (
	"github.com/golibs-starter/golib/web/constant"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_getOrNewCorrelationId_WhenRequestContainsRequestId_ShouldReturnItsRequestID(t *testing.T) {
	r, _ := http.NewRequest("GET", "/test", nil)
	r.Header.Set(constant.HeaderCorrelationId, "test-request-id")
	reqId := getOrNewCorrelationId(r)
	assert.Equal(t, "test-request-id", reqId)
}

func Test_getOrNewCorrelationId_WhenRequestNotContainsRequestId_ShouldReturnNewRequestID(t *testing.T) {
	r, _ := http.NewRequest("GET", "/test", nil)
	reqId := getOrNewCorrelationId(r)
	assert.NotEmpty(t, reqId)
}
