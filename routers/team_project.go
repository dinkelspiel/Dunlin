package routers

import (
	"database/sql"
	"net/http"

	"github.com/dinkelspiel/cdn/dao"
	"github.com/dinkelspiel/cdn/models"
	"github.com/dinkelspiel/cdn/services"
	"github.com/dinkelspiel/cdn/storage"
	"github.com/gin-gonic/gin"
)

func TeamProjectRouter(v1 *gin.RouterGroup, db *sql.DB) {
	team := v1.Group("/teams/:teamSlug/projects/:projectSlug")
	team.GET("", func(c *gin.Context) {
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

		c.JSON(http.StatusOK, gin.H{
			"message":     "Found team project",
			"teamProject": models.SerializeTeamProject(*teamProject),
		})
	})

	team.GET("/files/*filepath", func(c *gin.Context) {
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

		filepath := c.Param("filepath")

		files, err := services.GetTeamProjectFiles(*teamProject, filepath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fileList := []gin.H{}

		for _, file := range files {
			fileList = append(fileList, storage.SerializeFSItem(file))
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Found files",
			"files":   fileList,
		})
	})
}
