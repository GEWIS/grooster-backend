package models

import (
	"time"
)

type Values []string

type Roster struct {
	BaseModel

	Name string `json:"name"`

	RosterShift []RosterShift `json:"rosterShift" gorm:"foreignKey:RosterID"`

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

	Roster *Roster `json:"-" gorm:"foreignKey:RosterID;constraint:fk_rosters_roster_shift,OnDelete:CASCADE;"`

	Order uint `json:"order"`

	ShiftGroupID *uint `json:"shiftGroupId" gorm:"default:null"`

	ShiftGroup *ShiftGroup `json:"-" gorm:"foreignKey:ShiftGroupID;constraint:OnDelete:SET NULL;"`
} // @name RosterShift

type RosterAnswer struct {
	BaseModel

	UserID uint `json:"userId" gorm:"uniqueIndex:user_answer_idx"`

	RosterID uint `json:"rosterId" gorm:"uniqueIndex:user_answer_idx"`

	Roster *Roster `json:"roster" gorm:"foreignKey:RosterID;constraint:fk_roster_answers_roster,OnDelete:CASCADE;"`

	RosterShiftID uint `json:"rosterShiftId" gorm:"uniqueIndex:user_answer_idx;"`

	RosterShift *RosterShift `json:"-" gorm:"foreignKey:RosterShiftID;constraint:OnDelete:CASCADE;"`

	Value string `json:"value"`
} // @name RosterAnswer

type SavedShift struct {
	BaseModel

	RosterID uint `json:"rosterId"`

	RosterShiftID uint `json:"rosterShiftId" gorm:"constraint:OnDelete:CASCADE;"`

	RosterShift *RosterShift `json:"rosterShift" gorm:"foreignKey:RosterShiftID;constraint:OnDelete:CASCADE;"`

	Users []*User `json:"users" gorm:"many2many:user_shift_saved;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
} // @name SavedShift

type SavedShiftOrdering struct {
	ShiftName string `json:"shiftName"`

	Users []*User `json:"users"`
} // @name SavedShiftOrdering

type RosterTemplate struct {
	BaseModel

	OrganID uint `json:"organId"`

	Name string `json:"name"`

	Shifts []RosterTemplateShift `json:"shifts" gorm:"foreignKey:TemplateID;constraint:OnDelete:CASCADE;"`
} // @name RosterTemplate

type RosterTemplateShift struct {
	BaseModel

	TemplateID uint `json:"templateId"`

	Template *RosterTemplate `json:"-" gorm:"foreignKey:TemplateID;constraint:OnDelete:CASCADE;"`

	ShiftName string `json:"shiftName"`

	ShiftGroupID *uint `json:"shiftGroupId" gorm:"default:null"`

	ShiftGroup *ShiftGroup `json:"-" gorm:"foreignKey:ShiftGroupID;constraint:OnDelete:SET NULL;"`
} // @name RosterTemplateShift

type RosterTemplateShiftPreference struct {
	BaseModel

	RosterTemplateShiftID uint `json:"rosterTemplateShiftID"`

	RosterTemplateShift *RosterTemplateShift `json:"-" gorm:"foreignKey:RosterTemplateShiftID;constraint:OnDelete:CASCADE;"`

	UserID uint `json:"userId"`

	User *User `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`

	Preference string `json:"value"`
} // @name RosterTemplateShiftPreference

type ShiftGroup struct {
	BaseModel

	OrganID uint `json:"organId" gorm:"uniqueIndex:organ_shift_group"`

	Organ Organ `json:"organ" gorm:"foreignKey:OrganID;constraint:OnDelete:CASCADE;"`

	Name string `json:"name" gorm:"uniqueIndex:organ_shift_group"`
} // @name ShiftGroup
