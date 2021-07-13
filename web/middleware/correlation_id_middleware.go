package middleware

import (
	"github.com/google/uuid"
	"gitlab.id.vin/vincart/golib/web/constant"
	"gitlab.id.vin/vincart/golib/web/log"
	"net/http"
)

// CorrelationId middleware responsible for inject correlationId to request attributes
// correlationId is usually sent in the request header by the client (see constant.HeaderCorrelationId),
// but sometimes it doesn't exist, we will generate it automatically
func CorrelationId() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			getOrCreateRequestAttributes(r).CorrelationId = getOrNewCorrelationId(r)
			next.ServeHTTP(w, r)
		})
	}
}

func getOrNewCorrelationId(r *http.Request) string {
	correlationId := r.Header.Get(constant.HeaderCorrelationId)
	if len(correlationId) > 0 {
		return correlationId
	}
	newCorrelationId, err := uuid.NewUUID()
	if err != nil {
		log.Error(r.Context(), "Cannot generate new correlation id with error [%v]", err)
		return ""
	}
	return newCorrelationId.String()
}
