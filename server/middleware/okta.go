package middleware

import (
	"letsquiz/config"
	"net/http"
)

// OktaAuth middleware to check authorization using Okta
func OktaAuth(next http.Handler) http.Handler {
	if !config.AppConfig.EnableOktaAuth {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		token := authHeader[len("Bearer "):]
		if !validateOktaToken(token) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Validate Okta token (this is a stub, replace with actual Okta validation logic)
func validateOktaToken(token string) bool {
	// TODO Okta's API to validate the token here
	// TODO Using the issuer and client ID from the config
	// TODO issuer := config.AppConfig.OktaIssuer
	// TODO clientID := config.AppConfig.OktaClientID

	return token != ""
}
