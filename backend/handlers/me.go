package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MeHandler handles the /me endpoint
func MeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":  "John Doe",
		"email": "john@example.com",
	})
}
