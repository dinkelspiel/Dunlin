package routers

import (
	"net/http"

	"github.com/dinkelspiel/cdn/db"
	"github.com/dinkelspiel/cdn/middleware"
	"github.com/dinkelspiel/cdn/models"
	"github.com/dinkelspiel/cdn/services"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

type CreateTeamBody struct {
	TeamName string `json:"teamName" binding:"required"`
}

func TeamsRouter(v1 *gin.RouterGroup, db *db.DB) {
	v1.POST("/teams", middleware.AuthMiddleware(db), func(c *gin.Context) {
		authUser, _ := c.MustGet("authUser").(models.User)
		var body CreateTeamBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		team := models.Team{
			Name:    body.TeamName,
			Slug:    slug.Make(body.TeamName),
			Owner:   &authUser,
			OwnerId: *authUser.Id,
		}

		_, err := services.CreateTeam(db, team)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Created Team",
		})
	})
}
