package routers_user

import (
	"net/http"

	"github.com/dinkelspiel/cdn/dao"
	"github.com/dinkelspiel/cdn/db"
	"github.com/dinkelspiel/cdn/middleware"
	"github.com/dinkelspiel/cdn/models"
	"github.com/gin-gonic/gin"
)

func UserTeamsRouter(v1 *gin.RouterGroup, db *db.DB) {
	user := v1.Group("/user")
	user.Use(middleware.AuthMiddleware(db))
	user.GET("/teams", func(c *gin.Context) {
		authUser, _ := c.MustGet("authUser").(models.User)

		teams, err := dao.GetTeamsByOwner(db, authUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		teamList := []gin.H{}

		for _, team := range *teams {
			teamList = append(teamList, models.SerializeTeam(team))
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Found teams",
			"teams":   teamList,
		})
	})
}
