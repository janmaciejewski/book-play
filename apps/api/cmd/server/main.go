package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/janmaciejewski/book-play/apps/api/internal/config"
	"github.com/janmaciejewski/book-play/apps/api/internal/middleware"
	"github.com/janmaciejewski/book-play/apps/api/internal/models"
	"github.com/janmaciejewski/book-play/apps/api/internal/modules/auth"
	"github.com/janmaciejewski/book-play/apps/api/internal/modules/chat"
	"github.com/janmaciejewski/book-play/apps/api/internal/modules/facility"
	"github.com/janmaciejewski/book-play/apps/api/internal/modules/mail"
	"github.com/janmaciejewski/book-play/apps/api/internal/modules/reservation"
	"github.com/janmaciejewski/book-play/apps/api/internal/modules/team"
	"github.com/janmaciejewski/book-play/apps/api/internal/modules/user"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := config.InitDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer config.CloseDatabase()

	if err := db.AutoMigrate(
		&models.User{}, &models.RefreshToken{}, &models.Facility{}, &models.FacilitySlot{},
		&models.Reservation{}, &models.Team{}, &models.TeamMember{},
		&models.TeamRecruitmentApplication{}, &models.EmailVerificationOTP{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed")

	if err := config.SeedDatabase(db); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	_, err = config.InitRedis(&cfg.Redis)
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
	}
	defer config.CloseRedis()

	mailService := mail.NewService(db)
	authService := auth.NewService(db)
	authHandler := auth.NewHandler(authService, mailService)

	facilityService := facility.NewService(db)
	facilityHandler := facility.NewHandler(facilityService)

	teamService := team.NewService(db)
	teamHandler := team.NewHandler(teamService)

	reservationService := reservation.NewService(db)
	reservationHandler := reservation.NewHandler(reservationService)

	userService := user.NewService(db)
	userHandler := user.NewHandler(userService)

	var chatService *chat.Service
	var chatHandler *chat.Handler
	if config.RedisClient != nil {
		chatService = chat.NewService(config.RedisClient)
		chatHandler = chat.NewHandler(chatService)
	}

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"name": "Book-Play API", "version": "1.0.0", "status": "running"})
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := router.Group("/api/v1")
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", authHandler.RefreshToken)
			authGroup.POST("/logout", authHandler.Logout)
			authGroup.POST("/send-otp", authHandler.SendOTP)
			authGroup.POST("/verify-otp", authHandler.VerifyOTP)
			authGroup.POST("/reset-password", authHandler.ResetPassword)
		}

		v1.GET("/facilities", facilityHandler.GetAll)
		v1.GET("/facilities/mine", middleware.JWTAuth(), func(c *gin.Context) {
			facilityHandler.GetMine(c)
		})
		v1.GET("/facilities/:id", facilityHandler.GetByID)
		v1.GET("/facilities/:id/availability", facilityHandler.GetAvailability)

		protected := v1.Group("")
		protected.Use(middleware.JWTAuth())
		{
			protected.GET("/auth/me", authHandler.GetMe)
			protected.POST("/facilities", facilityHandler.Create)
			protected.PUT("/facilities/:id", facilityHandler.UpdateFacility)
			protected.PUT("/facilities/:id/slots", facilityHandler.UpdateSlots)
			protected.PUT("/facilities/:id/close", facilityHandler.ToggleClose)

			protected.GET("/reservations", reservationHandler.GetAll)
			protected.GET("/reservations/:id", reservationHandler.GetByID)
			protected.POST("/reservations", reservationHandler.Create)
			protected.PUT("/reservations/:id/cancel", reservationHandler.Cancel)
			protected.GET("/facilities/my/reservations", reservationHandler.GetForFacilityOwner)
			protected.PUT("/reservations/:id/status", reservationHandler.UpdateStatus)

			protected.GET("/teams", teamHandler.GetAll)
			protected.GET("/my-teams", teamHandler.GetMyTeams)
			protected.GET("/teams/:id", teamHandler.GetByID)
			protected.POST("/teams", teamHandler.Create)
			protected.PUT("/teams/:id", teamHandler.Update)
			protected.POST("/teams/:id/logo", teamHandler.UploadLogo)
			protected.POST("/teams/:id/members", teamHandler.AddMember)
			protected.DELETE("/teams/:id/members/:memberId", teamHandler.RemoveMember)
			protected.PUT("/teams/:id/members/:memberId/role", teamHandler.UpdateMemberRole)
			protected.GET("/teams/:id/search-users", teamHandler.SearchUsers)
			protected.PUT("/teams/:id/recruitment", teamHandler.ToggleRecruitment)
			protected.POST("/teams/:id/apply", teamHandler.ApplyRecruitment)
			protected.GET("/teams/:id/applications", teamHandler.GetApplications)
			protected.PUT("/teams/:id/applications/:appId", teamHandler.HandleApplication)
			protected.DELETE("/teams/:id", teamHandler.DeleteTeam)

			protected.GET("/users", userHandler.GetAll)
			protected.GET("/users/:id", userHandler.GetProfile)
			protected.PUT("/users/:id", userHandler.UpdateProfile)
			protected.PUT("/users/:id/role", userHandler.UpdateRole)
			protected.POST("/users/:id/avatar", userHandler.UploadAvatar)
			protected.DELETE("/users/:id", userHandler.Delete)

			// Endpointy czatu – tylko jeśli Redis jest dostępny
			if chatHandler != nil {
				protected.GET("/teams/:id/chat", chatHandler.GetMessages)
				protected.POST("/teams/:id/chat", chatHandler.SendMessage)
			}
		}
	}

	router.Static("/uploads", "./uploads")

	port := cfg.App.Port
	if port == "" {
		port = "8080"
	}
	// Uruchamia scheduler przypomnień w tle
	go func() {
		log.Println("📧 Reservation reminder scheduler started (checking every hour)")
		// Uruchamia od razu przy starcie dla testów
		log.Println("Running initial reminder check...")
		if err := mailService.ProcessReminders(); err != nil {
			log.Printf("Warning: Initial reminder check failed: %v", err)
		}
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			log.Println("Running scheduled reminder check...")
			if err := mailService.ProcessReminders(); err != nil {
				log.Printf("Warning: Reminder check failed: %v", err)
			}
		}
	}()

	fmt.Printf("Server starting on port %s...\n", port)
	if err := router.Run(":" + port); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start server: %v\n", err)
		os.Exit(1)
	}
}
