package handlers

import (
	"calhacks/api/db"
	"calhacks/api/lib"
	"calhacks/api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListIssues(c *gin.Context) {
	var req ScanRequest

	// Bind the JSON body to the struct
	if err := c.BindJSON(&req); err != nil {
		// If there is an error, return a 400 Bad Request response
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	token := c.GetHeader("Authorization")

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	// Remove the "Bearer " prefix from the token
	token = token[len("Bearer "):]

	var repository models.Repo
	result := db.DB.Where("owner = ? AND name = ?", req.Owner, req.Repo).First(&repository)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	if repository.RescanSecurity {
		c.JSON(http.StatusConflict, gin.H{"error": "Scan in progress"})
		go func() {
			lib.GetOrUpdateSecurityIssues(req.Owner, req.Repo, token)
		}()
		return
	}

	var issues []models.SecurityIssue
	result = db.DB.Where("repository_id = ?", repository.ID).Find(&issues)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	for i := range issues {
		issues[i].Repository = repository
	}

	c.JSON(http.StatusOK, issues)
}
