package models

// User model
// This model defines a user which can input date into a roster
type User struct {
	BaseModel

	Name string `json:"name"`

	GEWISID uint `json:"gewis_id" gorm:"uniqueIndex:idx_name"`

	Organs []Organ `json:"organs" gorm:"many2many:user_organs;"`
} // @name User

type UserCreateRequest struct {
	Name string

	GEWISID uint

	Organs []Organ
}
