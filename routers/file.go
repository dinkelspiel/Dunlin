package routers

import (
	"database/sql"
	"net/http"

	"github.com/dinkelspiel/cdn/dao"
	"github.com/dinkelspiel/cdn/services"
	"github.com/gin-gonic/gin"
)

func FileRouter(r *gin.RouterGroup, db *sql.DB) {
	r.GET("/files/:teamSlug/:projectSlug/*filepath", func(c *gin.Context) {
		teamSlug := c.Param("teamSlug")

		team, err := dao.GetTeamBySlug(db, teamSlug)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if team == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No team found with slug"})
			return
		}

		teamProjectSlug := c.Param("projectSlug")

		teamProject, err := dao.GetTeamProjectInTeamBySlug(db, *team, teamProjectSlug)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if teamProject == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No team project found with slug in team"})
			return
		}

		filePath := c.Param("filepath")

		path, err := services.GetFilePathToFileInTeamProject(*teamProject, filePath)

		c.Header("Content-Description", "File Transfer")
		c.File(path)
	})
}
