package roster

import (
	"GEWIS-Rooster/internal/models"
	"time"
)

type CreateRequest struct {
	Name string `json:"name"`

	Date time.Time `json:"date"`

	OrganID uint `json:"organId"`

	Shifts []string `json:"shifts"`

	TemplateID *uint `json:"templateId"`
} // @name RosterCreateRequest

type UpdateRequest struct {
	Name *string `json:"name"`

	Date *time.Time `json:"date"`

	Saved *bool `json:"saved"`
} // @name RosterUpdateRequest

type ShiftCreateRequest struct {
	Name string `json:"name"`

	RosterID uint `json:"rosterId"`
} // @name ShiftCreateRequest

type ShiftUpdateRequest struct {
	Order *int `json:"order"`
} // @name ShiftUpdateRequest

type AnswerCreateRequest struct {
	UserID uint `json:"userId"`

	RosterID uint `json:"rosterId"`

	RosterShiftID uint `json:"rosterShiftId"`

	Value string `json:"value"`
} // @name AnswerCreateRequest

type AnswerUpdateRequest struct {
	Value string `json:"value"`
} // @name AnswerUpdateRequest

type SavedShiftUpdateRequest struct {
	UserIDs []uint `json:"users"`
} // @name SavedShiftUpdateRequest

type SavedShiftResponse struct {
	SavedShifts []*models.SavedShift `json:"savedShifts"`

	SavedShiftOrdering []*models.SavedShiftOrdering `json:"savedShiftOrdering"`
} // @name SavedShiftResponse

type FilterParams struct {
	ID      *uint      `form:"id"`
	Date    *time.Time `form:"date" time_format:"2006-01-02"`
	OrganID *uint      `form:"organId"`
} // @name RosterFilterParams

type TemplateCreateRequest struct {
	OrganID uint `json:"organId"`

	Name string `json:"name"`

	Shifts []string `json:"shifts"`
} // @name TemplateCreateRequest

type TemplateFilterParams struct {
	OrganID *uint `form:"organId"`
} // @name TemplateFilterParams

// TODO updating roster templates does not yet work

type TemplateUpdateParams struct {
	Name string `json:"name"`

	Shifts []string `json:"shifts"`
} // @name TemplateUpdateParams

type TemplateShiftPreferenceCreateRequest struct {
	UserID uint `json:"userId"`

	RosterTemplateShiftID uint `json:"rosterTemplateShiftID"`

	Preference string `json:"preference"`
} // @name TemplateShiftPreferenceCreateRequest

type TemplateShiftPreferenceFilterParams struct {
	UserID uint `form:"userId"`

	TemplateID uint `form:"templateId"`
} // @name TemplateShiftPreferenceFilterParams

type TemplateShiftPreferenceUpdateRequest struct {
	Preference string `json:"preference"`
} // @name TemplateShiftPreferenceUpdateRequest
