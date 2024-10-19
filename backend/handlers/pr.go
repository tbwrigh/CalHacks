package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PrHandler handles the /pr endpoint
func PrHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"title":       "PR #123",
		"description": "Fixes a bug in the application.",
		"status":      "Open",
	})
}
