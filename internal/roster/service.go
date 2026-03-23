package roster

import (
	"GEWIS-Rooster/internal/models"
	"GEWIS-Rooster/internal/user"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"time"
)

type Service interface {
	RosterManager
	ShiftManager
	TemplateManager

	FillRosterPreferences(uint) ([]*models.RosterAnswer, error)

	SaveRoster(uint) error
	UpdateSavedShift(uint, *SavedShiftUpdateRequest) (*models.SavedShift, error)
	GetSavedRoster(uint) ([]*models.SavedShift, []*models.SavedShiftOrdering, error)

	CreateShiftGroup(ShiftGroupCreateRequest) (*models.ShiftGroup, error)
	GetShiftGroups(ShiftGroupFilterParams) (*[]models.ShiftGroup, error)
	GetShiftGroup(uint) (*models.ShiftGroup, error)

	UpdateShiftGroupPriority(groupID uint, params GroupUpdatePriorityParam) (*models.ShiftGroupPriority, error)
}

type UserProvider interface {
	Get(*user.FilterParams) ([]*models.User, error)
}

type service struct {
	db *gorm.DB
	u  UserProvider
}

func NewRosterService(db *gorm.DB, userService UserProvider) Service {
	return &service{db: db, u: userService}
}

func (s *service) FillRosterPreferences(rosterID uint) ([]*models.RosterAnswer, error) {
	filter := &FilterParams{
		ID: &rosterID,
	}

	rosters, err := s.GetRosters(filter)

	if err != nil {
		return nil, err
	}

	if len(rosters) == 0 || len(rosters) > 1 {
		return nil, errors.New("only one roster should be found")
	}

	toFillRoster := rosters[0]

	if toFillRoster.TemplateID == nil {
		return nil, errors.New("roster shift has no linked template")
	}

	userFilter := user.FilterParams{
		OrganID: &toFillRoster.OrganID,
	}
	users, err := s.u.Get(&userFilter)

	if err != nil {
		return nil, err
	}

	userIDs := make([]uint, len(users))
	for i, getUser := range users {
		userIDs[i] = getUser.ID
	}

	var preferences []models.RosterTemplateShiftPreference
	err = s.db.Preload("RosterTemplateShift").
		Joins("JOIN roster_template_shifts ON roster_template_shifts.id = roster_template_shift_preferences.roster_template_shift_id").
		Where("roster_template_shift_preferences.user_id IN ? AND roster_template_shifts.template_id = ?", userIDs, toFillRoster.TemplateID).
		Find(&preferences).Error
	if err != nil {
		return nil, err
	}

	nameToShiftID := make(map[string]uint)

	for _, shift := range toFillRoster.RosterShift {
		nameToShiftID[shift.Name] = shift.ID
	}

	newAnswers := make([]*models.RosterAnswer, len(preferences))

	for prefID, pref := range preferences {
		if newShiftID, exists := nameToShiftID[pref.RosterTemplateShift.ShiftName]; exists {
			var existingAnswer models.RosterAnswer

			err := s.db.Where("user_id = ? AND roster_id = ? AND roster_shift_id = ?",
				pref.UserID, toFillRoster.ID, newShiftID).First(&existingAnswer).Error

			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}

			if errors.Is(err, gorm.ErrRecordNotFound) {
				answer := models.RosterAnswer{
					UserID:        pref.UserID,
					RosterID:      toFillRoster.ID,
					RosterShiftID: newShiftID,
					Value:         pref.Preference,
				}

				if err := s.db.Create(&answer).Error; err != nil {
					return nil, err
				}
				newAnswers[prefID] = &answer
			} else {
				newAnswers[prefID] = &existingAnswer
			}
		}
	}

	return newAnswers, nil
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

func (s *service) CreateShiftGroup(params ShiftGroupCreateRequest) (*models.ShiftGroup, error) {
	shiftGroup := models.ShiftGroup{
		OrganID: params.OrganID,
		Name:    params.Name,
	}

	if err := s.db.Create(&shiftGroup).Error; err != nil {
		return nil, err
	}

	return &shiftGroup, nil
}

func (s *service) GetShiftGroups(filters ShiftGroupFilterParams) (*[]models.ShiftGroup, error) {
	var shiftGroups []models.ShiftGroup

	db := s.db.Model(&models.ShiftGroup{})

	db = db.Where("organ_id = ?", filters.OrganID)

	if err := db.Find(&shiftGroups).Error; err != nil {
		return nil, err
	}

	return &shiftGroups, nil
}

func (s *service) GetShiftGroup(ID uint) (*models.ShiftGroup, error) {
	var shiftGroup models.ShiftGroup

	if err := s.db.First(&shiftGroup, ID).Error; err != nil {
		return nil, err
	}

	return &shiftGroup, nil
}

func (s *service) UpdateShiftGroupPriority(groupID uint, params GroupUpdatePriorityParam) (*models.ShiftGroupPriority, error) {
	newRecord := models.ShiftGroupPriority{
		UserID:       params.UserID,
		ShiftGroupID: groupID,
		Priority:     params.Priority,
	}

	err := s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "shift_group_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"priority", "updated_at"}),
	}).Create(&newRecord).Error

	if err != nil {
		return nil, err
	}

	err = s.db.Where("user_id = ? AND shift_group_id = ?", params.UserID, groupID).First(&newRecord).Error

	if err != nil {
		return nil, err
	}

	return &newRecord, err
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

		// Get the latest shift from users to check when they were last assigned
		// It first checks by groups and if no group is assigned it checks on name
		err := s.db.Table("users AS u").
			Select(`
				u.*, 
				MAX(r.date) AS last_date, 
				COALESCE(MAX(sgp.priority), 1) AS group_priority
    		`).
			Joins("JOIN user_organs AS uo ON u.id = uo.user_id").
			Joins("JOIN roster_shifts AS target_rs ON target_rs.name = ?", savedShift.RosterShift.Name).
			Joins(`LEFT JOIN shift_group_priorities AS sgp ON 
				sgp.user_id = u.id AND 
				sgp.shift_group_id = target_rs.shift_group_id`).
			Joins(`LEFT JOIN roster_shifts AS rs ON (
				(target_rs.shift_group_id IS NOT NULL AND rs.shift_group_id = target_rs.shift_group_id) OR 
				(target_rs.shift_group_id IS NULL AND rs.name = target_rs.name)
			)`).
			Joins("LEFT JOIN user_shift_saved AS uss ON uss.user_id = u.id").
			Joins("LEFT JOIN saved_shifts AS ss ON ss.roster_shift_id = rs.id AND ss.id = uss.saved_shift_id").
			Joins("LEFT JOIN rosters AS r ON r.id = ss.roster_id").
			Where("uo.organ_id = ?", organID).
			Group("u.id").
			Order("group_priority DESC, last_date ASC").
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
