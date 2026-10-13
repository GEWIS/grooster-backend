package models

import (
	"time"
)

// BaseModel replaces gorm.Model but with JSON tags
type BaseModel struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
