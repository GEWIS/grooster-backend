package roster

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func (h *Handler) registerShiftRoutes(g *gin.RouterGroup) {
	shiftGroup := g.Group("/shift")
	{
		shiftGroup.POST("", h.CreateRosterShift)
		shiftGroup.PATCH("/:id", h.UpdateRosterShift)
		shiftGroup.DELETE("/:id", h.DeleteRosterShift)
	}

	answerGroup := shiftGroup.Group("/answer")
	{
		answerGroup.POST("", h.CreateRosterAnswer)
		answerGroup.PATCH("/:id", h.UpdateRosterAnswer)
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

	if err := c.ShouldBindJSON(&param); err != nil {
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
