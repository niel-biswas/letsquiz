package middleware

import (
	"letsquiz/config"
	"net/http"
	"sync"
	"time"
)

// RateLimiter middleware to limit the number of requests
func RateLimiter(next http.Handler) http.Handler {
	var (
		mu      sync.Mutex
		visits  = make(map[string]int)
		limit   = config.DbConfig.RateLimit // limit of requests per minute per IP
		timeout = time.Minute
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		mu.Lock()
		visits[ip]++
		count := visits[ip]
		mu.Unlock()

		if count > limit {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		time.AfterFunc(timeout, func() {
			mu.Lock()
			defer mu.Unlock()
			visits[ip]--
		})

		next.ServeHTTP(w, r)
	})
}
