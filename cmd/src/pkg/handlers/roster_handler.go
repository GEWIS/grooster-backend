package handlers

import (
	"GEWIS-Rooster/cmd/src/pkg/models"
	"GEWIS-Rooster/cmd/src/pkg/services"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type RosterHandler struct {
	rosterService services.RosterServiceInterface
}

func NewRosterHandler(rosterService services.RosterServiceInterface, rg *gin.RouterGroup) *RosterHandler {
	h := &RosterHandler{rosterService: rosterService}

	g := rg.Group("/roster")

	g.POST("", h.CreateRoster)
	g.GET("", h.GetRosters)
	g.GET(":id", h.GetRoster)
	g.PATCH("/:id", h.UpdateRoster)
	g.DELETE("/:id", h.DeleteRoster)

	g.POST("/shift", h.CreateRosterShift)
	g.PATCH("/shift", h.UpdateRosterShift)
	g.DELETE("/shift/:id", h.DeleteRosterShift)
	g.POST("/answer", h.CreateRosterAnswer)
	g.PATCH("/answer/:id", h.UpdateRosterAnswer)

	g.POST("/:id/save", h.SaveRoster)
	g.PATCH("/saved-shift/:id", h.UpdateSavedShift)
	g.GET("/saved-shift/:id", h.GetSavedRoster)

	g.POST("/template", h.CreateRosterTemplate)
	g.GET("/template", h.GetRosterTemplates)
	g.GET("/template/:id", h.GetRosterTemplate)
	g.PUT("/template/:id", h.UpdateRosterTemplate)
	g.DELETE("/template/:id", h.DeleteRosterTemplate)

	return h
}

// CreateRoster
//
//	@Summary	CreateRoster a new roster
//	@Security	BearerAuth
//	@Tags		Roster
//	@Accept		json
//	@Produce	json
//	@Param		createParams	body		models.RosterCreateRequest	true	"Roster input"
//	@Success	200				{object}	models.Roster
//	@Failure	400				{string}	string
//	@ID			createRoster
//	@Router		/roster [post]
func (h *RosterHandler) CreateRoster(c *gin.Context) {
	var param *models.RosterCreateRequest

	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	createdRoster, err := h.rosterService.CreateRoster(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdRoster)
}

// GetRosters
//
//	@Summary	Get all rosters or query by date and organ
//	@Security	BearerAuth
//	@Tags		Roster
//	@Accept		json
//	@Produce	json
//	@Param		date	query		string	false	"Date filter (ISO format)"
//	@Param		organId	query		uint	false	"Organ ID filter"
//	@Success	200		{array}		models.Roster
//	@Failure	400		{string}	string
//	@ID			getRosters
//	@Router		/roster [get]
func (h *RosterHandler) GetRosters(c *gin.Context) {
	var params models.RosterFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	rosters, err := h.rosterService.GetRosters(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, rosters)
}

// GetRoster
//
//	@Summary	Get a specific roster by id
//	@Security	BearerAuth
//	@Tags		Roster
//	@Accept		json
//	@Produce	json
//
//	@Param		id	path		uint	true	"Roster ID"
//
//	@Success	200	{object}	models.Roster
//	@Failure	400	{string}	string
//	@Failure	404	{string}	string
//	@ID			getRoster
//	@Router		/roster/{id} [get]
func (h *RosterHandler) GetRoster(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	id := uint(id64)

	params := &models.RosterFilterParams{
		ID: &id,
	}
	roster, err := h.rosterService.GetRosters(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, roster)
}

// UpdateRoster
//
//	@Summary	Update a roster
//	@Security	BearerAuth
//	@Tags		Roster
//	@Accept		json
//	@Produce	json
//	@Param		id				path		uint						true	"Roster ID"
//	@Param		updateParams	body		models.RosterUpdateRequest	true	"Roster input"
//	@Success	200				{object}	models.Roster
//	@Failure	400				{string}	string
//	@ID			updateRoster
//	@Router		/roster/{id} [patch]
func (h *RosterHandler) UpdateRoster(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid roster ID"})
		return
	}

	var updateParams models.RosterUpdateRequest
	if err := c.ShouldBindJSON(&updateParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updatedRoster, err := h.rosterService.UpdateRoster(uint(id), &updateParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedRoster)
}

// DeleteRoster
//
//	@Summary	DeleteRoster a roster
//	@Security	BearerAuth
//	@Tags		Roster
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"Roster ID"
//	@Success	200	{string}	string
//	@Failure	400	{string}	string
//	@Failure	404	{string}	string
//	@ID			deleteRoster
//	@Router		/roster/{id} [delete]
func (h *RosterHandler) DeleteRoster(c *gin.Context) {
	rosterId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid roster ID"})
		return
	}

	if err := h.rosterService.DeleteRoster(uint(rosterId)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Roster deleted",
	})
}

