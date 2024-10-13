package models

import "time"

type UserAnswer struct {
	ID             int       `gorm:"primaryKey" json:"id"`
	AttemptID      int       `gorm:"not null" json:"attempt_id"`
	QuestionID     int       `gorm:"not null" json:"question_id"`
	ChosenAnswerID int       `gorm:"not null" json:"chosen_answer_id"`
	IsCorrect      bool      `gorm:"type:tinyint(1)" json:"is_correct"`
	AnsweredDate   time.Time `gorm:"type:date" json:"answered_date"`
}
