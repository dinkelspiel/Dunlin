package routers

import (
	"net/http"

	"github.com/dinkelspiel/cdn/db"
	"github.com/dinkelspiel/cdn/middleware"
	"github.com/dinkelspiel/cdn/services"
	"github.com/gin-gonic/gin"
)

func Map[T any, U any](input []T, mapper func(T) U) []U {
	result := make([]U, len(input))
	for i, v := range input {
		result[i] = mapper(v)
	}
	return result
}

func StatisticsRouter(v1 *gin.RouterGroup, db *db.DB) {
	statistics := v1.Group("/statistics")

	statistics.GET("", middleware.AuthMiddleware(db), func(c *gin.Context) {
		diskStats := services.GetHostDiskStats()
		teamProjectsSizes, err := services.GetSizeOfTeamProjectsOnDisk(db)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Retrieved statistics successfully",
			"data": gin.H{
				"disk": services.SerializeDiskStats(diskStats),
				"teamProjectSizes": Map(*teamProjectsSizes, func(teamProjectsSize services.TeamProjectSize) gin.H {
					return services.SerializeTeamProjectSize(teamProjectsSize)
				}),
			},
		})
	})
}
