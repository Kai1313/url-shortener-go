package handlers

import (
    "net/http"
    "strings"

    "github.com/Kai1313/url-shortener-fullstack/backend/internal/models"
    "github.com/Kai1313/url-shortener-fullstack/backend/internal/repository"
    "github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
    urlRepo repository.URLRepository
}

func NewAnalyticsHandler(urlRepo repository.URLRepository) *AnalyticsHandler {
    return &AnalyticsHandler{
        urlRepo: urlRepo,
    }
}

func (h *AnalyticsHandler) GetStats(c *gin.Context) {
    code := c.Param("code")

    url, err := h.urlRepo.FindByCode(code)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
        return
    }

    // Get analytics data
    analytics, err := h.urlRepo.GetAnalytics(url.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch analytics"})
        return
    }

    todayClicks, _ := h.urlRepo.GetTodayClicks(url.ID)
    weeklyClicks, _ := h.urlRepo.GetWeeklyClicks(url.ID)

    // Parse analytics for device, OS, browser breakdown
    devices := make(map[string]int)
    os := make(map[string]int)
    browsers := make(map[string]int)
    referers := make(map[string]int)

    for _, a := range analytics {
        // Device
        if a.Device != "" {
            devices[a.Device]++
        } else {
            devices["Unknown"]++
        }

        // OS
        if a.OS != "" {
            os[a.OS]++
        }

        // Browser
        if a.Browser != "" {
            browsers[a.Browser]++
        }

        // Referer
        if a.Referer != "" {
            host := extractHost(a.Referer)
            if host != "" {
                referers[host]++
            }
        }
    }

    response := models.StatsResponse{
        TotalClicks:  url.Clicks,
        TodayClicks:  todayClicks,
        WeeklyClicks: weeklyClicks,
        Devices:      devices,
        OS:           os,
        Browsers:     browsers,
        Referers:     referers,
    }

    c.JSON(http.StatusOK, response)
}

func extractHost(referer string) string {
    // Remove protocol
    if strings.HasPrefix(referer, "http://") {
        referer = strings.TrimPrefix(referer, "http://")
    } else if strings.HasPrefix(referer, "https://") {
        referer = strings.TrimPrefix(referer, "https://")
    }

    // Get host part
    parts := strings.Split(referer, "/")
    if len(parts) > 0 {
        return parts[0]
    }
    return referer
}