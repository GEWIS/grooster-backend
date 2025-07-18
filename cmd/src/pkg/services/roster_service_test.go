package services

import (
	"GEWIS-Rooster/cmd/seeder/seeder"
	database "GEWIS-Rooster/cmd/src/pkg"
	"GEWIS-Rooster/cmd/src/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
	"time"
)

type TestRosterSuite struct {
	suite.Suite
	db      *gorm.DB
	service RosterService
}

func (suite *TestRosterSuite) SetupTest() {
	db := database.ConnectDB(":memory:")
	seeder.Seeder(db)
	suite.db = db
	suite.service = RosterService{db: db}
}

func (suite *TestRosterSuite) TestCreateRoster_ValidInput() {
	params := models.RosterCreateRequest{
		Name:    "Valid Name",
		Date:    time.Now().Add(25 * time.Hour),
		OrganID: 1,
	}

	roster, err := suite.service.CreateRoster(&params)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), roster)
	assert.Equal(suite.T(), params.Name, roster.Name)
	assert.Equal(suite.T(), params.Date.Local(), roster.Date.Local())
	assert.Equal(suite.T(), params.OrganID, roster.OrganID)
}

func (suite *TestRosterSuite) TestCreateRoster_EmptyName() {
	params := models.RosterCreateRequest{
		Name:    "",
		Date:    time.Now().Add(25 * time.Hour),
		OrganID: 1,
	}

	roster, err := suite.service.CreateRoster(&params)
	assert.Error(suite.T(), err, "Expected error for empty name")
	assert.Nil(suite.T(), roster)
}

func (suite *TestRosterSuite) TestCreateRoster_ZeroDate() {
	params := models.RosterCreateRequest{
		Name:    "Valid Name",
		Date:    time.Time{},
		OrganID: 1,
	}

	roster, err := suite.service.CreateRoster(&params)
	assert.Error(suite.T(), err, "Expected error for zero date")
	assert.Nil(suite.T(), roster)
}

func (suite *TestRosterSuite) TestCreateRoster_InvalidOrganID() {
	params := models.RosterCreateRequest{
		Name:    "Valid Name",
		Date:    time.Now().Add(25 * time.Hour),
		OrganID: 0,
	}

	roster, err := suite.service.CreateRoster(&params)
	assert.Error(suite.T(), err, "Expected error for invalid organ ID")
	assert.Nil(suite.T(), roster)
}

func (suite *TestRosterSuite) TestCreateRoster_WithShifts() {
	var organ models.Organ
	suite.db.First(&organ)

	shift := []string{"Shift 1", "Shift 42"}

	params := models.RosterCreateRequest{
		Name:    "Valid Name",
		Date:    time.Now().Add(25 * time.Hour),
		OrganID: organ.ID,
		Shifts:  shift,
	}

	roster, err := suite.service.CreateRoster(&params)
	assert.NoError(suite.T(), err)

	var rosterShiftNames []string
	for _, rs := range roster.RosterShift {
		rosterShiftNames = append(rosterShiftNames, rs.Name)
	}

	assert.ElementsMatch(suite.T(), shift, rosterShiftNames)
}

func (suite *TestRosterSuite) TestGetRosters_All() {
	cParams := models.RosterCreateRequest{
		Name:    "Valid Name",
		Date:    time.Now().Add(25 * time.Hour),
		OrganID: 1,
	}
	_, _ = suite.service.CreateRoster(&cParams)
	var params models.RosterFilterParams
	rosters, err := suite.service.GetRosters(&params)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), rosters)

	var foundRoster *models.Roster
	for _, r := range rosters {
		if r.Name == cParams.Name {
			foundRoster = r
			break
		}
	}
	assert.NotNil(suite.T(), foundRoster, "Expected to find a roster with name %q", cParams.Name)
	if foundRoster == nil {
		return
	}

	assert.NotNil(suite.T(), foundRoster.RosterShift, "RosterShift should not be nil")
	assert.NotNil(suite.T(), foundRoster.RosterAnswer, "RosterAnswer should not be nil")
	assert.NotNil(suite.T(), foundRoster.Organ, "Organ should be loaded")
	assert.Equal(suite.T(), cParams.OrganID, foundRoster.OrganID)
}