// CreateRosterShift
//
//	@Summary	Create a new roster shift
//	@Security	BearerAuth
//	@Tags		Roster Shift
//	@Accept		json
//	@Produce	json
//	@Param		createParams	body		models.RosterShiftCreateRequest	true	"Roster shift input"
//	@Success	200				{object}	models.RosterShift
//	@Failure	400				{string}	string
//	@ID			createRosterShift
//	@Router		/roster/shift [post]
func (h *RosterHandler) CreateRosterShift(c *gin.Context) {
	var param *models.RosterShiftCreateRequest

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
//		@Param     updateParams   body      models.RosterShiftUpdateRequest    true   "Update input"
//		@Success   200            {object}  models.RosterShift
//		@Failure   400            {string}  string
//	 	@Failure 404 {string} string
//		@ID        updateRosterShift
//		@Router    /roster/shift/{id} [patch]
func (h *RosterHandler) UpdateRosterShift(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shift ID"})
		return
	}

	var param models.RosterShiftUpdateRequest
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
func (h *RosterHandler) DeleteRosterShift(c *gin.Context) {
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
//	@Param		createParams	body		models.RosterAnswerCreateRequest	true	"Roster answer input"
//	@Success	200				{object}	models.RosterAnswer
//	@Failure	400				{string}	string
//	@ID			createRosterAnswer
//	@Router		/roster/answer [post]
func (h *RosterHandler) CreateRosterAnswer(c *gin.Context) {
	var param *models.RosterAnswerCreateRequest

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
//	@Param		updateParams	body		models.RosterAnswerUpdateRequest	true	"New answer value"
//	@Success	200				{object}	models.RosterAnswer
//	@Failure	400				{string}	json
//	@Failure	404				{string}	json
//	@ID			updateRosterAnswer
//	@Router		/roster/answer/{id} [patch]
func (h *RosterHandler) UpdateRosterAnswer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid roster answer ID"})
		return
	}

	var updateParams models.RosterAnswerUpdateRequest
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

// SaveRoster
//
//	@Summary	Save a specific roster
//	@Security	BearerAuth
//	@Tags		Saved Shift
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"Roster ID"
//	@Success	200	{string}	json
//	@Failure	400	{string}	json
//	@Failure	404	{string}	json
//	@ID			rosterSave
//	@Router		/roster/{id}/save [post]
func (h *RosterHandler) SaveRoster(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid roster ID"})
		return
	}

	err = h.rosterService.SaveRoster(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "SavedShift not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update shift"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// UpdateSavedShift
//
//	@Summary	Update a specific saved shift
//	@Security	BearerAuth
//	@Tags		Saved Shift
//	@Accept		json
//	@Produce	json
//	@Param		id				path		int								true	"SavedShift ID"
//	@Param		updateParams	body		models.SavedShiftUpdateRequest	true	"Update data"
//	@Success	200				{object}	models.SavedShift
//	@Failure	400				{string}	string	"Invalid request"
//	@Failure	404				{string}	string	"SavedShift not found"
//	@ID			updateSavedShift
//	@Router		/roster/saved-shift/{id} [patch]
func (h *RosterHandler) UpdateSavedShift(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid saved shift ID"})
		return
	}

	var updateParams models.SavedShiftUpdateRequest
	if err := c.ShouldBindJSON(&updateParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	updatedShift, err := h.rosterService.UpdateSavedShift(uint(id), &updateParams)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "SavedShift not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update shift"})
		}
		return
	}

	c.JSON(http.StatusOK, updatedShift)
}

