package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logger middleware to log each request
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("Method: %s, URI: %s, Time: %v", r.Method, r.RequestURI, time.Since(start))
	})
}
