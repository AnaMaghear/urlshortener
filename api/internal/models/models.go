package models

import "time"

// ShortURL table
type ShortURL struct {
	ID          uint      `gorm:"primaryKey"`
	Code        string    `gorm:"size:16;uniqueIndex;not null"`
	OriginalURL string    `gorm:"type:text;not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	ExpiresAt   *time.Time
}

// ClickEvent table
type ClickEvent struct {
	ID         uint      `gorm:"primaryKey"`
	ShortURLID uint      `gorm:"index;not null"`
	ClickedAt  time.Time `gorm:"autoCreateTime"`
	IPAddress  string
	UserAgent  string
	Country    string
}
