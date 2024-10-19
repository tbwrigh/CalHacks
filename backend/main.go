package main

import (
	"calhacks/api/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Use handlers for the endpoints
	router.GET("/me", handlers.MeHandler)
	router.GET("/pr", handlers.PrHandler)

	router.Run(":8080")
}
