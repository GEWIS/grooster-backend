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

func (h *Handler) registerShiftRoutes(g *gin.RouterGroup, db *gorm.DB) {
	shiftGroup := g.Group("/shift")
	{
		shiftGroup.POST("", requireShiftOrganRoleBody(db, models.RoleAdmin), h.CreateRosterShift)
		shiftGroup.PATCH("/:id", requireShiftOrganRoleParam(db, "id", models.RoleAdmin), h.UpdateRosterShift)
		shiftGroup.DELETE("/:id", requireShiftOrganRoleParam(db, "id", models.RoleAdmin), h.DeleteRosterShift)
	}

	answerGroup := g.Group("/answer")
	{
		answerGroup.POST("", requireShiftOrganRoleBody(db, models.RoleMember), h.CreateRosterAnswer)
		answerGroup.PATCH("/:id", requireShiftAnswerOrganRoleParam(db, "id", models.RoleMember), h.UpdateRosterAnswer)
	}

}

// CreateRosterShift
//
//	@Summary	Create a new roster shift
//	@Security	BearerAuth
//	@Tags		Roster Shift
//	@Accept		json
//	@Produce	json
//	@Param		createParams	body		ShiftCreateRequest	true	"Roster shift input"
//	@Success	200				{object}	models.RosterShift
//	@Failure	400				{string}	string
//	@ID			createRosterShift
//	@Router		/roster/shift [post]
func (h *Handler) CreateRosterShift(c *gin.Context) {
	var param *ShiftCreateRequest

	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	createdRosterShift, err := h.rosterService.CreateRosterShift(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdRosterShift)
}

// UpdateRosterShift
//
//		@Summary   Update a roster shift
//		@Security  BearerAuth
//		@Tags      Roster Shift
//		@Accept    json
//		@Produce   json
//		@Param     id             path      int                               true   "Roster Shift ID"
//		@Param     updateParams   body      ShiftUpdateRequest    true   "Update input"
//		@Success   200            {object}  models.RosterShift
//		@Failure   400            {string}  string
//	 	@Failure 404 {string} string
//		@ID        updateRosterShift
//		@Router    /roster/shift/{id} [patch]
func (h *Handler) UpdateRosterShift(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shift ID"})
		return
	}

	var param ShiftUpdateRequest
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updatedShift, err := h.rosterService.UpdateRosterShift(uint(id), &param)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Roster shift not found"})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedShift)
}

// DeleteRosterShift
//
//	@Summary	Deletes a roster shift
//	@Security	BearerAuth
//	@Tags		Roster Shift
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"Roster Answer ID"
//	@Success	200	{string}	string
//	@Failure	400	{string}	string
//	@ID			deleteRosterShift
//	@Router		/roster/shift/{id} [delete]
func (h *Handler) DeleteRosterShift(c *gin.Context) {
	rosterId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid roster ID"})
		return
	}

	if err := h.rosterService.DeleteRosterShift(uint(rosterId)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Roster deleted",
	})
}

// CreateRosterAnswer
//
//	@Summary	Create a new roster shift answer
//	@Security	BearerAuth
//	@Tags		Roster Answer
//	@Accept		json
//	@Produce	json
//	@Param		createParams	body		AnswerCreateRequest	true	"Roster answer input"
//	@Success	200				{object}	models.RosterAnswer
//	@Failure	400				{string}	string
//	@ID			createRosterAnswer
//	@Router		/roster/answer [post]
func (h *Handler) CreateRosterAnswer(c *gin.Context) {
	var param *AnswerCreateRequest

	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	createdAnswer, err := h.rosterService.CreateRosterAnswer(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdAnswer)
}

// UpdateRosterAnswer
//
//	@Summary	Updates a roster answer with the new value
//	@Security	BearerAuth
//	@Tags		Roster Answer
//	@Accept		json
//	@Produce	json
//	@Param		id				path		int									true	"Roster Answer ID"
//	@Param		updateParams	body		AnswerUpdateRequest	true	"New answer value"
//	@Success	200				{object}	models.RosterAnswer
//	@Failure	400				{string}	json
//	@Failure	404				{string}	json
//	@ID			updateRosterAnswer
//	@Router		/roster/answer/{id} [patch]
func (h *Handler) UpdateRosterAnswer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid roster answer ID"})
		return
	}

	var updateParams AnswerUpdateRequest
	if err := c.ShouldBindJSON(&updateParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updatedAnswer, err := h.rosterService.UpdateRosterAnswer(uint(id), &updateParams)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedAnswer)
}

// requireRosterOrganRoleParam validates the existence of a roster by its ID
// from the URL parameters and ensures the current user has the required
// minimum role within that roster's organization.
func requireShiftOrganRoleParam(db *gorm.DB, paramStr string, minRole models.OrganRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		shiftID := c.Param(paramStr)
		if shiftID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": paramStr + " is required"})
			return
		}

		var shift models.RosterShift
		if err := db.Preload("Roster").First(&shift, "id = ?", shiftID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Roster not found"})
			return
		}

		checkAccess(c, db, shift.Roster.OrganID, minRole)
	}
}

// requireShiftAnswerOrganRoleParam validates the existence of a roster by its ID
// from the URL parameters and ensures the current user has the required
// minimum role within that roster's organization.
func requireShiftAnswerOrganRoleParam(db *gorm.DB, paramStr string, minRole models.OrganRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		answerID := c.Param(paramStr)
		if answerID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": paramStr + " is required"})
			return
		}

		var answer models.RosterAnswer
		if err := db.Preload("Roster").First(&answer, "id = ?", answerID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Answer not found"})
			return
		}

		checkAccess(c, db, answer.Roster.OrganID, minRole)
	}
}

func requireShiftOrganRoleBody(db *gorm.DB, minRole models.OrganRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			RosterID uint `json:"rosterID"`
		}
		if err := c.ShouldBindBodyWith(&body, binding.JSON); err != nil || body.RosterID == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Valid rosterID is required in body"})
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
