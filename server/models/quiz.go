package models

import "time"

type Quiz struct {
	ID               int       `gorm:"primaryKey" json:"id"`
	Title            string    `gorm:"type:varchar(100);not null" json:"title"`
	Description      string    `gorm:"type:varchar(300)" json:"description"`
	ContentURL       string    `gorm:"type:varchar(2083)" json:"content_url"`
	CategoryID       int       `gorm:"not null" json:"category_id"`
	CreatorID        int       `gorm:"not null" json:"creator_id"`
	CreationDate     time.Time `gorm:"type:date" json:"creation_date"`
	LastModifiedDate time.Time `gorm:"type:date" json:"last_modified_date"`
	TimeLimitInMins  int       `gorm:"not null" json:"time_limit_in_mins"`
	Points           int       `gorm:"not null" json:"points"`
	DifficultyLevel  string    `gorm:"not null" json:"difficulty_level"`
	HintExplanation  string    `gorm:"not null" json:"hint_explanation"`
	QuestionCount    int       `gorm:"not null" json:"question_count"`
	IsActive         bool      `gorm:"type:tinyint(1)" json:"is_active"`
}
