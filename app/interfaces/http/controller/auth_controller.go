package controller

import (
	"encoding/json"
	"net/http"

	"cacto-cms/app/application/auth"
	"cacto-cms/app/interfaces/http/middleware"
	"cacto-cms/app/shared/errors"
	"cacto-cms/app/shared/validation"
	"cacto-cms/config"
)

// AuthController handles authentication requests
type AuthController struct {
	authService *auth.Service
	config      *config.Config
}

// NewAuthController creates a new auth controller
func NewAuthController(authService *auth.Service, cfg *config.Config) *AuthController {
	return &AuthController{
		authService: authService,
		config:      cfg,
	}
}

// Login handles user login
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.ErrorResponse(w, errors.NewBadRequest("Invalid request body"), c.config)
		return
	}

	// Validate request
	if err := validation.ValidateStruct(&req); err != nil {
		middleware.ErrorResponse(w, err, c.config)
		return
	}

	// Login
	response, err := c.authService.Login(&req)
	if err != nil {
		middleware.ErrorResponse(w, err, c.config)
		return
	}

	// Set cookie (Secure flag based on HTTPS config)
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    response.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   c.config.GetCookieSecure(), // Based on HTTPS config
		SameSite: http.SameSiteStrictMode,
	})

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Register handles user registration
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req auth.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.ErrorResponse(w, errors.NewBadRequest("Invalid request body"), c.config)
		return
	}

	// Validate request
	if err := validation.ValidateStruct(&req); err != nil {
		middleware.ErrorResponse(w, err, c.config)
		return
	}

	// Register
	user, err := c.authService.Register(&req)
	if err != nil {
		middleware.ErrorResponse(w, err, c.config)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Logout handles user logout
func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   c.config.GetCookieSecure(),
		MaxAge:   -1,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}
