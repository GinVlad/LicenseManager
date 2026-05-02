// @title           LicenseManager API
// @version         1.0
// @description     Self-hosted license management server. Admin routes require Bearer JWT; client routes require X-App-Key header.
// @host            localhost:8080
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in     header
// @name   Authorization
// @description Type "Bearer" followed by a space and the JWT token from /admin/login.

// @securityDefinitions.apikey AppApiKey
// @in     header
// @name   X-App-Key
// @description API key issued per application (shown on app creation).

package main

import (
	"context"
	"log"
	"net/http"
	"time"

	_ "github.com/cheo/licensemanager/docs"
	"github.com/cheo/licensemanager/internal/config"
	"github.com/cheo/licensemanager/internal/database"
	"github.com/cheo/licensemanager/internal/handlers"
	"github.com/cheo/licensemanager/internal/middleware"
	"github.com/cheo/licensemanager/internal/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	// Seed admin account if password is set and no admins exist
	if cfg.AdminPass != "" {
		var count int
		db.QueryRow(context.Background(), "SELECT COUNT(*) FROM admins").Scan(&count)
		if count == 0 {
			hash, _ := bcrypt.GenerateFromPassword([]byte(cfg.AdminPass), 12)
			db.Exec(context.Background(),
				"INSERT INTO admins (email, password_hash) VALUES ($1, $2) ON CONFLICT DO NOTHING",
				cfg.AdminEmail, string(hash))
			log.Printf("Admin account created: %s", cfg.AdminEmail)
		}
	}

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(
		gin.Recovery(),
		gin.Logger(),
		middleware.SecureHeaders(),
		middleware.CORS(cfg.CORSOrigins),
		middleware.RateLimit(120, time.Minute),
	)

	licenseSvc := services.NewLicenseService(db)
	licenseH := handlers.NewLicenseHandler(licenseSvc)
	adminH := handlers.NewAdminHandler(db, licenseSvc, cfg.JWTSecret)

	// Public license validation (called by client apps)
	public := r.Group("/api/v1/license")
	public.Use(middleware.AppAPIKey(db))
	{
		public.POST("/validate", licenseH.Validate)
		public.DELETE("/hwid", licenseH.RemoveHWID)
	}

	// Admin routes (JWT protected)
	admin := r.Group("/api/v1/admin")
	{
		admin.POST("/login", adminH.Login)

		protected := admin.Group("")
		protected.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			protected.GET("/apps", adminH.ListApps)
			protected.POST("/apps", adminH.CreateApp)
			protected.PUT("/apps/:id", adminH.UpdateApp)
			protected.DELETE("/apps/:id", adminH.DeleteApp)
			protected.GET("/licenses", adminH.ListLicenses)
			protected.POST("/licenses", adminH.CreateLicense)
			protected.PUT("/licenses/:id", adminH.UpdateLicense)
			protected.DELETE("/licenses/:id", adminH.DeleteLicense)
			protected.GET("/licenses/:id/hwids", adminH.ListHWIDs)
			protected.POST("/licenses/:id/hwids", adminH.AddHWID)
			protected.DELETE("/hwids/:id", adminH.DeleteHWID)
		}
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now().UTC()})
	})

	// Swagger UI — http://localhost:8080/swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Printf("LicenseManager starting on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server: %v", err)
	}
}
