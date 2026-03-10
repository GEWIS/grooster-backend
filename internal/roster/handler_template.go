package roster

import (
	"GEWIS-Rooster/internal/models"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func (h *Handler) registerTemplateRoutes(g *gin.RouterGroup, db *gorm.DB) {
	templateGroup := g.Group("/template")
	{
		templateGroup.POST("", requireRosterOrganRoleBody(db, models.RoleAdmin), h.CreateRosterTemplate)
		templateGroup.GET("", h.GetRosterTemplates)
		templateGroup.GET("/:id", h.GetRosterTemplate)
		templateGroup.PUT("/:id", requireTemplateOrganRoleParam(db, "id", models.RoleAdmin), h.UpdateRosterTemplate)
		templateGroup.DELETE("/:id", requireTemplateOrganRoleParam(db, "id", models.RoleAdmin), h.DeleteRosterTemplate)

		templateGroup.PATCH("/shift/:id", h.UpdateRosterTemplateShift)

		templateGroup.POST("/shift-preference", h.CreateRosterTemplateShiftPreference)
		templateGroup.GET("/shift-preference", h.GetRosterTemplateShiftPreferences)
		templateGroup.PATCH("/shift-preference/:id", h.UpdateRosterTemplateShiftPreference)
	}
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

// UpdateRosterTemplateShift
//
//	@Summary   Updates a roster template shift by ID
//	@Security  BearerAuth
//	@Tags      Roster
//	@Accept    json
//	@Produce   json
//	@Param     id     path      int    true   "Shift ID"
//	@Param     params body      TemplateShiftUpdateRequest  true "Update params"
//	@Success   200    {object}   models.RosterTemplateShift
//	@Failure   400    {string}   string
//	@Failure   404    {string}   string
//	@ID        updateRosterTemplateShift
//	@Router    /roster/template/shift/{id} [patch]
func (h *Handler) UpdateRosterTemplateShift(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updateParams TemplateShiftUpdateRequest
	if err := c.ShouldBindJSON(&updateParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	shift, err := h.rosterService.UpdateRosterTemplateShift(uint(id64), &updateParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shift)
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

// requireRosterOrganRoleParam validates the existence of a roster by its ID
// from the URL parameters and ensures the current user has the required
// minimum role within that roster's organization.
func requireTemplateOrganRoleParam(db *gorm.DB, paramStr string, minRole models.OrganRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		templateID := c.Param(paramStr)
		if templateID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": paramStr + " is required"})
			return
		}

		var template models.RosterTemplate
		if err := db.First(&template, "id = ?", templateID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Roster not found"})
			return
		}

		checkAccess(c, db, template.OrganID, minRole)
	}
}
