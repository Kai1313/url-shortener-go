package models

import (
    "time"
)

type Analytics struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    URLID     uint      `gorm:"not null;index" json:"url_id"`
    IPAddress string    `json:"ip_address"`
    UserAgent string    `json:"user_agent"`
    Device    string    `json:"device"`
    OS        string    `json:"os"`
    Browser   string    `json:"browser"`
    Referer   string    `json:"referer"`
    Country   string    `json:"country"`
    City      string    `json:"city"`
    CreatedAt time.Time `json:"created_at"`
    URL       URL       `gorm:"foreignKey:URLID" json:"url,omitempty"`
}

type AnalyticsRequest struct {
    IPAddress string `json:"ip_address"`
    UserAgent string `json:"user_agent"`
    Referer   string `json:"referer"`
}

type StatsResponse struct {
    TotalClicks  int               `json:"total_clicks"`
    TodayClicks  int               `json:"today_clicks"`
    WeeklyClicks []DailyClickCount `json:"weekly_clicks"`
    Devices      map[string]int    `json:"devices"`
    OS           map[string]int    `json:"os"`
    Browsers     map[string]int    `json:"browsers"`
    Referers     map[string]int    `json:"referers"`
}

type DailyClickCount struct {
    Date   string `json:"date"`
    Clicks int    `json:"clicks"`
}