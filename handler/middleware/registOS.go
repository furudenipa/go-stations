package middleware

import (
	"context"
	"net/http"

	"github.com/mileusna/useragent"
)

type OSContextKey string

func WithOS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua := useragent.Parse(r.UserAgent())
		ctx := context.WithValue(r.Context(), OSContextKey("os"), ua.OS)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}
