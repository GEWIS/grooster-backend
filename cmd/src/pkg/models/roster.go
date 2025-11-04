package models

import (
	"gorm.io/datatypes"
	"time"
)

type Values []string

type Roster struct {
	BaseModel

	Name string `json:"name"`

	RosterShift []RosterShift `json:"rosterShift"`

	RosterAnswer []RosterAnswer `json:"rosterAnswer" gorm:"foreignKey:RosterID"`

	Values Values `json:"values" gorm:"serializer:json"`

	OrganID uint `json:"organId"`

	Organ Organ `json:"organ" gorm:"foreignKey:OrganID"`

	Date time.Time `json:"date"`

	Saved bool `json:"saved" gorm:"default:false"`
} // @name Roster

type RosterShift struct {
	BaseModel

	Name string `json:"name"`

	RosterID uint `json:"rosterId"`
} // @name RosterShift

type RosterAnswer struct {
	BaseModel

	UserID uint `json:"userId" gorm:"uniqueIndex:user_answer_idx"`

	RosterID uint `json:"rosterId" gorm:"uniqueIndex:user_answer_idx"`

	RosterShiftID uint `json:"rosterShiftId" gorm:"uniqueIndex:user_answer_idx;constraint:OnDelete:CASCADE;"`

	Value string `json:"value"`
} // @name RosterAnswer

type SavedShift struct {
	BaseModel

	RosterID uint `json:"rosterId"`

	RosterShiftID uint `json:"rosterShiftId"`

	RosterShift *RosterShift `json:"rosterShift"`

	Users []*User `json:"users" gorm:"many2many:user_shift_saved;"`
} // @name SavedShift

type SavedShiftOrdering struct {
	ShiftName string `json:"shiftName"`

	Users []*User `json:"users"`
}

type RosterTemplate struct {
	BaseModel

	OrganID uint `json:"organId"`

	Name string `json:"name"`

	Shifts datatypes.JSONSlice[string] `gorm:"serializer:json" json:"shifts"`
} // @name RosterTemplate

type RosterCreateRequest struct {
	Name string `json:"name"`

	Date time.Time `json:"date"`

	OrganID uint `json:"organId"`

	Shifts []string `json:"shifts"`
} // @name RosterCreateRequest

type RosterUpdateRequest struct {
	Name *string `json:"name"`

	Date *time.Time `json:"date"`
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

type RosterAnswerUpdateRequest struct {
	Value string `json:"value"`
} // @name RosterAnswerUpdateRequest

type SavedShiftUpdateRequest struct {
	UserIDs []uint `json:"users"`
} // @name SavedShiftUpdateRequest

type SavedShiftResponse struct {
	SavedShifts []*SavedShift `json:"savedShifts"`

	SavedShiftOrdering []*SavedShiftOrdering `json:"savedShiftOrdering"`
}

type RosterFilterParams struct {
	ID      *uint      `form:"id"`
	Date    *time.Time `form:"date" time_format:"2006-01-02"`
	OrganID *uint      `form:"organId"`
} // @name RosterFilterParams

type RosterTemplateCreateRequest struct {
	OrganID uint `json:"organId"`

	Name string `json:"name"`

	Shifts []string `json:"shifts"`
} // @name RosterTemplateCreateRequest

type RosterTemplateFilterParams struct {
	OrganID *uint `form:"organId"`
} // @name RosterTemplateFilterParams

type RosterTemplateUpdateParams struct {
	Name string `json:"name"`

	Shifts []string `json:"shifts"`
} // @name RosterTemplateUpdateParams
