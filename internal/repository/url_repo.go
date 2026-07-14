package repository

import (
    "context"
    "time"

    "github.com/Kai1313/url-shortener-fullstack/backend/internal/models"
    "github.com/redis/go-redis/v9"
    "gorm.io/gorm"
)

type URLRepository interface {
    Create(url *models.URL) error
    FindByCode(code string) (*models.URL, error)
    FindByID(id uint) (*models.URL, error)
    FindByUserID(userID uint) ([]models.URL, error)
    Update(url *models.URL) error
    Delete(id uint) error
    IncrementClicks(code string) error
    CacheURL(code string, url string) error
    GetCachedURL(code string) (string, error)
    DeleteCache(code string) error
    GetAnalytics(urlID uint) ([]models.Analytics, error)
    GetTodayClicks(urlID uint) (int, error)
    GetWeeklyClicks(urlID uint) ([]models.DailyClickCount, error)
}

type urlRepository struct {
    db    *gorm.DB
    redis *redis.Client
}

func NewURLRepository(db *gorm.DB, redis *redis.Client) URLRepository {
    return &urlRepository{
        db:    db,
        redis: redis,
    }
}

func (r *urlRepository) Create(url *models.URL) error {
    return r.db.Create(url).Error
}

func (r *urlRepository) FindByCode(code string) (*models.URL, error) {
    var url models.URL
    err := r.db.Where("short_code = ? AND is_active = ?", code, true).First(&url).Error
    if err != nil {
        return nil, err
    }
    return &url, nil
}

func (r *urlRepository) FindByID(id uint) (*models.URL, error) {
    var url models.URL
    err := r.db.First(&url, id).Error
    if err != nil {
        return nil, err
    }
    return &url, nil
}

func (r *urlRepository) FindByUserID(userID uint) ([]models.URL, error) {
    var urls []models.URL
    err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&urls).Error
    return urls, err
}

func (r *urlRepository) Update(url *models.URL) error {
    return r.db.Save(url).Error
}

func (r *urlRepository) Delete(id uint) error {
    return r.db.Delete(&models.URL{}, id).Error
}

func (r *urlRepository) IncrementClicks(code string) error {
    return r.db.Model(&models.URL{}).Where("short_code = ?", code).Update("clicks", gorm.Expr("clicks + ?", 1)).Error
}

func (r *urlRepository) CacheURL(code string, url string) error {
    ctx := context.Background()
    return r.redis.Set(ctx, "url:"+code, url, 1*time.Hour).Err()
}

func (r *urlRepository) GetCachedURL(code string) (string, error) {
    ctx := context.Background()
    return r.redis.Get(ctx, "url:"+code).Result()
}

func (r *urlRepository) DeleteCache(code string) error {
    ctx := context.Background()
    return r.redis.Del(ctx, "url:"+code).Err()
}

func (r *urlRepository) GetAnalytics(urlID uint) ([]models.Analytics, error) {
    var analytics []models.Analytics
    err := r.db.Where("url_id = ?", urlID).Order("created_at DESC").Find(&analytics).Error
    return analytics, err
}

func (r *urlRepository) GetTodayClicks(urlID uint) (int, error) {
    var count int64
    today := time.Now().Truncate(24 * time.Hour)
    err := r.db.Model(&models.Analytics{}).
        Where("url_id = ? AND created_at >= ?", urlID, today).
        Count(&count).Error
    return int(count), err
}

func (r *urlRepository) GetWeeklyClicks(urlID uint) ([]models.DailyClickCount, error) {
    var results []models.DailyClickCount

    // Get last 7 days of data
    query := `
        SELECT DATE(created_at) as date, COUNT(*) as clicks
        FROM analytics
        WHERE url_id = ? AND created_at >= NOW() - INTERVAL '7 days'
        GROUP BY DATE(created_at)
        ORDER BY date ASC
    `

    err := r.db.Raw(query, urlID).Scan(&results).Error
    if err != nil {
        return nil, err
    }

    // Fill in missing days with 0
    var allDays []models.DailyClickCount
    for i := 6; i >= 0; i-- {
        date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
        found := false
        for _, r := range results {
            if r.Date == date {
                allDays = append(allDays, r)
                found = true
                break
            }
        }
        if !found {
            allDays = append(allDays, models.DailyClickCount{Date: date, Clicks: 0})
        }
    }

    return allDays, nil
}