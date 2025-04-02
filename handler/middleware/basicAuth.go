package middleware

import (
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler/auth"
	"github.com/TechBowl-japan/go-stations/handler/auth/basic"
)

func BasicAuth(h http.Handler, c auth.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rUserID, rPass, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// if auth.BasicAuthenticater()
		if !basic.Authenticate(rUserID, rPass, auth.NewConfigFromEnv()) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}
