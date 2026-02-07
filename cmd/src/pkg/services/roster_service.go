package services

import (
	"GEWIS-Rooster/cmd/src/pkg/models"
	"errors"
	"fmt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"slices"
	"time"
)

type RosterServiceInterface interface {
	CreateRoster(*models.RosterCreateRequest) (*models.Roster, error)
	GetRosters(*models.RosterFilterParams) ([]*models.Roster, error)
	UpdateRoster(uint, *models.RosterUpdateRequest) (*models.Roster, error)
	DeleteRoster(ID uint) error

	CreateRosterShift(*models.RosterShiftCreateRequest) (*models.RosterShift, error)
	UpdateRosterShift(uint, *models.RosterShiftUpdateRequest) (*models.RosterShift, error)
	DeleteRosterShift(ID uint) error
	CreateRosterAnswer(*models.RosterAnswerCreateRequest) (*models.RosterAnswer, error)
	UpdateRosterAnswer(uint, *models.RosterAnswerUpdateRequest) (*models.RosterAnswer, error)

	SaveRoster(uint) error
	UpdateSavedShift(uint, *models.SavedShiftUpdateRequest) (*models.SavedShift, error)
	GetSavedRoster(uint) ([]*models.SavedShift, []*models.SavedShiftOrdering, error)

	CreateRosterTemplate(*models.RosterTemplateCreateRequest) (*models.RosterTemplate, error)
	GetRosterTemplate(uint) (*models.RosterTemplate, error)
	GetRosterTemplates(*models.RosterTemplateFilterParams) ([]*models.RosterTemplate, error)
	UpdateRosterTemplate(uint, *models.RosterTemplateUpdateParams) (*models.RosterTemplate, error)
	DeleteRosterTemplate(ID uint) error

	CreateRosterTemplateShiftPreference(models.TemplateShiftPreferenceCreateRequest) (*models.RosterTemplateShiftPreference, error)
	UpdateRosterTemplateShiftPreference(uint, models.TemplateShiftPreferenceUpdateRequest) (*models.RosterTemplateShiftPreference, error)
}

type RosterService struct {
	db *gorm.DB
}

func NewRosterService(db *gorm.DB) *RosterService {
	return &RosterService{db: db}
}

func (s *RosterService) CreateRoster(params *models.RosterCreateRequest) (*models.Roster, error) {
	var users []models.User
	var values = models.Values{"J", "X", "L", "N"} //TODO Change this to input values

	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}
	if err := s.db.Find(&models.Organ{}, params.OrganID).Error; err != nil {
		return nil, err
	}
	if !isTodayOrLater(params.Date) {
		return nil, errors.New("date must be after the current date")
	}
	if params.Name == "" {
		return nil, errors.New("name is required")
	}

	roster := models.Roster{
		Name:    params.Name,
		Date:    params.Date,
		OrganID: params.OrganID,
		Values:  values,
	}

	if err := s.db.Create(&roster).Error; err != nil {
		return nil, err
	}

	if params.Shifts != nil && len(params.Shifts) > 0 {
		for index, shift := range params.Shifts {
			rosterShift := &models.RosterShift{
				Name:     shift,
				RosterID: roster.ID,
				Order:    uint(index),
			}

			if err := s.db.Create(&rosterShift).Error; err != nil {
				return nil, err
			}
		}
	}

	if params.TemplateID != nil {
		
	}

	if err := s.db.Preload("Organ").Preload("RosterShift").First(&roster, roster.ID).Error; err != nil {
		return nil, err
	}

	return &roster, nil
}

func (s *RosterService) GetRosters(params *models.RosterFilterParams) ([]*models.Roster, error) {
	db := s.db.Model(&models.Roster{})

	if params.ID != nil {
		db = db.Where("id = ?", *params.ID)
	}

	if params.Date != nil {
		db = db.Where("date = ?", *params.Date)
	}

	if params.OrganID != nil {
		db = db.Where("organ_id = ?", *params.OrganID)
	}

	db.
		Preload("RosterShift").
		Preload("RosterAnswer").
		Preload("Organ")

	var rosters []*models.Roster
	if err := db.Find(&rosters).Error; err != nil {
		return nil, err
	}

	return rosters, nil
}

