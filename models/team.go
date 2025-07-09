package models

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Team struct {
	Id        *int64
	Name      string
	Slug      string
	OwnerId   int64
	Owner     *User
	UpdatedAt *time.Time
	CreatedAt *time.Time
}

func SerializeTeam(team Team) gin.H {
	return gin.H{
		"id":        team.Id,
		"name":      team.Name,
		"slug":      team.Slug,
		"ownerId":   team.OwnerId,
		"owner":     SerializeUser(*team.Owner),
		"updatedAt": team.UpdatedAt,
		"createdAt": team.CreatedAt,
	}
}
