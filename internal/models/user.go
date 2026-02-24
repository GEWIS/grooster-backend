package models

// User model
// This model defines a user which can input date into a roster
type User struct {
	BaseModel

	Name string `json:"name"`

	GEWISID uint `json:"gewis_id" gorm:"uniqueIndex:idx_name"`

	Organs []Organ `json:"organs" gorm:"many2many:user_organs;"`

	Shifts []*SavedShift `gorm:"many2many:user_shift_saved;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
} // @name User