func (suite *TestRosterSuite) TestGetRosters_FilterByID() {
	cParams := models.RosterCreateRequest{
		Name:    "RosterByID",
		Date:    time.Now().Add(25 * time.Hour),
		OrganID: 1,
	}
	created, err := suite.service.CreateRoster(&cParams)
	assert.NoError(suite.T(), err)

	params := models.RosterFilterParams{
		ID: &created.ID,
	}
	rosters, err := suite.service.GetRosters(&params)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), rosters, 1)
	assert.Equal(suite.T(), created.ID, rosters[0].ID)
	assert.Equal(suite.T(), cParams.Name, rosters[0].Name)
}

func (suite *TestRosterSuite) TestGetRosters_FilterByDate() {
	targetDate := time.Now().Add(25 * time.Hour)

	cParams := models.RosterCreateRequest{
		Name:    "RosterByDate",
		Date:    targetDate,
		OrganID: 1,
	}
	_, err := suite.service.CreateRoster(&cParams)
	assert.NoError(suite.T(), err)

	params := models.RosterFilterParams{
		Date: &targetDate,
	}
	rosters, err := suite.service.GetRosters(&params)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), rosters)

	found := false
	for _, r := range rosters {
		if r.Date.Equal(targetDate) && r.Name == cParams.Name {
			found = true
			break
		}
	}
	assert.True(suite.T(), found, "Expected to find a roster with date %v and name %q", targetDate, cParams.Name)
}

func (suite *TestRosterSuite) TestGetRosters_FilterByOrganID() {
	organID := uint(1)
	cParams := models.RosterCreateRequest{
		Name:    "RosterByOrgan",
		Date:    time.Now().Add(25 * time.Hour),
		OrganID: organID,
	}
	_, err := suite.service.CreateRoster(&cParams)
	assert.NoError(suite.T(), err)

	params := models.RosterFilterParams{
		OrganID: &organID,
	}
	rosters, err := suite.service.GetRosters(&params)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), rosters)

	for _, r := range rosters {
		assert.Equal(suite.T(), organID, r.OrganID)
	}
}

func (suite *TestRosterSuite) TestGetRosters_FilterByMultipleFields() {
	organID := uint(1)
	targetDate := time.Now().Add(25 * time.Hour)

	cParams := models.RosterCreateRequest{
		Name:    "RosterMultiFilter",
		Date:    targetDate,
		OrganID: organID,
	}
	_, err := suite.service.CreateRoster(&cParams)
	assert.NoError(suite.T(), err)

	params := models.RosterFilterParams{
		Date:    &targetDate,
		OrganID: &organID,
	}
	rosters, err := suite.service.GetRosters(&params)
	assert.NoError(suite.T(), err)

	found := false
	for _, r := range rosters {
		if r.Date.Equal(targetDate) && r.OrganID == organID && r.Name == cParams.Name {
			found = true
			break
		}
	}
	assert.True(suite.T(), found, "Expected to find a roster with date %v, organID %d and name %q", targetDate, organID, cParams.Name)
}

func (suite *TestRosterSuite) TestGetRosters_ReturnEmpty() {
	var roster models.Roster
	suite.db.Last(&roster)
	assert.NotEmpty(suite.T(), roster)

	noExistID := roster.ID + 1

	params := models.RosterFilterParams{
		ID: &noExistID,
	}
	rosters, _ := suite.service.GetRosters(&params)

	assert.Empty(suite.T(), rosters)
}

func (suite *TestRosterSuite) TestUpdateRoster_Valid() {
	var roster *models.Roster
	suite.db.First(&roster)

	name := "New Name"
	params := models.RosterUpdateRequest{
		Name: &name,
	}

	roster, err := suite.service.UpdateRoster(roster.ID, &params)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), name, roster.Name)
}

func (suite *TestRosterSuite) TestUpdateRoster_OnlyDate() {
	var roster *models.Roster
	suite.db.First(&roster)

	newDate := time.Now().Add(48 * time.Hour)
	params := models.RosterUpdateRequest{
		Date: &newDate,
	}

	updatedRoster, err := suite.service.UpdateRoster(roster.ID, &params)
	assert.NoError(suite.T(), err)
	assert.WithinDuration(suite.T(), newDate, updatedRoster.Date, time.Second)
}

func (suite *TestRosterSuite) TestUpdateRoster_InvalidDate() {
	var roster *models.Roster
	suite.db.First(&roster)

	newDate := time.Now().Add(-24 * time.Hour)
	params := models.RosterUpdateRequest{
		Date: &newDate,
	}

	_, err := suite.service.UpdateRoster(roster.ID, &params)
	assert.Error(suite.T(), err)
}

