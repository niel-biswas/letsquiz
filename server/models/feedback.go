package models

import "time"

type Feedback struct {
	ID           int       `gorm:"primaryKey" json:"id"`
	UserID       int       `gorm:"not null" json:"user_id"`
	QuizID       int       `gorm:"not null" json:"quiz_id"`
	Feedback     string    `gorm:"type:varchar(2000)" json:"feedback"`
	TicketID     string    `gorm:"type:varchar(45)" json:"ticket_id"`
	CreationDate time.Time `gorm:"type:date" json:"creation_date"`
}
