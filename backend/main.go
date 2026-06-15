package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"ocean-survey-backend/config"
	"ocean-survey-backend/database"
	"ocean-survey-backend/middleware"
	"ocean-survey-backend/router"
	"ocean-survey-backend/utils"
	"ocean-survey-backend/websocket"
)

var DB *gorm.DB
var RDB *redis.Client

func main() {
	config.Load()

	initDB()
	initRedis()

	hub := websocket.NewHub()
	go hub.Run()

	r := gin.Default()

	r.Use(middleware.CORS())

	r.GET("/health", func(c *gin.Context) {
		utils.Success(c, gin.H{"status": "ok"})
	})

	router.RegisterRoutes(r, DB, hub)

	addr := fmt.Sprintf(":%s", config.AppConfig.ServerPort)
	log.Printf("Server starting on port %s", config.AppConfig.ServerPort)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initDB() {
	var err error
	DB, err = database.InitPostgres(config.AppConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connection established")
}

func initRedis() {
	var err error
	RDB, err = database.InitRedis(config.AppConfig)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Redis connection established")
}
