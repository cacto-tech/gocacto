package controller

import (
	"encoding/json"
	"net/http"

	"cacto-cms/app/application/auth"
	"cacto-cms/app/interfaces/http/middleware"
	"cacto-cms/app/interfaces/templates/admin"
	"cacto-cms/app/shared/errors"
	"cacto-cms/app/shared/validation"
	"cacto-cms/config"
)

// AdminController handles admin-related requests
type AdminController struct {
	authService *auth.Service
	baseURL     string
	config      *config.Config
}

// NewAdminController creates a new admin controller
func NewAdminController(authService *auth.Service, baseURL string, cfg *config.Config) *AdminController {
	return &AdminController{
		authService: authService,
		baseURL:     baseURL,
		config:      cfg,
	}
}

// ShowLogin displays the admin login page
func (c *AdminController) ShowLogin(w http.ResponseWriter, r *http.Request) {
	// Check if already authenticated
	if userID, ok := middleware.GetUserID(r.Context()); ok && userID > 0 {
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	admin.Login().Render(r.Context(), w)
}

// HandleLogin handles admin login (both form and JSON)
func (c *AdminController) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req auth.LoginRequest

	// Check Content-Type to determine if it's JSON or form data
	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		// JSON request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			if isAPIRequest(r) {
				middleware.ErrorResponse(w, errors.NewBadRequest("Invalid request body"), c.config)
			} else {
				http.Error(w, "Invalid request", http.StatusBadRequest)
			}
			return
		}
	} else {
		// Form data
		req.Email = r.FormValue("email")
		req.Password = r.FormValue("password")
	}

	// Validate request
	if err := validation.ValidateStruct(&req); err != nil {
		if isAPIRequest(r) {
			middleware.ErrorResponse(w, err, c.config)
		} else {
			http.Error(w, "Validation failed", http.StatusBadRequest)
		}
		return
	}

	// Login
	response, err := c.authService.Login(&req)
	if err != nil {
		if isAPIRequest(r) {
			middleware.ErrorResponse(w, err, c.config)
		} else {
			http.Error(w, "Login failed", http.StatusUnauthorized)
		}
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

	// Check if API request
	if isAPIRequest(r) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		// Redirect to dashboard
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
	}
}

// ShowDashboard displays the admin dashboard
func (c *AdminController) ShowDashboard(w http.ResponseWriter, r *http.Request) {
	userEmail, _ := middleware.GetUserEmail(r.Context())
	userRole, _ := middleware.GetUserRole(r.Context())

	if isAPIRequest(r) {
		// Return JSON for API requests
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"email": userEmail,
			"role":  userRole,
		})
		return
	}

	// Return HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	admin.Dashboard(userEmail, userRole).Render(r.Context(), w)
}

// HandleLogout handles admin logout
func (c *AdminController) HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   c.config.GetCookieSecure(),
		MaxAge:   -1,
	})

	if isAPIRequest(r) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
	} else {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
	}
}

// isAPIRequest checks if request is an API request (JSON preferred)
func isAPIRequest(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	return accept == "application/json" || r.Header.Get("Content-Type") == "application/json"
}
