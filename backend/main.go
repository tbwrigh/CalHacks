package main

import (
	"calhacks/api/db"
	"calhacks/api/handlers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database
	db.Init()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Allow your frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Use handlers for the endpoints
	router.GET("/me", handlers.MeHandler)
	router.GET("/pr", handlers.PrHandler)
	router.GET("/repo", handlers.RepoHandler)
	router.POST("/scan/start", handlers.StartScanHandler)
	router.POST("/scan/status", handlers.GetScanStatusHandler)
	router.POST("/scan/results", handlers.GetScanResultsHandler)
	router.POST("/install", handlers.GetInstallHandler)

	router.Run(":8080")
}
