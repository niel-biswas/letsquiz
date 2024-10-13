package models

import "time"

type Leaderboard struct {
	ID           int       `gorm:"primaryKey" json:"id"`
	UserID       int       `gorm:"not null" json:"user_id"`
	QuizID       int       `gorm:"not null" json:"quiz_id"`
	Score        float64   `gorm:"type:decimal(10,2)" json:"score"`
	UserRank     int       `gorm:"type:smallint" json:"user_rank"`
	CreationDate time.Time `gorm:"type:date" json:"creation_date"`
}
