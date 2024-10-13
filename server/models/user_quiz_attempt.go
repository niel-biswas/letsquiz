package models

import "time"

type UserQuizAttempt struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	UserID    int       `gorm:"not null" json:"user_id"`
	QuizID    int       `gorm:"not null" json:"quiz_id"`
	Score     float64   `gorm:"type:decimal(10,2)" json:"score"`
	StartTime time.Time `gorm:"type:datetime" json:"start_time"`
	EndTime   time.Time `gorm:"type:datetime" json:"end_time"`
}