func (s *RosterService) UpdateRoster(id uint, params *models.RosterUpdateRequest) (*models.Roster, error) {
	var roster *models.Roster

	if err := s.db.First(&roster, id).Error; err != nil {
		return nil, err
	}

	if params.Date != nil && !isTodayOrLater(*params.Date) {
		return nil, errors.New("date must be after the current date")
	}

	if params.Date != nil {
		roster.Date = *params.Date
	}
	if params.Name != nil {
		roster.Name = *params.Name
	}
	if params.Saved != nil {
		roster.Saved = *params.Saved
	}

	if err := s.db.Save(&roster).Error; err != nil {
		return nil, err
	}

	return roster, nil
}

func (s *RosterService) DeleteRoster(ID uint) error {
	var roster models.Roster
	if err := s.db.Delete(&roster, ID).Error; err != nil {
		return err
	}

	return nil
}

func (s *RosterService) CreateRosterShift(createParams *models.RosterShiftCreateRequest) (*models.RosterShift, error) {
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

func (s *RosterService) UpdateRosterShift(ID uint, updateParams *models.RosterShiftUpdateRequest) (*models.RosterShift, error) {
	var rosterShift models.RosterShift

	if err := s.db.First(&rosterShift, ID).Error; err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	if updateParams.Order != nil {
		updates["order"] = updateParams.Order
	}

	if err := s.db.Model(&rosterShift).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &rosterShift, nil
}

func (s *RosterService) DeleteRosterShift(ID uint) error {
	var rosterShift *models.RosterShift

	if err := s.db.Delete(&rosterShift, ID).Error; err != nil {
		return err
	}

	return nil
}

func (s *RosterService) CreateRosterAnswer(params *models.RosterAnswerCreateRequest) (*models.RosterAnswer, error) {
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

func (s *RosterService) UpdateRosterAnswer(ID uint, updateParams *models.RosterAnswerUpdateRequest) (*models.RosterAnswer, error) {
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

func (s *RosterService) SaveRoster(ID uint) error {
	var roster *models.Roster
	if err := s.db.Preload("RosterShift").First(&roster, ID).Error; err != nil {
		return err
	}

	for _, shift := range roster.RosterShift {
		var existing models.SavedShift
		err := s.db.Where("roster_id = ? AND roster_shift_id = ?", roster.ID, shift.ID).First(&existing).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := s.createSavedShift(roster.ID, &shift); err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	roster.Saved = true
	if err := s.db.Save(&roster).Error; err != nil {
		return err
	}

	return nil
}

func (s *RosterService) GetSavedRoster(ID uint) ([]*models.SavedShift, []*models.SavedShiftOrdering, error) {
	var savedShifts []*models.SavedShift
	if err := s.db.Preload(clause.Associations).Where("roster_id = ?", ID).Find(&savedShifts).Error; err != nil {
		return nil, nil, err
	}

	savedShiftOrdering, err := s.getSavedShiftOrdering(savedShifts)

	if err != nil {
		return nil, nil, err
	}

	return savedShifts, savedShiftOrdering, nil
}

func (s *RosterService) UpdateSavedShift(ID uint, updateParams *models.SavedShiftUpdateRequest) (*models.SavedShift, error) {
	var saved *models.SavedShift
	if err := s.db.Preload("Users").First(&saved, ID).Error; err != nil {
		return nil, err
	}

	if updateParams.UserIDs != nil {
		var users []*models.User
		if err := s.db.Where("ID IN ?", updateParams.UserIDs).Find(&users).Error; err != nil {

			return nil, err
		}
		// Replace existing users with the new set
		if err := s.db.Model(&saved).Association("Users").Replace(users); err != nil {
			return nil, err
		}
		// Reload associations to get fresh data
		if err := s.db.Preload("Users").Preload("RosterShift").First(&saved, ID).Error; err != nil {
			return nil, err
		}
	}

	return saved, nil
}

func (s *RosterService) CreateRosterTemplate(params *models.RosterTemplateCreateRequest) (*models.RosterTemplate, error) {
	var organ models.Organ

	if err := s.db.First(&organ, params.OrganID).Error; err != nil {
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
		OrganID: organ.ID,
		Name:    params.Name,
		Shifts:  shifts,
	}

	if err := s.db.Create(&template).Error; err != nil {
		return nil, err
	}

	return &template, nil
}

func (s *RosterService) GetRosterTemplate(ID uint) (*models.RosterTemplate, error) {
	var template models.RosterTemplate
	if err := s.db.Preload("Shifts").First(&template, ID).Error; err != nil {
		return nil, err
	}

	return &template, nil
}

func (s *RosterService) GetRosterTemplates(params *models.RosterTemplateFilterParams) ([]*models.RosterTemplate, error) {
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

func (s *RosterService) UpdateRosterTemplate(id uint, params *models.RosterTemplateUpdateParams) (*models.RosterTemplate, error) {
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

func (s *RosterService) DeleteRosterTemplate(ID uint) error {
	result := s.db.Where("id = ?", ID).Delete(&models.RosterTemplate{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("roster template with ID %d not found", ID)
	}
	return nil
}

func (s *RosterService) CreateRosterTemplateShiftPreference(params models.TemplateShiftPreferenceCreateRequest) (*models.RosterTemplateShiftPreference, error) {
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

func (s *RosterService) UpdateRosterTemplateShiftPreference(id uint, params models.TemplateShiftPreferenceUpdateRequest) (*models.RosterTemplateShiftPreference, error) {
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

func (s *RosterService) createSavedShift(rID uint, shift *models.RosterShift) error {
	var savedShift = models.SavedShift{
		RosterID:    rID,
		RosterShift: shift,
		Users:       []*models.User{},
	}

	if err := s.db.Create(&savedShift).Error; err != nil {
		return err
	}
	return nil
}

func (s *RosterService) getSavedShiftOrdering(savedShifts []*models.SavedShift) ([]*models.SavedShiftOrdering, error) {
	var orderings []*models.SavedShiftOrdering

	for _, savedShift := range savedShifts {
		users := []*models.User{}

		var organID uint
		if err := s.db.Model(&models.Roster{}).
			Select("organ_id").
			Where("id = ?", savedShift.RosterID).
			Scan(&organID).Error; err != nil {
			return nil, err
		}

		err := s.db.Table("users AS u").
			Select("u.*, MAX(r.date) AS last_date").
			Joins("JOIN user_organs AS uo ON u.id = uo.user_id").
			Joins("LEFT JOIN user_shift_saved AS uss ON uss.user_id = u.id").
			Joins("LEFT JOIN roster_shifts AS rs ON rs.name = ? ", savedShift.RosterShift.Name).
			Joins("LEFT JOIN saved_shifts AS ss ON ss.roster_shift_id = rs.id AND ss.id = uss.saved_shift_id").
			Joins("LEFT JOIN rosters AS r ON r.id = ss.roster_id").
			Where("uo.organ_id = ?", organID).
			Group("u.id").
			Order("last_date ASC"). // Removed NULLS FIRST for MariaDB compatibility
			Scan(&users).Error

		if err != nil {
			log.Println("Error:", err)
		}

		orderings = append(orderings, &models.SavedShiftOrdering{
			ShiftName: savedShift.RosterShift.Name,
			Users:     users,
		})
	}

	return orderings, nil
}

func isTodayOrLater(date time.Time) bool {
	now := time.Now().In(date.Location())

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, date.Location())
	inputDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	return !inputDate.Before(today)
}
