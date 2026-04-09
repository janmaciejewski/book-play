package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/janmaciejewski/book-play/apps/api/internal/config"
	"github.com/janmaciejewski/book-play/apps/api/internal/middleware"
	"github.com/janmaciejewski/book-play/apps/api/internal/models"
	"github.com/janmaciejewski/book-play/apps/api/internal/modules/auth"
	"github.com/janmaciejewski/book-play/apps/api/internal/modules/facility"
	"github.com/janmaciejewski/book-play/apps/api/internal/modules/reservation"
	"github.com/janmaciejewski/book-play/apps/api/internal/modules/team"
	"github.com/janmaciejewski/book-play/apps/api/internal/modules/user"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := config.InitDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer config.CloseDatabase()

	// Auto-migrate models
	if err := db.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.Facility{},
		&models.FacilitySlot{},
		&models.Reservation{},
		&models.Team{},
		&models.TeamMember{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed")

	// Seed database with initial data
	if err := config.SeedDatabase(db); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	// Initialize Redis (optional in development)
	_, err = config.InitRedis(&cfg.Redis)
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
	}
	defer config.CloseRedis()

	// Initialize services
	authService := auth.NewService(db)
	authHandler := auth.NewHandler(authService)

	facilityService := facility.NewService(db)
	facilityHandler := facility.NewHandler(facilityService)

	teamService := team.NewService(db)
	teamHandler := team.NewHandler(teamService)

	reservationService := reservation.NewService(db)
	reservationHandler := reservation.NewHandler(reservationService)

	userService := user.NewService(db)
	userHandler := user.NewHandler(userService)

	// Setup Gin router
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Apply middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	// Root route - API info
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":    "Book-Play API",
			"version": "1.0.0",
			"status":  "running",
			"endpoints": gin.H{
				"health": "/health",
				"api":    "/api/v1",
			},
		})
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"version": "1.0.0",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", authHandler.RefreshToken)
			authGroup.POST("/logout", authHandler.Logout)
		}

		// Public facility routes
		v1.GET("/facilities", facilityHandler.GetAll)
		v1.GET("/facilities/:id", facilityHandler.GetByID)
		v1.GET("/facilities/:id/availability", facilityHandler.GetAvailability)

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.JWTAuth())
		{
			protected.GET("/auth/me", authHandler.GetMe)

			// Facility routes (protected - create only)
			protected.POST("/facilities", facilityHandler.Create)

			// Reservation routes
			protected.GET("/reservations", reservationHandler.GetAll)
			protected.GET("/reservations/:id", reservationHandler.GetByID)
			protected.POST("/reservations", reservationHandler.Create)
			protected.PUT("/reservations/:id/cancel", reservationHandler.Cancel)

			// Team routes
			protected.GET("/teams", teamHandler.GetAll)
			protected.GET("/teams/:id", teamHandler.GetByID)
			protected.POST("/teams", teamHandler.Create)

			// User routes (admin only)
			protected.GET("/users", userHandler.GetAll)
			protected.PUT("/users/:id/role", userHandler.UpdateRole)
			protected.DELETE("/users/:id", userHandler.Delete)
		}
	}

	// Start server
	port := cfg.App.Port
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s...\n", port)
	if err := router.Run(":" + port); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start server: %v\n", err)
		os.Exit(1)
	}
}
