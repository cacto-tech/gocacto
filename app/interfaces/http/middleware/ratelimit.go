package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/httprate"
)

// RateLimit creates a rate limiting middleware
func RateLimit(requestsPerMinute int) func(http.Handler) http.Handler {
	return httprate.Limit(
		requestsPerMinute,
		1*time.Minute,
		httprate.WithKeyFuncs(func(r *http.Request) (string, error) {
			// Rate limit by IP address
			ip, err := httprate.KeyByIP(r)
			if err != nil {
				return "", err
			}
			return ip, nil
		}),
	)
}

// RateLimitAuth creates stricter rate limiting for auth endpoints
func RateLimitAuth() func(http.Handler) http.Handler {
	// Stricter limit for auth: 5 requests per minute
	return RateLimit(5)
}

// RateLimitAPI creates rate limiting for API endpoints
func RateLimitAPI() func(http.Handler) http.Handler {
	// API limit: 60 requests per minute
	return RateLimit(60)
}
