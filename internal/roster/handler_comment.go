package roster

import (
	"GEWIS-Rooster/internal/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func (h *Handler) registerCommentRoutes(g *gin.RouterGroup, db *gorm.DB) {
	commentGroup := g.Group("/comment")
	{
		commentGroup.POST("", requireCommentOrganRoleBody(db, models.RoleMember), h.CreateRosterComment)
		commentGroup.GET("", requireCommentOrganRoleQuery(db, "rosterId", models.RoleMember), h.GetRosterComments)
	}
}

// CreateRosterComment
//
//	@Summary	Create a new roster comment
//	@Security	BearerAuth
//	@Tags		Roster Comment
//	@Accept		json
//	@Produce	json
//	@Param		createParams	body		CommentCreateRequest	true	"Roster comment input"
//	@Success	201				{object}	models.RosterComment
//	@Failure	400				{string}	string
//	@Failure	403				{string}	string
//	@Failure	404				{string}	string
//	@ID			createRosterComment
//	@Router		/roster/comment [post]
func (h *Handler) CreateRosterComment(c *gin.Context) {
	var param *CommentCreateRequest

	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	callerID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if callerID.(uint) != param.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only comment on your own roster entry"})
		return
	}

	createdComment, err := h.rosterService.CreateRosterComment(param)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Roster not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdComment)
}

// GetRosterComments
//
//	@Summary	Get all comments for a roster
//	@Security	BearerAuth
//	@Tags		Roster Comment
//	@Accept		json
//	@Produce	json
//	@Param		rosterId	query		int	true	"Roster ID"
//	@Success	200			{array}		models.RosterComment
//	@Failure	400			{string}	string
//	@Failure	404			{string}	string
//	@ID			getRosterComments
//	@Router		/roster/comment [get]
func (h *Handler) GetRosterComments(c *gin.Context) {
	rosterID, err := strconv.ParseUint(c.Query("rosterId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid roster ID"})
		return
	}

	comments, err := h.rosterService.GetRosterComments(uint(rosterID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// requireCommentOrganRoleBody validates the existence of a roster referenced
// by rosterId in the JSON body and ensures the current user has the required
// minimum role within that roster's organization.
func requireCommentOrganRoleBody(db *gorm.DB, minRole models.OrganRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			RosterID uint `json:"rosterId"`
		}
		if err := c.ShouldBindBodyWith(&body, binding.JSON); err != nil || body.RosterID == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Valid rosterId is required in body"})
			return
		}

		var roster models.Roster
		if err := db.First(&roster, "id = ?", body.RosterID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Roster not found"})
			return
		}
		checkAccess(c, db, roster.OrganID, minRole)
	}
}

// requireCommentOrganRoleQuery validates the existence of a roster referenced
// by the given query parameter and ensures the current user has the required
// minimum role within that roster's organization.
func requireCommentOrganRoleQuery(db *gorm.DB, queryStr string, minRole models.OrganRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		rosterID := c.Query(queryStr)
		if rosterID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": queryStr + " is required"})
			return
		}

		var roster models.Roster
		if err := db.First(&roster, "id = ?", rosterID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Roster not found"})
			return
		}

		checkAccess(c, db, roster.OrganID, minRole)
	}
}
