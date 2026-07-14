package handlers

import (
    "net/http"

    "github.com/Kai1313/url-shortener-fullstack/backend/internal/config"
    "github.com/Kai1313/url-shortener-fullstack/backend/internal/models"
    "github.com/Kai1313/url-shortener-fullstack/backend/internal/repository"
    "github.com/Kai1313/url-shortener-fullstack/backend/internal/utils"
    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    userRepo repository.UserRepository
    config   *config.Config
}

func NewAuthHandler(userRepo repository.UserRepository, cfg *config.Config) *AuthHandler {
    return &AuthHandler{
        userRepo: userRepo,
        config:   cfg,
    }
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req models.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Check if user exists
    existing, _ := h.userRepo.FindByEmail(req.Email)
    if existing != nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
        return
   }

    // Hash password
    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    // Create user
    user := &models.User{
        Email:    req.Email,
        Password: hashedPassword,
        Name:     req.Name,
    }

    if err := h.userRepo.Create(user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    // Generate token
    token, err := utils.GenerateToken(user.ID, h.config.JWTSecret, h.config.JWTExpiryHours)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusCreated, models.AuthResponse{
        Token: token,
        User:  *user,
    })
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req models.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Find user
    user, err := h.userRepo.FindByEmail(req.Email)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // Check password
    if !utils.CheckPassword(req.Password, user.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // Generate token
    token, err := utils.GenerateToken(user.ID, h.config.JWTSecret, h.config.JWTExpiryHours)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, models.AuthResponse{
        Token: token,
        User:  *user,
    })
}

func (h *AuthHandler) Logout(c *gin.Context) {
    // JWT is stateless, but we could blacklist the token here if needed
    c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}