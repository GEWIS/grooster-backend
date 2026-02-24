package roster

import (
	_ "GEWIS-Rooster/internal/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type Handler struct {
	rosterService Service
}

func NewRosterHandler(rosterService Service, rg *gin.RouterGroup) *Handler {
	h := &Handler{rosterService: rosterService}

	g := rg.Group("/roster")

	g.POST("", h.CreateRoster)
	g.GET("", h.GetRosters)
	g.GET(":id", h.GetRoster)
	g.PATCH("/:id", h.UpdateRoster)
	g.DELETE("/:id", h.DeleteRoster)

	g.POST("/shift", h.CreateRosterShift)
	g.PATCH("/shift/:id", h.UpdateRosterShift)
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

	g.POST("/template/shift-preference", h.CreateRosterTemplateShiftPreference)
	g.GET("/template/shift-preference", h.GetRosterTemplateShiftPreferences)
	g.PATCH("/template/shift-preference/{id}", h.UpdateRosterTemplateShiftPreference)

	return h
}

// CreateRoster
//
//	@Summary	CreateRoster a new roster
//	@Security	BearerAuth
//	@Tags		Roster
//	@Accept		json
//	@Produce	json
//	@Param		createParams	body		RosterCreateRequest	true	"Roster input"
//	@Success	200				{object}	models.Roster
//	@Failure	400				{string}	string
//	@ID			createRoster
//	@Router		/roster [post]
func (h *Handler) CreateRoster(c *gin.Context) {
	var param *CreateRequest

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
//		@Summary	Get all rosters or query by date and organ
//		@Security	BearerAuth
//		@Tags		Roster
//		@Accept		json
//		@Produce	json
//	 @Param     filter  query     RosterFilterParams  false  "Filter parameters"
//		@Success	200		{array}		models.Roster
//		@Failure	400		{string}	string
//		@ID			getRosters
//		@Router		/roster [get]
func (h *Handler) GetRosters(c *gin.Context) {
	var params FilterParams
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
func (h *Handler) GetRoster(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	id := uint(id64)

	params := &FilterParams{
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
//	@Param		updateParams	body		RosterUpdateRequest	true	"Roster input"
//	@Success	200				{object}	models.Roster
//	@Failure	400				{string}	string
//	@ID			updateRoster
//	@Router		/roster/{id} [patch]
func (h *Handler) UpdateRoster(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid roster ID"})
		return
	}

	var updateParams UpdateRequest
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
func (h *Handler) DeleteRoster(c *gin.Context) {
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
func (h *Handler) SaveRoster(c *gin.Context) {
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
//	@Param		updateParams	body		SavedShiftUpdateRequest	true	"Update data"
//	@Success	200				{object}	models.SavedShift
//	@Failure	400				{string}	string	"Invalid request"
//	@Failure	404				{string}	string	"SavedShift not found"
//	@ID			updateSavedShift
//	@Router		/roster/saved-shift/{id} [patch]
func (h *Handler) UpdateSavedShift(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid saved shift ID"})
		return
	}

	var updateParams SavedShiftUpdateRequest
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
//	@Success	200	{object}		SavedShiftResponse	"Saved Shift Response"
//	@Failure	400	{string}	string				"Invalid request"
//	@Failure	404	{string}	string				"SavedShift not found"
//	@ID			getSavedRoster
//	@Router		/roster/saved-shift/{id} [get]
func (h *Handler) GetSavedRoster(c *gin.Context) {
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

	response := SavedShiftResponse{
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
//	@Param		params	body		TemplateCreateRequest					false	"Template Params"
//	@Success	200	{array}		models.RosterTemplate	"Created Template"
//	@Failure	400	{string}	string				"Invalid request"
//	@ID			createRosterTemplate
//	@Router		/roster/template [post]
func (h *Handler) CreateRosterTemplate(c *gin.Context) {
	var param *TemplateCreateRequest

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
//	@Param		params	query		TemplateFilterParams	false	"Date filter (ISO format)"
//	@Success	200		{array}		models.RosterTemplate
//	@Failure	400		{string}	string
//	@ID			getRosterTemplates
//	@Router		/roster/template [get]
func (h *Handler) GetRosterTemplates(c *gin.Context) {
	var params TemplateFilterParams
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
func (h *Handler) GetRosterTemplate(c *gin.Context) {
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
//	@Param		params	body		TemplateUpdateParams	false "Update params"
//	@Success	200		{object}		models.RosterTemplate
//	@Failure	400		{string}	string
//	@Failure 	404 	{string} 	string
//	@ID			updateRosterTemplate
//	@Router		/roster/template/{id} [put]
func (h *Handler) UpdateRosterTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	id := uint(id64)

	var updateParams TemplateUpdateParams
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
func (h *Handler) DeleteRosterTemplate(c *gin.Context) {
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

// CreateRosterTemplateShiftPreference
//
//	@Summary   Creates a roster template shift preference
//	@Security  BearerAuth
//	@Tags      Roster
//	@Accept    json
//	@Produce   json
//	@Param     params body      TemplateShiftPreferenceCreateRequest  true "Creation params"
//	@Success   201       {object}      models.RosterTemplateShiftPreference
//	@Failure   400       {string}   string
//	@ID        createRosterTemplateShiftPreference
//	@Router    /roster/template/shift-preference [post]
func (h *Handler) CreateRosterTemplateShiftPreference(c *gin.Context) {
	var params TemplateShiftPreferenceCreateRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	preference, err := h.rosterService.CreateRosterTemplateShiftPreference(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, preference)
}

// GetRosterTemplateShiftPreferences
//
//	@Summary   Gets shift preferences filtered by user and template
//	@Security  BearerAuth
//	@Tags      Roster
//	@Accept    json
//	@Produce   json
//	@Param     userId      query    int     true  "User ID"
//	@Param     templateId  query    int     true  "Template ID"
//	@Success   200         {array}  models.RosterTemplateShiftPreference
//	@Failure   400         {string} string
//	@ID        getRosterTemplateShiftPreferences
//	@Router    /roster/template/shift-preference [get]
func (h *Handler) GetRosterTemplateShiftPreferences(c *gin.Context) {
	var params TemplateShiftPreferenceFilterParams

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters: " + err.Error()})
		return
	}

	preferences, err := h.rosterService.GetRosterTemplateShiftPreferences(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch preferences: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, preferences)
}

// UpdateRosterTemplateShiftPreference
//
//	@Summary   Updates a roster template shift preference by ID
//	@Security  BearerAuth
//	@Tags      Roster
//	@Accept    json
//	@Produce   json
//	@Param     id     path      int    true   "Preference ID"
//	@Param     params body      TemplateShiftPreferenceUpdateRequest  true "Update params"
//	@Success   200       {object}      models.RosterTemplateShiftPreference
//	@Failure   400       {string}   string
//	@Failure   404       {string}   string
//	@ID        updateRosterTemplateShiftPreference
//	@Router    /roster/template/shift-preference/{id} [patch]
func (h *Handler) UpdateRosterTemplateShiftPreference(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	id := uint(id64)

	var params TemplateShiftPreferenceUpdateRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	updatedPreference, err := h.rosterService.UpdateRosterTemplateShiftPreference(id, params)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Shift preference not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedPreference)
}
