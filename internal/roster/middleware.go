package roster

import (
	"GEWIS-Rooster/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func checkAccess(c *gin.Context, db *gorm.DB, organID interface{}, minRole models.OrganRole) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var userOrgan models.UserOrgan

	if err := db.Where("user_id = ? AND organ_id = ?", userID, organID).First(&userOrgan).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not a member of this organ"})
		return
	}

	if models.RoleWeights[userOrgan.Role] < models.RoleWeights[minRole] {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient organ permissions"})
		return
	}

	c.Next()
}
