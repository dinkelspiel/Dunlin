package models

import (
	"time"

	"github.com/gin-gonic/gin"
)

type TeamProject struct {
	Id        *int64
	Name      string
	Slug      string
	TeamId    int64
	Team      *Team
	UpdatedAt *time.Time
	CreatedAt *time.Time
}

func SerializeTeamProject(teamProject TeamProject) gin.H {
	return gin.H{
		"id":        teamProject.Id,
		"name":      teamProject.Name,
		"slug":      teamProject.Slug,
		"teamId":    teamProject.TeamId,
		"team":      SerializeTeam(*teamProject.Team),
		"updatedAt": teamProject.UpdatedAt,
		"createdAt": teamProject.CreatedAt,
	}
}
