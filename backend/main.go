package main

import (
	"calhacks/api/db"
	"calhacks/api/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database
	db.Init()

	router := gin.Default()

	// Use handlers for the endpoints
	router.GET("/me", handlers.MeHandler)
	router.GET("/pr", handlers.PrHandler)
	router.GET("/repo", handlers.RepoHandler)
	router.POST("/scan/start", handlers.StartScanHandler)
	router.POST("/scan/status", handlers.GetScanStatusHandler)
	router.POST("/scan/results", handlers.GetScanResultsHandler)

	router.Run(":8080")
}
