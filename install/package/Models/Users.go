package Models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           int    `json:"id" gorm:"size:36;not null;uniqueIndex;primary_key"`
	Name         string `json:"name"`
	Email        string `json:"email" gorm:"size:100;not null"`
	Password     string `json:"password" gorm:"size:100;not null"`
	PasswordTemp string `json:"password_temp"`
	Code         string `json:"code"`
	Phone        string `json:"phone"`
	Role         string `json:"role"`
	Active       int    `json:"active"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}

func (b *User) TableName() string {
	return "users"
}
