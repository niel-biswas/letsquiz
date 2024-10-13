package models

type Category struct {
	ID          int    `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(70);not null" json:"name"`
	Description string `gorm:"type:varchar(300)" json:"description"`
}
