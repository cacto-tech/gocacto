package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"
	"time"
)

var (
	csrfTokens = make(map[string]time.Time)
	csrfMutex  sync.RWMutex
	csrfTTL    = 24 * time.Hour
)

// CSRFProtection provides CSRF token generation and validation
type CSRFProtection struct {
	exemptPaths []string
}

// NewCSRFProtection creates a new CSRF protection middleware
func NewCSRFProtection(exemptPaths []string) *CSRFProtection {
	// Clean up expired tokens periodically
	go cleanupExpiredTokens()
	
	return &CSRFProtection{
		exemptPaths: exemptPaths,
	}
}

// GenerateToken generates a new CSRF token
func GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	
	token := base64.URLEncoding.EncodeToString(bytes)
	
	csrfMutex.Lock()
	csrfTokens[token] = time.Now().Add(csrfTTL)
	csrfMutex.Unlock()
	
	return token, nil
}

// ValidateToken validates a CSRF token
func ValidateToken(token string) bool {
	csrfMutex.RLock()
	expires, exists := csrfTokens[token]
	csrfMutex.RUnlock()
	
	if !exists {
		return false
	}
	
	if time.Now().After(expires) {
		csrfMutex.Lock()
		delete(csrfTokens, token)
		csrfMutex.Unlock()
		return false
	}
	
	return true
}

// Middleware validates CSRF tokens for state-changing requests
func (c *CSRFProtection) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip CSRF check for exempt paths
		for _, path := range c.exemptPaths {
			if r.URL.Path == path {
				next.ServeHTTP(w, r)
				return
			}
		}
		
		// Only check state-changing methods
		if r.Method == "GET" || r.Method == "HEAD" || r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}
		
		// Get token from header or form
		token := r.Header.Get("X-CSRF-Token")
		if token == "" {
			token = r.FormValue("csrf_token")
		}
		
		// Also check cookie (double submit cookie pattern)
		cookieToken, err := r.Cookie("csrf_token")
		if err == nil && cookieToken.Value != "" {
			// If header token matches cookie token, it's valid
			if token == cookieToken.Value && ValidateToken(token) {
				next.ServeHTTP(w, r)
				return
			}
		}
		
		// Validate token
		if token == "" || !ValidateToken(token) {
			http.Error(w, "Invalid CSRF token", http.StatusForbidden)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// cleanupExpiredTokens periodically removes expired tokens
func cleanupExpiredTokens() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	
	for range ticker.C {
		csrfMutex.Lock()
		now := time.Now()
		for token, expires := range csrfTokens {
			if now.After(expires) {
				delete(csrfTokens, token)
			}
		}
		csrfMutex.Unlock()
	}
}
