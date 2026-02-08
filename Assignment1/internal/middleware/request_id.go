package middleware

import (
	"net/http"
	"strconv"
	"time"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := strconv.FormatInt(time.Now().UnixNano(), 10)
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r)
	})
}
