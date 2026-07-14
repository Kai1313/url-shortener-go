package models

import (
    "time"

    "gorm.io/gorm"
)

type URL struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    OriginalURL string         `gorm:"not null" json:"original_url"`
    ShortCode   string         `gorm:"uniqueIndex;not null" json:"short_code"`
    UserID      *uint          `json:"user_id"`
    Clicks      int            `gorm:"default:0" json:"clicks"`
    ExpiresAt   *time.Time     `json:"expires_at"`
    IsActive    bool           `gorm:"default:true" json:"is_active"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
    User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Analytics   []Analytics    `gorm:"foreignKey:URLID" json:"analytics,omitempty"`
}

type CreateURLRequest struct {
    OriginalURL string  `json:"original_url" binding:"required,url"`
    CustomCode  string  `json:"custom_code" binding:"omitempty,min=3,max=20,alphanum"`
    ExpiresAt   *string `json:"expires_at" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

type UpdateURLRequest struct {
    OriginalURL string  `json:"original_url" binding:"omitempty,url"`
    IsActive    *bool   `json:"is_active"`
    ExpiresAt   *string `json:"expires_at" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

type URLResponse struct {
    ID          uint       `json:"id"`
    OriginalURL string     `json:"original_url"`
    ShortURL    string     `json:"short_url"`
    ShortCode   string     `json:"short_code"`
    Clicks      int        `json:"clicks"`
    ExpiresAt   *time.Time `json:"expires_at"`
    IsActive    bool       `json:"is_active"`
    CreatedAt   time.Time  `json:"created_at"`
}