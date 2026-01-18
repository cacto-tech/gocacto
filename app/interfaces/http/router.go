package http

import (
	"net/http"

	"cacto-cms/app/interfaces/http/controller"
	"cacto-cms/app/interfaces/http/middleware"
	"cacto-cms/app/shared/auth"
	"cacto-cms/config"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
)

// Router wraps the HTTP router
type Router struct {
	*chi.Mux
}

// NewRouter creates a new router with all routes and middleware
func NewRouter(
	pageController *controller.PageController,
	authController *controller.AuthController,
	adminController *controller.AdminController,
	jwtManager *auth.JWTManager,
	cfg *config.Config,
) *Router {
	r := chi.NewRouter()

	// Global middleware (order matters!)
	r.Use(chimw.StripSlashes)
	r.Use(middleware.PathTraversalPrevention)
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.CORSRestricted(cfg.AllowedOrigins))
	r.Use(middleware.Logger)
	r.Use(middleware.Recovery)
	r.Use(middleware.ErrorHandler(cfg))

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./web/uploads"))))

	// Public routes
	r.Get("/", pageController.ShowHome)
	r.Get("/{slug}", pageController.ShowPage)

	// Auth routes (API) - with rate limiting
	r.Group(func(r chi.Router) {
		r.Use(middleware.RateLimitAuth())
		r.Post("/api/auth/login", authController.Login)
		r.Post("/api/auth/register", authController.Register)
		r.Post("/api/auth/logout", authController.Logout)
	})

	// Admin login (public) - with rate limiting
	r.Group(func(r chi.Router) {
		r.Use(middleware.RateLimitAuth())
		r.Get("/admin/login", adminController.ShowLogin)
		r.Post("/admin/login", adminController.HandleLogin)
	})

	// Protected routes (require authentication)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(jwtManager))
		r.Use(middleware.RequireAuth)

		// Admin routes (require admin/editor role)
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole("admin", "editor"))

			r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "/admin/dashboard", http.StatusFound)
			})

			r.Get("/admin/dashboard", adminController.ShowDashboard)
			r.Get("/admin/logout", adminController.HandleLogout)
			r.Post("/admin/logout", adminController.HandleLogout)
		})
	})

	// Sitemap
	r.Get("/sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/static/sitemap.xml")
	})

	return &Router{r}
}
