package models

import (
    "time"

    "gorm.io/gorm"
)

type User struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Email     string         `gorm:"uniqueIndex;not null" json:"email"`
    Password  string         `gorm:"not null" json:"-"`
    Name      string         `json:"name"`
    URLs      []URL          `gorm:"foreignKey:UserID" json:"urls,omitempty"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type RegisterRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    Name     string `json:"name" binding:"required"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
    Token string `json:"token"`
    User  User   `json:"user"`
}