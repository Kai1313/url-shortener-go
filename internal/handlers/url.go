package handlers

import (
    "net/http"
    "time"

    "github.com/Kai1313/url-shortener-fullstack/backend/internal/config"
    "github.com/Kai1313/url-shortener-fullstack/backend/internal/models"
    "github.com/Kai1313/url-shortener-fullstack/backend/internal/repository"
    "github.com/Kai1313/url-shortener-fullstack/backend/internal/utils"
    "github.com/gin-gonic/gin"
)

type URLHandler struct {
    urlRepo repository.URLRepository
    config  *config.Config
}

func NewURLHandler(urlRepo repository.URLRepository, cfg *config.Config) *URLHandler {
    return &URLHandler{
        urlRepo: urlRepo,
        config:  cfg,
    }
}

func (h *URLHandler) Create(c *gin.Context) {
    var req models.CreateURLRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Generate short code
    shortCode := req.CustomCode
    if shortCode == "" {
        shortCode = utils.GenerateShortCode(6)
    }

    // Check if code exists
    existing, _ := h.urlRepo.FindByCode(shortCode)
    if existing != nil {
        // If custom code was provided, return error
        if req.CustomCode != "" {
            c.JSON(http.StatusConflict, gin.H{"error": "Custom code already taken"})
            return
        }
        // Otherwise generate a new one
        shortCode = utils.GenerateShortCode(8)
    }

    // Parse expiry
    var expiresAt *time.Time
    if req.ExpiresAt != nil && *req.ExpiresAt != "" {
        t, err := time.Parse(time.RFC3339, *req.ExpiresAt)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expiry date format"})
            return
        }
        expiresAt = &t
    }

    // Get user ID from context (set by auth middleware)
    userID, exists := c.Get("userID")
    var userIDPtr *uint
    if exists {
        id := userID.(uint)
        userIDPtr = &id
    }

    // Create URL
    url := &models.URL{
        OriginalURL: req.OriginalURL,
        ShortCode:   shortCode,
        UserID:      userIDPtr,
        ExpiresAt:   expiresAt,
        IsActive:    true,
    }

    if err := h.urlRepo.Create(url); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create URL"})
        return
    }

    // Cache the URL
    _ = h.urlRepo.CacheURL(shortCode, req.OriginalURL)

    response := models.URLResponse{
        ID:          url.ID,
        OriginalURL: url.OriginalURL,
        ShortURL:    h.config.BaseURL + "/" + url.ShortCode,
        ShortCode:   url.ShortCode,
        Clicks:      url.Clicks,
        ExpiresAt:   url.ExpiresAt,
        IsActive:    url.IsActive,
        CreatedAt:   url.CreatedAt,
    }

    c.JSON(http.StatusCreated, response)
}

func (h *URLHandler) Redirect(c *gin.Context) {
    code := c.Param("code")

    // Check cache first
    cachedURL, err := h.urlRepo.GetCachedURL(code)
    if err == nil && cachedURL != "" {
        // Increment clicks async
        go h.urlRepo.IncrementClicks(code)
        c.Redirect(http.StatusMovedPermanently, cachedURL)
        return
    }

    // Find in database
    url, err := h.urlRepo.FindByCode(code)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
        return
    }

    // Check if expired
    if url.ExpiresAt != nil && url.ExpiresAt.Before(time.Now()) {
        c.JSON(http.StatusGone, gin.H{"error": "URL has expired"})
        return
    }

    // Check if active
    if !url.IsActive {
        c.JSON(http.StatusGone, gin.H{"error": "URL has been deactivated"})
        return
    }

    // Cache for future requests
    _ = h.urlRepo.CacheURL(code, url.OriginalURL)

    // Increment clicks
    go h.urlRepo.IncrementClicks(code)

    c.Redirect(http.StatusMovedPermanently, url.OriginalURL)
}

func (h *URLHandler) List(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    urls, err := h.urlRepo.FindByUserID(userID.(uint))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch URLs"})
        return
    }

    var response []models.URLResponse
    for _, url := range urls {
        response = append(response, models.URLResponse{
            ID:          url.ID,
            OriginalURL: url.OriginalURL,
            ShortURL:    h.config.BaseURL + "/" + url.ShortCode,
            ShortCode:   url.ShortCode,
            Clicks:      url.Clicks,
            ExpiresAt:   url.ExpiresAt,
            IsActive:    url.IsActive,
            CreatedAt:   url.CreatedAt,
        })
    }

    c.JSON(http.StatusOK, response)
}

func (h *URLHandler) GetByCode(c *gin.Context) {
    code := c.Param("code")
    userID, _ := c.Get("userID")

    url, err := h.urlRepo.FindByCode(code)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
        return
    }

    // Check ownership
    if url.UserID == nil || *url.UserID != userID.(uint) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
        return
    }

    response := models.URLResponse{
        ID:          url.ID,
        OriginalURL: url.OriginalURL,
        ShortURL:    h.config.BaseURL + "/" + url.ShortCode,
        ShortCode:   url.ShortCode,
        Clicks:      url.Clicks,
        ExpiresAt:   url.ExpiresAt,
        IsActive:    url.IsActive,
        CreatedAt:   url.CreatedAt,
    }

    c.JSON(http.StatusOK, response)
}

func (h *URLHandler) Update(c *gin.Context) {
    code := c.Param("code")
    userID, _ := c.Get("userID")

    var req models.UpdateURLRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    url, err := h.urlRepo.FindByCode(code)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
        return
    }

    // Check ownership
    if url.UserID == nil || *url.UserID != userID.(uint) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
        return
    }

    // Update fields
    if req.OriginalURL != "" {
        url.OriginalURL = req.OriginalURL
        // Update cache
        _ = h.urlRepo.CacheURL(code, req.OriginalURL)
    }

    if req.IsActive != nil {
        url.IsActive = *req.IsActive
    }

    if req.ExpiresAt != nil && *req.ExpiresAt != "" {
        t, err := time.Parse(time.RFC3339, *req.ExpiresAt)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expiry date format"})
            return
        }
        url.ExpiresAt = &t
    }

    if err := h.urlRepo.Update(url); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update URL"})
        return
    }

    response := models.URLResponse{
        ID:          url.ID,
        OriginalURL: url.OriginalURL,
        ShortURL:    h.config.BaseURL + "/" + url.ShortCode,
        ShortCode:   url.ShortCode,
        Clicks:      url.Clicks,
        ExpiresAt:   url.ExpiresAt,
        IsActive:    url.IsActive,
        CreatedAt:   url.CreatedAt,
    }

    c.JSON(http.StatusOK, response)
}

func (h *URLHandler) Delete(c *gin.Context) {
    code := c.Param("code")
    userID, _ := c.Get("userID")

    url, err := h.urlRepo.FindByCode(code)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
        return
    }

    // Check ownership
    if url.UserID == nil || *url.UserID != userID.(uint) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
        return
    }

    if err := h.urlRepo.Delete(url.ID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete URL"})
        return
    }

    // Delete cache
    _ = h.urlRepo.DeleteCache(code)

    c.JSON(http.StatusOK, gin.H{"message": "URL deleted successfully"})
}