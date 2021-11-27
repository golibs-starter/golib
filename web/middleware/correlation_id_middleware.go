package middleware

import (
	"github.com/google/uuid"
	"gitlab.com/golibs-starter/golib/web/constant"
	"gitlab.com/golibs-starter/golib/web/context"
	"gitlab.com/golibs-starter/golib/web/log"
	"net/http"
)

// CorrelationId middleware responsible to inject correlationId to request attributes
// correlationId is usually sent in the request header by the client (see constant.HeaderCorrelationId),
// but sometimes it doesn't exist, we will generate it automatically
func CorrelationId() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			context.GetOrCreateRequestAttributes(r).CorrelationId = getOrNewCorrelationId(r)
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
