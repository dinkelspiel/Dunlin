package middleware

import (
	"net/http"
	"strings"

	"github.com/dinkelspiel/cdn/dao"
	"github.com/dinkelspiel/cdn/db"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(db *db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		cookie, err := c.Cookie("cdn-session-token")
		if err == nil {
			token = cookie
		} else {
			authHeader := c.GetHeader("Authorization")
			if len(authHeader) > 7 && strings.ToLower(authHeader[:7]) == "bearer " {
				token = authHeader[7:]
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Missing session token in cookie or Authorization header"})
				c.Abort()
				return
			}
		}

		userSession, err := dao.GetUserSessionByToken(db, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if userSession == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No user session found with token"})
			c.Abort()
			return
		}

		c.Set("authUser", *userSession.User)

		c.Next()
	}
}
