package roster

import (
	"GEWIS-Rooster/internal/models"
	"github.com/stretchr/testify/assert"
)

func (suite *TestRosterSuite) TestCreateRosterComment_Valid() {
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

	answer := models.RosterAnswer{
		UserID:        1,
		RosterID:      roster.ID,
		RosterShiftID: shift.ID,
		Value:         "yes",
	}
	suite.db.Create(&answer)

	createParams := &CommentCreateRequest{
		RosterID: roster.ID,
		UserID:   1,
		Comment:  "I can't work this shift",
	}

	comment, err := suite.service.CreateRosterComment(createParams)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), comment)
	assert.Equal(suite.T(), createParams.Comment, comment.Comment)
	assert.Equal(suite.T(), createParams.UserID, comment.UserID)
	assert.Equal(suite.T(), roster.ID, comment.RosterID)
}

func (suite *TestRosterSuite) TestCreateRosterComment_UpdatesExisting() {
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

	suite.db.Create(&models.RosterAnswer{
		UserID:        1,
		RosterID:      roster.ID,
		RosterShiftID: shift.ID,
		Value:         "yes",
	})

	first, err := suite.service.CreateRosterComment(&CommentCreateRequest{
		RosterID: roster.ID,
		UserID:   1,
		Comment:  "first comment",
	})
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), first)

	second, err := suite.service.CreateRosterComment(&CommentCreateRequest{
		RosterID: roster.ID,
		UserID:   1,
		Comment:  "updated comment",
	})
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), second)

	assert.Equal(suite.T(), first.ID, second.ID, "expected the existing comment to be updated, not duplicated")
	assert.Equal(suite.T(), "updated comment", second.Comment)

	comments, err := suite.service.GetRosterComments(roster.ID)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), comments, 1, "expected only one comment to remain for this user+roster")
}

func (suite *TestRosterSuite) TestCreateRosterComment_UserNotInOrgan() {
	organ := models.Organ{Name: "Unlinked Organ"}
	suite.db.Create(&organ)

	roster := models.Roster{
		Name:    "Test Roster",
		Values:  []string{"yes", "no"},
		OrganID: organ.ID,
	}
	suite.db.Create(&roster)

	createParams := &CommentCreateRequest{
		RosterID: roster.ID,
		UserID:   1,
		Comment:  "I can't work this shift",
	}

	comment, err := suite.service.CreateRosterComment(createParams)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), comment)
	assert.Contains(suite.T(), err.Error(), "not part of this organ")
}

func (suite *TestRosterSuite) TestCreateRosterComment_RosterNotFound() {
	createParams := &CommentCreateRequest{
		RosterID: 9999,
		UserID:   1,
		Comment:  "I can't work this shift",
	}

	comment, err := suite.service.CreateRosterComment(createParams)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), comment)
	assert.Contains(suite.T(), err.Error(), "roster not found")
}

func (suite *TestRosterSuite) TestGetRosterComments_FiltersByRoster() {
	rosterOne := models.Roster{
		Name:    "Roster One",
		Values:  []string{"yes", "no"},
		OrganID: uint(1),
	}
	suite.db.Create(&rosterOne)

	rosterTwo := models.Roster{
		Name:    "Roster Two",
		Values:  []string{"yes", "no"},
		OrganID: uint(1),
	}
	suite.db.Create(&rosterTwo)

	_, err := suite.service.CreateRosterComment(&CommentCreateRequest{RosterID: rosterOne.ID, UserID: 1, Comment: "comment one"})
	assert.NoError(suite.T(), err)
	_, err = suite.service.CreateRosterComment(&CommentCreateRequest{RosterID: rosterTwo.ID, UserID: 1, Comment: "comment two"})
	assert.NoError(suite.T(), err)

	comments, err := suite.service.GetRosterComments(rosterOne.ID)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), comments, 1)
	assert.Equal(suite.T(), "comment one", comments[0].Comment)
}
