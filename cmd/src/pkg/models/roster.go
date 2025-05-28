package models

import (
	"time"
)

// Roster
// @Description Roster info
// @Description with a gorm model, name, users, shifts and values
type Roster struct {
	BaseModel

	Name string `json:"name"`

	Users []*User `json:"users" gorm:"many2many:user_roster;"`

	RosterShift []*RosterShift `json:"rosterShift"`

	RosterAnswer []*RosterAnswer `json:"rosterAnswer" gorm:"foreignKey:RosterID"`

	Values Values `json:"values" gorm:"serializer:json"`

	Organ uint `json:"organ" gorm:"foreignKey:OrganID"`

	Date time.Time `json:"date"`

	Saved bool `json:"saved" gorm:"default:false"`
} // @name Roster

// RosterShift
// @Description One column of a roster
type RosterShift struct {
	BaseModel

	Name string `json:"name"`

	RosterID uint `json:"rosterId"`
} // @name RosterShift

// RosterAnswer
// @Description An answer from a user to a shift.
type RosterAnswer struct {
	BaseModel

	UserID uint `json:"userId" gorm:"uniqueIndex:user_answer_idx"`

	RosterID uint `json:"rosterId" gorm:"uniqueIndex:user_answer_idx"`

	RosterShiftID uint `json:"rosterShiftId" gorm:"uniqueIndex:user_answer_idx;constraint:OnDelete:CASCADE;"`

	Value string `json:"value"`
} // @name RosterAnswer

// SavedShift
// @Description A saved roster
type SavedShift struct {
	BaseModel

	RosterID uint `json:"rosterId"`

	RosterShiftID uint `json:"rosterShiftId"`

	RosterShift *RosterShift `json:"rosterShift"`

	Users []*User `json:"users" gorm:"many2many:user_shift_saved;"`
} // @name SavedShift

type Values []string

// RosterCreateRequest
// @Description Roster create request
type RosterCreateRequest struct {
	// Name the name of the new roster
	Name string `json:"name"`
	// Date the date that this roster will take place
	Date time.Time `json:"date"`
} // @name RosterCreateRequest

type RosterUpdateRequest struct {
	UserIDs *[]uint `json:"userIds"`
} // @name RosterUpdateRequest

type RosterShiftCreateRequest struct {
	Name string `json:"name"`

	RosterID uint `json:"rosterId"`
} // @name RosterShiftCreateRequest

type RosterAnswerCreateRequest struct {
	UserID uint `json:"userId"`

	RosterID uint `json:"rosterId"`

	RosterShiftID uint `json:"rosterShiftId"`

	Value string `json:"value"`
} // @name RosterAnswerCreateRequest

// RosterAnswerUpdateRequest
// @Description Request to update a certain roster answer
type RosterAnswerUpdateRequest struct {
	Value string `json:"value"`
} // @name RosterAnswerUpdateRequest

// SavedShiftUpdateRequest
// @Description Request to update a certain saved shift
type SavedShiftUpdateRequest struct {
	UserIDs []uint `json:"users"`
} // @name SavedShiftUpdateRequest
