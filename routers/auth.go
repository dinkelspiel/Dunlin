package routers

import (
	"database/sql"
	"math/rand/v2"
	"net/http"
	"strconv"

	"github.com/dinkelspiel/cdn/dao"
	"github.com/dinkelspiel/cdn/models"
	"github.com/gin-gonic/gin"
)

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

type AuthSendCodeBody struct {
	Email string `json:"email" binding:"required"`
}

type AuthVerifyCodeBody struct {
	Code int64 `json:"code" binding:"required"`
}

func AuthRouter(v1 *gin.RouterGroup, db *sql.DB) {
	auth := v1.Group("/auth")
	auth.POST("/send-code", func(c *gin.Context) {
		var body AuthSendCodeBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := dao.GetUserByEmail(db, body.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if user == nil {
			c.JSON(http.StatusCreated, gin.H{
				"message": "If an account exists for the email address you provided, a verification code has been sent. Please check your inbox (and spam folder) for further instructions.",
			})
			return
		}

		code := randRange(10000, 99999)
		auth_code := models.UserAuthCode{
			Code:   code,
			UserId: *user.Id,
			Used:   false,
		}
		dao.CreateUserAuthCode(db, auth_code)

		// TODO: Send email here

		c.JSON(http.StatusCreated, gin.H{
			"message": "If an account exists for the email address you provided, a verification code has been sent. Please check your inbox (and spam folder) for further instructions. " + strconv.Itoa(auth_code.Code),
		})
	})

	auth.POST("/verify-code", func(c *gin.Context) {
		var body AuthVerifyCodeBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		code, err := dao.GetUserAuthCodeByCode(db, body.Code)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = dao.UpdateUserAuthCodeToUsed(db, *code)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user_session, err := dao.CreateUserSession(db, *code.User)

		c.SetCookie("cdn-session-token", user_session.Token, 999999, "/", "localhost", false, true)

		c.JSON(http.StatusCreated, gin.H{
			"message": "Log in successful",
			"token":   user_session.Token,
		})
	})

	auth.GET("/check-session", func(c *gin.Context) {
		cookie, err := c.Cookie("cdn-session-token")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user_session, err := dao.GetUserSessionByToken(db, cookie)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "Check Successful",
			"id":       user_session.User.Id,
			"username": user_session.User.Username,
			"email":    user_session.User.Email,
		})
	})
}
