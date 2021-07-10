package middleware

import (
	"gitlab.id.vin/vincart/golib/web/context"
	"net/http"
)

func AdvancedResponseWriter() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(context.NewAdvancedResponseWriter(w), r)
		})
	}
}
