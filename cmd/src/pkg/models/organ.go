package models

// Organ
// @Description An organ that users can be part of.
type Organ struct {
	BaseModel

	Name string `gorm:"uniqueIndex"`
} // @name Organ
