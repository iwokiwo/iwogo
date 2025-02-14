package entity

import "time"

type User struct {
	ID           int    `json:"id" gorm:"size:36;not null;uniqueIndex;primary_key"`
	Name         string `json:"name"`
	Email        string `json:"email" gorm:"size:100;not null"`
	Password     string `json:"password" gorm:"size:100;not null"`
	PasswordTemp string `json:"password_temp"`
	Phone        string `json:"phone"`
	Code         string `json:"code"`
	Active       int    `json:"active"`
	Role         string `json:"role"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
