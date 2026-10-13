package organ

import (
	"GEWIS-Rooster/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func requireRosterOrganMemberRoleParam(db *gorm.DB, minRole models.OrganRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		organID := c.Param("id")
		if organID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id" + " is required"})
			return
		}

		val, exists := c.Get("userID")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		userID, ok := val.(uint)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
			return
		}

		var organ models.Organ
		if err := db.First(&organ, "id = ?", organID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Roster not found"})
			return
		}

		var user models.User
		if err := db.First(&user, "id = ?", userID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
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
}
