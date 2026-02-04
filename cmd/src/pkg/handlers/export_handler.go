package handlers

import (
	"GEWIS-Rooster/cmd/src/pkg/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ExportHandler struct {
	exportService services.ExportServiceInterface
}

func NewExportHandler(exportService services.ExportServiceInterface, rg *gin.RouterGroup) *ExportHandler {
	h := &ExportHandler{exportService: exportService}

	g := rg.Group("/export")

	g.GET("/roster/:id", h.AssignmentToPng)

	return h
}

// AssignmentToPng
//
// @Summary      Export roster assignments as PNG
// @Description  Generates and downloads a PNG image containing the shift assignments for a specific roster.
// @Security     BearerAuth
// @Tags         Export
// @Produce      image/png
// @Param        id   path      uint  true  "Roster ID"
// @Success      200  {file}    binary
// @Failure      400  {object}  map[string]string "Invalid ID format"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /export/roster/{id} [get]
func (h *ExportHandler) AssignmentToPng(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 32) // Use 32 for uint compatibility
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	imgBytes, err := h.exportService.AssignmentsToPng(uint(id64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fileName := fmt.Sprintf("roster-%d.png", id64)

	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Data(http.StatusOK, "image/png", imgBytes)
}
