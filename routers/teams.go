package routers

import (
	"database/sql"
	"net/http"

	"github.com/dinkelspiel/cdn/dao"
	"github.com/dinkelspiel/cdn/models"
	"github.com/gin-gonic/gin"
)

func TeamsRouter(v1 *gin.RouterGroup, db *sql.DB) {
	// teams := v1.Group("/teams")
	v1.GET("/teams", func(c *gin.Context) {
		slug := c.Query("slug")

		team, err := dao.GetTeamBySlug(db, slug)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if team == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No team found with slug"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Found team",
			"team":    models.SerializeTeam(*team),
		})
	})
}
