package models

import "time"

type Question struct {
	ID                  int       `gorm:"primaryKey" json:"id"`
	QuizID              int       `gorm:"not null" json:"quiz_id"`
	Text                string    `gorm:"type:varchar(350);not null" json:"text"`
	Type                string    `gorm:"type:varchar(30)" json:"type"`
	HintExplanation     string    `gorm:"type:varchar(200)" json:"hint_explanation"`
	DifficultyLevel     string    `gorm:"type:varchar(10)" json:"difficulty_level"`
	Points              float64   `gorm:"type:decimal(10,2)" json:"points"`
	MultiChoiceAnsLimit int       `gorm:"type:int" json:"multi_choice_ans_limit"`
	CreationDate        time.Time `gorm:"type:date" json:"creation_date"`
	LastModifiedDate    time.Time `gorm:"type:date" json:"last_modified_date"`
}