func (suite *TestRosterSuite) TestUpdateRoster_NameAndDate() {
	var roster *models.Roster
	suite.db.First(&roster)

	newName := "Updated Roster"
	newDate := time.Now().Add(72 * time.Hour)
	params := models.RosterUpdateRequest{
		Name: &newName,
		Date: &newDate,
	}

	updatedRoster, err := suite.service.UpdateRoster(roster.ID, &params)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), newName, updatedRoster.Name)
	assert.WithinDuration(suite.T(), newDate, updatedRoster.Date, time.Second)
}

func (suite *TestRosterSuite) TestUpdateRoster_NoFields() {
	var roster *models.Roster
	suite.db.First(&roster)

	originalName := roster.Name
	originalDate := roster.Date

	params := models.RosterUpdateRequest{}

	updatedRoster, err := suite.service.UpdateRoster(roster.ID, &params)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), originalName, updatedRoster.Name)
	assert.WithinDuration(suite.T(), originalDate, updatedRoster.Date, time.Second)
}

func (suite *TestRosterSuite) TestUpdateRoster_NotFound() {
	var roster *models.Roster
	suite.db.Last(&roster)
	nonExistID := roster.ID + 1

	name := "Should Not Exist"
	params := models.RosterUpdateRequest{
		Name: &name,
	}

	updatedRoster, err := suite.service.UpdateRoster(nonExistID, &params)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), updatedRoster)
}

func (suite *TestRosterSuite) TestDeleteRoster_Valid() {
	cParams := models.RosterCreateRequest{
		Name:    "Valid Name",
		Date:    time.Now().Add(25 * time.Hour),
		OrganID: 1,
	}
	roster, err := suite.service.CreateRoster(&cParams)
	assert.NoError(suite.T(), err)

	err = suite.service.DeleteRoster(roster.ID)
	assert.NoError(suite.T(), err)

	params := models.RosterFilterParams{ID: &roster.ID}
	rosters, err := suite.service.GetRosters(&params)
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), rosters)
}

func (suite *TestRosterSuite) TestCreateRosterShift_Valid() {
	roster := models.Roster{
		Name:    "Test Roster",
		Values:  []string{"yes", "no"},
		OrganID: 1,
	}
	suite.db.Create(&roster)

	createParams := &models.RosterShiftCreateRequest{
		Name:     "Morning Shift",
		RosterID: roster.ID,
	}

	shift, err := suite.service.CreateRosterShift(createParams)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), shift)
	assert.Equal(suite.T(), createParams.Name, shift.Name)
	assert.Equal(suite.T(), createParams.RosterID, shift.RosterID)
}

func (suite *TestRosterSuite) TestCreateRosterShift_RosterNotFound() {
	createParams := &models.RosterShiftCreateRequest{
		Name:     "Evening Shift",
		RosterID: 9999,
	}

	shift, err := suite.service.CreateRosterShift(createParams)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), shift)
	assert.Contains(suite.T(), err.Error(), "roster not found")
}

func (suite *TestRosterSuite) TestDeleteRosterShift_Valid() {
	roster := models.Roster{
		Name:    "Roster for Deletion",
		Values:  []string{"yes", "no"},
		OrganID: 1,
	}
	suite.db.Create(&roster)

	shift := models.RosterShift{
		Name:     "Shift to Delete",
		RosterID: roster.ID,
	}
	suite.db.Create(&shift)

	err := suite.service.DeleteRosterShift(shift.ID)

	assert.NoError(suite.T(), err)

	var deletedShift models.RosterShift
	result := suite.db.First(&deletedShift, shift.ID)
	assert.Error(suite.T(), result.Error)
	assert.Equal(suite.T(), gorm.ErrRecordNotFound, result.Error)
}

func (suite *TestRosterSuite) TestDeleteRosterShift_NotFound() {
	nonExistentID := uint(999999)
	err := suite.service.DeleteRosterShift(nonExistentID)

	assert.NoError(suite.T(), err)
}

func (suite *TestRosterSuite) TestCreateRosterAnswer_Valid() {

	roster := models.Roster{
		Name:    "Test Roster",
		Values:  []string{"yes", "no"},
		OrganID: uint(1),
	}
	suite.db.Create(&roster)

	shift := models.RosterShift{
		RosterID: roster.ID,
	}
	suite.db.Create(&shift)

	createParams := &models.RosterAnswerCreateRequest{
		UserID:        1,
		RosterID:      roster.ID,
		RosterShiftID: shift.ID,
		Value:         "yes",
	}

	answer, err := suite.service.CreateRosterAnswer(createParams)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), answer)
	assert.Equal(suite.T(), "yes", answer.Value)
	assert.Equal(suite.T(), createParams.UserID, answer.UserID)
}

