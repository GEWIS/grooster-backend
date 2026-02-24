package models

// Organ
// @Description An organ that users can be part of.
type Organ struct {
	BaseModel

	Name string `json:"name" gorm:"uniqueIndex"`

	Users []*User `json:"users" gorm:"many2many:user_organs;"`
} // @name Organ

type UserOrgan struct {
	UserID uint `gorm:"primaryKey"`

	OrganID uint `gorm:"primaryKey"`

	Username string `json:"username" gorm:"size:255"`
} // @name UserOrgan
