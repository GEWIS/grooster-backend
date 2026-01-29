package models

import (
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

	TemplateID *uint `json:"templateId" gorm:"foreignKey:TemplateID"`
} // @name Roster

type RosterShift struct {
	BaseModel

	Name string `json:"name"`

	RosterID uint `json:"rosterId"`

	Order uint `json:"order"`
} // @name RosterShift

type RosterAnswer struct {
	BaseModel

	UserID uint `json:"userId" gorm:"uniqueIndex:user_answer_idx"`

	RosterID uint `json:"rosterId" gorm:"uniqueIndex:user_answer_idx"`

	RosterShiftID uint `json:"rosterShiftId" gorm:"uniqueIndex:user_answer_idx;"`

	RosterShift *RosterShift `json:"-" gorm:"foreignKey:RosterShiftID;constraint:OnDelete:CASCADE;"`

	Value string `json:"value"`
} // @name RosterAnswer

type SavedShift struct {
	BaseModel

	RosterID uint `json:"rosterId"`

	RosterShiftID uint `json:"rosterShiftId" gorm:"constraint:OnDelete:CASCADE;"`

	RosterShift *RosterShift `json:"rosterShift" gorm:"foreignKey:RosterShiftID;constraint:OnDelete:CASCADE;"`

	Users []*User `json:"users" gorm:"many2many:user_shift_saved;"`
} // @name SavedShift

type SavedShiftOrdering struct {
	ShiftName string `json:"shiftName"`

	Users []*User `json:"users"`
} // @name SavedShiftOrdering

type RosterTemplate struct {
	BaseModel

	OrganID uint `json:"organId"`

	Name string `json:"name"`

	Shifts []RosterTemplateShift `json:"shifts" gorm:"foreignKey:TemplateID"`
} // @name RosterTemplate

type RosterTemplateShift struct {
	BaseModel

	TemplateID uint `json:"templateId"`

	ShiftName string `json:"shiftName"`
} // @name RosterTemplateShift

type RosterCreateRequest struct {
	Name string `json:"name"`

	Date time.Time `json:"date"`

	OrganID uint `json:"organId"`

	Shifts []string `json:"shifts"`
} // @name RosterCreateRequest

type RosterUpdateRequest struct {
	Name *string `json:"name"`

	Date *time.Time `json:"date"`

	Saved *bool `json:"saved"`
} // @name RosterUpdateRequest

type RosterShiftCreateRequest struct {
	Name string `json:"name"`

	RosterID uint `json:"rosterId"`
} // @name RosterShiftCreateRequest

type RosterShiftUpdateRequest struct {
	Order *int `json:"order"`
} // @name RosterShiftUpdateRequest

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
} // @name SavedShiftResponse

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

// TODO updating roster templates does not yet work

type RosterTemplateUpdateParams struct {
	Name string `json:"name"`

	Shifts []string `json:"shifts"`
} // @name RosterTemplateUpdateParams