func (suite *TestRosterSuite) TestCreateRosterAnswer_InvalidValue() {
	roster := models.Roster{
		Name:    "Test Roster",
		Values:  []string{"yes", "no"},
		OrganID: 1,
	}
	suite.db.Create(&roster)

	shift := models.RosterShift{
		RosterID: roster.ID,
	}
	suite.db.Create(&shift)

	createParams := &models.RosterAnswerCreateRequest{
		UserID:        1,
		RosterID:      roster.ID,
		RosterShiftID: shift.ID,
		Value:         "maybe",
	}

	answer, err := suite.service.CreateRosterAnswer(createParams)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), answer)
	assert.Contains(suite.T(), err.Error(), "is not a valid value")
}

func (suite *TestRosterSuite) TestCreateRosterAnswer_RosterNotFound() {
	createParams := &models.RosterAnswerCreateRequest{
		UserID:        1,
		RosterID:      9999,
		RosterShiftID: 1,
		Value:         "yes",
	}

	answer, err := suite.service.CreateRosterAnswer(createParams)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), answer)
	assert.Contains(suite.T(), err.Error(), "roster not found")
}

func (suite *TestRosterSuite) TestCreateRosterAnswer_RosterShiftNotFound() {
	roster := models.Roster{
		Name:    "Test Roster",
		Values:  []string{"yes"},
		OrganID: 1,
	}
	suite.db.Create(&roster)

	createParams := &models.RosterAnswerCreateRequest{
		UserID:        1,
		RosterID:      roster.ID,
		RosterShiftID: 9999,
		Value:         "yes",
	}

	answer, err := suite.service.CreateRosterAnswer(createParams)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), answer)
	assert.Contains(suite.T(), err.Error(), "roster shift not found")
}

func (suite *TestRosterSuite) TestUpdateRosterAnswer_Valid() {
	var answer *models.RosterAnswer
	suite.db.First(&answer)

	updateParams := &models.RosterAnswerUpdateRequest{
		Value: "new value",
	}

	updatedAnswer, err := suite.service.UpdateRosterAnswer(answer.ID, updateParams)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), updatedAnswer)
	assert.Equal(suite.T(), "new value", updatedAnswer.Value)
}

func (suite *TestRosterSuite) TestUpdateRosterAnswer_NotFound() {
	nonExistentID := uint(99999)
	updateParams := &models.RosterAnswerUpdateRequest{
		Value: "new value",
	}

	updatedAnswer, err := suite.service.UpdateRosterAnswer(nonExistentID, updateParams)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), updatedAnswer)
}

func (suite *TestRosterSuite) TestSaveRoster_Success() {
	roster := models.Roster{
		Name:    "Test Roster",
		OrganID: 1,
	}
	suite.db.Create(&roster)

	shifts := []models.RosterShift{
		{RosterID: roster.ID, Name: "Shift 1"},
		{RosterID: roster.ID, Name: "Shift 2"},
	}
	for _, shift := range shifts {
		suite.db.Create(&shift)
	}

	err := suite.service.SaveRoster(roster.ID)

	assert.NoError(suite.T(), err)

	var savedRoster models.Roster
	err = suite.db.First(&savedRoster, roster.ID).Error
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), savedRoster.Saved)
}

func (suite *TestRosterSuite) TestSaveRoster_RosterNotFound() {
	err := suite.service.SaveRoster(99999)
	assert.Error(suite.T(), err)
}

func (suite *TestRosterSuite) TestUpdateSavedShift_Success() {
	var user models.User
	suite.db.First(&user)

	err := suite.service.SaveRoster(1)
	assert.NoError(suite.T(), err)

	savedShifts, _ := suite.service.GetSavedRoster(1)

	updateParams := &models.SavedShiftUpdateRequest{
		UserIDs: []uint{user.ID},
	}

	updatedSavedShift, err := suite.service.UpdateSavedShift(savedShifts[0].ID, updateParams)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), updatedSavedShift)

	var users []models.User
	err = suite.db.Model(updatedSavedShift).Association("Users").Find(&users)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), users, 1)
	assert.Equal(suite.T(), user.ID, users[0].ID)
}

