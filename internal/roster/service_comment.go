package roster

import (
	"GEWIS-Rooster/internal/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
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

	if err := s.db.Create(&comment).Error; err != nil {
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
