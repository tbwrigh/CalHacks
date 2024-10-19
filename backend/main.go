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

	router.Run(":8080")
}