// GetSavedRoster
//
//	@Summary	Get all saved shifts for a specific roster
//	@Security	BearerAuth
//	@Tags		Saved Shift
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int					true	"Roster ID"
//	@Success	200	{object}		models.SavedShiftResponse	"Saved Shift Response"
//	@Failure	400	{string}	string				"Invalid request"
//	@Failure	404	{string}	string				"SavedShift not found"
//	@ID			getSavedRoster
//	@Router		/roster/saved-shift/{id} [get]
func (h *RosterHandler) GetSavedRoster(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid roster ID"})
		return
	}

	savedShifts, savedShiftOrdering, err := h.rosterService.GetSavedRoster(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := models.SavedShiftResponse{
		SavedShifts:        savedShifts,
		SavedShiftOrdering: savedShiftOrdering,
	}
	// Log the entire struct as a field called "response"
	log.Debug().Interface("response", response).Msg("Sending saved roster response")

	c.JSON(http.StatusOK, response)
}

// CreateRosterTemplate
//
//	@Summary	Creates a template of a roster by defining the name of the shifts
//	@Security	BearerAuth
//	@Tags		Roster
//	@Accept		json
//	@Produce	json
//	@Param		params	body		models.RosterTemplateCreateRequest					false	"Template Params"
//	@Success	200	{array}		models.RosterTemplate	"Created Template"
//	@Failure	400	{string}	string				"Invalid request"
//	@ID			createRosterTemplate
//	@Router		/roster/template [post]
func (h *RosterHandler) CreateRosterTemplate(c *gin.Context) {
	var param *models.RosterTemplateCreateRequest

	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	createdTemplate, err := h.rosterService.CreateRosterTemplate(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdTemplate)
}

// GetRosterTemplates
//
//	@Summary	Get all rosters templates or query by organ ID
//	@Security	BearerAuth
//	@Tags		Roster
//	@Accept		json
//	@Produce	json
//	@Param		params	query		models.RosterTemplateFilterParams	false	"Date filter (ISO format)"
//	@Success	200		{array}		models.RosterTemplate
//	@Failure	400		{string}	string
//	@ID			getRosterTemplates
//	@Router		/roster/template [get]
func (h *RosterHandler) GetRosterTemplates(c *gin.Context) {
	var params models.RosterTemplateFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	rosterTemplates, err := h.rosterService.GetRosterTemplates(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, rosterTemplates)
}

// GetRosterTemplate
//
//	@Summary	Get a roster template by ID
//	@Security	BearerAuth
//	@Tags		Roster
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"Template ID"
//	@Success	200		{object}		models.RosterTemplate
//	@Failure 	404 	{string} 	string
//	@Failure	400		{string}	string
//	@ID			getRosterTemplate
//	@Router		/roster/template/{id} [get]
func (h *RosterHandler) GetRosterTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	id := uint(id64)

	rosterTemplate, err := h.rosterService.GetRosterTemplate(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, rosterTemplate)
}

// UpdateRosterTemplate
//
//	@Summary	Updates a roster template by ID
//	@Security	BearerAuth
//	@Tags		Roster
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"Template ID"
//	@Param		params	body		models.RosterTemplateUpdateParams	false "Update params"
//	@Success	200		{object}		models.RosterTemplate
//	@Failure	400		{string}	string
//	@Failure 	404 	{string} 	string
//	@ID			updateRosterTemplate
//	@Router		/roster/template/{id} [put]
func (h *RosterHandler) UpdateRosterTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	id := uint(id64)

	var updateParams models.RosterTemplateUpdateParams
	if err := c.ShouldBindJSON(&updateParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	updatedTemplate, err := h.rosterService.UpdateRosterTemplate(id, &updateParams)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Roster template not found"})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedTemplate)
}

// DeleteRosterTemplate
//
//	@Summary	Deletes a roster template by ID
//	@Security	BearerAuth
//	@Tags		Roster
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"Template ID"
//	@Success	200		{string}	string
//	@Failure 	404 	{string} 	string
//	@Failure	400		{string}	string
//	@ID			deleteRosterTemplate
//	@Router		/roster/template/{id} [delete]
func (h *RosterHandler) DeleteRosterTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	id := uint(id64)

	err = h.rosterService.DeleteRosterTemplate(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}
