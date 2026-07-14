package main

import (
    "log"
    "os"

    "github.com/Kai1313/url-shortener-fullstack/backend/internal/config"
    "github.com/Kai1313/url-shortener-fullstack/backend/internal/database"
    "github.com/Kai1313/url-shortener-fullstack/backend/internal/handlers"
    "github.com/Kai1313/url-shortener-fullstack/backend/internal/middleware"
    "github.com/Kai1313/url-shortener-fullstack/backend/internal/models"
    "github.com/Kai1313/url-shortener-fullstack/backend/internal/repository"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func main() {
    // Load configuration
    cfg := config.LoadConfig()

    // Initialize database
    db, err := database.InitDB(cfg)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Auto migrate
    if err := db.AutoMigrate(&models.User{}, &models.URL{}, &models.Analytics{}); err != nil {
        log.Fatal("Failed to migrate database:", err)
    }

    // Initialize Redis
    redisClient := database.InitRedis(cfg)

    // Initialize repositories
    urlRepo := repository.NewURLRepository(db, redisClient)
    userRepo := repository.NewUserRepository(db)

    // Initialize handlers
    authHandler := handlers.NewAuthHandler(userRepo, cfg)
    urlHandler := handlers.NewURLHandler(urlRepo, cfg)
    analyticsHandler := handlers.NewAnalyticsHandler(urlRepo)

    // Setup router
    router := gin.Default()

    // CORS middleware
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

    // Public routes
    router.GET("/:code", urlHandler.Redirect)
    router.GET("/api/urls/:code/stats", analyticsHandler.GetStats)

    // Auth routes
    auth := router.Group("/api/auth")
    {
        auth.POST("/register", authHandler.Register)
        auth.POST("/login", authHandler.Login)
        auth.POST("/logout", authHandler.Logout)
    }

    // Protected routes
    api := router.Group("/api")
    api.Use(middleware.AuthMiddleware(cfg))
    {
        api.POST("/urls", urlHandler.Create)
        api.GET("/urls", urlHandler.List)
        api.GET("/urls/:code", urlHandler.GetByCode)
        api.DELETE("/urls/:code", urlHandler.Delete)
        api.PUT("/urls/:code", urlHandler.Update)
    }

    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server starting on port %s", port)
    if err := router.Run(":" + port); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}