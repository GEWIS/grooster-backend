package roster

import (
	"GEWIS-Rooster/internal/models"
	"errors"
	"fmt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"slices"
	"time"
)

type Service interface {
	CreateRoster(*CreateRequest) (*models.Roster, error)
	GetRosters(*FilterParams) ([]*models.Roster, error)
	UpdateRoster(uint, *UpdateRequest) (*models.Roster, error)
	DeleteRoster(ID uint) error

	CreateRosterShift(*ShiftCreateRequest) (*models.RosterShift, error)
	UpdateRosterShift(uint, *ShiftUpdateRequest) (*models.RosterShift, error)
	DeleteRosterShift(ID uint) error
	CreateRosterAnswer(*AnswerCreateRequest) (*models.RosterAnswer, error)
	UpdateRosterAnswer(uint, *AnswerUpdateRequest) (*models.RosterAnswer, error)

	SaveRoster(uint) error
	UpdateSavedShift(uint, *SavedShiftUpdateRequest) (*models.SavedShift, error)
	GetSavedRoster(uint) ([]*models.SavedShift, []*models.SavedShiftOrdering, error)

	CreateRosterTemplate(*TemplateCreateRequest) (*models.RosterTemplate, error)
	GetRosterTemplate(uint) (*models.RosterTemplate, error)
	GetRosterTemplates(*TemplateFilterParams) ([]*models.RosterTemplate, error)
	UpdateRosterTemplate(uint, *TemplateUpdateParams) (*models.RosterTemplate, error)
	DeleteRosterTemplate(ID uint) error

	CreateRosterTemplateShiftPreference(TemplateShiftPreferenceCreateRequest) (*models.RosterTemplateShiftPreference, error)
	GetRosterTemplateShiftPreferences(TemplateShiftPreferenceFilterParams) ([]models.RosterTemplateShiftPreference, error)
	UpdateRosterTemplateShiftPreference(uint, TemplateShiftPreferenceUpdateRequest) (*models.RosterTemplateShiftPreference, error)
}

type service struct {
	db *gorm.DB
}

func NewRosterService(db *gorm.DB) Service {
	return &service{db: db}
}

func (s *service) CreateRoster(params *CreateRequest) (*models.Roster, error) {
	var users []models.User
	var values = models.Values{"J", "X", "L", "N"} //TODO Change this to input values

	err := s.db.Joins("JOIN user_organs ON user_organs.user_id = users.id").
		Where("user_organs.organ_id = ?", params.OrganID).
		Find(&users).Error

	if err != nil {
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

	nameToShiftID := make(map[string]uint)

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

			nameToShiftID[rosterShift.Name] = rosterShift.ID
		}
	}

	if params.TemplateID != nil {
		userIDs := make([]uint, len(users))
		for i, getUser := range users {
			userIDs[i] = getUser.ID
		}

		var preferences []models.RosterTemplateShiftPreference
		err := s.db.Preload("RosterTemplateShift").
			Joins("JOIN roster_template_shifts ON roster_template_shifts.id = roster_template_shift_preferences.roster_template_shift_id").
			Where("roster_template_shift_preferences.user_id IN ? AND roster_template_shifts.template_id = ?", userIDs, *params.TemplateID).
			Find(&preferences).Error
		if err != nil {
			return nil, err
		}

		for _, pref := range preferences {
			if newShiftID, exists := nameToShiftID[pref.RosterTemplateShift.ShiftName]; exists {
				answer := models.RosterAnswer{
					UserID:        pref.UserID,
					RosterID:      roster.ID,
					RosterShiftID: newShiftID,
					Value:         pref.Preference,
				}

				if err := s.db.Create(&answer).Error; err != nil {
					return nil, err
				}
			}
		}
	}

	if err := s.db.Preload("Organ").Preload("RosterShift").First(&roster, roster.ID).Error; err != nil {
		return nil, err
	}

	return &roster, nil
}

func (s *service) GetRosters(params *FilterParams) ([]*models.Roster, error) {
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

func (s *service) UpdateRoster(id uint, params *UpdateRequest) (*models.Roster, error) {
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

func (s *service) DeleteRoster(ID uint) error {
	var roster models.Roster
	if err := s.db.Delete(&roster, ID).Error; err != nil {
		return err
	}

	return nil
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

func (s *service) SaveRoster(ID uint) error {
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

func (s *service) GetSavedRoster(ID uint) ([]*models.SavedShift, []*models.SavedShiftOrdering, error) {
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

func (s *service) UpdateSavedShift(ID uint, updateParams *SavedShiftUpdateRequest) (*models.SavedShift, error) {
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

func (s *service) createSavedShift(rID uint, shift *models.RosterShift) error {
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

func (s *service) getSavedShiftOrdering(savedShifts []*models.SavedShift) ([]*models.SavedShiftOrdering, error) {
	var orderings []*models.SavedShiftOrdering

	for _, savedShift := range savedShifts {
		var users []*models.User

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
