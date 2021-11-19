package middleware

import (
	assert "github.com/stretchr/testify/require"
	"gitlab.id.vin/vincart/common/golib/pubsub"
	"gitlab.id.vin/vincart/common/golib/web/constant"
	"gitlab.id.vin/vincart/common/golib/web/context"
	"gitlab.id.vin/vincart/common/golib/web/event"
	"net/http"
	"testing"
)

type mockEventPublisher struct {
	Event pubsub.Event
}

func (m *mockEventPublisher) Publish(event pubsub.Event) {
	m.Event = event
}

type mockResponseWriter struct {
}

func (d mockResponseWriter) Header() http.Header {
	return map[string][]string{}
}

func (d mockResponseWriter) Write(bytes []byte) (int, error) {
	return 0, nil
}

func (d mockResponseWriter) WriteHeader(statusCode int) {
}

type dummyTestRequestContextHandler struct {
	writer         http.ResponseWriter
	request        *http.Request
	responseStatus int
}

func (d *dummyTestRequestContextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(d.responseStatus)
	d.writer = w
	d.request = r
}

func TestRequestContext_ShouldAttachAttributesToTheRequest(t *testing.T) {
	publisher := &mockEventPublisher{}
	pubsub.ReplaceGlobal(publisher)

	next := dummyTestRequestContextHandler{responseStatus: http.StatusOK}
	handler := RequestContext()
	assert.NotNil(t, handler)

	internalHandler := handler(&next)
	assert.NotNil(t, internalHandler)

	handlerFunc, ok := internalHandler.(http.HandlerFunc)
	assert.True(t, ok)

	r, _ := http.NewRequest("GET", "/test?q=keyword", nil)
	r.Header.Set(constant.HeaderUserAgent, "FAKE-UA")
	r.Header.Set(constant.HeaderClientIpAddress, "1.1.1.1")
	r.Header.Set(constant.HeaderDeviceId, "fake-device-id")
	r.Header.Set(constant.HeaderDeviceSessionId, "fake-device-session-id")
	r.Header.Set(constant.HeaderServiceClientName, "fake-caller-service-name")
	handlerFunc(context.NewAdvancedResponseWriter(&mockResponseWriter{}), r)

	val := r.Context().Value(constant.ContextReqAttribute)
	assert.NotNil(t, val)
	assert.IsType(t, &context.RequestAttributes{}, val)

	requestAttr := val.(*context.RequestAttributes)
	assert.Equal(t, http.StatusOK, requestAttr.StatusCode)
	assert.Equal(t, "GET", requestAttr.Method)
	assert.Equal(t, "/test", requestAttr.Uri)
	assert.Equal(t, "q=keyword", requestAttr.Query)
	assert.Equal(t, "/test?q=keyword", requestAttr.Url)
	assert.Equal(t, "FAKE-UA", requestAttr.UserAgent)
	assert.Equal(t, "1.1.1.1", requestAttr.ClientIpAddress)
	assert.Equal(t, "1.1.1.1", requestAttr.ClientIpAddress)
	assert.Equal(t, "fake-device-id", requestAttr.DeviceId)
	assert.Equal(t, "fake-device-session-id", requestAttr.DeviceSessionId)
	assert.Equal(t, "fake-caller-service-name", requestAttr.CallerId)
	assert.NotNil(t, requestAttr.SecurityAttributes)

	assert.NotNil(t, publisher.Event)
	assert.IsType(t, &event.RequestCompletedEvent{}, publisher.Event)
	requestCompletedEvent := publisher.Event.(*event.RequestCompletedEvent)
	assert.IsType(t, event.RequestCompletedPayload{}, requestCompletedEvent.Payload())
	payload := requestCompletedEvent.Payload().(event.RequestCompletedPayload)
	assert.Equal(t, http.StatusOK, payload.Status)
	assert.NotZero(t, payload.ExecutionTime)
	assert.Equal(t, "/test", payload.Uri)
	assert.Equal(t, "q=keyword", payload.Query)
	assert.Equal(t, "/test?q=keyword", payload.Url)
	assert.Equal(t, "FAKE-UA", payload.UserAgent)
	assert.Equal(t, "1.1.1.1", payload.ClientIpAddress)
	assert.Equal(t, "1.1.1.1", payload.ClientIpAddress)
	assert.Equal(t, "fake-device-id", payload.DeviceId)
	assert.Equal(t, "fake-device-session-id", payload.DeviceSessionId)
	assert.Equal(t, "fake-caller-service-name", payload.CallerId)
	assert.Empty(t, payload.UserId)
	assert.Empty(t, payload.TechnicalUsername)
}

func TestRequestContext_WhenReturnBadRequest_ShouldAttachRequestAttributesCorrectly(t *testing.T) {
	publisher := &mockEventPublisher{}
	pubsub.ReplaceGlobal(publisher)

	next := dummyTestRequestContextHandler{responseStatus: http.StatusBadRequest}
	handler := RequestContext()
	assert.NotNil(t, handler)

	internalHandler := handler(&next)
	assert.NotNil(t, internalHandler)

	handlerFunc, ok := internalHandler.(http.HandlerFunc)
	assert.True(t, ok)

	r, _ := http.NewRequest("GET", "/test", nil)
	handlerFunc(context.NewAdvancedResponseWriter(&mockResponseWriter{}), r)

	val := r.Context().Value(constant.ContextReqAttribute)
	assert.NotNil(t, val)
	requestAttr, ok := val.(*context.RequestAttributes)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, requestAttr.StatusCode)

	assert.NotNil(t, publisher.Event)
	assert.IsType(t, &event.RequestCompletedEvent{}, publisher.Event)
	requestCompletedEvent := publisher.Event.(*event.RequestCompletedEvent)
	assert.IsType(t, event.RequestCompletedPayload{}, requestCompletedEvent.Payload())
	payload := requestCompletedEvent.Payload().(event.RequestCompletedPayload)
	assert.Equal(t, http.StatusBadRequest, payload.Status)
}
