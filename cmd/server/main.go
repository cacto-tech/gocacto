package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	authservice "cacto-cms/app/application/auth"
	"cacto-cms/app/application/component"
	"cacto-cms/app/application/page"
	userservice "cacto-cms/app/application/user"
	"cacto-cms/app/infrastructure/database"
	componentpersistence "cacto-cms/app/infrastructure/persistence/component"
	pagepersistence "cacto-cms/app/infrastructure/persistence/page"
	userpersistence "cacto-cms/app/infrastructure/persistence/user"
	httphandlers "cacto-cms/app/interfaces/http"
	"cacto-cms/app/interfaces/http/controller"
	"cacto-cms/app/shared/auth"
	"cacto-cms/app/shared/seo"
	"cacto-cms/app/shared/sitemap"
	"cacto-cms/config"
)

func main() {
	// Load config
	cfg := config.Load()

	// Initialize database
	db, err := database.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Seeders are now run via artisan CLI: go run ./cmd/artisan migrate:fresh --seed

	// Initialize repositories
	pageRepo := pagepersistence.NewRepository(db.DB)
	componentRepo := componentpersistence.NewRepository(db.DB)
	userRepo := userpersistence.NewRepository(db.DB)

	// Initialize services
	pageService := page.NewService(pageRepo)
	componentService := component.NewService(componentRepo)
	userService := userservice.NewService(userRepo)

	// Initialize auth
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiration)
	authService := authservice.NewService(userService, cfg.JWTSecret, cfg.JWTExpiration)

	// Initialize SEO manager
	seoManager := seo.NewManager(cfg.BaseURL, cfg.SiteName, cfg.SiteDescription)

	// Initialize sitemap generator
	sitemapPath := "./web/static/sitemap.xml"
	sitemapGen := sitemap.NewGenerator(cfg.BaseURL, sitemapPath, pageRepo)
	sitemapGen.ScheduleDaily()
	log.Println("📍 Sitemap generator scheduled")

	// Initialize controllers
	pageController := controller.NewPageController(
		cfg.BaseURL,
		pageService,
		componentService,
		seoManager,
	)
	
	authController := controller.NewAuthController(authService, cfg)
	adminController := controller.NewAdminController(authService, cfg.BaseURL, cfg)

	// Setup router
	router := httphandlers.NewRouter(pageController, authController, adminController, jwtManager, cfg)

	// Start server
	addr := ":" + cfg.ServerPort
	log.Printf("🚀 Server starting on %s", addr)
	log.Printf("📂 Database: %s", cfg.DBPath)
	log.Printf("🌐 Visit: http://localhost%s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func init() {
	// Create necessary directories
	dirs := []string{
		"./web/static",
		"./web/uploads",
		"./web/static/css",
		"./web/static/js",
		"./web/static/images",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Warning: Could not create directory %s: %v", dir, err)
		}
	}

	fmt.Println(`
	╔═════════════════════════════════════════╗
	║             Cacto CMS v1.0.0            ║
	║   Performance-Oriented Enterprise CMS   ║
	╚═════════════════════════════════════════╝
	`)
}
