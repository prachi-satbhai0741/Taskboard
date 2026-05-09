package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prachi-satbhai0741/Taskboard/db"
	"github.com/prachi-satbhai0741/Taskboard/handlers"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("  No .env file found, reading from system environment")
	}

	db.Connect()

	r := gin.Default()

	r.GET("/health", handlers.HealthCheck)
	r.GET("/tasks", handlers.GetTasks)
	r.POST("/tasks", handlers.CreateTask)
	r.PUT("/tasks/:id", handlers.UpdateTask)
	r.DELETE("/tasks/:id", handlers.DeleteTask)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf(" Taskboard running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf(" Server failed to start: %v", err)
	}
}
// phase 3 test