func (suite *TestRosterSuite) TestUpdateSavedShift_UserLoadFailure() {
	err := suite.service.SaveRoster(1)
	assert.NoError(suite.T(), err)

	savedShifts, _ := suite.service.GetSavedRoster(1)

	updateParams := &models.SavedShiftUpdateRequest{
		UserIDs: []uint{99999, 99998},
	}

	updatedSavedShift, err := suite.service.UpdateSavedShift(savedShifts[0].ID, updateParams)
	assert.NoError(suite.T(), err)

	var users []models.User
	err = suite.db.Model(updatedSavedShift).Association("Users").Find(&users)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), users, 0)
}

func (suite *TestRosterSuite) TestRosterTemplateCreate_Valid() {
	var organ models.Organ
	suite.db.First(&organ)

	expectedShifts := []string{"Shift 1", "Shift 42"}

	params := models.RosterTemplateCreateRequest{
		OrganID: organ.ID,
		Shifts:  expectedShifts,
	}

	template, err := suite.service.CreateRosterTemplate(&params)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), template)

	assert.Equal(suite.T(), organ.ID, template.OrganID)
	assert.ElementsMatch(suite.T(), expectedShifts, template.Shifts)
}

func (suite *TestRosterSuite) TestRosterTemplateCreate_InValidOrgan() {
	var organ models.Organ
	suite.db.Last(&organ)

	expectedShifts := []string{"Shift 1", "Shift 42"}

	params := models.RosterTemplateCreateRequest{
		OrganID: uint(999),
		Shifts:  expectedShifts,
	}

	template, err := suite.service.CreateRosterTemplate(&params)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), template)
}

func (suite *TestRosterSuite) TestRosterTemplateCreate_NoShifts() {
	var organ models.Organ
	suite.db.Last(&organ)

	expectedShifts := []string{}

	params := models.RosterTemplateCreateRequest{
		OrganID: uint(999),
		Shifts:  expectedShifts,
	}

	template, err := suite.service.CreateRosterTemplate(&params)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), template)
}

func (suite *TestRosterSuite) TestRosterTemplateGet_OneValid() {
	var template models.RosterTemplate

	suite.db.First(&template)

	newTemplate, err := suite.service.GetRosterTemplate(template.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), template.ID, newTemplate.ID)
}

func (suite *TestRosterSuite) TestRosterTemplateGet_OneInValid() {
	invalidID := uint(9999)

	newTemplate, err := suite.service.GetRosterTemplate(invalidID)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), newTemplate)
}

func (suite *TestRosterSuite) TestRosterTemplateGet_All() {
	var count int64
	suite.db.Model(&models.RosterTemplate{}).Count(&count)

	templates, err := suite.service.GetRosterTemplates(nil)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), templates)
	assert.Len(suite.T(), templates, int(count))
}

func (suite *TestRosterSuite) TestRosterTemplateGet_ByOrganID() {
	var template *models.RosterTemplate
	suite.db.First(&template)

	params := models.RosterTemplateFilterParams{OrganID: &template.OrganID}

	templates, err := suite.service.GetRosterTemplates(&params)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), templates)

	assert.Equal(suite.T(), template.ID, templates[0].ID)
	assert.Equal(suite.T(), template.OrganID, templates[0].OrganID)
}

func (suite *TestRosterSuite) TestRosterTemplateGet_AllInvalid() {
	var template *models.RosterTemplate
	suite.db.First(&template)

	ID := uint(9999)
	params := models.RosterTemplateFilterParams{OrganID: &ID}

	templates, err := suite.service.GetRosterTemplates(&params)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), len(templates), 0)
}

func (suite *TestRosterSuite) TestRosterTemplateDelete_Valid() {
	var template *models.RosterTemplate
	suite.db.First(&template)

	err := suite.service.DeleteRosterTemplate(template.ID)
	assert.NoError(suite.T(), err)

	var newTemplate *models.RosterTemplate
	suite.db.First(&newTemplate)

	assert.NotEqual(suite.T(), template.ID, newTemplate.ID)
}

func (suite *TestRosterSuite) TestRosterTemplateDelete_InValid() {
	ID := uint(9999)

	err := suite.service.DeleteRosterTemplate(ID)
	assert.Error(suite.T(), err)
}

func TestRosterService(t *testing.T) {
	suite.Run(t, new(TestRosterSuite))
}
