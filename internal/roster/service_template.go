package roster

import (
	"GEWIS-Rooster/internal/models"
	"errors"
	"fmt"
	"gorm.io/datatypes"
)

type TemplateManager interface {
	CreateRosterTemplate(*TemplateCreateRequest) (*models.RosterTemplate, error)
	GetRosterTemplate(uint) (*models.RosterTemplate, error)
	GetRosterTemplates(*TemplateFilterParams) ([]*models.RosterTemplate, error)
	UpdateRosterTemplate(uint, *TemplateUpdateParams) (*models.RosterTemplate, error)
	DeleteRosterTemplate(ID uint) error

	UpdateRosterTemplateShift(uint, *TemplateShiftUpdateRequest) (*models.RosterTemplateShift, error)

	CreateRosterTemplateShiftPreference(TemplateShiftPreferenceCreateRequest) (*models.RosterTemplateShiftPreference, error)
	GetRosterTemplateShiftPreferences(TemplateShiftPreferenceFilterParams) ([]models.RosterTemplateShiftPreference, error)
	UpdateRosterTemplateShiftPreference(uint, TemplateShiftPreferenceUpdateRequest) (*models.RosterTemplateShiftPreference, error)
}

func (s *service) CreateRosterTemplate(params *TemplateCreateRequest) (*models.RosterTemplate, error) {
	var userOrgan models.Organ

	if err := s.db.First(&userOrgan, params.OrganID).Error; err != nil {
		return nil, err
	}

	if len(params.Shifts) == 0 {
		return nil, errors.New("no shifts were given")
	}

	shifts := make([]models.RosterTemplateShift, len(params.Shifts))
	for i, name := range params.Shifts {
		shifts[i] = models.RosterTemplateShift{
			ShiftName: name,
		}
	}

	template := models.RosterTemplate{
		OrganID: userOrgan.ID,
		Name:    params.Name,
		Shifts:  shifts,
	}

	if err := s.db.Create(&template).Error; err != nil {
		return nil, err
	}

	return &template, nil
}

func (s *service) GetRosterTemplate(ID uint) (*models.RosterTemplate, error) {
	var template models.RosterTemplate
	if err := s.db.Preload("Shifts").First(&template, ID).Error; err != nil {
		return nil, err
	}

	return &template, nil
}

func (s *service) GetRosterTemplates(params *TemplateFilterParams) ([]*models.RosterTemplate, error) {
	var templates []*models.RosterTemplate
	db := s.db.Model(&models.RosterTemplate{}).Preload("Shifts")

	if params != nil {
		if params.OrganID != nil {
			db.Where("organ_id = ?", params.OrganID)
		}
	}

	if err := db.Find(&templates).Error; err != nil {
		return nil, err
	}

	return templates, nil
}

func (s *service) UpdateRosterTemplate(id uint, params *TemplateUpdateParams) (*models.RosterTemplate, error) {
	var template models.RosterTemplate
	if err := s.db.First(&template, id).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"name":   params.Name,
		"shifts": datatypes.JSONSlice[string](params.Shifts),
	}

	if err := s.db.Model(&template).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &template, nil
}

func (s *service) DeleteRosterTemplate(ID uint) error {
	result := s.db.Where("id = ?", ID).Delete(&models.RosterTemplate{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("roster template with ID %d not found", ID)
	}
	return nil
}

func (s *service) UpdateRosterTemplateShift(ID uint, updateParams *TemplateShiftUpdateRequest) (*models.RosterTemplateShift, error) {
	var templateShift models.RosterTemplateShift

	if err := s.db.First(&templateShift, ID).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"shift_group_id": updateParams.ShiftGroupID,
	}

	if err := s.db.Model(&templateShift).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &templateShift, nil
}

func (s *service) CreateRosterTemplateShiftPreference(params TemplateShiftPreferenceCreateRequest) (*models.RosterTemplateShiftPreference, error) {
	templateShiftPreference := models.RosterTemplateShiftPreference{
		UserID:                params.UserID,
		RosterTemplateShiftID: params.RosterTemplateShiftID,
		Preference:            params.Preference,
	}

	if err := s.db.Create(&templateShiftPreference).Error; err != nil {
		return nil, err
	}

	return &templateShiftPreference, nil
}

func (s *service) GetRosterTemplateShiftPreferences(params TemplateShiftPreferenceFilterParams) ([]models.RosterTemplateShiftPreference, error) {
	var preferences []models.RosterTemplateShiftPreference

	query := s.db.Model(&models.RosterTemplateShiftPreference{})

	query = query.Where("user_id = ?", params.UserID)

	query = query.Joins("JOIN roster_template_shifts ON roster_template_shifts.id = roster_template_shift_preferences.roster_template_shift_id").
		Where("roster_template_shifts.template_id = ?", params.TemplateID)

	if err := query.Find(&preferences).Error; err != nil {
		return nil, err
	}

	return preferences, nil
}

func (s *service) UpdateRosterTemplateShiftPreference(id uint, params TemplateShiftPreferenceUpdateRequest) (*models.RosterTemplateShiftPreference, error) {
	var templateShiftPreference models.RosterTemplateShiftPreference
	if err := s.db.First(&templateShiftPreference, id).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"preference": params.Preference,
	}

	if err := s.db.Model(&templateShiftPreference).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &templateShiftPreference, nil
}
