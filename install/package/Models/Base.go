package Models

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ID        uint       `json:"id" gorm:"size:36;not null;uniqueIndex;primary_key" form:"id"`
	CreatedAt *time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt *time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	DeletedAt gorm.DeletedAt
}
