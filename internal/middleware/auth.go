package middleware

import (
    "net/http"
    "strings"

    "github.com/Kai1313/url-shortener-fullstack/backend/internal/config"
    "github.com/Kai1313/url-shortener-fullstack/backend/internal/utils"
    "github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        // Bearer token format
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
            c.Abort()
            return
        }

        token := parts[1]

        // Validate token
        userID, err := utils.ValidateToken(token, cfg.JWTSecret)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
            c.Abort()
            return
        }

        // Set user ID in context
        c.Set("userID", userID)
        c.Next()
    }
}