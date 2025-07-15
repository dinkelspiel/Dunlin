package routers

import (
	"math/rand/v2"
	"net/http"
	"strconv"

	"github.com/dinkelspiel/cdn/dao"
	"github.com/dinkelspiel/cdn/db"
	"github.com/dinkelspiel/cdn/middleware"
	"github.com/dinkelspiel/cdn/models"
	"github.com/dinkelspiel/cdn/services"
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

func AuthRouter(v1 *gin.RouterGroup, db *db.DB) {
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
		authCode := models.UserAuthCode{
			Code:   code,
			UserId: *user.Id,
			Used:   false,
		}
		dao.CreateUserAuthCode(db, authCode)

		// TODO: Send email here

		c.JSON(http.StatusCreated, gin.H{
			"message": "If an account exists for the email address you provided, a verification code has been sent. Please check your inbox (and spam folder) for further instructions. " + strconv.Itoa(authCode.Code),
		})
	})

	auth.POST("/verify-code", func(c *gin.Context) {
		var body AuthVerifyCodeBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		code, err := dao.GetUnusedUserAuthCodeByCode(db, body.Code)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if code == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No unused code found"})
			return
		}

		err = dao.UpdateUserAuthCodeToUsed(db, *code)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userSession, err := dao.CreateUserSession(db, *code.User)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		config, _ := services.LoadConfig()

		var secure bool
		if config.AppUrl == "localhost" {
			secure = false
		} else {
			secure = true
		}

		c.SetCookie("cdn-session-token", userSession.Token, 999999, "/", config.AppUrl, secure, true)

		c.JSON(http.StatusCreated, gin.H{
			"message": "Log in successful",
			"token":   userSession.Token,
		})
	})

	auth.GET("/check-session", middleware.AuthMiddleware(db), func(c *gin.Context) {
		authUser, _ := c.MustGet("authUser").(models.User)

		c.JSON(http.StatusOK, gin.H{
			"message":  "Check Successful",
			"id":       authUser.Id,
			"username": authUser.Username,
			"email":    authUser.Email,
		})
	})
}
