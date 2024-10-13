package models

import "time"

type User struct {
	ID               int       `gorm:"primaryKey" json:"id"`
	UserName         string    `gorm:"type:varchar(30);not null" json:"user_name"`
	UserFullName     string    `gorm:"type:varchar(64);not null" json:"user_full_name"`
	Email            string    `gorm:"type:varchar(320);not null" json:"email"`
	RegistrationDate time.Time `gorm:"type:date" json:"registration_date"`
	LastLoginDate    time.Time `gorm:"type:date" json:"last_login_date"`
	IsActive         bool      `gorm:"type:tinyint(1)" json:"is_active"`
	LastModifiedDate time.Time `gorm:"type:date" json:"last_modified_date"`
}
