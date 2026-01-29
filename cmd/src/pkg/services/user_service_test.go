package services

import (
	"GEWIS-Rooster/cmd/seeder/seeder"
	database "GEWIS-Rooster/cmd/src/pkg"
	"GEWIS-Rooster/cmd/src/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type TestUserSuite struct {
	suite.Suite
	db      *gorm.DB
	service UserService
}

func (suite *TestUserSuite) SetupTest() {
	db := database.ConnectDB(":memory:")
	seeder.Seeder(db)
	suite.db = db
	suite.service = UserService{db: db}
}

func (suite *TestUserSuite) TestCreateUser_ValidInput() {
	var organ models.Organ
	suite.db.First(&organ)

	params := models.UserCreateRequest{
		Name:    "Test User",
		GEWISID: uint(10),
		Organs:  []models.Organ{organ},
	}
	user, err := suite.service.Create(&params)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), params.Name, user.Name)
	assert.Equal(suite.T(), params.GEWISID, user.GEWISID)
	assert.Equal(suite.T(), params.Organs[0].Name, user.Organs[0].Name)
}

func (suite *TestUserSuite) TestCreateUser_WithoutOrgans() {
	params := models.UserCreateRequest{
		Name:    "Test User Without Organs",
		GEWISID: uint(11),
		Organs:  []models.Organ{},
	}
	user, err := suite.service.Create(&params)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), params.Name, user.Name)
	assert.Equal(suite.T(), 0, len(user.Organs))
}

func (suite *TestUserSuite) TestCreateUser_DuplicateGEWISID() {
	var organ models.Organ
	suite.db.First(&organ)

	params := models.UserCreateRequest{
		Name:    "First User",
		GEWISID: uint(13),
		Organs:  []models.Organ{organ},
	}
	_, err := suite.service.Create(&params)
	assert.NoError(suite.T(), err)

	duplicateParams := models.UserCreateRequest{
		Name:    "Duplicate User",
		GEWISID: uint(13),
		Organs:  []models.Organ{organ},
	}
	user, err := suite.service.Create(&duplicateParams)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
}

func (suite *TestUserSuite) TestCreateUser_EmptyName() {
	var organ models.Organ
	suite.db.First(&organ)

	params := models.UserCreateRequest{
		Name:    "",
		GEWISID: uint(14),
		Organs:  []models.Organ{organ},
	}
	user, err := suite.service.Create(&params)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
}

func (suite *TestUserSuite) TestGetUser_ValidInput() {
	var users []*models.User
	err := suite.db.Find(&users).Error
	assert.NoError(suite.T(), err)

	getUsers, err := suite.service.GetUsers(nil)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), getUsers)
	assert.Equal(suite.T(), len(users), len(getUsers))
}

func (suite *TestUserSuite) TestGetUser_ByID() {
	var organ models.Organ
	err := suite.db.First(&organ).Error
	assert.NoError(suite.T(), err)

	req := &models.UserCreateRequest{
		Name:    "Test User",
		GEWISID: uint(10),
		Organs:  []models.Organ{organ},
	}

	user, err := suite.service.Create(req)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)

	filters := &models.UserFilterParams{ID: &user.ID}
	users, err := suite.service.GetUsers(filters)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), users, 1)
	assert.Equal(suite.T(), user.ID, users[0].ID)
}

func (suite *TestUserSuite) TestGetUser_ByGEWISID() {
	gewisID := uint(456)
	var organ models.Organ
	err := suite.db.First(&organ).Error
	assert.NoError(suite.T(), err)

	req := &models.UserCreateRequest{
		Name:    "Test User",
		GEWISID: gewisID,
		Organs:  []models.Organ{organ},
	}

	user, err := suite.service.Create(req)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)

	filters := &models.UserFilterParams{GEWISID: &gewisID}
	users, err := suite.service.GetUsers(filters)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), users, 1)
	assert.Equal(suite.T(), gewisID, users[0].GEWISID)
}

func (suite *TestUserSuite) TestGetUser_ByOrganID() {
	var organ models.Organ
	err := suite.db.First(&organ).Error
	assert.NoError(suite.T(), err)

	req := &models.UserCreateRequest{
		Name:    "Test User",
		GEWISID: uint(10),
		Organs:  []models.Organ{organ},
	}

	user, err := suite.service.Create(req)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)

	organID := organ.ID

	filters := &models.UserFilterParams{OrganID: &organID}
	users, err := suite.service.GetUsers(filters)

	assert.NoError(suite.T(), err)
	found := false
	for _, organ := range users[0].Organs {
		if organ.ID == organID {
			found = true
			break
		}
	}
	assert.True(suite.T(), found, "Expected user to be linked to organ ID %d", organID)
}

func (suite *TestUserSuite) TestGetUser_NoMatch() {
	nonExistentID := uint(99999)
	filters := &models.UserFilterParams{ID: &nonExistentID}

	users, err := suite.service.GetUsers(filters)

	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), users)
}

func (suite *TestUserSuite) TestDeleteUser_ValidInput() {
	var user models.User
	suite.db.First(&user)

	err := suite.service.Delete(user.ID)
	assert.NoError(suite.T(), err)

	var deletedUser models.User
	result := suite.db.First(&deletedUser, user.ID)
	assert.Equal(suite.T(), gorm.ErrRecordNotFound, result.Error)
}

func (suite *TestUserSuite) TestDeleteUser_NotFound() {
	var user models.User
	suite.db.Last(&user)

	err := suite.service.Delete(user.ID + 1)
	assert.Error(suite.T(), err)
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(TestUserSuite))
}
