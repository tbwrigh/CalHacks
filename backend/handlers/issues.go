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

type FixRequest struct {
	ID uint `json:"id"`
}

func GetFixIssueHandler(c *gin.Context) {
	var req FixRequest

	// Bind the JSON body to the struct
	if err := c.BindJSON(&req); err != nil {
		// If there is an error, return a 400 Bad Request response
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	var issue models.SecurityIssue
	result := db.DB.Where("id = ?", req.ID).First(&issue)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Issue not found"})
		return
	}

	token := c.GetHeader("Authorization")

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	// Remove the "Bearer " prefix from the token
	token = token[len("Bearer "):]

	go func() {
		var repository models.Repo
		result := db.DB.Where("id = ?", issue.RepositoryID).First(&repository)

		if result.Error != nil {
			return
		}

		code, err := lib.GetCode(repository.Owner, repository.Name, issue.Path, token)

		if err != nil {
			return
		}

		newCode, err := lib.FixSecurity(issue.StartLine, issue.EndLine, issue.FullDescription, code)

		if err != nil {
			return
		}

		fileChange := lib.FileChange{
			Path:    issue.Path,
			Content: newCode,
			Message: "Auto-fix security issue",
			NewFile: false,
		}

		err = lib.CreatePR(repository.Owner, repository.Name, []lib.FileChange{fileChange}, token)

		if err != nil {
			return
		}

		issue.FixSuggested = true
		db.DB.Save(&issue)
	}()

	c.JSON(http.StatusOK, issue)
}
