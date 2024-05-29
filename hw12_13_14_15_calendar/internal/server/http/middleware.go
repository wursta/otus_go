package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := fmt.Sprintf(
			"%s [%s] %s %s %s \"%s\"",
			r.RemoteAddr,
			time.Now().Format(time.RFC3339),
			r.Method,
			r.URL.Path,
			r.Proto,
			r.Header.Get("User-Agent"),
		)
		fmt.Println(s)

		next.ServeHTTP(w, r)
	})
}
