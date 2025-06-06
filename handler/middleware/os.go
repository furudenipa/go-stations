package middleware

import (
	"fmt"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler/context"
	"github.com/mileusna/useragent"
)

func WithOS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua := useragent.Parse(r.UserAgent())
		fmt.Println(ua.Name)
		ctx := context.WithOS(r.Context(), ua.OS)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}
