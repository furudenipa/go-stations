package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/TechBowl-japan/go-stations/handler/context"
)

type log struct {
	Timestamp time.Time `json:"timestamp"`
	Latency   int64     `json:"latency"`
	Path      string    `json:"path"`
	OS        string    `json:"os"`
}

func Log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()
		h.ServeHTTP(w, r)
		log := log{
			Timestamp: timeStart,
			Latency:   time.Now().Sub(timeStart).Milliseconds(),
			Path:      r.URL.Path,
			OS:        context.OS(r.Context()),
		}
		err := json.NewEncoder(os.Stdout).Encode(log)
		if err != nil {
			fmt.Println("Log: failed to encode log, err =", err)
		}
	})
}
