package models

import "time"

type Answer struct {
	ID               int       `gorm:"primaryKey" json:"id"`
	QuestionID       int       `gorm:"not null" json:"question_id"`
	Text             string    `gorm:"type:varchar(350);not null" json:"text"`
	IsCorrect        bool      `gorm:"type:tinyint(1)" json:"is_correct"`
	CreationDate     time.Time `gorm:"type:date" json:"creation_date"`
	LastModifiedDate time.Time `gorm:"type:date" json:"last_modified_date"`
}
