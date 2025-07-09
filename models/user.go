package models

import (
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	Id        *int64
	Username  string
	Email     string
	UpdatedAt *time.Time
	CreatedAt *time.Time
}

func SerializeUser(user User) gin.H {
	return gin.H{
		"id":        user.Id,
		"username":  user.Username,
		"email":     user.Email,
		"updatedAt": user.UpdatedAt,
		"createdAt": user.CreatedAt,
	}
}
