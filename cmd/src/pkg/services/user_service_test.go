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
	var user models.User
	suite.db.First(&user)

	getUser, err := suite.service.GetUser(user.GEWISID)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), getUser)
	assert.Equal(suite.T(), user.Name, getUser.Name)
}

func (suite *TestUserSuite) TestGetUser_NotFound() {
	var user models.User
	suite.db.Last(&user)

	getUser, err := suite.service.GetUser(user.GEWISID + 1)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), getUser)
}

func (suite *TestUserSuite) TestDeleteUser_ValidInput() {
	var user models.User
	suite.db.First(&user)

	err := suite.service.Delete(user.ID)
	assert.NoError(suite.T(), err)

	var deletedUser models.User
	result := suite.db.First(&deletedUser, user.ID)
	assert.Error(suite.T(), result.Error)
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
