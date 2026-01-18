package config

import (
	"os"
	"strings"
	"time"
)

// Config holds application configuration
type Config struct {
	// Server
	ServerPort string
	BaseURL    string
	Environment string // development, production, test
	UseHTTPS   bool   // Whether HTTPS is enabled

	// Database
	DBPath string

	// File Storage
	UploadDir string
	MaxUploadSize int64 // in bytes

	// JWT
	JWTSecret     string
	JWTExpiration time.Duration

	// Site
	SiteName        string
	SiteDescription string

	// Security
	AllowedOrigins []string // CORS allowed origins
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	env := getEnv("ENV", "development")
	baseURL := getEnv("BASE_URL", "http://localhost:8080")
	useHTTPS := getEnvBool("USE_HTTPS", false)
	
	// If USE_HTTPS is not set, try to detect from BASE_URL
	if !useHTTPS {
		useHTTPS = strings.HasPrefix(baseURL, "https://")
	}
	
	cfg := &Config{
		ServerPort:     getEnv("PORT", "8080"),
		DBPath:          getEnv("DB_PATH", "./cacto.db"),
		BaseURL:         baseURL,
		UploadDir:       getEnv("UPLOAD_DIR", "./web/uploads"),
		MaxUploadSize:   10 * 1024 * 1024, // 10MB
		Environment:     env,
		UseHTTPS:        useHTTPS,
		JWTSecret:       getEnv("JWT_SECRET", generateDefaultSecret()),
		JWTExpiration:   24 * time.Hour,
		SiteName:        getEnv("SITE_NAME", "Cacto CMS"),
		SiteDescription: getEnv("SITE_DESCRIPTION", "Performance-focused enterprise CMS"),
		AllowedOrigins:  getAllowedOrigins(env, baseURL),
	}

	// Warn if using default secret in production
	if cfg.IsProduction() && cfg.JWTSecret == generateDefaultSecret() {
		// Log warning (will be handled by logger)
	}

	return cfg
}

// getAllowedOrigins returns allowed CORS origins based on environment
func getAllowedOrigins(env, baseURL string) []string {
	if env == "production" || env == "prod" {
		// In production, only allow the base URL origin
		if baseURL != "" {
			return []string{baseURL}
		}
		return []string{} // Empty means no CORS (or configure manually)
	}
	// Development: allow all
	return []string{"*"}
}

// IsDevelopment checks if environment is development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development" || c.Environment == "dev"
}

// IsProduction checks if environment is production
func (c *Config) IsProduction() bool {
	return c.Environment == "production" || c.Environment == "prod"
}

// GetCookieSecure returns whether cookies should be Secure (HTTPS only)
func (c *Config) GetCookieSecure() bool {
	return c.UseHTTPS || c.IsProduction()
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBool gets boolean environment variable
func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return strings.ToLower(value) == "true" || value == "1"
}

// generateDefaultSecret generates a default secret (should be overridden in production)
func generateDefaultSecret() string {
	// In production, this should be set via environment variable
	return "change-this-secret-in-production"
}
