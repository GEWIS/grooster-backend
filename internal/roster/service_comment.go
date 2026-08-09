package roster

import (
	"GEWIS-Rooster/internal/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommentManager interface {
	CreateRosterComment(*CommentCreateRequest) (*models.RosterComment, error)
	GetRosterComments(rosterID uint) ([]*models.RosterComment, error)
}

func (s *service) CreateRosterComment(params *CommentCreateRequest) (*models.RosterComment, error) {
	var roster *models.Roster
	if err := s.db.First(&roster, params.RosterID).Error; err != nil {
		return nil, fmt.Errorf("roster not found: %w", err)
	}

	var answer models.RosterAnswer
	if err := s.db.Where("user_id = ? AND roster_id = ?", params.UserID, params.RosterID).First(&answer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user is not part of this roster")
		}
		return nil, err
	}

	comment := models.RosterComment{
		RosterID: roster.ID,
		UserID:   params.UserID,
		Comment:  params.Comment,
	}

	// A user may only have one comment per roster, so re-submitting updates
	// the existing comment in place instead of creating a duplicate.
	err := s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "roster_id"}, {Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"comment", "updated_at"}),
	}).Create(&comment).Error
	if err != nil {
		return nil, err
	}

	if err := s.db.Where("roster_id = ? AND user_id = ?", roster.ID, params.UserID).First(&comment).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}

func (s *service) GetRosterComments(rosterID uint) ([]*models.RosterComment, error) {
	var comments []*models.RosterComment
	if err := s.db.Where("roster_id = ?", rosterID).Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}
