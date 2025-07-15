package routers

import (
	"net/http"

	"github.com/dinkelspiel/cdn/dao"
	"github.com/dinkelspiel/cdn/db"
	"github.com/dinkelspiel/cdn/models"
	"github.com/dinkelspiel/cdn/services"
	"github.com/gin-gonic/gin"
)

type SetupBody struct {
	AdminUsername string `json:"adminUsername" binding:"required"`
	AdminEmail    string `json:"adminEmail" binding:"required"`
}

func SetupRouter(v1 *gin.RouterGroup, db *db.DB) {
	v1.POST("/setup", func(c *gin.Context) {
		users_count, err := dao.GetAmountOfUsers(db)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if *users_count > 0 {
			c.AbortWithStatus(404)
			return
		}

		var body SetupBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user := models.User{
			Username: body.AdminUsername,
			Email:    body.AdminEmail,
		}
		_, err = services.RegisterUser(db, user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "OpenCDN successfully initialized.",
		})
	})
}
