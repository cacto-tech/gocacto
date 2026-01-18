package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	"cacto-cms/app/shared/errors"
	"cacto-cms/config"
)

// ErrorHandler handles errors and returns appropriate HTTP responses
func ErrorHandler(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a custom response writer to capture status code
			rr := &responseRecorder{
				ResponseWriter: w,
				statusCode:    http.StatusOK,
			}

			next.ServeHTTP(rr, r)

			// If status code indicates an error, handle it
			if rr.statusCode >= 400 {
				handleError(w, rr.statusCode, nil, cfg)
			}
		})
	}
}

// responseRecorder wraps http.ResponseWriter to capture status code
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

// handleError writes error response
func handleError(w http.ResponseWriter, statusCode int, err error, cfg *config.Config) {
	var appErr *errors.AppError

	if err != nil {
		appErr = errors.AsAppError(err)
	} else {
		// Create default error based on status code
		switch statusCode {
		case http.StatusNotFound:
			appErr = errors.ErrNotFound
		case http.StatusUnauthorized:
			appErr = errors.ErrUnauthorized
		case http.StatusForbidden:
			appErr = errors.ErrForbidden
		case http.StatusInternalServerError:
			appErr = errors.ErrInternal
		default:
			appErr = errors.New(errors.ErrCodeInternal, "An error occurred", statusCode)
		}
	}

	// Log error (always log, but don't expose details in production)
	if statusCode >= 500 {
		log.Printf("Error: %v", appErr)
	}

	// In production, don't expose internal error details
	errorMessage := appErr.Message
	if cfg.IsProduction() && statusCode >= 500 {
		errorMessage = "An internal error occurred"
	}

	// Write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.HTTPStatus)

	response := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    appErr.Code,
			"message": errorMessage,
		},
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode error response: %v", err)
	}
}

// ErrorResponse writes an error response with config
func ErrorResponse(w http.ResponseWriter, err error, cfg *config.Config) {
	appErr := errors.AsAppError(err)
	handleError(w, appErr.HTTPStatus, err, cfg)
}
