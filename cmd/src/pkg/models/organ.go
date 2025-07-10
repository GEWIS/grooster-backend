package models

// Organ
// @Description An organ that users can be part of.
type Organ struct {
	BaseModel

	Name string `gorm:"uniqueIndex"`

	Users []*User `json:"users" gorm:"many2many:user_organs;"`
} // @name Organ
