package routers

import (
	"net/http"

	"github.com/dinkelspiel/cdn/db"
	"github.com/dinkelspiel/cdn/models"
	"github.com/dinkelspiel/cdn/services"
	"github.com/gin-gonic/gin"
)

func TeamsRouter(v1 *gin.RouterGroup, db *db.DB) {
	v1.GET("/teams", func(c *gin.Context) {
		teamSlug := c.Query("slug")

		team, err := services.GetTeamBySlug(db, teamSlug)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Found team",
			"team":    models.SerializeTeam(*team),
		})
	})
}
