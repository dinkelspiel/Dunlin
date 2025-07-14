package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/dinkelspiel/cdn/routers"
	routers_user "github.com/dinkelspiel/cdn/routers/user"
	"github.com/dinkelspiel/cdn/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

type PostStatisticBody struct {
	ServerId     *int64  `json:"serverId"`
	ServerSecret *string `json:"serverSecret"`

	PlayerCount       int    `json:"playerCount"`
	GameVersion       string `json:"gameVersion"`
	ServerEnvironment string `json:"serverEnvironment"`
	OperatingSystem   string `json:"operatingSystem"`
	Arch              string `json:"arch"`
	JavaVersion       string `json:"javaVersion"`
}

func setupRouter(db *sql.DB) *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.HandleMethodNotAllowed = true

	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "HEAD", "POST", "PUT", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOrigins:     []string{"https://cdn.keii.dev/", "https://github.com", "http://localhost:5173"},
		MaxAge:           12 * time.Hour,
	}))

	r.StaticFile("/", "./static/index.html")

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	api := r.Group("/api")
	v1 := api.Group("/v1")

	// This route serves all the files
	routers.FileRouter(r.Group("/"), db)

	routers.AuthRouter(v1, db)
	routers.SetupRouter(v1, db)
	routers.TeamsRouter(v1, db)
	routers.TeamRouter(v1, db)
	routers.TeamProjectRouter(v1, db)
	routers_user.UserTeamsRouter(v1, db)

	return r
}

func main() {
	// Load Config
	config, err := services.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config: ", err)
		return
	}

	services.EnsureFoldersExist()

	// Load Database
	dsn := config.DatabaseUrl
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening database: ", err)
		return
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		log.Fatal("Error pinging database: ", err)
		return
	}

	r := setupRouter(db)

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
