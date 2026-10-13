package roster

import (
	"GEWIS-Rooster/internal/models"
	"fmt"
	"slices"
)

type ShiftManager interface {
	CreateRosterShift(*ShiftCreateRequest) (*models.RosterShift, error)
	UpdateRosterShift(uint, *ShiftUpdateRequest) (*models.RosterShift, error)
	DeleteRosterShift(ID uint) error

	CreateRosterAnswer(*AnswerCreateRequest) (*models.RosterAnswer, error)
	UpdateRosterAnswer(uint, *AnswerUpdateRequest) (*models.RosterAnswer, error)
}

func (s *service) CreateRosterShift(createParams *ShiftCreateRequest) (*models.RosterShift, error) {
	var roster *models.Roster
	if err := s.db.First(&roster, createParams.RosterID).Error; err != nil {
		return nil, fmt.Errorf("roster not found: %w", err)
	}

	var maxOrdering int
	err := s.db.Model(&models.RosterShift{}).
		Where("roster_id = ?", createParams.RosterID).
		Select("COALESCE(MAX(`order`), -1)").
		Row().Scan(&maxOrdering)

	if err != nil {
		return nil, err
	}

	rosterShift := models.RosterShift{
		Name:     createParams.Name,
		RosterID: createParams.RosterID,
		Order:    uint(maxOrdering + 1),
	}

	if err := s.db.Create(&rosterShift).Error; err != nil {
		return nil, err
	}

	return &rosterShift, nil
}

func (s *service) UpdateRosterShift(ID uint, updateParams *ShiftUpdateRequest) (*models.RosterShift, error) {
	var rosterShift models.RosterShift

	if err := s.db.First(&rosterShift, ID).Error; err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	if updateParams.Order != nil {
		updates["order"] = updateParams.Order
	}

	if updateParams.ShiftGroupID != nil {
		updates["shift_group_id"] = updateParams.ShiftGroupID
	}

	if err := s.db.Model(&rosterShift).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &rosterShift, nil
}

func (s *service) DeleteRosterShift(ID uint) error {
	var rosterShift *models.RosterShift

	if err := s.db.Delete(&rosterShift, ID).Error; err != nil {
		return err
	}

	return nil
}

func (s *service) CreateRosterAnswer(params *AnswerCreateRequest) (*models.RosterAnswer, error) {
	var roster *models.Roster
	if err := s.db.First(&roster, params.RosterID).Error; err != nil {
		return nil, fmt.Errorf("roster not found: %w", err)
	}

	var rosterShift *models.RosterShift
	if err := s.db.First(&rosterShift, params.RosterShiftID).Error; err != nil {
		return nil, fmt.Errorf("roster shift not found: %w", err)
	}

	if !slices.Contains(roster.Values, params.Value) {
		return nil, fmt.Errorf("%s is not a valid value for this roster", params.Value)
	}

	rosterAnswer := models.RosterAnswer{
		UserID:        params.UserID,
		RosterID:      roster.ID,
		RosterShiftID: params.RosterShiftID,
		Value:         params.Value,
	}

	if err := s.db.Create(&rosterAnswer).Error; err != nil {
		return nil, err
	}

	return &rosterAnswer, nil
}

func (s *service) UpdateRosterAnswer(ID uint, updateParams *AnswerUpdateRequest) (*models.RosterAnswer, error) {
	var answer *models.RosterAnswer

	if err := s.db.First(&answer, ID).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&answer).Updates(updateParams).Error; err != nil {
		return nil, err
	}
	if err := s.db.First(&answer, ID).Error; err != nil {
		return nil, err
	}
	return answer, nil
}
