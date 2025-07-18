package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/dinkelspiel/cdn/db"
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

func setupRouter(db *db.DB) *gin.Engine {
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
		AllowOrigins:     []string{"https://files.keii.dev/", "https://github.com", "http://localhost:5173"},
		MaxAge:           12 * time.Hour,
	}))

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

	r.Static("/assets", "./frontend/dist/assets")
	r.StaticFile("/favicon.ico", "./frontend/dist/favicon.ico")

	r.NoRoute(func(c *gin.Context) {
		indexPath := filepath.Join("./frontend/dist", "index.html")

		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			c.String(http.StatusNotFound, "index.html not found")
			return
		}

		c.File(indexPath)
	})

	return r
}

func main() {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Println("Running as:", u.Username)

	// Load Config
	config, err := services.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config: ", err)
		return
	}

	services.EnsureFoldersExists()

	// Load MariaDB Database
	dsn := config.MariaDBDatabaseUrl
	mariaDBClient, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening mariadb: ", err)
		return
	}
	defer mariaDBClient.Close()
	if err = mariaDBClient.Ping(); err != nil {
		log.Fatal("Error pinging mariadb: ", err)
		return
	}

	// Load Redis Database
	// redisClient := redis.NewClient(&redis.Options{
	// 	Addr: config.RedisDatabaseUrl,
	// })
	// defer redisClient.Close()

	// ctx := context.Background()
	// if err = redisClient.Ping(ctx).Err(); err != nil {
	// 	log.Fatal("Error pinging redis: ", err)
	// 	return
	// }

	db := db.DB{
		MariaDB: mariaDBClient,
		// Redis:   redisClient,
	}

	r := setupRouter(&db)

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
